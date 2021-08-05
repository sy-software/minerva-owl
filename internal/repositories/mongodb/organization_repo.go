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

type OrgRepo struct {
	db         *MongoDB
	collection *mongo.Collection
	config     *domain.Config
}

func NewOrgRepo(db *MongoDB, config *domain.Config) (*OrgRepo, error) {
	return &OrgRepo{
		db:         db,
		collection: db.client.Database("minerva").Collection("organizations"),
		config:     config,
	}, nil
}

func (repo *OrgRepo) List(skip int, limit int) ([]domain.Organization, error) {
	var orgs []domain.Organization
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	limit64 := int64(limit)
	skip64 := int64(skip)

	log.Debug().Msg("Listing elements")
	cur, err := repo.collection.Find(ctx, bson.D{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &skip64,
	})

	ctx, cancelFn = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	if err = cur.All(ctx, &orgs); err != nil {
		log.Debug().Err(err).Msg("List error")
		return orgs, err
	}

	log.Debug().Msgf("Found elements: %d", len(orgs))
	return orgs, err
}

func (repo *OrgRepo) Get(id string) (domain.Organization, error) {
	org := domain.Organization{}
	log.Debug().Msgf("Finding element with _id: %q", id)

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Debug().Err(err).Msg("Get error")
		return org, nil
	}

	result := repo.collection.FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	})

	if result.Err() == mongo.ErrNoDocuments {
		return org, ports.ErrItemNotFound{
			Id:    &id,
			Model: "Organization",
		}
	}

	if result.Err() != nil {
		log.Debug().Err(result.Err()).Msg("Get error")
		return org, result.Err()
	}
	err = result.Decode(&org)
	return org, err
}

func (repo *OrgRepo) Create(entity domain.Organization) (string, error) {
	log.Debug().Msgf("Saving: %v", entity)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	result, err := repo.collection.InsertOne(ctx, entity)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result.InsertedID), nil
}

func (repo *OrgRepo) Update(entity domain.Organization) error {
	log.Debug().Msgf("Saving: %v", entity)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(entity.Id)

	if err != nil {
		return err
	}

	result, err := repo.collection.UpdateOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "name", Value: entity.Name},
				primitive.E{Key: "description", Value: entity.Description},
				primitive.E{Key: "logo", Value: entity.Logo},
			},
		},
	})

	log.Debug().Msgf("Update result: %+v", result)

	return err
}

func (repo *OrgRepo) Delete(id string) error {
	log.Debug().Msgf("Deleting by id: %q", id)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := repo.collection.DeleteOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: objectId},
	})

	log.Debug().Msgf("Delete result: %+v", result)

	return err
}
