package repositories

import (
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx/v2"
	"github.com/sy-software/minerva-owl/internal/core/domain"
)

// Cassandra holds Cassandra DB related objects
type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocqlx.Session
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
		cluster.Consistency = gocql.Any
		cluster.ProtoVersion = 4
		cluster.ConnectTimeout = time.Second * config.ConnectTimeout
		cluster.Timeout = time.Second * config.ConnectTimeout
		cluster.NumConns = config.Connections

		// TODO: Pass a logger with our standard format
		// cluster.Logger = ....

		// TODO: Add authentication
		// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
		session, err := gocqlx.WrapSession(cluster.CreateSession())

		if err != nil {
			dbErr = err
			return
		}
		log.Info().Msg("Cassandra DB Session created")
		log.Info().Msg("Creating minerva Keyspace")
		// create keyspaces
		err = session.ExecStmt("CREATE KEYSPACE IF NOT EXISTS minerva WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 3};")
		if err != nil {
			log.Debug().Err(err).Msg("Error")
			dbErr = err
			return
		}

		log.Info().Msg("Keyspace created")

		instance = &Cassandra{
			cluster: cluster,
			session: &session,
		}
	})

	return instance, dbErr
}

func (cassandra *Cassandra) Close() {
	cassandra.session.Close()
}
