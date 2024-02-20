package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"github.com/labstack/echo/v4"
)

func (h *BaseHandler) Temporal(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   h.config.Temporal.UIHost,
		Scheme: "http",
	})

	originalDirector := proxy.Director

	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		// Add authentication token to be able to access temporal.
		// FIXME: Maybe cache it somehow so we don't generate it on every request?
		token, _ := jwt.NewTokenForUser(
			h.config.Jwt.PrivateKey,
			h.config.Jwt.PublicKey,
			cc.Principal.User.Email,
		)
		r.Header.Add("Authorization", "Bearer "+token)
	}

	proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}
