package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeature_Diff(t *testing.T) {
	f1 := Feature{
		Type:  "f1",
		Value: "f1",
		Gate: &Gate{
			Value:         "true",
			Groups:        []string{},
			Actors:        []string{"one"},
			ActorPercent:  10,
			PercentOfTime: 10,
		},
	}
	f2 := Feature{
		Type:  "f2",
		Value: "f2",
		Gate: &Gate{
			Value:         "false",
			Groups:        []string{"two"},
			Actors:        []string{"one", "three"},
			ActorPercent:  50,
			PercentOfTime: 60,
		},
	}

	fields := f1.Diff(&f2)
	assert.Equal(t, "f1 -> f2", fields["type"], "type did not match")
	assert.Equal(t, "f1 -> f2", fields["value"], "value did not match")
	assert.Equal(t, "true -> false", fields["gate_value"], "gate_value did not match")
	assert.Equal(t, "NO_VALUE -> two", fields["gate_groups"], "gate_groups did not match")
	assert.Equal(t, "one -> one,three", fields["gate_actors"], "gate_actors did not match")
	assert.Equal(t, "10 -> 50", fields["gate_actor_percent"], "gate_actor_percent did not match")
	assert.Equal(t, "10 -> 60", fields["gate_percent_of_time"], "gate_percent_of_time did not match")
}
