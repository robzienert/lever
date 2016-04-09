package cassandra

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
)

// Spec defines the specification for creating a new Cassandra connection.
type Spec struct {
	Keyspace    string
	Seeds       []string
	SSLCertPath string
	SSLKeyPath  string
	Username    string
	Password    string
}

// New returns a gocql Session for the given Cassandra cluster and keyspace.
// The function will continue to try to establish a connection until successful.
func New(spec Spec) chan *gocql.Session {
	logrus.WithFields(logrus.Fields{
		"keyspace": spec.Keyspace,
		"seeds":    spec.Seeds,
	}).Info("Establishing a connection with Cassandra cluster")

	ch := make(chan *gocql.Session, 1)
	go func() {
		defer close(ch)
		for {
			cluster := newClusterConfig(spec)
			if trySession(ch, cluster) {
				return
			}
		}
	}()
	return ch
}

func newClusterConfig(spec Spec) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(spec.Seeds...)
	cluster.Keyspace = spec.Keyspace
	cluster.CQLVersion = "3.2.0"
	cluster.ProtoVersion = 3
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = time.Second * 5
	if spec.SSLCertPath != "" && spec.SSLKeyPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CertPath:               spec.SSLCertPath,
			KeyPath:                spec.SSLKeyPath,
			EnableHostVerification: false,
		}
	}
	if spec.Username != "" && spec.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: spec.Username,
			Password: spec.Password,
		}
	}
	return cluster
}

func trySession(sessionCh chan<- *gocql.Session, cluster *gocql.ClusterConfig) bool {
	session, err := cluster.CreateSession()
	if err == nil && session != nil {
		sessionCh <- session
		return true
	}

	logrus.WithField("err", err).Error("Cannot connect to Cassandra: Retrying in 10 seconds")
	time.Sleep(10 * time.Second)
	return false
}
