package handlers

import (
	"math/rand"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type IndexData struct {
	*components.Base
	HealthChecks []*HealthCheck
}

type HealthCheck struct {
	Domain  string
	Healthy bool
	Uptime  string
	History []bool
}

func newMockHealthCheck(domain string) *HealthCheck {
	randBool := func() bool {
		return rand.Intn(2) == 1
	}

	var history []bool
	for i := 0; i < 90; i++ {
		history = append(history, randBool())
	}

	return &HealthCheck{
		Domain:  domain,
		Healthy: randBool(),
		Uptime:  "100",
		History: history,
	}
}

func (h *BaseHandler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.tmpl", &IndexData{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Status"),
			Navbar:       Pages,
		},
		HealthChecks: []*HealthCheck{
			newMockHealthCheck("example.com"),
			newMockHealthCheck("example.org"),
			newMockHealthCheck("example.net"),
			newMockHealthCheck("foo.example.net"),
		},
	})
}
