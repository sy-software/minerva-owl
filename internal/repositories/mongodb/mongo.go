package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

var mdbInstance *MongoDB
var mdbOnce sync.Once

// GetMongoDB creates or returns a singleton connection to a MongoDB instance
func GetMongoDB(config domain.MDBConfig) (*MongoDB, error) {
	var dbErr error
	mdbOnce.Do(func() {
		log.Info().Msg("Initializing Mongo DB connection")
		uri := "mongodb://"

		if len(config.Username) > 0 {
			uri += config.Username + ":" + config.Password + "@"
		}

		uri += fmt.Sprintf("%s:%d/", config.Host, config.Port)

		clientOpts := options.Client()
		maxPoolSize := uint64(config.MaxPoolSize)
		clientOpts.MaxPoolSize = &maxPoolSize

		client, err := mongo.NewClient(clientOpts.ApplyURI(uri))

		if err != nil {
			dbErr = err
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), config.ConnectTimeout*time.Second)
		defer cancelFn()
		err = client.Connect(ctx)

		if err != nil {
			dbErr = err
		}

		mdbInstance = &MongoDB{
			client: client,
		}

		log.Info().Msg("DB Connected")
	})

	return mdbInstance, dbErr
}

// Close should always be called when the process using the connection ends
func (mdb *MongoDB) Close() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	mdb.client.Disconnect(ctx)
}
