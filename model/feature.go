package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Feature NODOC
type Feature struct {
	Namespace   string    `json:"namespace,omitempty"`
	Key         string    `json:"key" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Value       string    `json:"value" binding:"required"`
	Gate        *Gate     `json:"gate" binding:"required"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// Diff returns a flatmap diff of two features, which can be used for auditing.
// It is expected that Namespace and Key do not change.
func (f *Feature) Diff(b *Feature) Fields {
	d := make(map[string]string, 0)
	if f.Type != b.Type {
		d["type"] = f.diffValue(f.Type, b.Type)
	}
	if f.Value != b.Value {
		d["value"] = f.diffValue(f.Value, b.Value)
	}
	if f.Gate.ActorPercent != b.Gate.ActorPercent {
		d["gate_actor_percent"] = f.diffValue(strconv.Itoa(f.Gate.ActorPercent), strconv.Itoa(b.Gate.ActorPercent))
	}
	fActors := strings.Join(f.Gate.Actors, ",")
	bActors := strings.Join(b.Gate.Actors, ",")
	if fActors != bActors {
		d["gate_actors"] = f.diffValue(fActors, bActors)
	}
	fGroups := strings.Join(f.Gate.Groups, ",")
	bGroups := strings.Join(b.Gate.Groups, ",")
	if fGroups != bGroups {
		d["gate_groups"] = f.diffValue(fGroups, bGroups)
	}
	if f.Gate.PercentOfTime != b.Gate.PercentOfTime {
		d["gate_percent_of_time"] = f.diffValue(strconv.Itoa(f.Gate.PercentOfTime), strconv.Itoa(b.Gate.PercentOfTime))
	}
	if f.Gate.Value != b.Gate.Value {
		d["gate_value"] = f.diffValue(f.Gate.Value, b.Gate.Value)
	}
	return d
}

func (f *Feature) diffValue(from string, to string) string {
	if from == "" {
		from = "NO_VALUE"
	}
	if to == "" {
		to = "NO_VALUE"
	}
	return fmt.Sprintf("%s -> %s", from, to)
}
