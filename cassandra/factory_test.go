package cassandra

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestNewClusterConfig(t *testing.T) {
	t.Parallel()

	spec := Spec{
		Keyspace:    "foo",
		Seeds:       []string{"1.1.1.1", "2.2.2.2"},
		SSLCertPath: "/cert",
		SSLKeyPath:  "/key",
		Username:    "user",
		Password:    "pass",
	}

	cluster := newClusterConfig(spec)
	assert.Equal(t, "foo", cluster.Keyspace)
	assert.Equal(t, gocql.Quorum, cluster.Consistency)
	assert.Equal(t, []string{"1.1.1.1", "2.2.2.2"}, cluster.Hosts)
}
