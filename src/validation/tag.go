package validation

import (
	"errors"
	"fmt"

	"github.com/feanor306/image_tagger/src/entities"
	"github.com/google/uuid"
)

// TAG_MAX_LENGTH is the maximum allowed size for Tag.Name
const TAG_MAX_LENGTH = 30

// ValidateTag will validate a single tag and set its Id if empty
func ValidateTag(tag *entities.Tag) error {
	if len(tag.Id) == 0 {
		tag.Id = uuid.New().String()
	}

	if len(tag.Name) == 0 {
		return errors.New("tag name can not be empty")
	}

	if len(tag.Name) > TAG_MAX_LENGTH {
		return errors.New(fmt.Sprintf("tag name can be max %d characters", TAG_MAX_LENGTH))
	}

	return nil
}
