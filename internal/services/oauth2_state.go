package services

import (
	"context"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"github.com/jmoiron/sqlx"
)

func CreateOAuth2State(ctx context.Context, db *sqlx.DB, oauth2State *models.OAuth2State) error {
	_, err := db.NamedExecContext(ctx,
		"INSERT INTO oauth2_states (state, expires_at) VALUES (:state, :expires_at)",
		oauth2State,
	)
	return err
}

func DeleteOAuth2State(ctx context.Context, db *sqlx.DB, state string) (deleted bool, err error) {
	res, err := db.ExecContext(ctx, "DELETE FROM oauth2_states WHERE state = $1 AND expires_at > strftime('%Y-%m-%dT%H:%M:%fZ')", state)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, err
}
