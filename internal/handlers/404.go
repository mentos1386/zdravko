package handlers

import (
	"net/http"

	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

func (h *BaseHandler) Error404(c echo.Context) error {
	return c.Render(http.StatusNotFound, "404.tmpl", &components.Base{
		NavbarActive: nil,
		Navbar:       Pages,
	})
}
