package managers

import (
	"io"

	"github.com/darylhjd/oams/backend/internal/database"
)

func ParseManagersFile(file io.Reader) ([]database.UpsertClassGroupManagerParams, error) {
	return []database.UpsertClassGroupManagerParams{}, nil
}
