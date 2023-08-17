package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/darylhjd/oams/backend/internal/env"
)

const (
	connectionNamespace = "database/connection"
)

const (
	sslMode     = "sslmode"
	sslRootCert = "sslrootcert"
)

type DB struct {
	// Conn holds the base connection to the database. Usually, one will not directly use this for interactions with
	// the database.
	Conn *sql.DB

	// queryable is used to support transactions.
	queryable qrm.Queryable
}

// WithTx returns a new DB object with following transaction. The new DB should not be used once the transaction has
// been committed.
func (d *DB) WithTx(tx *sql.Tx) *DB {
	return &DB{d.Conn, tx}
}

func (d *DB) Close() error {
	return d.Conn.Close()
}

// Connect and return an interface to the OAMS database.
func Connect(ctx context.Context) (*DB, error) {
	return ConnectDB(ctx, env.GetDatabaseName())
}

// ConnectDB is similar to Connect but allows you to specify a specific database to use.
func ConnectDB(ctx context.Context, dbName string) (*DB, error) {
	_, connString := GetConnectionProperties(dbName)

	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("%s - error creating connection pool: %w", connectionNamespace, err)
	}

	// Check that connection is successful.
	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s - could not ping database connection: %w", connectionNamespace, err)
	}

	return &DB{conn, conn}, nil
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
