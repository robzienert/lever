package cql

import (
	"github.com/robzienert/lever/model"
	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
)

type featureStore struct {
	session *gocql.Session
}

func (s *featureStore) Get(key string) (*model.Feature, error) {
	return s.one("SELECT * FROM features WHERE key = ?", key)
}

func (s *featureStore) GetByNamespace(namespace string, key string) (*model.Feature, error) {
	return s.one("SELECT * FROM features_namespaced WHERE namespace = ? AND key = ?", namespace, key)
}

func (s *featureStore) one(query string, args ...interface{}) (*model.Feature, error) {
	data := make(cqlResult, 0)
	err := s.session.Query(query, args...).MapScan(data)
	if err != nil && err.Error() != notFoundError {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"q":   query,
			"a":   args,
		}).Error("Could not execute CQL query")
		return nil, err
	}
	return marshalFeature(data), nil
}

func (s *featureStore) GetList() ([]*model.Feature, error) {
	return s.all("SELECT * FROM features")
}

func (s *featureStore) GetListByNamespace(namespace string) ([]*model.Feature, error) {
	return s.all("SELECT * FROM features_namespaced WHERE namespace = ?", namespace)
}

func (s *featureStore) all(query string, args ...interface{}) ([]*model.Feature, error) {
	iter := s.session.Query(query, args...).Iter()

	var all []*model.Feature
	data := make(cqlResult, 0)
	for iter.MapScan(data) {
		all = append(all, marshalFeature(data))
	}
	if err := iter.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"q":   query,
			"a":   args,
		}).Error("Could not execute CQL query")
		return nil, err
	}

	return all, nil
}

func (s *featureStore) Upsert(feature *model.Feature) (err error) {
	if feature.Namespace == "" {
		err = s.session.Query(
			`UPDATE features SET type = ?, value = ?, gate_value = ?, gate_groups = ?, gate_actors = ?,
gate_actor_percent = ?, gate_percent_of_time = ?, date_created = ?, last_updated = ?
WHERE key = ?`,
			feature.Type,
			feature.Value,
			feature.Gate.Value,
			feature.Gate.Groups,
			feature.Gate.Actors,
			feature.Gate.ActorPercent,
			feature.Gate.PercentOfTime,
			feature.DateCreated,
			feature.LastUpdated,
			feature.Key,
		).Exec()
	} else {
		err = s.session.Query(
			`UPDATE features_namespaced SET type = ?, value = ?, gate_value = ?, gate_groups = ?,
gate_actors = ?, gate_actor_percent = ?, gate_percent_of_time = ?, date_created = ?, last_updated = ?
WHERE namespace = ? AND key = ?`,
			feature.Type,
			feature.Value,
			feature.Gate.Value,
			feature.Gate.Groups,
			feature.Gate.Actors,
			feature.Gate.ActorPercent,
			feature.Gate.PercentOfTime,
			feature.DateCreated,
			feature.LastUpdated,
			feature.Namespace,
			feature.Key,
		).Exec()
	}
	return
}

func (s *featureStore) Delete(feature *model.Feature) (err error) {
	if feature.Namespace == "" {
		err = s.session.Query("DELETE FROM features WHERE key = ?", feature.Key).Exec()
	} else {
		err = s.session.Query("DELETE FROM features_namespaced WHERE namespace = ? AND key = ?", feature.Namespace, feature.Key).Exec()
	}
	return
}
