package store

import (
	"github.com/robzienert/lever/model"
	"golang.org/x/net/context"
)

// FeatureStore is the repository for interacting with the feature backends.
type FeatureStore interface {
	Get(string) (*model.Feature, error)
	GetByNamespace(string, string) (*model.Feature, error)
	GetList() ([]*model.Feature, error)
	GetListByNamespace(string) ([]*model.Feature, error)
	Upsert(*model.Feature) error
	Delete(*model.Feature) error
}

// GetFeature will proxy to the net.Context's feature storage backend to get
// an individual feature by key and optionally namespace.
func GetFeature(c context.Context, namespace string, key string) (*model.Feature, error) {
	if namespace == "" {
		return FromContext(c).Features().Get(key)
	}
	return FromContext(c).Features().GetByNamespace(namespace, key)
}

// GetFeatureList will proxy to the net.Context's feature storage backend to
// get a list of features, optionally by namespace.
func GetFeatureList(c context.Context, namespace string) ([]*model.Feature, error) {
	if namespace == "" {
		return FromContext(c).Features().GetList()
	}
	return FromContext(c).Features().GetListByNamespace(namespace)
}

// UpsertFeature will proxy the net.Context's feature storage to save a feature.
func UpsertFeature(c context.Context, feature *model.Feature) error {
	return FromContext(c).Features().Upsert(feature)
}

// DeleteFeature will proxy the net.Context's feature storage to delete features.
func DeleteFeature(c context.Context, feature *model.Feature) error {
	return FromContext(c).Features().Delete(feature)
}
