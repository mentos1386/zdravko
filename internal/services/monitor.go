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

type MonitorStatus string

const (
	MonitorStatusUnknown MonitorStatus = "UNKNOWN"
	MonitorStatusPaused  MonitorStatus = "PAUSED"
	MonitorStatusActive  MonitorStatus = "ACTIVE"
)

func getScheduleId(id string) string {
	return "monitor-" + id
}

func CountMonitors(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM monitors")
	return count, err
}

func GetMonitorStatus(ctx context.Context, temporal client.Client, id string) (MonitorStatus, error) {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getScheduleId(id))

	description, err := schedule.Describe(ctx)
	if err != nil {
		return MonitorStatusUnknown, err
	}

	if description.Schedule.State.Paused {
		return MonitorStatusPaused, nil
	}

	return MonitorStatusActive, nil
}

func SetMonitorStatus(ctx context.Context, temporal client.Client, id string, status MonitorStatus) error {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getScheduleId(id))

	if status == MonitorStatusActive {
		return schedule.Unpause(ctx, client.ScheduleUnpauseOptions{Note: "Unpaused by user"})
	}

	if status == MonitorStatusPaused {
		return schedule.Pause(ctx, client.SchedulePauseOptions{Note: "Paused by user"})
	}

	return nil
}

func CreateMonitor(ctx context.Context, db *sqlx.DB, monitor *models.Monitor) error {
	_, err := db.NamedExecContext(ctx,
		"INSERT INTO monitors (id, name, script, schedule) VALUES (:id, :name, :script, :schedule)",
		monitor,
	)
	return err
}

func UpdateMonitor(ctx context.Context, db *sqlx.DB, monitor *models.Monitor) error {
	_, err := db.NamedExecContext(ctx,
		"UPDATE monitors SET name=:name, script=:script, schedule=:schedule WHERE id=:id",
		monitor,
	)
	return err
}

func DeleteMonitor(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM monitors WHERE id=$1",
		id,
	)
	return err
}

func UpdateMonitorWorkerGroups(ctx context.Context, db *sqlx.DB, monitor *models.Monitor, workerGroups []*models.WorkerGroup) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		"DELETE FROM monitor_worker_groups WHERE monitor_id=$1",
		monitor.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, group := range workerGroups {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO monitor_worker_groups (monitor_id, worker_group_id) VALUES ($1, $2)",
			monitor.Id,
			group.Id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetMonitor(ctx context.Context, db *sqlx.DB, id string) (*models.Monitor, error) {
	monitor := &models.Monitor{}
	err := db.GetContext(ctx, monitor,
		"SELECT * FROM monitors WHERE id=$1",
		id,
	)
	return monitor, err
}

func GetMonitorWithWorkerGroups(ctx context.Context, db *sqlx.DB, id string) (*models.MonitorWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  monitors.id,
  monitors.name,
  monitors.script,
  monitors.schedule,
  monitors.created_at,
  monitors.updated_at,
  worker_groups.name as worker_group_name
FROM monitors
LEFT OUTER JOIN monitor_worker_groups ON monitors.id = monitor_worker_groups.monitor_id
LEFT OUTER JOIN worker_groups ON monitor_worker_groups.worker_group_id = worker_groups.id
WHERE monitors.id=$1
`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitor := &models.MonitorWithWorkerGroups{}

	for rows.Next() {
		var workerGroupName *string
		err = rows.Scan(
			&monitor.Id,
			&monitor.Name,
			&monitor.Script,
			&monitor.Schedule,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
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
		"SELECT * FROM monitors ORDER BY name",
	)
	return monitors, err
}

func GetMonitorsWithWorkerGroups(ctx context.Context, db *sqlx.DB) ([]*models.MonitorWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  monitors.id,
  monitors.name,
  monitors.script,
  monitors.schedule,
  monitors.created_at,
  monitors.updated_at,
  worker_groups.name as worker_group_name
FROM monitors
LEFT OUTER JOIN monitor_worker_groups ON monitors.id = monitor_worker_groups.monitor_id
LEFT OUTER JOIN worker_groups ON monitor_worker_groups.worker_group_id = worker_groups.id
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
			&monitor.Id,
			&monitor.Name,
			&monitor.Script,
			&monitor.Schedule,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
			&workerGroupName,
		)
		if err != nil {
			return nil, err
		}
		if workerGroupName != nil {
			workerGroups := []string{}
			if monitors[monitor.Id] != nil {
				workerGroups = monitors[monitor.Id].WorkerGroups
			}
			monitor.WorkerGroups = append(workerGroups, *workerGroupName)
		}
		monitors[monitor.Id] = monitor
	}

	return maps.Values(monitors), err
}

func DeleteMonitorSchedule(ctx context.Context, t client.Client, id string) error {
	schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(id))
	return schedule.Delete(ctx)
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
		workerGroupStrings[i] = group.Id
	}

	args := make([]interface{}, 1)
	args[0] = workflows.MonitorWorkflowParam{
		Script:         monitor.Script,
		MonitorId:      monitor.Id,
		WorkerGroupIds: workerGroupStrings,
	}

	options := client.ScheduleOptions{
		ID: getScheduleId(monitor.Id),
		Spec: client.ScheduleSpec{
			CronExpressions: []string{monitor.Schedule},
			Jitter:          time.Second * 10,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        getScheduleId(monitor.Id),
			Workflow:  workflows.NewWorkflows(nil).MonitorWorkflowDefinition,
			Args:      args,
			TaskQueue: "default",
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		},
	}

	schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(monitor.Id))

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
