package mock

import "github.com/robzienert/lever/model"

type BreadcrumbStore struct{}

func (s *BreadcrumbStore) GetList() ([]*model.Breadcrumb, error) {
	return nil, nil
}

func (s *BreadcrumbStore) Create(b *model.Breadcrumb) error {
	return nil
}
