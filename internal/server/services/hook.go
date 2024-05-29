package services

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database/models"
	internaltemporal "github.com/mentos1386/zdravko/internal/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func getHookScheduleId(id string) string {
	return "hook-" + id
}

func CountHooks(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM hooks")
	return count, err
}

func GetHookState(ctx context.Context, temporal client.Client, id string) (models.HookState, error) {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getHookScheduleId(id))

	description, err := schedule.Describe(ctx)
	if err != nil {
		return models.HookStateUnknown, err
	}

	if description.Schedule.State.Paused {
		return models.HookStatePaused, nil
	}

	return models.HookStateActive, nil
}

func SetHookState(ctx context.Context, temporal client.Client, id string, state models.HookState) error {
	schedule := temporal.ScheduleClient().GetHandle(ctx, getHookScheduleId(id))

	if state == models.HookStateActive {
		return schedule.Unpause(ctx, client.ScheduleUnpauseOptions{Note: "Unpaused by user"})
	}

	if state == models.HookStatePaused {
		return schedule.Pause(ctx, client.SchedulePauseOptions{Note: "Paused by user"})
	}

	return nil
}

func CreateHook(ctx context.Context, db *sqlx.DB, hook *models.Hook) error {
	_, err := db.NamedExecContext(ctx,
		`INSERT INTO hooks (id, name, script, schedule)
    VALUES (:id, :name, :script, :schedule)`,
		hook,
	)
	return err
}

func UpdateHook(ctx context.Context, db *sqlx.DB, hook *models.Hook) error {
	_, err := db.NamedExecContext(ctx,
		`UPDATE hooks SET script=:script, schedule=:schedule WHERE id=:id`,
		hook,
	)
	return err
}

func DeleteHook(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM hooks WHERE id=$1",
		id,
	)
	return err
}

func GetHook(ctx context.Context, db *sqlx.DB, id string) (*models.Hook, error) {
	hook := &models.Hook{}
	err := db.GetContext(ctx, hook,
		"SELECT * FROM hooks WHERE id=$1",
		id,
	)
	return hook, err
}

func GetHooks(ctx context.Context, db *sqlx.DB) ([]*models.Hook, error) {
	hooks := []*models.Hook{}
	err := db.SelectContext(ctx, &hooks,
		"SELECT * FROM hooks ORDER BY name",
	)
	return hooks, err
}

func DeleteHookSchedule(ctx context.Context, t client.Client, id string) error {
	schedule := t.ScheduleClient().GetHandle(ctx, getHookScheduleId(id))
	return schedule.Delete(ctx)
}

func CreateOrUpdateHookSchedule(
	ctx context.Context,
	t client.Client,
	hook *models.Hook,
) error {
	log.Println("Creating or Updating Hook Schedule")

	args := make([]interface{}, 1)
	args[0] = internaltemporal.WorkflowHookParam{
		Script: hook.Script,
		HookId: hook.Id,
	}

	options := client.ScheduleOptions{
		ID: getHookScheduleId(hook.Id),
		Spec: client.ScheduleSpec{
			CronExpressions: []string{hook.Schedule},
			Jitter:          time.Second * 10,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        getHookScheduleId(hook.Id),
			Workflow:  internaltemporal.WorkflowHookName,
			Args:      args,
			TaskQueue: internaltemporal.TEMPORAL_SERVER_QUEUE,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		},
	}

	schedule := t.ScheduleClient().GetHandle(ctx, getHookScheduleId(hook.Id))

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
