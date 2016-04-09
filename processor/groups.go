package processor

import (
	"strings"

	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/shared/strutil"
)

func groupsProcessor(g *model.Gate, value string) bool {
	for _, group := range strings.Split(value, ",") {
		if strutil.StringInSlice(group, g.Groups) {
			return true
		}
	}
	return false
}
