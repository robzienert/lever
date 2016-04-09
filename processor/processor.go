package processor

import (
	"fmt"

	"github.com/robzienert/lever/model"
)

type gateProcessor func(g *model.Gate, value string) bool

var gateProcessorMap = map[string]gateProcessor{
	model.BooleanGateType:         booleanProcessor,
	model.PercentOfActorsGateType: percentOfActorsProcessor,
	model.ActorsGateType:          actorsProcessor,
	model.GroupsGateType:          groupsProcessor,
	model.PercentOfTimeGateType:   percentOfTimeProcessor,
}

// ProcessGate will return the gate state of a feature given actors and groups.
func ProcessGate(g *model.Gate, actors string, groups string) (bool, error) {
	for _, gt := range g.Types() {
		f := gateProcessorMap[gt]
		if f == nil {
			return false, fmt.Errorf("could not load gate func for type: %s", gt)
		}

		var enabled bool
		switch gt {
		case model.ActorsGateType:
			enabled = f(g, actors)
		case model.PercentOfActorsGateType:
			enabled = f(g, actors)
		case model.GroupsGateType:
			enabled = f(g, groups)
		default:
			enabled = f(g, "")
		}
		if enabled {
			return true, nil
		}
	}
	return false, nil
}
