package handlers

import (
	"context"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"github.com/labstack/echo/v4"
)

type ApiV1WorkersConnectGETResponse struct {
	Endpoint string `json:"endpoint"`
	Group    string `json:"group"`
	Slug     string `json:"slug"`
}

func (h *BaseHandler) ApiV1WorkersConnectGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	response := ApiV1WorkersConnectGETResponse{
		Endpoint: h.config.Temporal.ServerHost,
		Group:    cc.Principal.Worker.Group,
		Slug:     cc.Principal.Worker.Slug,
	}

	return c.JSON(http.StatusOK, response)
}

// TODO: Can we instead get this from the Workflow outcome?
//
//	To somehow listen for the outcomes and then store them automatically.
func (h *BaseHandler) ApiV1HealthchecksHistoryPOST(c echo.Context) error {
	ctx := context.Background()

	slug := c.Param("slug")

	healthcheck, err := services.GetHealthcheckHttp(ctx, h.query, slug)
	if err != nil {
		return err
	}

	err = h.query.HealthcheckHttp.History.Model(healthcheck).Append(
		&models.HealthcheckHttpHistory{
			Status: "UP",
		})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}
