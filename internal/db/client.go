package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func NewDB(dbURL string, logger log.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Println("Could not initialize DB with given config variables")
		return nil, err
	}
	logger.Println("Connected to DB at: " + dbURL)
	return db, nil
}