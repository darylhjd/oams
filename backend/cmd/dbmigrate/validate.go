package main

import (
	"fmt"

	"github.com/darylhjd/oams/backend/internal/database"
)

// validateArguments that were provided to the application.
func validateArguments(args *arguments) error {
	// Check database name is valid.
	if args.name == "" {
		return fmt.Errorf("%s - database name cannot be empty", database.MigrationNamespace)
	}

	// Check validity of operations is correct.
	// One and only one operation must be specified.
	ops := 0
	if args.migrate {
		ops++
	}

	if args.create {
		ops++
	}

	if args.drop {
		ops++
	}

	if args.truncate {
		ops++
	}

	if ops == noOp {
		return fmt.Errorf("%s - no operation specified", database.MigrationNamespace)
	} else if ops > minOp {
		return fmt.Errorf("%s - more than one operation specified", database.MigrationNamespace)
	}

	// Check validity of options.
	ops = 0
	if args.version != noOp {
		ops++
	}

	if args.steps != noOp {
		ops++
	}

	if args.fullUp {
		ops++
	}

	if args.fullDown {
		ops++
	}

	// If migrate operation not selected but options set, return error.
	if !args.migrate && ops != noOp {
		return fmt.Errorf("%s - invalid options provided for operation", database.MigrationNamespace)
	}

	// If non-migration operation, we can exit early.
	if !args.migrate {
		return nil
	}

	if ops == noOp {
		return fmt.Errorf("%s - no options specified", database.MigrationNamespace)
	} else if ops > minOp {
		return fmt.Errorf("%s - more than one option specified", database.MigrationNamespace)
	}

	return nil
}
