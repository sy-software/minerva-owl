package repositories

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/sy-software/minerva-owl/internal/core/domain"
)

// Cassandra holds Cassandra DB related objects
type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var instance *Cassandra
var once sync.Once

// GetCassandra Gets a singleton connection with Cassandra DB
func GetCassandra(config domain.CDBConfig) *Cassandra {
	once.Do(func() {
		cluster := gocql.NewCluster(config.Host)
		cluster.Port = config.Port
		cluster.Consistency = gocql.Quorum
		cluster.ProtoVersion = 4
		cluster.ConnectTimeout = time.Second * config.ConnectTimeout
		cluster.Timeout = time.Second * config.ConnectTimeout
		cluster.NumConns = config.Connections

		// TODO: Add authentication
		// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
		session, err := cluster.CreateSession()

		if err != nil {
			// TODO: Handle errors correctly
			fmt.Printf("Error connection: %v", err)
		}

		// create keyspaces
		err = session.Query("CREATE KEYSPACE IF NOT EXISTS minerva WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 3};").Exec()
		if err != nil {
			fmt.Printf("Error creating keyspace: %v", err)
		}

		instance = &Cassandra{
			cluster: cluster,
			session: session,
		}
	})

	return instance
}

func (cassandra *Cassandra) Close() {
	cassandra.session.Close()
}
