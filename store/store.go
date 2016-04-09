package store

// Store is a repository for all the storage backends.
type Store interface {
	Name() string
	Breadcrumbs() BreadcrumbStore
	Features() FeatureStore
}

type store struct {
	name        string
	breadcrumbs BreadcrumbStore
	features    FeatureStore
}

func (s *store) Name() string                 { return s.name }
func (s *store) Breadcrumbs() BreadcrumbStore { return s.breadcrumbs }
func (s *store) Features() FeatureStore       { return s.features }

// New will create a new Store with the provided concrete backends.
func New(name string, breadcrumbs BreadcrumbStore, features FeatureStore) Store {
	return &store{name, breadcrumbs, features}
}
