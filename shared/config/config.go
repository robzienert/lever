package config

import (
	"errors"
	"strings"

	"github.com/robzienert/lever/shared/strutil"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Load the application config via Viper.
//
// Viper is a pretty powerful config library. We're doing a few things here.
//
// 1. File config support. While I'd like to deploy this Service out with
// Docker, I suspect we'll wind up doing the Ansible route first. Viper will
// auto-detect YAML, TOML or JSON files, so we just need to set "config" as the
// name, without the file extension. Kind of overkill, but easy.
//
// 2. Environment variable auto-binding. All of our configs will then be
// available like "LEVER_CASSANDRA_KEYSPACE".
//
// 3. Reading the config. Anywhere we need to get values we can just call into
// viper globally: `viper.GetString("cassandra_keyspace")`
func Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/lever/")

	viper.SetEnvPrefix("lever")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("releaseMode", true)
	viper.SetDefault("store", "cql")
	viper.SetDefault("http.addr", ":8500")
	viper.SetDefault("http.cert", "")
	viper.SetDefault("http.key", "")
	viper.SetDefault("cassandra.keyspace", "lever")
	viper.SetDefault("statsd.addr", "127.0.0.1:8125")
	viper.SetDefault("statsd.bufferLength", 100)

	viper.ReadInConfig()
}

// Validate the set configuration opts.
func Validate() error {
	validStores := []string{"cql", "memory"}
	if !strutil.StringInSlice(viper.GetString("store"), validStores) {
		return errors.New("invalid store config")
	}
	return nil
}

// SetReleaseFlag will adjust the logging levels for the application.
func SetReleaseFlag(isRelease bool) {
	if isRelease {
		logrus.Info("Running in \"release\" mode. Switch to \"debug\" mode in dev")
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	}
}
