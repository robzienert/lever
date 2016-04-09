package cql

import (
	"errors"
	"time"

	"github.com/robzienert/lever/model"
	"github.com/Sirupsen/logrus"
)

type cqlResult map[string]interface{}

func marshalBreadcrumb(d cqlResult) *model.Breadcrumb {
	if d == nil || len(d) == 0 {
		return nil
	}
	defer recoverMarshalPanic("breadcrumb", d)

	b := &model.Breadcrumb{}
	if v, ok := d["action"]; ok {
		b.Action = v.(string)
	}
	if v, ok := d["actor"]; ok {
		b.Actor = v.(string)
	}
	if v, ok := d["date_created"]; ok {
		b.DateCreated = v.(time.Time)
	}
	if v, ok := d["fields"]; ok {
		b.Fields = v.(model.Fields)
	}
	return b
}

func marshalFeature(d cqlResult) *model.Feature {
	if d == nil || len(d) == 0 {
		return nil
	}
	defer recoverMarshalPanic("feature", d)

	f := &model.Feature{
		Key:   d["key"].(string),
		Type:  d["type"].(string),
		Value: d["value"].(string),
		Gate: &model.Gate{
			Value:         d["gate_value"].(string),
			Groups:        d["gate_groups"].([]string),
			Actors:        d["gate_actors"].([]string),
			ActorPercent:  d["gate_actor_percent"].(int),
			PercentOfTime: d["gate_percent_of_time"].(int),
		},
		DateCreated: d["date_created"].(time.Time),
		LastUpdated: d["last_updated"].(time.Time),
	}
	if v, ok := d["namespace"]; ok {
		f.Namespace = v.(string)
	}
	return f
}

func recoverMarshalPanic(marshaler string, dat cqlResult) {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("unknwon panic")
		}
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"data": dat,
		}).Errorf("Recovered marshal error in %s", marshaler)
	}
}
