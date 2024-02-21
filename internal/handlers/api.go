package handlers

import (
	"context"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/pkg/api"
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

	var body api.ApiV1HealthchecksHistoryPOSTBody
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return err
	}

	healthcheck, err := services.GetHealthcheck(ctx, h.query, slug)
	if err != nil {
		return err
	}

	err = h.query.Healthcheck.History.Model(healthcheck).Append(
		&models.HealthcheckHistory{
			Status: body.Status,
			Note:   body.Note,
		})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}
