package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/darylhjd/oams/backend/internal/env"
)

const (
	sslMode     = "sslmode"
	sslRootCert = "sslrootcert"
)

// DB contains the database connection pool and the query interface to the database.
type DB struct {
	C *pgxpool.Pool
	Q *Queries
}

// NewTx creates a new transaction with corresponding Tx and Queries object.
// The caller is responsible for committing or rolling back the transaction.
// For nested transactions (save points), caller may provide an existing Queries object.
// Else, a normal transaction is created from the database connection pool.
func (d *DB) NewTx(ctx context.Context, q *Queries) (pgx.Tx, *Queries, error) {
	tx, err := d.C.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	txq := d.Q.WithTx(tx)
	if q != nil {
		txq = q.WithTx(tx)
	}

	return tx, txq, nil
}

// Close the database connection.
func (d *DB) Close() {
	d.C.Close()
}

// Connect and return an interface to the OAMS database.
func Connect(ctx context.Context) (*DB, error) {
	return ConnectDB(ctx, env.GetDatabaseName())
}

// ConnectDB is similar to Connect but allows you to specify a specific database to use.
func ConnectDB(ctx context.Context, dbName string) (*DB, error) {
	_, connString := GetConnectionProperties(dbName)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	// Check that connection is successful.
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{pool, New(pool)}, nil
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
