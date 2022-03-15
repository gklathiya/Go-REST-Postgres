package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifeTime = 5 * time.Minute

// ConnectSQL create database pool for Postgre
func ConnectSQL(dbs string) (*DB, error) {
	d, err := NewDatabase(dbs)
	if err != nil {
		panic(err) // Programme does not work after this
	}
	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetConnMaxIdleTime(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifeTime)
	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB try to ping databse connection
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dbs string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbs)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
