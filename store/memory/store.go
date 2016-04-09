package memory

import "github.com/robzienert/lever/store"

// Load will create a new memory store.
func Load() store.Store {
	return New()
}

// New will initialize a new memory store repository.
func New() store.Store {
	return store.New(
		"memory",
		&breadcrumbStore{},
		&featureStore{},
	)
}
