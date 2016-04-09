package cql

import (
	"github.com/gocql/gocql"
	"github.com/robzienert/cqlmigrate"
	"github.com/robzienert/http-healthcheck"
	monitor "github.com/robzienert/http-healthcheck/monitor/cassandra"
	"github.com/robzienert/lever/cassandra"
	"github.com/robzienert/lever/store"
)

// StoreSpec defines the arguments for creating a new CQL Store.
type StoreSpec struct {
	Keyspace   string
	Hosts      []string
	CertPath   string
	KeyPath    string
	Username   string
	Password   string
	Migrations []cqlmigrate.Spec
}

// ToCassandraSpec converts the Spec expected by the cassandra package.
func (s StoreSpec) ToCassandraSpec() cassandra.Spec {
	return cassandra.Spec{
		Keyspace:    s.Keyspace,
		Seeds:       s.Hosts,
		SSLCertPath: s.CertPath,
		SSLKeyPath:  s.KeyPath,
		Username:    s.Username,
		Password:    s.Password,
	}
}

// StoreResponse encapsulates all resulting objects from a CQL Store load.
//
// The session is returned so that its Closer impl can be deferred from the
// main method.
type StoreResponse struct {
	Store          store.Store
	Session        *gocql.Session
	HealthProvider healthcheck.Provider
}

// Load a new CQL storage backend given the CQL session.
//
// This function will block until the application is able to acquire a
// connection with Cassandra.
func Load(spec StoreSpec) (*StoreResponse, error) {
	session := <-cassandra.New(spec.ToCassandraSpec())
	healthProvider := monitor.NewHealthProvider(session)

	err := cassandra.RunMigrations(spec.Keyspace, session, spec.Migrations)
	if err != nil {
		return nil, err
	}

	return &StoreResponse{
		Store:          New(session),
		Session:        session,
		HealthProvider: healthProvider,
	}, nil
}

// New will initialize a new CQL storage repository.
func New(session *gocql.Session) store.Store {
	return store.New(
		"cql",
		&breadcrumbStore{session: session},
		&featureStore{session: session},
	)
}
