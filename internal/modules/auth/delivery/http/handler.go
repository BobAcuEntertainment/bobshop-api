package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	validator "github.com/go-playground/validator/v10"

	"bobshop/internal/modules/auth/application"
	"bobshop/internal/modules/auth/delivery/http/dto"
	"bobshop/internal/modules/auth/domain"
	"bobshop/internal/platform/response"
	"bobshop/internal/platform/web"
)

var (
	accessTokenName = "access_token"
)

type AuthHandler struct {
	authService   *application.AuthService
	validate      *validator.Validate
	cookieManager *web.CookieManager
}

func NewAuthHandler(authService *application.AuthService, cookieManager *web.CookieManager) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		validate:      validator.New(),
		cookieManager: cookieManager,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "mismatched fields", err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	err := h.authService.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			response.Conflict(c, "user already exists", err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Created(c, "User created successfully", nil)
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "mismatched fields", err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.BadRequest(c, "invalid fields", err)
		return
	}

	token, err := h.authService.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			response.Unauthorized(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	cookie := h.cookieManager.BuildCookie(accessTokenName, token, h.cookieManager.GetMaxAge())

	http.SetCookie(c.Writer, cookie)

	response.Success(c, http.StatusOK, "Signed in successfully", dto.SignInResponse{
		Email: req.Email,
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	cookie := h.cookieManager.BuildCookie(accessTokenName, "", -1)
	http.SetCookie(c.Writer, cookie)

	response.Success(c, http.StatusOK, "Signed out successfully", nil)
}
