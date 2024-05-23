package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/internal/server/services"
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
