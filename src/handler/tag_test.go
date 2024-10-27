package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/feanor306/image_tagger/src/entities"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestTagStore struct {
	Tags []entities.Tag
}

func (t TestTagStore) CreateTag(tag *entities.Tag) error {
	t.Tags = append(t.Tags, *tag)
	return nil
}

func (t TestTagStore) GetAllTags(count int) ([]entities.Tag, error) {
	return t.Tags, nil
}

func (t TestTagStore) GetTag(id string) (*entities.Tag, error) {
	for _, tg := range t.Tags {
		if tg.Id == id {
			return &tg, nil
		}
	}
	return nil, errors.New("tag not found")
}

func getTestTagHandler() TagHandler {
	return TagHandler{
		TestTagStore{
			Tags: make([]entities.Tag, 0, 100),
		},
	}
}

func TestHandleTagShowAll(t *testing.T) {
	tagHandler := getTestTagHandler()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, tagHandler.HandleTagShowAll(c))
	assert.Contains(t, rec.Body.String(), "<form hx-post=")
	assert.Contains(t, rec.Body.String(), "id=\"tag-list\"")
}

func TestHandleTagCreate(t *testing.T) {
	tagHandler := getTestTagHandler()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/tags", nil)
	req.Form = make(url.Values)
	req.Form.Set("name", "")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, tagHandler.HandleTagCreate(c))
	assert.Contains(t, rec.Body.String(), "tag name can not be empty")

	req = httptest.NewRequest(http.MethodPost, "/tags", nil)
	req.Form = make(url.Values)
	req.Form.Set("name", "Freddy")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	require.NoError(t, tagHandler.HandleTagCreate(c))
	assert.Contains(t, rec.Body.String(), "Freddy")
}
