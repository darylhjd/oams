package main

import (
	"flag"
	"fmt"

	"github.com/darylhjd/oams/backend/env"
)

// parseFlags for the programme.
func parseFlags() (*arguments, error) {
	defaultDatabaseName, err := env.GetDatabaseName()
	if err != nil {
		return nil, err
	}

	name := flag.String("name", defaultDatabaseName,
		"Name of the database to do operation on.\n"+
			"It is always good to specify which database you need explicitly to avoid any mistakes.\n"+
			"If no name is provided, then the default database will be used.")
	mig := flag.Bool("migrate", false,
		"Set this flag to tell the programme that you intend to do migrations on the specified database.\n"+
			"One of `migrate`, `create`, or `drop` must be set.\n"+
			"By default, this is false.")
	create := flag.Bool("create", false,
		"Set this flag to tell the programme that you intend to create a new database.\n"+
			"Specify the name of the database with the `name` flag.\n"+
			"One of `migrate`, `create`, or `drop` must be set.\n"+
			"By default, this is false.")
	drop := flag.Bool("drop", false,
		"Set this flag to tell the programme that you intend to drop a database.\n"+
			"Specify the name of the database with the `name` flag.\n"+
			"One of `migrate`, `create`, or `drop` must be set.\n"+
			"By default, this is false.")
	version := flag.Uint("version", noOp,
		"Version of the database to migrate to.\n"+
			fmt.Sprintf("Value must be a number greater than %d.\n", noOp)+
			"One of `version`, `steps`, `full-up`, or `full-down` must be set.\n"+
			"By default, this will do a no-op migration.")
	steps := flag.Int("steps", noOp,
		"Number of steps up or down to migrate the database.\n"+
			fmt.Sprintf("Value should not be %d.\n", noOp)+
			"One of `version`, `steps`, `full-up`, or `full-down` must be set.\n"+
			"By default, this will do a no-op migration.")
	fullUp := flag.Bool("full-up", false,
		"Set this flag to tell the programme to up-migrate the current database to latest specification.\n"+
			"One of `version`, `steps`, `full-up`, or `full-down` must be set.\n"+
			"By default, this will do a no-op migration.")
	fullDown := flag.Bool("full-down", false,
		"Set this flag to tell the programme to apply all down-migrations.\n"+
			"One of `version`, `steps`, `full-up`, or `full-down` must be set.\n"+
			"By default, this will do a no-op migration.")

	flag.Parse()

	return &arguments{
		name:     *name,
		migrate:  *mig,
		create:   *create,
		drop:     *drop,
		version:  *version,
		steps:    *steps,
		fullUp:   *fullUp,
		fullDown: *fullDown,
	}, nil
}
