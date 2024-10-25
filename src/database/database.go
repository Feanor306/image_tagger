package database

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/feanor306/image_tagger/src/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB handles database connection
type DB struct {
	conn *sql.DB
	sq   squirrel.StatementBuilderType
}

// dbInstance is a singleton db instance
var dbInstance *DB

// GetDatabase creates or returns the database instance
func GetDatabase(conf *config.Config) (*DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName)

	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	dbInstance = &DB{
		conn: dbConn,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(dbConn),
	}
	return dbInstance, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// InitDatabase will create all the tables and indices
func (db *DB) InitDatabase() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id text NOT NULL,
			name text NOT NULL,
			PRIMARY KEY (id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS tags_id_index ON tags USING btree (id text_pattern_ops);
		CREATE UNIQUE INDEX IF NOT EXISTS tags_name_index ON tags USING btree (name text_pattern_ops);

		CREATE TABLE IF NOT EXISTS media (
			id text NOT NULL,
			name text NOT NULL,
			filename text NOT NULL,
			PRIMARY KEY (id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS media_id_index ON media USING btree (filename text_pattern_ops);
		CREATE INDEX IF NOT EXISTS media_name_index ON media USING btree (name text_pattern_ops);
		CREATE INDEX IF NOT EXISTS media_filename_index ON media USING btree (filename text_pattern_ops);

		CREATE TABLE IF NOT EXISTS media_tags (
			media_id text NOT NULL,
			tag_id text NOT NULL,
			PRIMARY KEY (media_id, tag_id),
			FOREIGN KEY (media_id) REFERENCES media (id),
			FOREIGN KEY (tag_id) REFERENCES tags (id)
		);
		CREATE INDEX IF NOT EXISTS media_tag_media_id_index ON media_tags (media_id);
		CREATE INDEX IF NOT EXISTS media_tag_tag_id_index ON media_tags (tag_id);
	`)

	if err != nil {
		return err
	}
	return nil
}
