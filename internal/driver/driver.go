package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	SQL *sql.DB
}

var (
	databaseConnection = &Database{}
)

const (
	maxOpenDatabaseConnection = 10
	maxIdleDatabaseConnection = 5
	maxDatabaseLifeTime       = 5 * time.Minute
)

func ConnectSQL(dsn string) (*Database, error) {
	database, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	database.SetMaxOpenConns(maxOpenDatabaseConnection)
	database.SetMaxIdleConns(maxIdleDatabaseConnection)
	database.SetConnMaxLifetime(maxDatabaseLifeTime)

	databaseConnection.SQL = database

	err = checkDB(database)

	if err != nil {
		return nil, err
	}

	return databaseConnection, nil
}

func checkDB(database *sql.DB) error {
	err := database.Ping()

	if err != nil {
		return err
	}

	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
