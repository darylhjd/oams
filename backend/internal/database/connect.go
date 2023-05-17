package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"

	"github.com/darylhjd/oams/backend/internal/env"
)

const Namespace = "database"

const (
	sslMode     = "sslmode"
	sslRootCert = "sslrootcert"
)

// DB contains the database connection pool and the query interface to the database.
type DB struct {
	C *sql.DB
	Q *Queries
}

// Close the database connection.
func (d *DB) Close() error {
	return d.C.Close()
}

// Connect and return an interface to the OAMS database.
func Connect() (*DB, error) {
	return ConnectDB(env.GetDatabaseName())
}

// ConnectDB is similar to Connect but allows you to specify a specific database to use.
func ConnectDB(dbName string) (*DB, error) {
	driver, connString := GetConnectionProperties(dbName)

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
func GetConnectionProperties(dbName string) (driver, connString string) {
	driver = env.GetDatabaseType()

	// Set SSL mode.
	params := url.Values{}
	params.Set(sslMode, env.GetDatabaseSSLMode())
	params.Set(sslRootCert, env.GetDatabaseSSLRootCertLocation())

	return driver, (&url.URL{
		Scheme:   driver,
		User:     url.UserPassword(env.GetDatabaseUser(), env.GetDatabasePassword()),
		Host:     fmt.Sprintf("%s:%s", env.GetDatabaseHost(), env.GetDatabasePort()),
		Path:     dbName,
		RawQuery: params.Encode(),
	}).String()
}
