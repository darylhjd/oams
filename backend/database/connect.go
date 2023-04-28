package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"

	"github.com/darylhjd/oats/backend/env"
)

const Namespace = "database"

// DB contains the database connection pool and the query interface to the database.
type DB struct {
	Db *sql.DB
	Q  *Queries
}

// Close the database connection.
func (d *DB) Close() error {
	return d.Db.Close()
}

// Connect and return an interface to the Oats database.
func Connect() (*DB, error) {
	name, err := env.GetDatabaseName()
	if err != nil {
		return nil, err
	}

	return ConnectDB(name)
}

// ConnectDB is similar to Connect but allows you to specify a specific database to use.
func ConnectDB(dbName string) (*DB, error) {
	driver, connString, err := GetConnectionProperties(dbName)
	if err != nil {
		return nil, fmt.Errorf("%s - cannot parse connection properties: %w", Namespace, err)
	}

	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}

	// Check that connection is successful.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db, New(db)}, nil
}

// GetConnectionProperties returns the connection strings required to connect to a database.
func GetConnectionProperties(dbName string) (driver, connString string, err error) {
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

	// TODO: Set up SSL.
	params := url.Values{}
	params.Set("sslmode", "disable")

	return driver, (&url.URL{
		Scheme:   driver,
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     dbName,
		RawQuery: params.Encode(),
	}).String(), nil
}
