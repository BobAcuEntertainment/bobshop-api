package oauth

import (
	"app/env"
	"app/pkg/ecode"
	"app/pkg/enum"
	"app/pkg/tz"
	"app/store"
	"app/store/db"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Oauth struct {
	store *store.Store
}

func New(store *store.Store) *Oauth {
	return &Oauth{store}
}

func (s Oauth) Timezone(r *http.Request) (tz.Timezone, error) {
	name := r.Header.Get("Timezone")
	if len(name) == 0 {
		name = r.FormValue("timezone")
	}

	if len(name) == 0 {
		return tz.AsiaHoChiMinh, nil
	}

	timezone, err := time.LoadLocation(name)
	if err != nil {
		return tz.AsiaHoChiMinh, err
	}

	return tz.Timezone(timezone.String()), nil
}

func (s Oauth) NoAuth(r *http.Request) (*db.AuthSessionDto, error) {
	timezone, err := s.Timezone(r)
	if err != nil {
		return nil, ecode.InvalidTimezone
	}
	return &db.AuthSessionDto{Timezone: timezone}, nil
}

func (s Oauth) BearerAuth(r *http.Request) (*db.AuthSessionDto, error) {
	var (
		auth   = r.Header.Get("Authorization")
		prefix = "Bearer "
		access = ""
	)

	if len(auth) > 0 && strings.HasPrefix(auth, prefix) {
		access = auth[len(prefix):]
	} else {
		access = r.FormValue("access_token")
	}

	if len(access) == 0 {
		return nil, ecode.Unauthorized
	}

	timezone, err := s.Timezone(r)
	if err != nil {
		return nil, ecode.InvalidTimezone
	}

	return s.ValidateToken(r.Context(), access, timezone)
}

func (s Oauth) BasicAuth(r *http.Request) (*db.AuthSessionDto, error) {
	clientId, clientSecret, ok := r.BasicAuth()
	if !ok {
		return nil, ecode.Unauthorized
	}

	client, err := s.store.GetClient(r.Context(), clientId)
	if err != nil {
		return nil, ecode.Unauthorized.Stack(err)
	}

	if client.ClientSecret != clientSecret {
		return nil, ecode.Unauthorized
	}

	tenant, err := s.store.GetTenant(r.Context(), string(client.TenantId))
	if err != nil {
		return nil, ecode.Unauthorized.Stack(err)
	}

	if tenant.DataStatus == enum.DataStatusDisable {
		return nil, ecode.Forbidden
	}

	timezone, err := s.Timezone(r)
	if err != nil {
		return nil, ecode.InvalidTimezone
	}

	return &db.AuthSessionDto{Username: clientId, TenantId: client.TenantId, Timezone: timezone}, nil
}

func (s Oauth) GenerateToken(ctx context.Context, uid string) (*db.AuthTokenDto, error) {
	client, err := s.store.GetClient(ctx, env.ClientId)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	user, err := s.store.GetUser(ctx, uid)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	if user.DataStatus == enum.DataStatusDisable {
		return nil, ecode.Forbidden
	}

	tenant, err := s.store.GetTenant(ctx, string(user.TenantId))
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	if tenant.DataStatus == enum.DataStatusDisable {
		return nil, ecode.Forbidden
	}

	now := time.Now()
	key := []byte(client.ClientId + client.ClientSecret + client.SecureKey)
	dto := &db.AuthTokenDto{ExpiresIn: int64(time.Hour.Seconds()), TokenType: "Bearer"}

	access := &AccessClaims{
		BaseClaims: BaseClaims{
			Jti:       uuid.NewString(),
			Subject:   user.ID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiredAt: now.Add(time.Hour).Unix(),
		},
		JwtType: JwtTypeAccess,
		Version: user.VersionToken,
	}

	dto.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, access).SignedString(key)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	refresh := &RefreshClaims{
		BaseClaims: BaseClaims{
			Jti:       access.Jti,
			Subject:   user.ID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiredAt: now.Add(time.Hour * 24).Unix(),
		},
		JwtType: JwtTypeRefresh,
		Version: user.VersionToken,
	}

	dto.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refresh).SignedString(key)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	return dto, nil
}

func (s Oauth) RefreshToken(ctx context.Context, refresh string) (*db.AuthTokenDto, error) {
	client, err := s.store.GetClient(ctx, env.ClientId)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	key := []byte(client.ClientId + client.ClientSecret + client.SecureKey)

	claims := &RefreshClaims{}

	if _, err := jwt.ParseWithClaims(refresh, claims, func(token *jwt.Token) (any, error) { return key, nil }); err != nil || claims.JwtType != JwtTypeRefresh {
		return nil, ecode.InvalidToken
	}

	if claims.ExpiredAt < time.Now().Unix() || time.Now().Unix() < claims.NotBefore {
		return nil, ecode.InvalidToken
	}

	user, err := s.store.GetUser(ctx, claims.Subject)
	if err != nil {
		return nil, ecode.InternalServerError.Stack(err)
	}

	if user.VersionToken != claims.Version {
		return nil, ecode.InvalidToken
	}

	if s.IsRevoked(ctx, claims.Jti) {
		return nil, ecode.InvalidToken
	}

	tk, err := s.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	s.store.Rdb.Set(ctx, fmt.Sprintf("jti:%v", claims.Jti), true, time.Duration(claims.ExpiredAt-time.Now().Unix())*time.Second)

	return tk, nil
}

func (s Oauth) ValidateToken(ctx context.Context, access string, timezone tz.Timezone) (*db.AuthSessionDto, error) {
	client, err := s.store.GetClient(ctx, env.ClientId)
	if err != nil {
		return nil, ecode.Unauthorized.Stack(err)
	}

	key := []byte(client.ClientId + client.ClientSecret + client.SecureKey)

	claims := &AccessClaims{}

	if _, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (any, error) { return key, nil }); err != nil || claims.JwtType != JwtTypeAccess {
		return nil, ecode.Unauthorized
	}

	if claims.ExpiredAt < time.Now().Unix() || time.Now().Unix() < claims.NotBefore {
		return nil, ecode.Unauthorized
	}

	user, err := s.store.GetUser(ctx, claims.Subject)
	if err != nil {
		return nil, ecode.Unauthorized.Stack(err)
	}

	if user.VersionToken != claims.Version {
		return nil, ecode.Unauthorized
	}

	if s.IsRevoked(ctx, claims.Jti) {
		return nil, ecode.Unauthorized
	}

	session := &db.AuthSessionDto{
		Name:        user.Name,
		Phone:       user.Phone,
		Email:       user.Email,
		Username:    user.Username,
		UserId:      user.ID,
		TenantId:    user.TenantId,
		Permissions: user.Permissions,
		IsRoot:      user.IsRoot,
		IsTenant:    user.IsTenant,
		Timezone:    timezone,
		AccessToken: access,
	}

	if session.IsRoot {
		session.Permissions = enum.PermissionRootValues()
		if session.IsTenant {
			session.Permissions = enum.PermissionTenantValues()
		}
	}

	return session, nil
}

func (s Oauth) RevokeTokenByUser(ctx context.Context, uid string) error {
	if err := s.store.Db.User.IncrementVersionToken(ctx, uid); err != nil {
		return ecode.InternalServerError.Stack(err)
	}
	return s.store.DelUser(ctx, uid)
}

func (s Oauth) RevokeToken(ctx context.Context, access string) error {
	client, err := s.store.GetClient(ctx, env.ClientId)
	if err != nil {
		return ecode.InternalServerError.Stack(err)
	}

	key := []byte(client.ClientId + client.ClientSecret + client.SecureKey)

	claims := &AccessClaims{}

	if _, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (any, error) { return key, nil }); err != nil || claims.JwtType != JwtTypeAccess {
		return ecode.InvalidToken
	}

	return s.store.Rdb.Set(ctx, fmt.Sprintf("jti:%v", claims.Jti), true, time.Duration(claims.ExpiredAt-time.Now().Unix())*time.Second+time.Hour*23)
}

func (s Oauth) IsRevoked(ctx context.Context, jti string) bool {
	bytes, _ := s.store.Rdb.GetBytes(ctx, fmt.Sprintf("jti:%v", jti))
	return len(bytes) > 0
}
