package services

import (
	"context"
	"log"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func CreateHealthcheckHttp(ctx context.Context, db *gorm.DB, healthcheck *models.HealthcheckHttp) error {
	return db.WithContext(ctx).Create(healthcheck).Error
}

func GetHealthcheckHttp(ctx context.Context, q *query.Query, id uint) (*models.HealthcheckHttp, error) {
	log.Println("GetHealthcheckHttp")
	return q.HealthcheckHttp.WithContext(ctx).Where(
		q.HealthcheckHttp.ID.Eq(id),
	).First()
}

func StartHealthcheckHttp(ctx context.Context, t client.Client) error {
	log.Println("Starting HealthcheckHttp Workflow")

	args := make([]interface{}, 0)
	args = append(args, workflows.HealthcheckHttpWorkflowParam{Id: 1})

	_, err := t.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: "healthcheck-http-id",
		Spec: client.ScheduleSpec{
			CronExpressions: []string{"0 * * * *"},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        "healthcheck-http-id-workflow",
			Workflow:  workflows.HealthcheckHttpWorkflowDefinition,
			Args:      args,
			TaskQueue: "default",
		},
	})

	return err
}
