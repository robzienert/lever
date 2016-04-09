package cql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBreadcrumb_NilOnEmpty(t *testing.T) {
	assert.Nil(t, marshalBreadcrumb(cqlResult{}))
}

func TestMarshalFeature_NilOnEmpty(t *testing.T) {
	assert.Nil(t, marshalFeature(cqlResult{}))
}

func TestMarshalPanicRecovery(t *testing.T) {
	assert.NotPanics(t, func() {
		assert.Nil(t, marshalFeature(cqlResult{"foo": "bar"}))
	})
}
