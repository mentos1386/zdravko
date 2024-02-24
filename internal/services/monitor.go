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

func getScheduleId(monitor *models.Monitor, group string) string {
	return "monitor-" + monitor.Slug + "-" + group
}

func CreateMonitor(ctx context.Context, query *query.Query, monitor *models.Monitor) error {
	return query.Monitor.WithContext(ctx).Create(monitor)
}

func UpdateMonitor(ctx context.Context, q *query.Query, monitor *models.Monitor) error {
	_, err := q.Monitor.WithContext(ctx).Where(
		q.Monitor.Slug.Eq(monitor.Slug),
	).Updates(monitor)

	return err
}

func UpdateMonitorWorkerGroups(ctx context.Context, q *query.Query, monitor *models.Monitor, workerGroups []*models.WorkerGroup) error {
	return q.Monitor.WorkerGroups.Model(monitor).Replace(workerGroups...)
}

func GetMonitor(ctx context.Context, q *query.Query, slug string) (*models.Monitor, error) {
	return q.Monitor.WithContext(ctx).Where(
		q.Monitor.Slug.Eq(slug),
	).Preload(
		q.Monitor.WorkerGroups,
		q.Monitor.History,
	).First()
}

func GetMonitors(ctx context.Context, q *query.Query) ([]*models.Monitor, error) {
	return q.Monitor.WithContext(ctx).Preload(
		q.Monitor.History,
	).Preload(
		q.Monitor.WorkerGroups,
	).Find()
}

func CreateOrUpdateMonitorSchedule(
	ctx context.Context,
	t client.Client,
	monitor *models.Monitor,
) error {
	log.Println("Creating or Updating Monitor Schedule")

	args := make([]interface{}, 0)
	args = append(args, workflows.MonitorWorkflowParam{Script: monitor.Script, Slug: monitor.Slug})

	for _, group := range monitor.WorkerGroups {
		options := client.ScheduleOptions{
			ID: getScheduleId(monitor, group.Slug),
			Spec: client.ScheduleSpec{
				CronExpressions: []string{monitor.Schedule},
				Jitter:          time.Second * 10,
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        getScheduleId(monitor, group.Slug),
				Workflow:  workflows.NewWorkflows(nil).MonitorWorkflowDefinition,
				Args:      args,
				TaskQueue: group.Slug,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 3,
				},
			},
		}

		schedule := t.ScheduleClient().GetHandle(ctx, getScheduleId(monitor, group.Slug))

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

func CreateMonitorHistory(ctx context.Context, db *gorm.DB, monitorHistory *models.MonitorHistory) error {
	return db.WithContext(ctx).Create(monitorHistory).Error
}
