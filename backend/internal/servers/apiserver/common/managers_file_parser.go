package common

import (
	"errors"
	"fmt"
	"io"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/xuri/excelize/v2"
)

func ParseManagersFile(f io.Reader) ([]database.ProcessUpsertClassGroupManagerParams, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	sheets := file.GetSheetList()
	if len(sheets) != expectedManagersSheetCount {
		return nil, errors.New("invalid sheet count for manager file")
	}

	rows, err := file.GetRows(sheets[expectedManagersSheetCount-1])
	if err != nil {
		return nil, fmt.Errorf("cannot get sheet rows: %w", err)
	}

	if err = parseSanityCheck(rows); err != nil {
		return nil, fmt.Errorf("failed file sanity check: %w", err)
	}

	params, err := parseManagerData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed parsing manager data: %w", err)
	}

	return params, nil
}

func parseSanityCheck(rows [][]string) error {
	if len(rows) < expectedSanityCheckDataRows {
		return errors.New("unexpected number of rows for file")
	}

	if len(rows[expectedSanityCheckDataRows-1]) != expectedSanityCheckDataColumns {
		return errors.New("unexpected number of columns for file")
	}

	for idx, col := range rows[expectedSanityCheckDataRows-1] {
		if col != managersColumnNames[idx] {
			return fmt.Errorf("failed sanity check for header name: %s", col)
		}
	}

	return nil
}

func parseManagerData(rows [][]string) ([]database.ProcessUpsertClassGroupManagerParams, error) {
	params := make([]database.ProcessUpsertClassGroupManagerParams, 0, len(rows)-expectedSanityCheckDataRows)
	for index := expectedSanityCheckDataRows; index < len(rows); index++ {
		row := rows[index]

		if len(row) != expectedSanityCheckDataColumns {
			return nil, fmt.Errorf("unexpected number of data columns on row %d", index+1)
		}

		var (
			userId         string
			classCode      string
			classYear      int32
			classSemester  string
			classGroupName string
			classType      model.ClassType
			managingRole   model.ManagingRole
		)

		if _, err := fmt.Sscanf(row[userIdColumn], "%s", &userId); err != nil {
			return nil, fmt.Errorf("could not parse user id: %w", err)
		}

		if _, err := fmt.Sscanf(row[classCodeColumn], "%s", &classCode); err != nil {
			return nil, fmt.Errorf("could not parse class code: %w", err)
		}

		if _, err := fmt.Sscanf(row[classYearColumn], "%d", &classYear); err != nil {
			return nil, fmt.Errorf("could not parse class year: %w", err)
		}

		if _, err := fmt.Sscanf(row[classSemesterColumn], "%s", &classSemester); err != nil {
			return nil, fmt.Errorf("could not parse class semester: %w", err)
		}

		if _, err := fmt.Sscanf(row[classGroupNameColumn], "%s", &classGroupName); err != nil {
			return nil, fmt.Errorf("could not parse class group name: %w", err)
		}

		if _, err := fmt.Sscanf(row[classTypeColumn], "%s", &classType); err != nil {
			return nil, fmt.Errorf("could not parse class type %w", err)
		}

		if _, err := fmt.Sscanf(row[managingRowColumn], "%s", &managingRole); err != nil {
			return nil, fmt.Errorf("could not parse managing role: %w", err)
		}

		params = append(params, database.ProcessUpsertClassGroupManagerParams{
			UserID:         userId,
			ClassCode:      classCode,
			ClassYear:      classYear,
			ClassSemester:  classSemester,
			ClassGroupName: classGroupName,
			ClassType:      classType,
			ManagingRole:   managingRole,
		})
	}

	return params, nil
}
