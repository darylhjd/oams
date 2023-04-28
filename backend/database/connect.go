package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/darylhjd/oats/backend/env"
)

// Connect and return an interface to the Oats database.
func Connect() (*Queries, error) {
	driver, connString, err := GetConnectionProperties()
	if err != nil {
		return nil, fmt.Errorf("database - cannot parse connection properties: %w", err)
	}

	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}

	// Check that connection is successful.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return New(db), nil
}

func GetConnectionProperties() (driver, connString string, err error) {
	driver, err = env.GetDatabaseType()
	if err != nil {
		return "", "", err
	}

	user, err := env.GetDatabaseUser()
	if err != nil {
		return "", "", err
	}

	password, err := env.GetDatabasePassword()
	if err != nil {
		return "", "", err
	}

	host, err := env.GetDatabaseHost()
	if err != nil {
		return "", "", err
	}

	port, err := env.GetDatabasePort()
	if err != nil {
		return "", "", err
	}

	name, err := env.GetDatabaseName()
	if err != nil {
		return "", "", err
	}

	return driver, (&url.URL{
		Scheme: driver,
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}).String(), nil
}
