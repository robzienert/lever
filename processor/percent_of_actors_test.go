package processor

import (
	"fmt"
	"testing"

	"github.com/robzienert/lever/model"
	"github.com/stretchr/testify/assert"
)

var percentOfActorsTests = []struct {
	gate     *model.Gate
	value    string
	expected bool
}{
	{
		&model.Gate{Actors: []string{"one", "two"}},
		"one",
		false,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}, ActorPercent: 0},
		"one",
		false,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}, ActorPercent: 100},
		"three",
		false,
	},
	{
		&model.Gate{Actors: []string{"one", "two"}, ActorPercent: 100},
		"one",
		true,
	},
}

func TestPercentOfActorsProcessor(t *testing.T) {
	for i, tt := range percentOfActorsTests {
		actual := percentOfActorsProcessor(tt.gate, tt.value)
		if tt.expected {
			assert.True(t, actual, fmt.Sprintf("case %d", i+1))
		} else {
			assert.False(t, actual, fmt.Sprintf("case %d", i+1))
		}
	}
}
