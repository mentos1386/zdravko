package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database/models"
)

func CountTargets(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM targets")
	return count, err
}

func SetTargetState(ctx context.Context, db *sqlx.DB, id string, state models.TargetState) error {
	_, err := db.NamedExecContext(ctx,
		`UPDATE targets SET state=:state WHERE id=:id`,
		struct {
			Id    string
			State models.TargetState
		}{Id: id, State: state},
	)
	return err
}

func CreateTarget(ctx context.Context, db *sqlx.DB, target *models.Target) error {
	_, err := db.NamedExecContext(ctx,
		`INSERT INTO targets (id, name, "group", visibility, state, metadata) VALUES (:id, :name, :group, :visibility, :state, :metadata)`,
		target,
	)
	return err
}

func UpdateTarget(ctx context.Context, db *sqlx.DB, target *models.Target) error {
	_, err := db.NamedExecContext(ctx,
		`UPDATE targets SET visibility=:visibility, "group"=:group, metadata=:metadata WHERE id=:id`,
		target,
	)
	return err
}

func DeleteTarget(ctx context.Context, db *sqlx.DB, id string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM targets WHERE id=$1",
		id,
	)
	return err
}

func GetTarget(ctx context.Context, db *sqlx.DB, id string) (*models.Target, error) {
	target := &models.Target{}
	err := db.GetContext(ctx, target,
		"SELECT * FROM targets WHERE id=$1",
		id,
	)
	return target, err
}

func GetTargets(ctx context.Context, db *sqlx.DB) ([]*models.Target, error) {
	targets := []*models.Target{}
	err := db.SelectContext(ctx, &targets,
		"SELECT * FROM targets ORDER BY name",
	)
	return targets, err
}

func GetTargetsWithFilter(ctx context.Context, db *sqlx.DB, filter string) ([]*models.Target, error) {
	targets := []*models.Target{}
	err := db.SelectContext(ctx, &targets,
		"SELECT * FROM targets WHERE name ILIKE $1 ORDER BY name",
		"%"+filter+"%",
	)
	return targets, err
}
