package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/feanor306/image_tagger/src/config"
)

// db handles database connection
type db struct {
	conn *sql.DB
	conf *config.Config
}

// dbInstance is a singleton db instance
var dbInstance *db

// GetDatabase creates or returns the database instance
func GetDatabase(conf *config.Config) (*db, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName)

	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	dbInstance = &db{
		conn: dbConn,
		conf: conf,
	}
	return dbInstance, nil
}

// Close closes the database connection.
func (db *db) Close() error {
	log.Printf("Disconnected from database: %s", db.conf.DbName)
	return db.conn.Close()
}

// InitDatabase will create all the tables and indices
func (db *db) InitDatabase() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id text NOT NULL,
			name text NOT NULL,
			PRIMARY KEY (id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS tags_id_index ON tags USING btree (id text_pattern_ops);
		CREATE UNIQUE INDEX IF NOT EXISTS tags_name_index ON tags USING btree (name text_pattern_ops);

		CREATE TABLE IF NOT EXISTS images (
			id text NOT NULL,
			name text NOT NULL,
			filename text NOT NULL,
			PRIMARY KEY (id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS images_id_index ON images USING btree (filename text_pattern_ops);
		CREATE INDEX IF NOT EXISTS images_name_index ON images USING btree (name text_pattern_ops);
		CREATE INDEX IF NOT EXISTS images_filename_index ON images USING btree (filename text_pattern_ops);

		CREATE TABLE IF NOT EXISTS image_tags (
			image_id text NOT NULL,
			tag_id text NOT NULL,
			PRIMARY KEY (image_id, tag_id),
			FOREIGN KEY (image_id) REFERENCES images (id),
			FOREIGN KEY (tag_id) REFERENCES tags (id)
		);
		CREATE INDEX idx_image_tags_image_id ON image_tags (image_id);
		CREATE INDEX idx_image_tags_tag_id ON image_tags (tag_id);
	`)

	if err != nil {
		return err
	}
	return nil
}
