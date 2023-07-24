package tests

import (
	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/jackc/pgx/v5/pgtype"
)

func MockUpsertStudentsParams() database.UpsertStudentsParams {
	return database.UpsertStudentsParams{
		ID:   oauth2.MockIDTokenName,
		Name: "",
		Email: pgtype.Text{
			String: oauth2.MockAccountPreferredUsername,
			Valid:  true,
		},
	}
}

func StubStudent(createdAt, updatedAt pgtype.Timestamp) database.Student {
	return database.Student{
		ID:   oauth2.MockIDTokenName,
		Name: "",
		Email: pgtype.Text{
			String: oauth2.MockAccountPreferredUsername,
			Valid:  true,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
