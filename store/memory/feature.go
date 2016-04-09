package memory

import (
	"sync"

	"github.com/robzienert/lever/model"
)

type featureStore struct {
	features []*model.Feature
	lock     sync.RWMutex
}

func (s *featureStore) Get(key string) (*model.Feature, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.GetByNamespace("", key)
}

func (s *featureStore) GetByNamespace(namespace string, key string) (*model.Feature, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, f := range s.features {
		if f.Namespace == namespace && f.Key == key {
			return f, nil
		}
	}
	return nil, nil
}

func (s *featureStore) GetList() ([]*model.Feature, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.GetListByNamespace("")
}

func (s *featureStore) GetListByNamespace(namespace string) ([]*model.Feature, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var features []*model.Feature
	for _, f := range s.features {
		if f.Namespace == namespace {
			features = append(features, f)
		}
	}
	return features, nil
}

func (s *featureStore) Upsert(feature *model.Feature) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, f := range s.features {
		if f.Namespace == feature.Namespace && f.Key == feature.Key {
			s.features[i] = feature
			return nil
		}
	}
	s.features = append(s.features, feature)
	return nil
}

func (s *featureStore) Delete(feature *model.Feature) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, f := range s.features {
		if f.Namespace == feature.Namespace && f.Key == feature.Key {
			s.features = append(s.features[:i], s.features[i+1:]...)
		}
	}
	return nil
}
