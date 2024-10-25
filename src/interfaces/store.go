package interfaces

// Store contains database specific functionality
type Store interface {
	InitDatabase() error
	Close() error
}

// CompositeStore represents the full storage functionality
type CompositeStore interface {
	Store
	TagStore
	MediaStore
}
