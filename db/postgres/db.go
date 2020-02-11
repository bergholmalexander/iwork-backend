package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"go-api/db"
)

type PostgresDBStore struct {
	database *sql.DB
}

var CreateError = errors.New("create error")

func (p PostgresDBStore) Close() {
	p.database.Close()
}

func NewPostgresDataStore(dbUrl string) (*db.DataStore, error) {
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	dbStore := &PostgresDBStore{database: database}
	return &db.DataStore{
		WorkspaceProvider: dbStore,
		BookingProvider:   dbStore,
		OfferingProvider:  dbStore,
	}, nil
}