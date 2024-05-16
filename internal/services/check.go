package services

import (
	"context"
	"log"
	"sort"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"golang.org/x/exp/maps"
)

type CheckStatus string

const (
	CheckStatusUnknown CheckStatus = "UNKNOWN"
	CheckStatusPaused  CheckStatus = "PAUSED"
	CheckStatusActive  CheckStatus = "ACTIVE"
)

func getScheduleId(id string) string {
	return "check-" + id
}

func CountChecks(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM checks")
	return count, err
}

func GetCheckStatus(ctx context.Context, temporal client.Client, id string) (CheckStatus, error) {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getScheduleId(id))

	description, err := schedule.Describe(ctx)
	if err != nil {
		return CheckStatusUnknown, err
	}

	if description.Schedule.State.Paused {
		return CheckStatusPaused, nil
	}

	return CheckStatusActive, nil
}

func SetCheckStatus(ctx context.Context, temporal client.Client, id string, status CheckStatus) error {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getScheduleId(id))

	if status == CheckStatusActive {
		return schedule.Unpause(ctx, client.ScheduleUnpauseOptions{Note: "Unpaused by user"})
	}

	if status == CheckStatusPaused {
		return schedule.Pause(ctx, client.SchedulePauseOptions{Note: "Paused by user"})
	}

	return nil
}

func CreateCheck(ctx context.Context, db *sqlx.DB, check *models.Check) error {
	_, err := db.NamedExecContext(ctx,
		`INSERT INTO checks (id, name, "group", script, schedule) VALUES (:id, :name, :group, :script, :schedule)`,
		check,
	)
	return err
}

func UpdateCheck(ctx context.Context, db *sqlx.DB, check *models.Check) error {
	_, err := db.NamedExecContext(ctx,
		`UPDATE checks SET "group"=:group, script=:script, schedule=:schedule WHERE id=:id`,
		check,
	)
	return err
}

func DeleteCheck(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM checks WHERE id=$1",
		id,
	)
	return err
}

func UpdateCheckWorkerGroups(ctx context.Context, db *sqlx.DB, check *models.Check, workerGroups []*models.WorkerGroup) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		"DELETE FROM check_worker_groups WHERE check_id=$1",
		check.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, group := range workerGroups {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO check_worker_groups (check_id, worker_group_id) VALUES ($1, $2)",
			check.Id,
			group.Id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetCheck(ctx context.Context, db *sqlx.DB, id string) (*models.Check, error) {
	check := &models.Check{}
	err := db.GetContext(ctx, check,
		"SELECT * FROM checks WHERE id=$1",
		id,
	)
	return check, err
}

func GetCheckWithWorkerGroups(ctx context.Context, db *sqlx.DB, id string) (*models.CheckWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  checks.id,
  checks.name,
  checks."group",
  checks.script,
  checks.schedule,
  checks.created_at,
  checks.updated_at,
  worker_groups.name as worker_group_name
FROM checks
LEFT OUTER JOIN check_worker_groups ON checks.id = check_worker_groups.check_id
LEFT OUTER JOIN worker_groups ON check_worker_groups.worker_group_id = worker_groups.id
WHERE checks.id=$1
ORDER BY checks.name
`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	check := &models.CheckWithWorkerGroups{}

	for rows.Next() {
		var workerGroupName *string
		err = rows.Scan(
			&check.Id,
			&check.Name,
			&check.Group,
			&check.Script,
			&check.Schedule,
			&check.CreatedAt,
			&check.UpdatedAt,
			&workerGroupName,
		)
		if err != nil {
			return nil, err
		}
		if workerGroupName != nil {
			check.WorkerGroups = append(check.WorkerGroups, *workerGroupName)
		}
	}

	return check, err
}

func GetChecks(ctx context.Context, db *sqlx.DB) ([]*models.Check, error) {
	checks := []*models.Check{}
	err := db.SelectContext(ctx, &checks,
		"SELECT * FROM checks ORDER BY name",
	)
	return checks, err
}

func GetChecksWithWorkerGroups(ctx context.Context, db *sqlx.DB) ([]*models.CheckWithWorkerGroups, error) {
	rows, err := db.QueryContext(ctx,
		`
SELECT
  checks.id,
  checks.name,
  checks."group",
  checks.script,
  checks.schedule,
  checks.created_at,
  checks.updated_at,
  worker_groups.name as worker_group_name
FROM checks
LEFT OUTER JOIN check_worker_groups ON checks.id = check_worker_groups.check_id
LEFT OUTER JOIN worker_groups ON check_worker_groups.worker_group_id = worker_groups.id
ORDER BY checks.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := map[string]*models.CheckWithWorkerGroups{}

	for rows.Next() {
		check := &models.CheckWithWorkerGroups{}

		var workerGroupName *string
		err = rows.Scan(
			&check.Id,
			&check.Name,
			&check.Group,
			&check.Script,
			&check.Schedule,
			&check.CreatedAt,
			&check.UpdatedAt,
			&workerGroupName,
		)
		if err != nil {
			return nil, err
		}
		if workerGroupName != nil {
			workerGroups := []string{}
			if checks[check.Id] != nil {
				workerGroups = checks[check.Id].WorkerGroups
			}
			check.WorkerGroups = append(workerGroups, *workerGroupName)
		}
		checks[check.Id] = check
	}

	checksWithWorkerGroups := maps.Values(checks)
	sort.SliceStable(checksWithWorkerGroups, func(i, j int) bool {
		return checksWithWorkerGroups[i].Name < checksWithWorkerGroups[j].Name
	})

	return checksWithWorkerGroups, err
}

func DeleteCheckSchedule(ctx context.Context, t client.Client, id string) error {
	schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(id))
	return schedule.Delete(ctx)
}

func CreateOrUpdateCheckSchedule(
	ctx context.Context,
	t client.Client,
	check *models.Check,
	workerGroups []*models.WorkerGroup,
) error {
	log.Println("Creating or Updating Check Schedule")

	workerGroupStrings := make([]string, len(workerGroups))
	for i, group := range workerGroups {
		workerGroupStrings[i] = group.Id
	}

	args := make([]interface{}, 1)
	args[0] = workflows.CheckWorkflowParam{
		Script:         check.Script,
		CheckId:      check.Id,
		WorkerGroupIds: workerGroupStrings,
	}

	options := client.ScheduleOptions{
		ID: getScheduleId(check.Id),
		Spec: client.ScheduleSpec{
			CronExpressions: []string{check.Schedule},
			Jitter:          time.Second * 10,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        getScheduleId(check.Id),
			Workflow:  workflows.NewWorkflows(nil).CheckWorkflowDefinition,
			Args:      args,
			TaskQueue: "default",
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		},
	}

	schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(check.Id))

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
