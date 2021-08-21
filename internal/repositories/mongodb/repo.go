package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepo is an implementation of ports.Repo interface with MongoDB as datasource
type MongoRepo struct {
	db          *MongoDB
	collections map[string]*mongo.Collection
	config      *domain.Config
}

// NewMongoRepo creates an instance of MongoRepo
func NewMongoRepo(db *MongoDB, config *domain.Config) (*MongoRepo, error) {
	return &MongoRepo{
		db:          db,
		collections: map[string]*mongo.Collection{},
		config:      config,
	}, nil
}

// mongoGetCollection checks if we have a reference of a given collection, if no creates a new one and returns it
func (repo *MongoRepo) mongoGetCollection(collection string) *mongo.Collection {
	value, exists := repo.collections[collection]

	if exists {
		log.Debug().Msgf("Reusing connection for: %v", collection)
		return value
	} else {
		log.Debug().Msgf("Getting new collection connection for: %v", collection)
		value = repo.db.client.
			Database(repo.config.MongoDBConfig.DB).
			Collection(collection)

		repo.collections[collection] = value
		return value
	}
}

// List stores into results a list of items from the given collection applying the filters
// results must be a pointer to an Slice of an struct with bson tags for serialization
func (repo *MongoRepo) List(collection string, results interface{}, skip int, limit int, filters ...ports.Filter) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	limit64 := int64(limit)
	skip64 := int64(skip)

	dbFilters, err := formatFilters(filters)

	if err != nil {
		log.Debug().Err(err).Msgf("%v - List error", collection)
		return err
	}

	log.Debug().Msgf("%v - Listing elements", collection)
	cur, err := repo.mongoGetCollection(collection).Find(ctx, dbFilters, &options.FindOptions{
		Limit: &limit64,
		Skip:  &skip64,
	})

	if err != nil {
		log.Debug().Err(err).Msgf("%v - List error", collection)
		return err
	}

	ctx, cancelFn = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	if err = cur.All(ctx, results); err != nil {
		log.Debug().Err(err).Msgf("%v - List error", collection)
		return err
	}

	return err
}

// Get stores into result an item from collection with _id equals to id
// result must be a pointer to an instance of a struct with bson tags for serialization
func (repo *MongoRepo) Get(collection string, id string, result interface{}) error {
	log.Debug().Msgf("%v - Finding element with _id: %q", collection, id)

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Debug().Err(err).Msgf("%v - Get error", collection)
		return nil
	}

	rawResult := repo.mongoGetCollection(collection).FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	})

	if rawResult.Err() == mongo.ErrNoDocuments {
		return ports.ErrItemNotFound{
			Id:    &id,
			Model: collection,
		}
	}

	if rawResult.Err() != nil {
		log.Debug().Err(rawResult.Err()).Msgf("%v - Get error", collection)
		return rawResult.Err()
	}

	return rawResult.Decode(result)
}

// Get stores into result an item from collection matching the filters
// result must be a pointer to an instance of a struct with bson tags for serialization
func (repo *MongoRepo) GetOne(collection string, result interface{}, filters ...ports.Filter) error {
	log.Debug().Msgf("%v - Finding element with filters: %+v", collection, filters)

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	dbFilters, err := formatFilters(filters)

	if err != nil {
		log.Debug().Err(err).Msgf("%v - Get error", collection)
		return nil
	}

	rawResult := repo.mongoGetCollection(collection).FindOne(ctx, dbFilters)

	if rawResult.Err() == mongo.ErrNoDocuments {
		log.Debug().Err(rawResult.Err()).Msgf("%v - Get error", collection)
		return ports.ErrItemNotFound{
			Model: collection,
		}
	}

	if rawResult.Err() != nil {
		log.Debug().Err(rawResult.Err()).Msgf("%v - Get error", collection)
		return rawResult.Err()
	}

	return rawResult.Decode(result)
}

// Create saves the serialized version of entity into the collection
// entity must be an instance of a struct with bson tags for serialization
func (repo *MongoRepo) Create(collection string, entity interface{}) (string, error) {
	log.Debug().Msgf("%v - Saving: %v", collection, entity)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	result, err := repo.mongoGetCollection(collection).InsertOne(ctx, entity)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update saves the values of entity to the item with id from the collection
// entity must be an instance of a struct with bson tags for serialization
//
// If you whish to omit some fields from entity from saving you can pass the field
// names into the final omit parameter
func (repo *MongoRepo) Update(collection string, id string, entity interface{}, omit ...string) error {
	log.Debug().Msgf("%v - Saving: %v", collection, entity)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	bsonDoc, err := toBSONDoc(entity, omit...)
	if err != nil {
		return err
	}

	result, err := repo.mongoGetCollection(collection).UpdateOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	}, bson.D{
		primitive.E{
			Key:   "$set",
			Value: bsonDoc,
		},
	})

	log.Debug().Msgf("Update result: %+v", result)

	return err
}

// Delete removes item with id from collection
func (repo *MongoRepo) Delete(collection string, id string) error {
	log.Debug().Msgf("%v - Deleting by id: %q", collection, id)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := repo.mongoGetCollection(collection).DeleteOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	})

	log.Debug().Msgf("%v - Delete result: %+v", collection, result)

	return err
}

// toBSONDoc marshals the value of v into a bson.D and omits the fields matching a name from omit
func toBSONDoc(v interface{}, omit ...string) (bson.D, error) {
	// TODO: Support nested documents and arrays
	data, err := bson.Marshal(v)
	if err != nil {
		return bson.D{}, nil
	}

	var doc bson.D
	err = bson.Unmarshal(data, &doc)

	if len(omit) == 0 || err != nil {
		return doc, err
	}

	omitMap := map[string]bool{}
	for _, v := range omit {
		omitMap[v] = true
	}

	filtered := bson.D{}

	for _, field := range doc {
		_, shouldOmit := omitMap[field.Key]
		if field.Key != "_id" && !shouldOmit {
			filtered = append(filtered, field)
		}
	}

	return filtered, nil
}

// formatFilters takes a generic list of filters and converts them into a MongoDB query
func formatFilters(filters []ports.Filter) (bson.D, error) {
	result := bson.D{}
	// TODO: Cast ObjectId values to the right type
	for _, filter := range filters {
		switch filter.Name {
		case "$or", "$and":

			if values, ok := filter.Value.([]ports.Filter); ok {
				formated := []bson.E{}

				for _, v := range values {
					formated = append(formated, formatFilter(v))
				}

				result = append(result, bson.E{
					Key:   filter.Name,
					Value: formated,
				})
			} else {
				return result, fmt.Errorf("invalid value for operator %q: %+v", filter.Name, filter.Value)
			}
		default:
			result = append(result, formatFilter(filter))
		}
	}

	return result, nil
}

// formatFilters takes a generic filter and converts it into a MongoDB query
func formatFilter(f ports.Filter) bson.E {
	if value, ok := f.Value.(ports.Filter); ok {
		return bson.E{
			Key:   f.Name,
			Value: formatFilter(value),
		}
	} else {
		return bson.E{
			Key:   f.Name,
			Value: f.Value,
		}
	}
}
