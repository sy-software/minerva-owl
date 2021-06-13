package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/sy-software/minerva-owl/cmd/graphql/graph/generated"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, input model.NewOrganization) (*model.Organization, error) {
	return r.Handler.Create(input.Name, input.Description, input.Logo)
}

func (r *mutationResolver) UpdateOrganization(ctx context.Context, input model.UpdateOrganization) (*model.Organization, error) {
	return r.Handler.Update(input.ID, input.Name, input.Description, input.Logo)
}

func (r *mutationResolver) DeleteOrganization(ctx context.Context, id string) (*model.Organization, error) {
	return r.Handler.Delete(id)
}

func (r *queryResolver) Organizations(ctx context.Context) ([]*model.Organization, error) {
	return r.Handler.Query()
}

func (r *queryResolver) Organization(ctx context.Context, id string) (*model.Organization, error) {
	return r.Handler.QueryById(id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
