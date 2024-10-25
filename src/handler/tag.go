package handler

import (
	"github.com/feanor306/image_tagger/src/database"
	"github.com/feanor306/image_tagger/src/entities"
	"github.com/feanor306/image_tagger/src/validation"
	"github.com/feanor306/image_tagger/src/view/layout"
	"github.com/feanor306/image_tagger/src/view/tag"
	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	DB *database.DB
}

func (h TagHandler) HandleTagShowAll(c echo.Context) error {
	allTags, err := h.DB.GetAllTags(MAX_SIZE)
	if err != nil {
		return render(c, layout.Error(err))
	}

	return render(c, tag.TagsAll(allTags))
}

func (h TagHandler) HandleTagCreate(c echo.Context) error {
	t := entities.Tag{
		Name: c.FormValue("name"),
	}

	err := validation.ValidateTag(&t)
	if err != nil {
		return render(c, layout.Error(err))
	}

	err = h.DB.CreateTag(&t)
	if err != nil {
		return render(c, layout.Error(err))
	}

	return render(c, tag.TagPartial(t))
}
