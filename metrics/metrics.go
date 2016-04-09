package metrics

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// The standard metrics namespace
const namespace = "lever."

// Load will startup a new statsd client and hard-fail if it is not available in
// a production deployment.
//
// TODO In prod, should treat like a resource and retry.
func Load(addr string, bufferLength int) *statsd.Client {
	statsd, err := New(addr, bufferLength)
	if err != nil {
		if gin.IsDebugging() {
			logrus.WithField("err", err).Error("Could not create statsd client")
		} else {
			logrus.WithField("err", err).Fatal("Could not create statsd client")
		}
	}
	return statsd
}

// New creates a StatsD buffered client.
func New(addr string, bufferLength int) (*statsd.Client, error) {
	logrus.WithFields(logrus.Fields{
		"addr":   addr,
		"bufLen": bufferLength,
	}).Info("Starting StatsD client")

	c, err := statsd.NewBuffered(addr, bufferLength)
	if err != nil {
		return nil, err
	}
	c.Namespace = namespace

	return c, nil
}
