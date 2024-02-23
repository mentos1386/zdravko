package services

import (
	"context"
	"log"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"gorm.io/gorm"
)

func getScheduleId(healthcheck *models.Healthcheck, group string) string {
	return "healthcheck-" + healthcheck.Slug + "-" + group
}

func CreateHealthcheck(ctx context.Context, query *query.Query, healthcheck *models.Healthcheck) error {
	return query.Healthcheck.WithContext(ctx).Create(healthcheck)
}

func UpdateHealthcheck(ctx context.Context, q *query.Query, healthcheck *models.Healthcheck) error {
	_, err := q.Healthcheck.WithContext(ctx).Where(
		q.Healthcheck.Slug.Eq(healthcheck.Slug),
	).Updates(healthcheck)
	return err
}

func GetHealthcheck(ctx context.Context, q *query.Query, slug string) (*models.Healthcheck, error) {
	return q.Healthcheck.WithContext(ctx).Where(
		q.Healthcheck.Slug.Eq(slug),
		q.Healthcheck.DeletedAt.IsNull(),
	).Preload(
		q.Healthcheck.History,
	).First()
}

func GetHealthchecks(ctx context.Context, q *query.Query) ([]*models.Healthcheck, error) {
	return q.Healthcheck.WithContext(ctx).Preload(
		q.Healthcheck.History,
	).Where(
		q.Healthcheck.DeletedAt.IsNull(),
	).Find()
}

func CreateOrUpdateHealthcheckSchedule(ctx context.Context, t client.Client, healthcheck *models.Healthcheck) error {
	log.Println("Creating or Updating Healthcheck Schedule")

	args := make([]interface{}, 0)
	args = append(args, workflows.HealthcheckWorkflowParam{Script: healthcheck.Script, Slug: healthcheck.Slug})

	for _, group := range healthcheck.WorkerGroups {
		options := client.ScheduleOptions{
			ID: getScheduleId(healthcheck, group),
			//SearchAttributes: map[string]interface{}{
			//	"worker-group":     group,
			//	"healthcheck-slug": healthcheck.Slug,
			//},
			Spec: client.ScheduleSpec{
				CronExpressions: []string{healthcheck.Schedule},
				Jitter:          time.Second * 10,
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        getScheduleId(healthcheck, group),
				Workflow:  workflows.NewWorkflows(nil).HealthcheckWorkflowDefinition,
				Args:      args,
				TaskQueue: group,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 3,
				},
			},
		}

		schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(healthcheck, group))

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
	}

	return nil
}

func CreateHealthcheckHistory(ctx context.Context, db *gorm.DB, healthcheckHistory *models.HealthcheckHistory) error {
	return db.WithContext(ctx).Create(healthcheckHistory).Error
}
