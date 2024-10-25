package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/feanor306/image_tagger/src/entities"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// MEDIA_NAME_MAX_LENGTH is the maximum allowed size for Media.Name
const MEDIA_NAME_MAX_LENGTH = 30

// TagValue is used for unmarshalling tags
type TagValue struct {
	Value string `json:"value"`
}

// ValidateMedia will validate media params, set id if missing
// set public filename as well as set tags
func ValidateMedia(c echo.Context, media *entities.Media) error {
	id := uuid.New().String()
	if len(media.Id) == 0 {
		media.Id = id
	}

	if len(media.Name) == 0 {
		return errors.New("file name can not be empty")
	}

	if len(media.Name) > MEDIA_NAME_MAX_LENGTH {
		return errors.New(fmt.Sprintf("tag name can be max %d characters", MEDIA_NAME_MAX_LENGTH))
	}

	if len(media.Filename) == 0 {
		fileName, err := getPublicFileName(c, media.Id)
		if err != nil {
			return err
		}
		media.Filename = fileName
	}

	tags := c.FormValue("tags")
	if len(tags) == 0 {
		return errors.New("media must have at least one tag")
	}

	var tagValues []TagValue
	err := json.Unmarshal([]byte(tags), &tagValues)
	if err != nil {
		return err
	}

	var ts []entities.Tag
	for _, tag := range tagValues {
		ts = append(ts, entities.Tag{
			Name: tag.Value,
		})
	}

	media.Tags = ts

	return nil
}

// getPublicFileName will return publically accessible url path to file
func getPublicFileName(c echo.Context, id string) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	scheme := "http:/"
	if c.Request().TLS != nil {
		scheme = "https:/"
	}

	fns := strings.Split(file.Filename, ".")
	return strings.Join([]string{scheme, c.Request().Host, "files", fmt.Sprintf("%s.%s", id, fns[len(fns)-1])}, "/"), nil
}
