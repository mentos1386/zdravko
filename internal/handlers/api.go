package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/pkg/api"
	"github.com/labstack/echo/v4"
)

type ApiV1WorkersConnectGETResponse struct {
	Endpoint string `json:"endpoint"`
	Group    string `json:"group"`
}

func (h *BaseHandler) ApiV1WorkersConnectGET(c echo.Context) error {
	ctx := context.Background()
	cc := c.(AuthenticatedContext)

	workerGroup, err := services.GetWorkerGroup(ctx, h.db, cc.Principal.Worker.Group)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token invalid")
		}
		return err
	}

	response := ApiV1WorkersConnectGETResponse{
		Endpoint: h.config.Temporal.ServerHost,
		Group:    workerGroup.Id,
	}

	return c.JSON(http.StatusOK, response)
}

// TODO: Can we instead get this from the Workflow outcome?
//
//	To somehow listen for the outcomes and then store them automatically.
func (h *BaseHandler) ApiV1ChecksHistoryPOST(c echo.Context) error {
	ctx := context.Background()
	id := c.Param("id")

	var body api.ApiV1ChecksHistoryPOSTBody
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return err
	}

	_, err = services.GetCheck(ctx, h.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Check not found")
		}
		return err
	}

	err = services.AddHistoryForCheck(ctx, h.db, &models.CheckHistory{
		CheckId:     id,
		WorkerGroupId: body.WorkerGroupId,
		Status:        body.Status,
		Note:          body.Note,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}
