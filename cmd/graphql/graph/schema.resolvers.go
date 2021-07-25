package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/sy-software/minerva-owl/cmd/graphql/graph/generated"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
)

func (r *mutationResolver) CreateOrganization(ctx context.Context, input model.NewOrganization) (*model.Organization, error) {
	return r.OrgHandler.Create(input.Name, input.Description, input.Logo)
}

func (r *mutationResolver) UpdateOrganization(ctx context.Context, input model.UpdateOrganization) (*model.Organization, error) {
	return r.OrgHandler.Update(input.ID, input.Name, input.Description, input.Logo)
}

func (r *mutationResolver) DeleteOrganization(ctx context.Context, id string) (*model.Organization, error) {
	return r.OrgHandler.Delete(id)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	return r.UsrHandler.Create(input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	return r.UsrHandler.Update(input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	return r.UsrHandler.Delete(id)
}

func (r *queryResolver) Organizations(ctx context.Context, page *int, pageSize *int) ([]*model.Organization, error) {
	return r.OrgHandler.Query(page, pageSize)
}

func (r *queryResolver) Organization(ctx context.Context, id string) (*model.Organization, error) {
	return r.OrgHandler.QueryById(id)
}

func (r *queryResolver) Users(ctx context.Context, role *string, page *int, pageSize *int) ([]*model.User, error) {
	return r.UsrHandler.Query(role, page, pageSize)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.UsrHandler.QueryById(id)
}

func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.UsrHandler.QueryByUsername(username)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
