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

func CreateHealthcheck(ctx context.Context, db *gorm.DB, healthcheck *models.Healthcheck) error {
	return db.WithContext(ctx).Create(healthcheck).Error
}

func GetHealthcheck(ctx context.Context, q *query.Query, slug string) (*models.Healthcheck, error) {
	return q.Healthcheck.WithContext(ctx).Where(
		q.Healthcheck.Slug.Eq(slug),
		q.Healthcheck.DeletedAt.IsNull(),
	).First()
}

func GetHealthchecksWithHistory(ctx context.Context, q *query.Query) ([]*models.Healthcheck, error) {
	return q.Healthcheck.WithContext(ctx).Preload(
		q.Healthcheck.History,
	).Where(
		q.Healthcheck.DeletedAt.IsNull(),
	).Find()
}

func StartHealthcheck(ctx context.Context, t client.Client, healthcheck *models.Healthcheck) error {
	log.Println("Starting Healthcheck Workflow")

	args := make([]interface{}, 0)
	args = append(args, workflows.HealthcheckWorkflowParam{Script: healthcheck.Script, Slug: healthcheck.Slug})

	id := "healthcheck-" + healthcheck.Slug

	fakeWorkflows := workflows.NewWorkflows(nil)

	for _, group := range healthcheck.WorkerGroups {
		_, err := t.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: id + "-" + group,
			Spec: client.ScheduleSpec{
				CronExpressions: []string{healthcheck.Schedule},
				Jitter:          time.Second * 10,
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        id + "-" + group,
				Workflow:  fakeWorkflows.HealthcheckWorkflowDefinition,
				Args:      args,
				TaskQueue: group,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 3,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateHealthcheckHistory(ctx context.Context, db *gorm.DB, healthcheckHistory *models.HealthcheckHistory) error {
	return db.WithContext(ctx).Create(healthcheckHistory).Error
}
