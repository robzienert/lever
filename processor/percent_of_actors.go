package processor

import (
	"strings"

	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/shared/strutil"
	"github.com/spaolacci/murmur3"
)

func percentOfActorsProcessor(g *model.Gate, value string) bool {
	for _, actor := range strings.Split(value, ",") {
		if !strutil.StringInSlice(actor, g.Actors) {
			continue
		}
		if murmur3.Sum32([]byte(actor))%100 <= uint32(g.ActorPercent) {
			return true
		}
	}
	return false
}
