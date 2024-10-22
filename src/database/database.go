package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/feanor306/image_tagger/src/config"
)

type db struct {
	db   *sql.DB
	conf *config.Config
}

var dbInstance *db

// GetDatabase creates or returns the database instance
func GetDatabase(conf *config.Config) *db {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName)

	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	dbInstance = &db{
		db:   dbConn,
		conf: conf,
	}
	return dbInstance
}

// Close closes the database connection.
func (db *db) Close() error {
	log.Printf("Disconnected from database: %s", db.conf.DbName)
	return db.db.Close()
}
