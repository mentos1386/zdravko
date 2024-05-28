package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database/models"
)

type TargetHistory struct {
	*models.TargetHistory
	TargetName      string `db:"target_name"`
	WorkerGroupName string `db:"worker_group_name"`
	CheckName       string `db:"check_name"`
}

type TargetHistoryDateRange string

const (
	TargetHistoryDateRange60Days    TargetHistoryDateRange = "60_DAYS"
	TargetHistoryDateRange48Hours   TargetHistoryDateRange = "48_HOURS"
	TargetHistoryDateRange60Minutes TargetHistoryDateRange = "60_MINUTES"
)

func GetTargetHistoryForTarget(ctx context.Context, db *sqlx.DB, targetId string, dateRange TargetHistoryDateRange) ([]*TargetHistory, error) {
	dateRangeFilter := ""
	switch dateRange {
	case TargetHistoryDateRange60Days:
		dateRangeFilter = "AND strftime('%Y-%m-%dT%H:%M:%fZ', th.created_at) >= datetime('now', 'localtime', '-60 days')"
	case TargetHistoryDateRange48Hours:
		dateRangeFilter = "AND strftime('%Y-%m-%dT%H:%M:%fZ', th.created_at) >= datetime('now', 'localtime', '-48 hours')"
	case TargetHistoryDateRange60Minutes:
		dateRangeFilter = "AND strftime('%Y-%m-%dT%H:%M:%fZ', th.created_at) >= datetime('now', 'localtime', '-60 minutes')"
	}

	var targetHistory []*TargetHistory
	err := db.SelectContext(
		ctx,
		&targetHistory,
		fmt.Sprintf(`
    SELECT
      th.*,
      t.name AS target_name,
      wg.name AS worker_group_name,
      c.name AS check_name
    FROM target_histories th
      LEFT JOIN targets t ON th.target_id = t.id
      LEFT JOIN worker_groups wg ON th.worker_group_id = wg.id
      LEFT JOIN checks c ON th.check_id = c.id
    WHERE
      th.target_id = $1
      %s
    ORDER BY th.created_at DESC
  `, dateRangeFilter),
		targetId,
	)
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
