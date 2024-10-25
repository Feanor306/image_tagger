package interfaces

import "github.com/feanor306/image_tagger/src/entities"

// TagStore contains tag functionality
type TagStore interface {
	CreateTag(tag *entities.Tag) error
	GetAllTags(count int) ([]entities.Tag, error)
	GetTag(id string) (*entities.Tag, error)
}
