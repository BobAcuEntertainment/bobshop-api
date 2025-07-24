package infrastructure

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"bobshop/internal/platform/config"
)

type JwtTokenizer struct {
	cfg *config.JWTConfig
}

func NewJwtTokenizer(cfg *config.JWTConfig) *JwtTokenizer {
	return &JwtTokenizer{cfg: cfg}
}

func (j *JwtTokenizer) GenerateToken(userID string, role string) (string, error) {
	exp, err := time.ParseDuration(j.cfg.ExpirationHours)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(exp).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.Secret))
}

func (j *JwtTokenizer) ParseToken(token string) (map[string]any, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(j.cfg.Secret), nil
	})
	return claims, err
}
