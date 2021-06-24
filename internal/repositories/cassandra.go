package repositories

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"
)

type Cassandra struct {
	cluster gocql.ClusterConfig
	session gocql.Session
}

var instance *Cassandra
var once sync.Once

func GetCassandra() *Cassandra {
	once.Do(func() {
		cluster := gocql.NewCluster("127.0.0.1") //replace PublicIP with the IP addresses used by your cluster.
		cluster.Consistency = gocql.Quorum
		cluster.ProtoVersion = 4
		cluster.ConnectTimeout = time.Second * 10
		cluster.Timeout = time.Second * 10

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
			cluster: *cluster,
			session: *session,
		}
	})

	return instance
}
