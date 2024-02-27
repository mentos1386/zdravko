package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

func GetMonitorHistoryForMonitor(ctx context.Context, db *sqlx.DB, monitorSlug string) ([]*models.MonitorHistory, error) {
	var monitorHistory []*models.MonitorHistory
	err := db.SelectContext(ctx, &monitorHistory,
		"SELECT * FROM monitor_histories WHERE monitor_slug = $1 ORDER BY created_at DESC",
		monitorSlug,
	)
	return monitorHistory, err
}

func AddHistoryForMonitor(ctx context.Context, db *sqlx.DB, history *models.MonitorHistory) error {
	_, err := db.NamedExecContext(ctx,
		"INSERT INTO monitor_histories (monitor_slug, status, note) VALUES (:monitor_slug, :status, :note)",
		history,
	)
	return err
}
