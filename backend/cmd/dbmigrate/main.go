package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"

	"github.com/darylhjd/oats/backend/database"
	_ "github.com/darylhjd/oats/backend/env/autoloader"
)

type arguments struct {
	name string

	// Operations
	migrate bool
	create  bool
	drop    bool

	// Options
	version  uint
	steps    int
	fullUp   bool
	fullDown bool
}

const (
	noOp = iota
	minOp
)

func main() {
	args, err := parseFlags()
	if err != nil {
		log.Fatalf("%s - unable to parse cli commands: %s", database.MigrationNamespace, err)
	}

	if err = validateArguments(args); err != nil {
		log.Fatalf("%s - invalid arguments provided to programme: %s", database.MigrationNamespace, err)
	}

	migrator, err := database.NewMigrate(args.name, nil)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case args.migrate:
		err = migrateOp(args, migrator)
	case args.create:
		err = createOp(args, migrator)
	case args.drop:
		err = dropOp(args, migrator)
	default:
		log.Fatalf("%s - reached impossible operation", database.MigrationNamespace)
	}

	if err != nil {
		log.Fatalf("%s - error executing operation: %s", database.MigrationNamespace, err)
	}
}

// migrateOp specifies the operation for migrating the database.
func migrateOp(args *arguments, migrator *migrate.Migrate) error {
	var err error
	switch {
	case args.version != noOp:
		err = migrator.Migrate(args.version)
	case args.steps != noOp:
		err = migrator.Steps(args.steps)
	case args.fullUp:
		err = migrator.Up()
	case args.fullDown:
		err = migrator.Down()
	default:
		err = fmt.Errorf("%s - reached impossible migration sub-operation", database.MigrationNamespace)
	}

	source, db := migrator.Close()
	return errors.Join(source, db, err)
}

func createOp(args *arguments, _ *migrate.Migrate) error {
	return database.Create(args.name, false)
}

func dropOp(args *arguments, _ *migrate.Migrate) error {
	return database.Drop(args.name, true)
}
