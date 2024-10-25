package interfaces

import "github.com/feanor306/image_tagger/src/entities"

// MediaStore represents media storage functionality
type MediaStore interface {
	CreateMedia(media *entities.Media) error
	CreateMediaTags(media *entities.Media) error
	FindMedia(tag *entities.Tag, count int) ([]entities.Media, error)
}
