package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

type MonitorHistoryWithMonitor struct {
	*models.MonitorHistory
	MonitorName string `db:"monitor_name"`
	MonitorId   string `db:"monitor_id"`
}

func GetLastNMonitorHistory(ctx context.Context, db *sqlx.DB, n int) ([]*MonitorHistoryWithMonitor, error) {
	var monitorHistory []*MonitorHistoryWithMonitor
	err := db.SelectContext(ctx, &monitorHistory, `
    SELECT
      mh.*,
      wg.name AS worker_group_name,
      m.name AS monitor_name,
      m.id AS monitor_id
    FROM monitor_histories mh
      LEFT JOIN worker_groups wg ON mh.worker_group_id = wg.id
      LEFT JOIN monitor_worker_groups mwg ON mh.monitor_id = mwg.monitor_id
      LEFT JOIN monitors m ON mwg.monitor_id = m.id
    ORDER BY mh.created_at DESC
    LIMIT $1
    `, n)
	return monitorHistory, err
}

func GetMonitorHistoryForMonitor(ctx context.Context, db *sqlx.DB, monitorId string) ([]*models.MonitorHistory, error) {
	var monitorHistory []*models.MonitorHistory
	err := db.SelectContext(ctx, &monitorHistory, `
  SELECT
    mh.*,
    wg.name AS worker_group_name,
    wg.id AS worker_group_id
  FROM monitor_histories as mh
    LEFT JOIN worker_groups wg ON mh.worker_group_id = wg.id
    LEFT JOIN monitor_worker_groups mwg ON mh.monitor_id = mwg.monitor_id
  WHERE mh.monitor_id = $1
  ORDER BY mh.created_at DESC
  `, monitorId)
	return monitorHistory, err
}

func AddHistoryForMonitor(ctx context.Context, db *sqlx.DB, history *models.MonitorHistory) error {
	_, err := db.NamedExecContext(ctx,
		`
INSERT INTO monitor_histories (
  monitor_id,
  worker_group_id,
  status,
  note
) VALUES (
  :monitor_id,
  :worker_group_id,
  :status,
  :note
)`,
		history,
	)
	return err
}
