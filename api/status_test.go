package api

import (
	"errors"
	"testing"

	"github.com/robzienert/http-healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestMarshalHealthStatusResponse(t *testing.T) {
	status := healthcheck.Status{
		Healthy: false,
		Statuses: healthcheck.ProviderStatuses{
			"foo": nil,
			"bar": errors.New("Not OK"),
		},
	}

	actual := MarshalHealthStatusResponse(status)
	expected := HealthStatusResponse{
		Status: map[string]string{
			"foo": "OK",
			"bar": "Not OK",
		},
	}

	assert.Equal(t, expected, actual)
}
