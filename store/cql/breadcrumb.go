package cql

import (
	"github.com/robzienert/lever/model"
	"github.com/gocql/gocql"
)

type breadcrumbStore struct {
	session *gocql.Session
}

func (s *breadcrumbStore) GetList() ([]*model.Breadcrumb, error) {
	iter := s.session.Query("SELECT * FROM audit").Iter()

	var breadcrumbs []*model.Breadcrumb
	var result map[string]interface{}
	for iter.MapScan(result) {
		b := marshalBreadcrumb(result)
		breadcrumbs = append(breadcrumbs, b)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return breadcrumbs, nil
}

func (s *breadcrumbStore) Create(b *model.Breadcrumb) error {
	err := s.session.Query(
		`INSERT INTO audit (action, actor, fields, date_created) VALUES (?, ?, ?, ?)`,
		b.Action, b.Actor, b.Fields, b.DateCreated).Exec()
	if err != nil {
		return err
	}
	return nil
}
