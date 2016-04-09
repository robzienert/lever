package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"golang.org/x/net/context"
)

// Key represents the context value for the metrics module.
const Key = "metrics"

// FromContext returns the statsd client instance from net.Context.
func FromContext(c context.Context) *statsd.Client {
	v := c.Value(Key)
	if v != nil {
		return v.(*statsd.Client)
	}
	return nil
}

// WithTiming is a convenience function to record execution times for arbitrary
// blocks of code.
//
// StatsD may not be available (like running locally or in tests). The function
// will silently fallback to only running the function in this case.
func WithTiming(c context.Context, name string, fn func()) {
	statsd := FromContext(c)
	if statsd != nil {
		start := time.Now()
		fn()
		statsd.TimeInMilliseconds(name, float64(time.Now().Sub(start).Nanoseconds()/1000), nil, 1)
	} else {
		fn()
	}
}
