package memory

import (
	"sync"

	"github.com/robzienert/lever/model"
)

type breadcrumbStore struct {
	breadcrumbs []*model.Breadcrumb
	lock        sync.RWMutex
}

func (s *breadcrumbStore) GetList() ([]*model.Breadcrumb, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.breadcrumbs, nil
}

func (s *breadcrumbStore) Create(b *model.Breadcrumb) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.breadcrumbs = append(s.breadcrumbs, b)
	return nil
}
