package model

import (
	"sort"

	"github.com/Sirupsen/logrus"
)

// BooleanGateType is a binary gate.
//
// ActorsGateType will enable if the given actor(s) are in a known, approved
// list of actors.
//
// GroupsGateType will enable if the given group(s) are in the known, approved
// list of actors.
//
// PercentOfActorsGateType will first validate that the actor is allowed, then
// hash the actor to see if they are within a percentage of enabled users in
// that list.
//
// PercentOfTimeGateType will enable the gate a percentage of the time.
const (
	BooleanGateType         = "boolean"
	ActorsGateType          = "actors"
	GroupsGateType          = "groups"
	PercentOfActorsGateType = "percentOfActors"
	PercentOfTimeGateType   = "percentOfTime"
)

// Gate NODOC
type Gate struct {
	Value         string   `json:"value,omitempty"`
	Groups        []string `json:"groups,omitempty"`
	Actors        []string `json:"actors,omitempty"`
	ActorPercent  int      `json:"actorPercent,omitempty"`
	PercentOfTime int      `json:"percentOfTime,omitempty"`
}

// Types returns a slice of all types of the gate, in order of evaluation
// precedence.
func (g *Gate) Types() (types []string) {
	if g.Value == "true" || g.Value == "false" {
		types = append(types, BooleanGateType)
	}
	if len(g.Actors) > 0 {
		if g.ActorPercent > 0 {
			types = append(types, PercentOfActorsGateType)
		} else {
			types = append(types, ActorsGateType)
		}
	}
	if len(g.Groups) > 0 {
		types = append(types, GroupsGateType)
	}
	if g.PercentOfTime > 0 {
		types = append(types, PercentOfTimeGateType)
	}
	sort.Sort(byPrecedence(types))
	return
}

// From highest precedence to lowest.
var gatePrecedence = []string{
	BooleanGateType,
	GroupsGateType,
	ActorsGateType,
	PercentOfActorsGateType,
	PercentOfTimeGateType,
}

type byPrecedence []string

func (s byPrecedence) Len() int {
	return len(s)
}

func (s byPrecedence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPrecedence) Less(i, j int) bool {
	return s.pos(s[i]) < s.pos(s[j])
}

func (s byPrecedence) pos(val string) int {
	for i, g := range gatePrecedence {
		if g == val {
			return i
		}
	}
	logrus.WithField("val", val).Error("Unknown value for byPrecedence sorter")
	return 0
}
