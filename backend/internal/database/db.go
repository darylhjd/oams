package database

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/qrm"
)

// DB provides an interface to the database connection.
type DB struct {
	// Conn holds the base connection to the database. Usually, one will not directly use this for interactions with
	// the database.
	Conn *sql.DB

	// qe is used to support transactions.
	qe queryExecutable
}

// queryExecutable allows for the use of transaction-enabled DB.
type queryExecutable interface {
	qrm.Queryable
	qrm.Executable
}

// AsTx returns a new DB and *sql.Tx object, with the new DB making use of the following transaction. The caller should
// commit or rollback as required, and the provided DB should not be used once a commit or rollback is done.
func (d *DB) AsTx(ctx context.Context, opts *sql.TxOptions) (*DB, *sql.Tx, error) {
	tx, err := d.Conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return &DB{d.Conn, tx}, tx, nil
}

func (d *DB) Close() error {
	return d.Conn.Close()
}
