package repositories

import (
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
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
func GetCassandra(config domain.CDBConfig) (*Cassandra, error) {
	var dbErr error
	once.Do(func() {
		log.Info().Msg("Initializing Cassandra DB connection")
		cluster := gocql.NewCluster(config.Host)
		cluster.Port = config.Port
		cluster.Consistency = gocql.Quorum
		cluster.ProtoVersion = 4
		cluster.ConnectTimeout = time.Second * config.ConnectTimeout
		cluster.Timeout = time.Second * config.ConnectTimeout
		cluster.NumConns = config.Connections

		// TODO: Pass a logger with our standard format
		// cluster.Logger = ....

		// TODO: Add authentication
		// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
		session, err := cluster.CreateSession()

		if err != nil {
			dbErr = err
			return
		}

		// create keyspaces
		err = session.Query("CREATE KEYSPACE IF NOT EXISTS minerva WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 3};").Exec()
		if err != nil {
			dbErr = err
			return
		}

		instance = &Cassandra{
			cluster: cluster,
			session: session,
		}
	})

	return instance, dbErr
}

func (cassandra *Cassandra) Close() {
	cassandra.session.Close()
}
