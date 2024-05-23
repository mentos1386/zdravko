package services

import (
	"context"

	"github.com/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

type TargetHistory struct {
	*models.TargetHistory
	TargetName string `db:"target_name"`
}

func GetLastNTargetHistory(ctx context.Context, db *sqlx.DB, n int) ([]*TargetHistory, error) {
	var targetHistory []*TargetHistory
	err := db.SelectContext(ctx, &targetHistory, `
    SELECT
      th.*,
      t.name AS target_name
    FROM target_histories th
      LEFT JOIN targets t ON th.target_id = t.id
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
      t.name AS target_name
    FROM target_histories th
      LEFT JOIN targets t ON th.target_id = t.id
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
  status,
  note
) VALUES (
  :target_id,
  :status,
  :note
)`,
		history,
	)
	return err
}
