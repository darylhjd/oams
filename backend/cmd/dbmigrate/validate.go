package main

import (
	"fmt"

	"github.com/darylhjd/oats/backend/database"
)

// validateArguments that were provided to the application.
func validateArguments(args *arguments) error {
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

	if ops == noOp {
		return fmt.Errorf("%s - no operation specified", database.MigrationNamespace)
	} else if ops > minOp {
		return fmt.Errorf("%s - more than one operation specified", database.MigrationNamespace)
	}

	// Check validity of options.
	// If migrate operation not selected but options set, return error.
	if !args.migrate && (args.version != noOp || args.steps != noOp) {
		return fmt.Errorf("%s - invalid options provided for operation", database.MigrationNamespace)
	}

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

	if ops == noOp {
		return fmt.Errorf("%s - no options specified", database.MigrationNamespace)
	} else if ops >= minOp {
		return fmt.Errorf("%s - more than one option specified", database.MigrationNamespace)
	}

	return nil
}
