package processor

import (
	"testing"

	"github.com/robzienert/lever/model"
	"github.com/stretchr/testify/assert"
)

var booleanTests = []struct {
	gate     *model.Gate
	value    string
	expected bool
}{
	{
		&model.Gate{Value: "true"},
		"one",
		true,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}},
		"three",
		false,
	},
	{
		&model.Gate{Value: "blah blah"},
		"",
		false,
	},
}

func TestBooleanProcessor(t *testing.T) {
	for _, tt := range booleanTests {
		actual := booleanProcessor(tt.gate, tt.value)
		if tt.expected {
			assert.True(t, actual)
		} else {
			assert.False(t, actual)
		}
	}
}
