package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

type CheckHistoryWithCheck struct {
	*models.CheckHistory
	CheckName string `db:"check_name"`
	CheckId   string `db:"check_id"`
}

func GetLastNCheckHistory(ctx context.Context, db *sqlx.DB, n int) ([]*CheckHistoryWithCheck, error) {
	var checkHistory []*CheckHistoryWithCheck
	err := db.SelectContext(ctx, &checkHistory, `
    SELECT
      mh.*,
      wg.name AS worker_group_name,
      m.name AS check_name,
      m.id AS check_id
    FROM check_histories mh
      LEFT JOIN worker_groups wg ON mh.worker_group_id = wg.id
      LEFT JOIN check_worker_groups mwg ON mh.check_id = mwg.check_id
      LEFT JOIN checks m ON mwg.check_id = m.id
    ORDER BY mh.created_at DESC
    LIMIT $1
    `, n)
	return checkHistory, err
}

func GetCheckHistoryForCheck(ctx context.Context, db *sqlx.DB, checkId string) ([]*models.CheckHistory, error) {
	var checkHistory []*models.CheckHistory
	err := db.SelectContext(ctx, &checkHistory, `
  SELECT
    mh.*,
    wg.name AS worker_group_name,
    wg.id AS worker_group_id
  FROM check_histories as mh
    LEFT JOIN worker_groups wg ON mh.worker_group_id = wg.id
    LEFT JOIN check_worker_groups mwg ON mh.check_id = mwg.check_id
  WHERE mh.check_id = $1
  ORDER BY mh.created_at DESC
  `, checkId)
	return checkHistory, err
}

func AddHistoryForCheck(ctx context.Context, db *sqlx.DB, history *models.CheckHistory) error {
	_, err := db.NamedExecContext(ctx,
		`
INSERT INTO check_histories (
  check_id,
  worker_group_id,
  status,
  note
) VALUES (
  :check_id,
  :worker_group_id,
  :status,
  :note
)`,
		history,
	)
	return err
}
