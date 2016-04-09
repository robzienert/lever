package processor

import (
	"testing"

	"github.com/robzienert/lever/model"
	"github.com/stretchr/testify/assert"
)

var actorsTests = []struct {
	gate     *model.Gate
	value    string
	expected bool
}{
	{
		&model.Gate{Value: "true"},
		"one",
		false,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}},
		"three",
		false,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}},
		"one",
		true,
	},
}

func TestActorsProcessor(t *testing.T) {
	for _, tt := range actorsTests {
		actual := actorsProcessor(tt.gate, tt.value)
		if tt.expected {
			assert.True(t, actual)
		} else {
			assert.False(t, actual)
		}
	}
}
