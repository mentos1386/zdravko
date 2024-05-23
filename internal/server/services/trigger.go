package services

import (
	"context"

	"github.com/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

func CountTriggers(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM triggers")
	return count, err
}

func CreateTrigger(ctx context.Context, db *sqlx.DB, trigger *models.Trigger) error {
	_, err := db.NamedExecContext(ctx,
		`INSERT INTO triggers (id, name, script) VALUES (:id, :name, :script)`,
		trigger,
	)
	return err
}

func UpdateTrigger(ctx context.Context, db *sqlx.DB, trigger *models.Trigger) error {
	_, err := db.NamedExecContext(ctx,
		`UPDATE triggers SET script=:script WHERE id=:id`,
		trigger,
	)
	return err
}

func DeleteTrigger(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM triggers WHERE id=$1",
		id,
	)
	return err
}

func GetTrigger(ctx context.Context, db *sqlx.DB, id string) (*models.Trigger, error) {
	trigger := &models.Trigger{}
	err := db.GetContext(ctx, trigger,
		"SELECT * FROM triggers WHERE id=$1",
		id,
	)
	return trigger, err
}

func GetTriggers(ctx context.Context, db *sqlx.DB) ([]*models.Trigger, error) {
	triggers := []*models.Trigger{}
	err := db.SelectContext(ctx, &triggers,
		"SELECT * FROM triggers ORDER BY name",
	)
	return triggers, err
}
