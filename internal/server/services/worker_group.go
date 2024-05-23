package services

import (
	"context"

	"github.com/mentos1386/zdravko/database/models"
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

func GetWorkerGroupsWithChecks(ctx context.Context, db *sqlx.DB) ([]*models.WorkerGroupWithChecks, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.id,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  checks.name as check_name
FROM worker_groups
LEFT OUTER JOIN check_worker_groups ON worker_groups.id = check_worker_groups.worker_group_id
LEFT OUTER JOIN checks ON check_worker_groups.check_id = checks.id
ORDER BY worker_groups.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workerGroups := map[string]*models.WorkerGroupWithChecks{}

	for rows.Next() {
		workerGroup := &models.WorkerGroupWithChecks{}

		var checkName *string
		err = rows.Scan(
			&workerGroup.Id,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&checkName,
		)
		if err != nil {
			return nil, err
		}

		if checkName != nil {
			checks := []string{}
			if workerGroups[workerGroup.Id] != nil {
				checks = workerGroups[workerGroup.Id].Checks
			}
			workerGroup.Checks = append(checks, *checkName)
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

func GetWorkerGroupWithChecks(ctx context.Context, db *sqlx.DB, id string) (*models.WorkerGroupWithChecks, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  worker_groups.id,
  worker_groups.name,
  worker_groups.created_at,
  worker_groups.updated_at,
  checks.name as check_name
FROM worker_groups
LEFT OUTER JOIN check_worker_groups ON worker_groups.id = check_worker_groups.worker_group_id
LEFT OUTER JOIN checks ON check_worker_groups.check_id = checks.id
WHERE worker_groups.id=$1
`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workerGroup := &models.WorkerGroupWithChecks{}

	for rows.Next() {
		var checkName *string
		err = rows.Scan(
			&workerGroup.Id,
			&workerGroup.Name,
			&workerGroup.CreatedAt,
			&workerGroup.UpdatedAt,
			&checkName,
		)
		if err != nil {
			return nil, err
		}
		if checkName != nil {
			workerGroup.Checks = append(workerGroup.Checks, *checkName)
		}
	}

	return workerGroup, err
}
