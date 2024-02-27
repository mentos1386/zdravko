package services

import (
	"context"
	"log"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"golang.org/x/exp/maps"
)

func getScheduleId(monitor *models.Monitor) string {
	return "monitor-" + monitor.Slug
}

func CreateMonitor(ctx context.Context, db *sqlx.DB, monitor *models.Monitor) error {
	_, err := db.NamedExecContext(ctx,
		"INSERT INTO monitors (slug, name, script, schedule) VALUES (:slug, :name, :script, :schedule)",
		monitor,
	)
	return err
}

func UpdateMonitor(ctx context.Context, db *sqlx.DB, monitor *models.Monitor) error {
	_, err := db.NamedExecContext(ctx,
		"UPDATE monitors SET name=:name, script=:script, schedule=:schedule WHERE slug=:slug",
		monitor,
	)
	return err
}

func UpdateMonitorWorkerGroups(ctx context.Context, db *sqlx.DB, monitor *models.Monitor, workerGroups []*models.WorkerGroup) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		"DELETE FROM monitor_worker_groups WHERE monitor_slug=$1",
		monitor.Slug,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, group := range workerGroups {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO monitor_worker_groups (monitor_slug, worker_group_slug) VALUES ($1, $2)",
			monitor.Slug,
			group.Slug,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetMonitor(ctx context.Context, db *sqlx.DB, slug string) (*models.Monitor, error) {
	monitor := &models.Monitor{}
	err := db.GetContext(ctx, monitor,
		"SELECT * FROM monitors WHERE slug=$1 AND deleted_at IS NULL",
		slug,
	)
	return monitor, err
}

func GetMonitorWithWorkerGroups(ctx context.Context, db *sqlx.DB, slug string) (*models.MonitorWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  monitors.slug,
  monitors.name,
  monitors.script,
  monitors.schedule,
  monitors.created_at,
  monitors.updated_at,
  monitors.deleted_at,
  worker_groups.name as worker_group_name
FROM monitors
LEFT OUTER JOIN monitor_worker_groups ON monitors.slug = monitor_worker_groups.monitor_slug
LEFT OUTER JOIN worker_groups ON monitor_worker_groups.worker_group_slug = worker_groups.slug
WHERE monitors.slug=$1 AND monitors.deleted_at IS NULL
`,
		slug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitor := &models.MonitorWithWorkerGroups{}

	for rows.Next() {
		var workerGroupName *string
		err = rows.Scan(
			&monitor.Slug,
			&monitor.Name,
			&monitor.Script,
			&monitor.Schedule,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
			&monitor.DeletedAt,
			&workerGroupName,
		)
		if err != nil {
			return nil, err
		}
		if workerGroupName != nil {
			monitor.WorkerGroups = append(monitor.WorkerGroups, *workerGroupName)
		}
	}

	return monitor, err
}

func GetMonitors(ctx context.Context, db *sqlx.DB) ([]*models.Monitor, error) {
	monitors := []*models.Monitor{}
	err := db.SelectContext(ctx, &monitors,
		"SELECT * FROM monitors WHERE deleted_at IS NULL ORDER BY name",
	)
	return monitors, err
}

func GetMonitorsWithWorkerGroups(ctx context.Context, db *sqlx.DB) ([]*models.MonitorWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  monitors.slug,
  monitors.name,
  monitors.script,
  monitors.schedule,
  monitors.created_at,
  monitors.updated_at,
  monitors.deleted_at,
  worker_groups.name as worker_group_name
FROM monitors
LEFT OUTER JOIN monitor_worker_groups ON monitors.slug = monitor_worker_groups.monitor_slug
LEFT OUTER JOIN worker_groups ON monitor_worker_groups.worker_group_slug = worker_groups.slug
WHERE monitors.deleted_at IS NULL
ORDER BY monitors.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitors := map[string]*models.MonitorWithWorkerGroups{}

	for rows.Next() {
		monitor := &models.MonitorWithWorkerGroups{}

		var workerGroupName *string
		err = rows.Scan(
			&monitor.Slug,
			&monitor.Name,
			&monitor.Script,
			&monitor.Schedule,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
			&monitor.DeletedAt,
			&workerGroupName,
		)
		if err != nil {
			return nil, err
		}
		if workerGroupName != nil {
			workerGroups := []string{}
			if monitors[monitor.Slug] != nil {
				workerGroups = monitors[monitor.Slug].WorkerGroups
			}
			monitor.WorkerGroups = append(workerGroups, *workerGroupName)
		}
		monitors[monitor.Slug] = monitor
	}

	return maps.Values(monitors), err
}

func CreateOrUpdateMonitorSchedule(
	ctx context.Context,
	t client.Client,
	monitor *models.Monitor,
	workerGroups []*models.WorkerGroup,
) error {
	log.Println("Creating or Updating Monitor Schedule")

	workerGroupStrings := make([]string, len(workerGroups))
	for i, group := range workerGroups {
		workerGroupStrings[i] = group.Slug
	}

	args := make([]interface{}, 1)
	args[0] = workflows.MonitorWorkflowParam{
		Script:       monitor.Script,
		Slug:         monitor.Slug,
		WorkerGroups: workerGroupStrings,
	}

	options := client.ScheduleOptions{
		ID: getScheduleId(monitor),
		Spec: client.ScheduleSpec{
			CronExpressions: []string{monitor.Schedule},
			Jitter:          time.Second * 10,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        getScheduleId(monitor),
			Workflow:  workflows.NewWorkflows(nil).MonitorWorkflowDefinition,
			Args:      args,
			TaskQueue: "default",
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		},
	}

	schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(monitor))

	// If exists, we update
	_, err := schedule.Describe(ctx)
	if err == nil {
		err = schedule.Update(ctx, client.ScheduleUpdateOptions{
			DoUpdate: func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
				return &client.ScheduleUpdate{
					Schedule: &client.Schedule{
						Spec:   &options.Spec,
						Action: options.Action,
						Policy: input.Description.Schedule.Policy,
						State:  input.Description.Schedule.State,
					},
				}, nil
			},
		})
		if err != nil {
			return err
		}
	} else {
		schedule, err = t.ScheduleClient().Create(ctx, options)
		if err != nil {
			return err
		}
	}

	err = schedule.Trigger(ctx, client.ScheduleTriggerOptions{})
	if err != nil {
		return err
	}

	return nil
}
