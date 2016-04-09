package processor

import (
	"math/rand"

	"github.com/robzienert/lever/model"
)

func percentOfTimeProcessor(g *model.Gate, value string) bool {
	return rand.Intn(100) >= g.PercentOfTime
}
