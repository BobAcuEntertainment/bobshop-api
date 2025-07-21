package web

import (
	"net/http"

	"bobshop/internal/platform/config"
)

type CookieManager struct {
	cfg *config.CookieConfig
}

func NewCookieManager(cfg *config.CookieConfig) *CookieManager {
	return &CookieManager{cfg: cfg}
}

func (c *CookieManager) getSameSite() http.SameSite {
	if c.cfg.SameSite == "none" {
		return http.SameSiteNoneMode
	}
	return http.SameSiteStrictMode
}

func (c *CookieManager) GetMaxAge() int {
	return c.cfg.MaxAge
}

func (c *CookieManager) BuildCookie(name, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		HttpOnly: c.cfg.HttpOnly,
		Secure:   c.cfg.Secure,
		SameSite: c.getSameSite(),
	}
}
