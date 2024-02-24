package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"gorm.io/gorm/clause"
)

func GetOrCreateWorkerGroup(ctx context.Context, q *query.Query, workerGroup models.WorkerGroup) (*models.WorkerGroup, error) {
	tx := q.Begin()

	if err := tx.WorkerGroup.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&workerGroup); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	wg, err := tx.WorkerGroup.WithContext(ctx).Where(
		q.WorkerGroup.Slug.Eq(workerGroup.Slug),
	).First()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return wg, tx.Commit()
}

func GetWorkerGroups(ctx context.Context, q *query.Query) ([]*models.WorkerGroup, error) {
	return q.WorkerGroup.WithContext(ctx).Preload(q.WorkerGroup.Monitors).Find()
}

func GetWorkerGroupsBySlug(ctx context.Context, q *query.Query, slugs []string) ([]*models.WorkerGroup, error) {
	return q.WorkerGroup.WithContext(ctx).Where(
		q.WorkerGroup.Slug.In(slugs...),
	).Find()
}

func GetWorkerGroup(ctx context.Context, q *query.Query, slug string) (*models.WorkerGroup, error) {
	return q.WorkerGroup.WithContext(ctx).Where(
		q.WorkerGroup.Slug.Eq(slug),
	).First()
}
