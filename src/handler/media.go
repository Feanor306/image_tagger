package handler

import (
	"strings"

	"github.com/feanor306/image_tagger/src/database"
	"github.com/feanor306/image_tagger/src/entities"
	"github.com/feanor306/image_tagger/src/validation"
	"github.com/feanor306/image_tagger/src/view/layout"
	"github.com/feanor306/image_tagger/src/view/media"
	"github.com/labstack/echo/v4"
)

// MediaHandler is a handler that manages media requests
type MediaHandler struct {
	DB *database.DB
}

// HandleMediaShowAll will render all media that contains certain tag
func (h MediaHandler) HandleMediaShowAll(c echo.Context) error {
	tag := entities.Tag{
		Id: c.QueryParam("tag"),
	}
	mediaPlural, err := h.DB.FindMedia(&tag, MAX_SIZE)

	if err != nil {
		return render(c, layout.Error(err))
	}

	return render(c, media.MediaByTag(mediaPlural))
}

// HandleMediaNew will render media creation form
func (h MediaHandler) HandleMediaNew(c echo.Context) error {
	tags, err := h.DB.GetAllTags(100)
	if err != nil {
		return render(c, layout.Error(err))
	}

	tagNames := make([]string, 0, len(tags))
	for _, t := range tags {
		tagNames = append(tagNames, t.Name)
	}

	return render(c, media.NewMedia(strings.Join(tagNames, ",")))
}

// HandleMediaCreate will handle media creation requests
// involving validation, file upload and saving to database
func (h MediaHandler) HandleMediaCreate(c echo.Context) error {
	m := entities.Media{
		Name: c.FormValue("name"),
	}

	err := validation.ValidateMedia(c, &m)
	if err != nil {
		return render(c, layout.Error(err))
	}

	err = handleFileUpload(c, m.Id)
	if err != nil {
		return render(c, layout.Error(err))
	}

	err = h.DB.CreateMedia(&m)
	if err != nil {
		return render(c, layout.Error(err))
	}

	return render(c, media.MediaPartial(m))
}
