package model

import (
	"fmt"
	"time"
)

// Fields allow defining arbitrary data with a breadcrumb
type Fields map[string]string

// Breadcrumb is an individual action caused by an external force on the system.
type Breadcrumb struct {
	Action      string    `json:"action"`
	Actor       string    `json:"actor"`
	DateCreated time.Time `json:"dateCreated"`
	Fields      Fields    `json:"fields"`
}

// NewBreadcrumb is the primary factory for building new Breadcrumbs.
func NewBreadcrumb(action string, actor string) *Breadcrumb {
	return &Breadcrumb{
		Action:      action,
		Actor:       actor,
		DateCreated: time.Now().UTC(),
		Fields:      Fields{},
	}
}

// WithField allows defining a single data field with the breadcrumb.
func (b *Breadcrumb) WithField(key string, value string) *Breadcrumb {
	b.Fields[key] = value
	return b
}

// WithFields allows defining n-Fields with a breadcrumb.
func (b *Breadcrumb) WithFields(fields Fields) *Breadcrumb {
	b.Fields = fields
	return b
}

// ToEvent converts a breadcrumb into a text blob suitable for reporting into
// DataDog's Event stream.
func (b *Breadcrumb) ToEvent() string {
	var fields string
	for k, v := range b.Fields {
		fields += fmt.Sprintf("%s: %s\n", k, v)
	}
	return fmt.Sprintf("Action: %s\nBy: %s\n%s", b.Action, b.Actor, fields)
}
