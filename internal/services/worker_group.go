package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"golang.org/x/exp/maps"
)

func GetActiveWorkers(ctx context.Context, workerGroupSlug string, temporal client.Client) ([]string, error) {
	response, err := temporal.DescribeTaskQueue(ctx, workerGroupSlug, enums.TASK_QUEUE_TYPE_ACTIVITY)
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
		"INSERT INTO worker_groups (slug, name) VALUES (:slug, :name)",
		workerGroup,
	)
	return err
}

func GetWorkerGroups(ctx context.Context, db *sqlx.DB) ([]*models.WorkerGroup, error) {
	var workerGroups []*models.WorkerGroup
	err := db.SelectContext(ctx, &workerGroups,
		"SELECT * FROM worker_groups WHERE deleted_at IS NULL ORDER BY name",
	)
	return workerGroups, err
}

func GetWorkerGroupsWithMonitors(ctx context.Context, db *sqlx.DB) ([]*models.WorkerGroupWithMonitors, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.slug,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  worker_groups.deleted_at,
  monitors.name as monitor_name
FROM worker_groups
LEFT OUTER JOIN monitor_worker_groups ON worker_groups.slug = monitor_worker_groups.worker_group_slug
LEFT OUTER JOIN monitors ON monitor_worker_groups.monitor_slug = monitors.slug
WHERE worker_groups.deleted_at IS NULL AND monitors.deleted_at IS NULL
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
			&workerGroup.Slug,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&workerGroup.DeletedAt,
			&monitorName,
		)
		if err != nil {
			return nil, err
		}

		if monitorName != nil {
			monitors := []string{}
			if workerGroups[workerGroup.Slug] != nil {
				monitors = workerGroups[workerGroup.Slug].Monitors
			}
			workerGroup.Monitors = append(monitors, *monitorName)
		}

		workerGroups[workerGroup.Slug] = workerGroup
	}

	return maps.Values(workerGroups), err
}

func GetWorkerGroupsBySlug(ctx context.Context, db *sqlx.DB, slugs []string) ([]*models.WorkerGroup, error) {
	var workerGroups []*models.WorkerGroup
	err := db.SelectContext(ctx, &workerGroups,
		"SELECT * FROM worker_groups WHERE slug = ANY($1) AND deleted_at IS NULL",
		slugs,
	)
	return workerGroups, err
}

func GetWorkerGroup(ctx context.Context, db *sqlx.DB, slug string) (*models.WorkerGroup, error) {
	var workerGroup models.WorkerGroup
	err := db.GetContext(ctx, &workerGroup,
		"SELECT * FROM worker_groups WHERE slug = $1 AND deleted_at IS NULL",
		slug,
	)
	return &workerGroup, err
}

func GetWorkerGroupWithMonitors(ctx context.Context, db *sqlx.DB, slug string) (*models.WorkerGroupWithMonitors, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.slug,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  worker_groups.deleted_at,
  monitors.name as monitor_name
FROM worker_groups
LEFT OUTER JOIN monitor_worker_groups ON worker_groups.slug = monitor_worker_groups.worker_group_slug
LEFT OUTER JOIN monitors ON monitor_worker_groups.monitor_slug = monitors.slug
WHERE worker_groups.slug=$1 AND worker_groups.deleted_at IS NULL AND monitors.deleted_at IS NULL
`,
		slug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workerGroup := &models.WorkerGroupWithMonitors{}

	for rows.Next() {
		var monitorName *string
		err = rows.Scan(
			&workerGroup.Slug,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&workerGroup.DeletedAt,
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
