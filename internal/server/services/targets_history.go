package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database/models"
)

type TargetHistory struct {
	*models.TargetHistory
	TargetName      string `db:"target_name"`
	WorkerGroupName string `db:"worker_group_name"`
	CheckName       string `db:"check_name"`
}

func GetLastNTargetHistory(ctx context.Context, db *sqlx.DB, n int) ([]*TargetHistory, error) {
	var targetHistory []*TargetHistory
	err := db.SelectContext(ctx, &targetHistory, `
    SELECT
      th.*,
      t.name AS target_name,
      wg.name AS worker_group_name,
      c.name AS check_name
    FROM target_histories th
      LEFT JOIN targets t ON th.target_id = t.id
      LEFT JOIN worker_groups wg ON th.worker_group_id = wg.id
      LEFT JOIN checks c ON th.check_id = c.id
    WHERE th.target_id = $1
    ORDER BY th.created_at DESC
    LIMIT $1
    `, n)
	return targetHistory, err
}

func GetTargetHistoryForTarget(ctx context.Context, db *sqlx.DB, targetId string) ([]*TargetHistory, error) {
	var targetHistory []*TargetHistory
	err := db.SelectContext(ctx, &targetHistory, `
    SELECT
      th.*,
      t.name AS target_name,
      wg.name AS worker_group_name,
      c.name AS check_name
    FROM target_histories th
      LEFT JOIN targets t ON th.target_id = t.id
      LEFT JOIN worker_groups wg ON th.worker_group_id = wg.id
      LEFT JOIN checks c ON th.check_id = c.id
    WHERE th.target_id = $1
    ORDER BY th.created_at DESC
  `, targetId)
	return targetHistory, err
}

func AddHistoryForTarget(ctx context.Context, db *sqlx.DB, history *models.TargetHistory) error {
	_, err := db.NamedExecContext(ctx,
		`
INSERT INTO target_histories (
  target_id,
  worker_group_id,
  check_id,
  status,
  note
) VALUES (
  :target_id,
  :worker_group_id,
  :check_id,
  :status,
  :note
)`,
		history,
	)
	return err
}
