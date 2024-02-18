package services

import (
	"context"
	"log"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"gorm.io/gorm"
)

func CreateWorker(ctx context.Context, db *gorm.DB, worker *models.Worker) error {
	return db.WithContext(ctx).Create(worker).Error
}

func GetWorker(ctx context.Context, q *query.Query, slug string) (*models.Worker, error) {
	log.Println("GetWorker")
	return q.Worker.WithContext(ctx).Where(
		q.Worker.Slug.Eq(slug),
	).First()
}
