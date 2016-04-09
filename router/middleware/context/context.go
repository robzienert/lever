package context

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/robzienert/http-healthcheck"
	"github.com/robzienert/lever/metrics"
	"github.com/robzienert/lever/store"
	"github.com/satori/go.uuid"
)

// SetHealthMonitor will set the given health monitor into the net.Context.
func SetHealthMonitor(m *healthcheck.Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(healthcheck.Key, m)
		c.Next()
	}
}

// SetStatsD will set the given statsd client into the net.Context.
func SetStatsD(s *statsd.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(metrics.Key, s)
		c.Next()
	}
}

// SetStore will set the storage backend into the net.Context.
func SetStore(s store.Store) gin.HandlerFunc {
	logrus.Infof("Using storage backend: %s", s.Name())
	return func(c *gin.Context) {
		store.ToContext(c, s)
		c.Next()
	}
}

// SetRequestUUID will search for an X-ST-CORRELATION header and set a
// request-level correlation ID into the net.Context. If no header is found, a
// new UUID will be generated.
func SetRequestUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.Request.Header.Get("X-ST-CORRELATION")
		if u == "" {
			u = uuid.NewV4().String()
		}
		contextLogger := logrus.WithField("uuid", u)
		c.Set("log", contextLogger)
		c.Set("uuid", u)
	}
}
