package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/feanor306/image_tagger/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDatabase(t *testing.T) *DB {
	dbConn, err := sql.Open("pgx", "postgres://test:test@localhost:5432/test")
	require.NoError(t, err, "creating DB conn should not fail")

	testDb := &DB{
		conn: dbConn,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(dbConn),
	}

	err = testDb.InitDatabase()
	require.NoError(t, err, "initializing db should not fail")

	return testDb
}

func closeTestDatabase(t *testing.T, db *DB) {
	err := db.Close()
	require.NoError(t, err, "closing DB should not fail")
}

func TestCreateTag(t *testing.T) {
	db := getTestDatabase(t)
	defer closeTestDatabase(t, db)

	cases := []struct {
		tag entities.Tag
		err error
		msg string
	}{
		{
			tag: entities.Tag{
				Id:   "One",
				Name: "Frodo",
			},
			err: nil,
			msg: "no error on first tag",
		},
		{
			tag: entities.Tag{
				Id:   "One",
				Name: "Gandalf",
			},
			err: errors.New("ERROR: duplicate key value violates unique constraint \"tags_pkey\" (SQLSTATE 23505)"),
			msg: "duplicate tag id",
		},
		{
			tag: entities.Tag{
				Id:   "Two",
				Name: "Frodo",
			},
			err: errors.New("ERROR: duplicate key value violates unique constraint \"tags_name_index\" (SQLSTATE 23505)"),
			msg: "duplicate tag name",
		},
	}

	for _, c := range cases {
		err := db.CreateTag(&c.tag)
		if err == nil {
			assert.Equal(t, c.err, err, c.msg)
		} else {
			assert.Equal(t, c.err.Error(), err.Error(), c.msg)
		}
	}
}

func TestGetTag(t *testing.T) {
	db := getTestDatabase(t)
	defer closeTestDatabase(t, db)

	existingId := "prince"
	tag := entities.Tag{
		Id:   existingId,
		Name: "Legolas",
	}
	err := db.CreateTag(&tag)
	require.NoError(t, err, "creating tag should not fail")

	existingTag, err := db.GetTag(existingId)
	require.NoError(t, err, "finding tag should not fail")
	assert.Equal(t, tag, *existingTag, "db tag should match created tag")

	nonExistingTag, err := db.GetTag("non-existing-id")
	require.Nil(t, nonExistingTag, "should not find tag by non existing id")
	assert.Equal(t, errors.New("tag doesn't exist in database"), err, "trying to find non existing tag should return err")
}

func TestGetAllTags(t *testing.T) {
	db := getTestDatabase(t)
	defer closeTestDatabase(t, db)

	tags := []entities.Tag{
		{
			Id:   "noldor-king-1",
			Name: "Finwe",
		},
		{
			Id:   "noldor-king-2",
			Name: "Feanor",
		},
		{
			Id:   "noldor-king-3",
			Name: "Fingolfin",
		},
	}

	for _, tg := range tags {
		err := db.CreateTag(&tg)
		require.NoError(t, err, "creating tag should not fail")
	}

	tagsDb, err := db.GetAllTags(99)
	require.NoError(t, err, "finding all tags should not fail")
	assert.GreaterOrEqual(t, len(tagsDb), len(tags), "find all tags should find inserted tags")
}

func TestCreateMedia(t *testing.T) {
	db := getTestDatabase(t)
	defer closeTestDatabase(t, db)

	existingTag := entities.Tag{
		Id:   "dragon-1",
		Name: "Glaurung",
	}
	nonExistingTag := entities.Tag{
		Id:   "i-dont-exist",
		Name: "Imagination",
	}
	err := db.CreateTag(&existingTag)
	require.Nil(t, err, "creating tag should not fail")

	cases := []struct {
		media entities.Media
		err   error
		msg   string
	}{
		{
			media: entities.Media{
				Id:       "numenor-island",
				Name:     "Numenor",
				Filename: "numenor.jpg",
				Tags:     []entities.Tag{existingTag},
			},
			err: nil,
			msg: "creating valid media with valid tag",
		},
		{
			media: entities.Media{
				Id:       "numenor-island",
				Name:     "Invalid Media 1",
				Filename: "invalid_media.jpg",
				Tags:     []entities.Tag{existingTag},
			},
			err: errors.New("ERROR: duplicate key value violates unique constraint \"media_pkey\" (SQLSTATE 23505)"),
			msg: "media must have unique id",
		},
		{
			media: entities.Media{
				Id:       "invalid-media-id-1",
				Name:     "Invalid Media Name",
				Filename: "Invalid Media Filename",
				Tags:     []entities.Tag{existingTag, nonExistingTag},
			},
			err: errors.New("not all tags provided for media exist - create them first"),
			msg: "media must have valid existing tags",
		},
	}

	for _, cs := range cases {
		err := db.CreateMedia(&cs.media)
		if err == nil {
			assert.Equal(t, cs.err, err, cs.msg)
		} else {
			assert.Equal(t, cs.err.Error(), err.Error(), cs.msg)
		}
	}
}
