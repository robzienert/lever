package mock

import "github.com/robzienert/lever/store"

func LoadFeatureStore(featureStore *FeatureStore) store.Store {
	return store.New("mock", &BreadcrumbStore{}, featureStore)
}
