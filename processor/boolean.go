package processor

import "github.com/robzienert/lever/model"

func booleanProcessor(g *model.Gate, value string) bool {
	return g.Value == "true"
}
