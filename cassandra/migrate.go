package cassandra

import (
	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
	"github.com/robzienert/cqlmigrate"
)

// RunMigrations is a convenience wrapper around the migration package.
func RunMigrations(keyspace string, session *gocql.Session, migrations []cqlmigrate.Spec) error {
	runner := cqlmigrate.New(&cqlmigrate.Config{
		Keyspace: keyspace,
		Session:  session,
	})
	ok, err := runner.Run(migrations)
	if ok {
		logrus.Info("Successfully ran Cassandra migrations")
		return nil
	}
	return err
}
