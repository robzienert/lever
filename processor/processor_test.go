package processor

import (
	"testing"

	"github.com/robzienert/lever/model"
	"github.com/stretchr/testify/assert"
)

func TestProcessGate(t *testing.T) {
	gate := &model.Gate{Value: "true"}
	enabled, err := ProcessGate(gate, "", "")
	assert.NoError(t, err)
	assert.True(t, enabled)
}
