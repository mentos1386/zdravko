package services

import (
	"context"
	"log"

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
	log.Println("GetHealthcheck")
	return q.Healthcheck.WithContext(ctx).Where(
		q.Healthcheck.Slug.Eq(slug),
	).First()
}

func StartHealthcheck(ctx context.Context, t client.Client, healthcheck *models.Healthcheck) error {
	log.Println("Starting Healthcheck Workflow")

	args := make([]interface{}, 0)
	args = append(args, workflows.HealthcheckWorkflowParam{Script: healthcheck.Script})

	id := "healthcheck-" + healthcheck.Slug

	for _, group := range healthcheck.WorkerGroups {
		_, err := t.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: id + "-" + group,
			Spec: client.ScheduleSpec{
				CronExpressions: []string{healthcheck.Schedule},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        id + "-" + group,
				Workflow:  workflows.HealthcheckWorkflowDefinition,
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
