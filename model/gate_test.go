package model

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var gateTypesTests = []struct {
	gate     Gate
	expected []string
}{
	{
		Gate{Value: "true", Groups: []string{"one"}, Actors: []string{"one"}, ActorPercent: 1, PercentOfTime: 1},
		[]string{BooleanGateType, GroupsGateType, PercentOfActorsGateType, PercentOfTimeGateType},
	},
	{
		Gate{Value: "invalid", Groups: []string{"one"}, Actors: []string{"one"}},
		[]string{GroupsGateType, ActorsGateType},
	},
}

func TestGate_Types(t *testing.T) {
	for _, tt := range gateTypesTests {
		assert.Equal(t, tt.expected, tt.gate.Types())
	}
}

var byPrecedenceTests = []struct {
	input    []string
	expected []string
}{
	{[]string{BooleanGateType, GroupsGateType}, []string{BooleanGateType, GroupsGateType}},
	{[]string{GroupsGateType, BooleanGateType}, []string{BooleanGateType, GroupsGateType}},
	{[]string{PercentOfTimeGateType, PercentOfActorsGateType}, []string{PercentOfActorsGateType, PercentOfTimeGateType}},
}

func TestByPrecedenceSorter(t *testing.T) {
	for _, tt := range byPrecedenceTests {
		sort.Sort(byPrecedence(tt.input))
		assert.Equal(t, tt.expected, tt.input)
	}
}
