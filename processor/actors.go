package processor

import (
	"strings"

	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/shared/strutil"
)

func actorsProcessor(g *model.Gate, value string) bool {
	for _, actor := range strings.Split(value, ",") {
		if strutil.StringInSlice(actor, g.Actors) {
			return true
		}
	}
	return false
}
