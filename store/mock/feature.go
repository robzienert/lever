package mock

import "github.com/robzienert/lever/model"

type FeatureStore struct {
	GetFn     func(namespace string, key string) (*model.Feature, error)
	GetListFn func(namespace string) ([]*model.Feature, error)
	UpsertFn  func(feature *model.Feature) error
	DeleteFn  func(feature *model.Feature) error
}

func (s *FeatureStore) Get(key string) (*model.Feature, error) {
	return s.GetFn("", key)
}

func (s *FeatureStore) GetByNamespace(namespace string, key string) (*model.Feature, error) {
	return s.GetFn(namespace, key)
}

func (s *FeatureStore) GetList() ([]*model.Feature, error) {
	return s.GetListFn("")
}

func (s *FeatureStore) GetListByNamespace(namespace string) ([]*model.Feature, error) {
	return s.GetListFn(namespace)
}

func (s *FeatureStore) Upsert(feature *model.Feature) error {
	return s.UpsertFn(feature)
}

func (s *FeatureStore) Delete(feature *model.Feature) error {
	return s.DeleteFn(feature)
}
