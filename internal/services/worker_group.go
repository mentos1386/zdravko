package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"golang.org/x/exp/maps"
)

func CountWorkerGroups(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM worker_groups")
	return count, err
}

func GetActiveWorkers(ctx context.Context, workerGroupId string, temporal client.Client) ([]string, error) {
	response, err := temporal.DescribeTaskQueue(ctx, workerGroupId, enums.TASK_QUEUE_TYPE_ACTIVITY)
	if err != nil {
		return make([]string, 0), err
	}

	workers := make([]string, len(response.Pollers))
	for i, poller := range response.Pollers {
		workers[i] = poller.Identity
	}

	return workers, nil
}

func CreateWorkerGroup(ctx context.Context, db *sqlx.DB, workerGroup *models.WorkerGroup) error {
	_, err := db.NamedExecContext(ctx,
		"INSERT INTO worker_groups (id, name) VALUES (:id, :name)",
		workerGroup,
	)
	return err
}

func DeleteWorkerGroup(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM worker_groups WHERE id = $1",
		id,
	)
	return err
}

func GetWorkerGroups(ctx context.Context, db *sqlx.DB) ([]*models.WorkerGroup, error) {
	var workerGroups []*models.WorkerGroup
	err := db.SelectContext(ctx, &workerGroups,
		"SELECT * FROM worker_groups ORDER BY name",
	)
	return workerGroups, err
}

func GetWorkerGroupsWithMonitors(ctx context.Context, db *sqlx.DB) ([]*models.WorkerGroupWithMonitors, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.id,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  monitors.name as monitor_name
FROM worker_groups
LEFT OUTER JOIN monitor_worker_groups ON worker_groups.id = monitor_worker_groups.worker_group_id
LEFT OUTER JOIN monitors ON monitor_worker_groups.monitor_id = monitors.id
ORDER BY worker_groups.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workerGroups := map[string]*models.WorkerGroupWithMonitors{}

	for rows.Next() {
		workerGroup := &models.WorkerGroupWithMonitors{}

		var monitorName *string
		err = rows.Scan(
			&workerGroup.Id,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&monitorName,
		)
		if err != nil {
			return nil, err
		}

		if monitorName != nil {
			monitors := []string{}
			if workerGroups[workerGroup.Id] != nil {
				monitors = workerGroups[workerGroup.Id].Monitors
			}
			workerGroup.Monitors = append(monitors, *monitorName)
		}

		workerGroups[workerGroup.Id] = workerGroup
	}

	return maps.Values(workerGroups), err
}

func GetWorkerGroupsById(ctx context.Context, db *sqlx.DB, ids []string) ([]*models.WorkerGroup, error) {
	var workerGroups []*models.WorkerGroup
	err := db.SelectContext(ctx, &workerGroups,
		"SELECT * FROM worker_groups WHERE id = ANY($1)",
		ids,
	)
	return workerGroups, err
}

func GetWorkerGroup(ctx context.Context, db *sqlx.DB, id string) (*models.WorkerGroup, error) {
	var workerGroup models.WorkerGroup
	err := db.GetContext(ctx, &workerGroup,
		"SELECT * FROM worker_groups WHERE id = $1",
		id,
	)
	return &workerGroup, err
}

func GetWorkerGroupWithMonitors(ctx context.Context, db *sqlx.DB, id string) (*models.WorkerGroupWithMonitors, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.id,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  monitors.name as monitor_name
FROM worker_groups
LEFT OUTER JOIN monitor_worker_groups ON worker_groups.id = monitor_worker_groups.worker_group_id
LEFT OUTER JOIN monitors ON monitor_worker_groups.monitor_id = monitors.id
WHERE worker_groups.id=$1
`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workerGroup := &models.WorkerGroupWithMonitors{}

	for rows.Next() {
		var monitorName *string
		err = rows.Scan(
			&workerGroup.Id,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&monitorName,
		)
		if err != nil {
			return nil, err
		}
		if monitorName != nil {
			workerGroup.Monitors = append(workerGroup.Monitors, *monitorName)
		}
	}

	return workerGroup, err
}
