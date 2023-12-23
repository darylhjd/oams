package database

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (d *DB) UpdateUserSignature(ctx context.Context, id string, newSignature string) error {
	hash, err := argon2id.CreateHash(newSignature, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	stmt := UserSignatures.INSERT(
		UserSignatures.UserID,
		UserSignatures.Signature,
	).MODEL(
		model.UserSignature{
			UserID:    id,
			Signature: hash,
		},
	).ON_CONFLICT(
		UserSignatures.UserID,
	).DO_UPDATE(
		SET(
			UserSignatures.Signature.SET(UserSignatures.EXCLUDED.Signature),
		),
	)
	_, err = stmt.ExecContext(ctx, d.qe)
	return err
}
