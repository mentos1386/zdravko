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

func CreateHealthcheckHttp(ctx context.Context, db *gorm.DB, healthcheck *models.HealthcheckHttp) error {
	return db.WithContext(ctx).Create(healthcheck).Error
}

func GetHealthcheckHttp(ctx context.Context, q *query.Query, slug string) (*models.HealthcheckHttp, error) {
	log.Println("GetHealthcheckHttp")
	return q.HealthcheckHttp.WithContext(ctx).Where(
		q.HealthcheckHttp.Slug.Eq(slug),
	).First()
}

func StartHealthcheckHttp(ctx context.Context, t client.Client, healthcheckHttp *models.HealthcheckHttp) error {
	log.Println("Starting HealthcheckHttp Workflow")

	args := make([]interface{}, 0)
	args = append(args, workflows.HealthcheckHttpWorkflowParam{Url: healthcheckHttp.Url, Method: healthcheckHttp.Method})

	for _, group := range healthcheckHttp.WorkerGroups {
		_, err := t.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: "healthcheck-http-" + healthcheckHttp.Slug,
			Spec: client.ScheduleSpec{
				CronExpressions: []string{healthcheckHttp.Schedule},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        "healthcheck-http-" + healthcheckHttp.Slug,
				Workflow:  workflows.HealthcheckHttpWorkflowDefinition,
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
