package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/feanor306/image_tagger/src/config"
	"github.com/feanor306/image_tagger/src/entities"

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
		CREATE UNIQUE INDEX IF NOT EXISTS media_id_index ON media USING btree (id text_pattern_ops);
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

// CreateTag will write a single Tag to database
func (db *DB) CreateTag(tag *entities.Tag) error {
	_, err := db.sq.
		Insert("tags").
		Columns("id", "name").
		Values(tag.Id, tag.Name).
		Exec()

	return err
}

// CreateMediaTags will create join table records for media_tags
func (db *DB) CreateMediaTags(media *entities.Media) error {
	tagNames := make([]string, 0, len(media.Tags))
	for _, tag := range media.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	rows, err := db.sq.
		Select("id", "name").
		From("tags").
		Where(squirrel.Eq{"name": tagNames}).
		Query()

	if err != nil {
		return err
	}

	dbTags := make([]entities.Tag, 0, len(media.Tags))
	for rows.Next() {
		var tag entities.Tag
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil {
			return err
		}
		dbTags = append(dbTags, tag)
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	if len(dbTags) != len(media.Tags) {
		return errors.New("not all tags provided for media exist - create them first")
	}

	insertSql := db.sq.
		Insert("media_tags").
		Columns("media_id", "tag_id")

	for _, tag := range dbTags {
		insertSql = insertSql.Values(media.Id, tag.Id)
	}

	_, err = insertSql.Exec()
	return err
}

// GetAllTags will retrieve all tags from database
func (db *DB) GetAllTags(count int) ([]entities.Tag, error) {
	rows, err := db.sq.
		Select("id", "name").
		From("tags").
		Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	tags := make([]entities.Tag, 0, count)
	for rows.Next() {
		var tag entities.Tag
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// GetTag will return a single tag by id
func (db *DB) GetTag(id string) (*entities.Tag, error) {
	var tag entities.Tag
	err := db.sq.
		Select("id", "name").
		From("tags").
		Where(squirrel.Eq{"id": id}).
		QueryRow().
		Scan(&tag.Id, &tag.Name)

	if len(tag.Id) == 0 {
		return nil, errors.New("tag doesn't exist in database")
	}

	return &tag, err
}

// CreateMedia will create media in database
func (db *DB) CreateMedia(media *entities.Media) error {
	_, err := db.sq.
		Insert("media").
		Columns("id", "name", "filename").
		Values(media.Id, media.Name, media.Filename).
		Exec()

	if err != nil {
		return err
	}

	return db.CreateMediaTags(media)
}

// GetMediaIdsByTag will return media id-s
// that are related to a tag by its id
func (db *DB) GetMediaIdsByTag(tag *entities.Tag, count int) ([]string, error) {
	rows, err := db.sq.
		Select("media_id").
		From("media_tags").
		Where(squirrel.Eq{"tag_id": tag.Id}).
		Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	result := make([]string, 0, count)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return result, err
		}
		result = append(result, id)
	}

	return result, rows.Err()
}

// FindMedia will return all media
// that is connected to a tag by id
// with all media props set
func (db *DB) FindMedia(tag *entities.Tag, count int) ([]entities.Media, error) {
	// split up in 2 queries because of squirrel query nesting issues
	mediaIdPlural, err := db.GetMediaIdsByTag(tag, count)

	if err != nil {
		return nil, err
	}

	rows, err := db.sq.
		Select("media.*", "tags.*").
		From("media").
		Join("media_tags ON media.id = media_tags.media_id").
		Join("tags ON tags.id = media_tags.tag_id").
		Where(squirrel.Eq{"media_tags.media_id": mediaIdPlural}).
		Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	mediaResult := make([]entities.Media, 0, count)
	for rows.Next() {
		var mediaRow entities.Media
		var tag entities.Tag
		if err := rows.Scan(&mediaRow.Id, &mediaRow.Name, &mediaRow.Filename, &tag.Id, &tag.Name); err != nil {
			return mediaResult, err
		}

		exists := -1
		for idx, mr := range mediaResult {
			if mr.Id == mediaRow.Id {
				exists = idx
			}
		}

		if exists >= 0 && exists < len(mediaResult) {
			mediaResult[exists].Tags = append(mediaResult[exists].Tags, tag)
		} else {
			mediaRow.Tags = []entities.Tag{tag}
			mediaResult = append(mediaResult, mediaRow)
		}
	}

	return mediaResult, rows.Err()
}
