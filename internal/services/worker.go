package services

import (
	"context"
	"log"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
)

func CreateWorker(ctx context.Context, q *query.Query, worker *models.Worker) error {
	return q.Worker.WithContext(ctx).Create(worker)
}

func GetWorker(ctx context.Context, q *query.Query, slug string) (*models.Worker, error) {
	log.Println("GetWorker")
	return q.Worker.WithContext(ctx).Where(
		q.Worker.Slug.Eq(slug),
	).First()
}
