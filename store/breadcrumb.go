package store

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/robzienert/lever/metrics"
	"github.com/robzienert/lever/model"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

// BreadcrumbStore is the repsitory for interacting with audit backends.
type BreadcrumbStore interface {
	GetList() ([]*model.Breadcrumb, error)
	Create(*model.Breadcrumb) error
}

// GetBreadcrumbList will proxy to the net.Context's breadcrumb storage backend
// to return a full list of audit breadcrumbs.
func GetBreadcrumbList(c context.Context) ([]*model.Breadcrumb, error) {
	return FromContext(c).Breadcrumbs().GetList()
}

// SaveBreadcrumb will proxy to the net.Context's breadcrumb storage backend
// to create a new audit breadcrumb.
func SaveBreadcrumb(c context.Context, b *model.Breadcrumb) error {
	// Important bit: Since SaveBreadcrumb is called via goroutines, we need to
	// make sure that we're using a non-global instance of the logger or data
	// races will happen.
	logrus.New().WithField("b", *b).Info("Creating breadcrumb")

	if viper.GetBool("releaseMode") && metrics.FromContext(c) != nil {
		// TODO It might be valuable to parse the feature key, if we wind up
		// standardizing on key format that would allow us to easily correlate to
		// shards, etc.
		metrics.FromContext(c).Event(&statsd.Event{
			Title:     "Dynamic Feature Change",
			Text:      b.ToEvent(),
			Priority:  statsd.Low,
			AlertType: statsd.Info,
		})
	}

	return FromContext(c).Breadcrumbs().Create(b)
}
