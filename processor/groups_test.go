package processor

import (
	"testing"

	"github.com/robzienert/lever/model"
	"github.com/stretchr/testify/assert"
)

var groupsTests = []struct {
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
		&model.Gate{Groups: []string{"one", "two"}},
		"three",
		false,
	},
	{
		&model.Gate{Groups: []string{"one", "two"}},
		"one",
		true,
	},
}

func TestGroupsProcessor(t *testing.T) {
	for _, tt := range groupsTests {
		actual := groupsProcessor(tt.gate, tt.value)
		if tt.expected {
			assert.True(t, actual)
		} else {
			assert.False(t, actual)
		}
	}
}
