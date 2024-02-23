package handlers

import (
	"net/http"

	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

func (h *BaseHandler) Incidents(c echo.Context) error {
	return c.Render(http.StatusOK, "incidents.tmpl", &components.Base{
		NavbarActive: GetPageByTitle(Pages, "Incidents"),
		Navbar:       Pages,
	})
}
