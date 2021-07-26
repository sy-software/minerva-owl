package handlers

import (
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/utils"
)

// UserGraphqlHandler works as adapter between GraphQL endpoints and a UserService
type UserGraphqlHandler struct {
	service service.UserService
}

// NewUserGraphqlHandler creates an instance of UserGraphqlHandler
func NewUserGraphqlHandler(service service.UserService) *UserGraphqlHandler {
	return &UserGraphqlHandler{
		service: service,
	}
}

// Create saves a new user into a repository
func (handler *UserGraphqlHandler) Create(input model.NewUser) (*model.User, error) {
	domainUser, err := handler.service.Create(
		input.Name,
		input.Username,
		utils.CoalesceStr(input.Picture, ""),
		input.Role,
		input.Provider,
		input.TokenID,
		input.Status,
	)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), err
}

// Update saves changes into an exiting User
func (handler *UserGraphqlHandler) Update(input model.UpdateUser) (*model.User, error) {
	domainUser, err := handler.service.Update(*graphQLUpdateToUser(&input))

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), err
}

// Delete removes a User with the provided id
func (handler *UserGraphqlHandler) Delete(id string) (*model.User, error) {
	out, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	err = handler.service.Delete(id, false)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&out), nil
}

// Query returns a paginated list of Users that can be filtered by role
func (handler *UserGraphqlHandler) Query(role *string, page *int, pageSize *int) ([]*model.User, error) {
	output := []*model.User{}
	var users []domain.User
	var err error
	if role != nil {
		users, err = handler.service.ListByRole(*role, page, pageSize)
		if err != nil {
			return output, err
		}
	} else {
		users, err = handler.service.List(page, pageSize)
		if err != nil {
			return output, err
		}
	}

	for _, u := range users {
		output = append(output, userToGraphQL(&u))
	}

	return output, nil
}

// QueryById returns the User with the provided id
func (handler *UserGraphqlHandler) QueryById(id string) (*model.User, error) {
	domainUser, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), nil
}

// QueryById returns the User with the provided username
func (handler *UserGraphqlHandler) QueryByUsername(username string) (*model.User, error) {
	domainUser, err := handler.service.GetByUsername(username)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), nil
}

// userToGraphQL converts the internal User model into the GraphQL version
func userToGraphQL(source *domain.User) *model.User {
	return &model.User{
		ID:         source.Id,
		Name:       source.Name,
		Username:   source.Username,
		Picture:    &source.Picture,
		Role:       source.Role,
		Provider:   source.Provider,
		TokenID:    source.TokenID,
		CreateDate: source.CreateDate,
		UpdateDate: source.UpdateDate,
		Status:     source.Status,
	}
}

// graphQLUpdateToUser converts the GraphQL user model into our internal model
func graphQLUpdateToUser(source *model.UpdateUser) *domain.User {
	return &domain.User{
		Id:       source.ID,
		Name:     source.Name,
		Username: source.Username,
		Picture:  utils.CoalesceStr(source.Picture, ""),
		Role:     source.Role,
		Provider: source.Provider,
		TokenID:  source.TokenID,
		Status:   source.Status,
	}
}
