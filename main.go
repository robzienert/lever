package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/robzienert/gin-middleware/header"
	"github.com/robzienert/gin-middleware/oauth"
	"github.com/robzienert/http-healthcheck"
	"github.com/robzienert/lever/metrics"
	"github.com/robzienert/lever/router"
	"github.com/robzienert/lever/router/middleware/context"
	"github.com/robzienert/lever/shared/config"
	"github.com/robzienert/lever/shared/server"
	"github.com/robzienert/lever/store"
	"github.com/robzienert/lever/store/cql"
	"github.com/robzienert/lever/store/memory"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
)

func init() {
	config.Load()
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()
	config.SetReleaseFlag(viper.GetBool("releaseMode"))

	if err := config.Validate(); err != nil {
		logrus.WithField("err", err).Fatal("Config value is not valid")
	}

	statsd := metrics.Load(viper.GetString("statsd.addr"), viper.GetInt("statsd.bufferLength"))
	{
		defer statsd.Close()
	}

	var healthProviders []healthcheck.Provider

	var backendStore store.Store
	if viper.GetString("store") == "cql" {
		cqlStoreResp, err := cql.Load(cql.StoreSpec{
			Keyspace:   viper.GetString("cassandra.keyspace"),
			Hosts:      viper.GetStringSlice("cassandra.hosts"),
			CertPath:   viper.GetString("cassandra.ssl.certPath"),
			KeyPath:    viper.GetString("cassandra.ssl.keyPath"),
			Username:   viper.GetString("cassandra.auth.username"),
			Password:   viper.GetString("cassandra.auth.password"),
			Migrations: cassandraMigrations,
		})
		if err != nil {
			if cqlStoreResp != nil && cqlStoreResp.Session != nil {
				cqlStoreResp.Session.Close()
			}
			logrus.WithField("err", err).Fatal("Error running Cassandra migrations")
		}
		defer cqlStoreResp.Session.Close()
		healthProviders = append(healthProviders, cqlStoreResp.HealthProvider)
		backendStore = cqlStoreResp.Store
	} else {
		backendStore = memory.Load()
	}

	healthMonitor := healthcheck.New(healthcheck.DefaultSupervisor, healthProviders...)
	{
		defer healthMonitor.Close()
		healthMonitor.Start()
	}

	tokenValidator := oauth.NewSpringSecTokenValidator(
		oauth.SpringSecTokenValidatorSpec{
			Host:     viper.GetString("oauth.host"),
			User:     viper.GetString("oauth.user"),
			Password: viper.GetString("oauth.password"),
		},
	)

	server := server.Load(viper.GetString("http.addr"), viper.GetString("http.cert"), viper.GetString("http.key"))
	server.Run(router.Load(
		tokenValidator,
		header.Version(header.VersionHeader, version),
		context.SetStatsD(statsd),
		context.SetHealthMonitor(healthMonitor),
		context.SetStore(backendStore),
	))
}
