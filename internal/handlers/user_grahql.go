package handlers

import (
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/utils"
)

type UserGraphqlHandler struct {
	service service.UserService
}

func NewUserGraphqlHandler(service service.UserService) *UserGraphqlHandler {
	return &UserGraphqlHandler{
		service: service,
	}
}

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

func (handler *UserGraphqlHandler) Update(input model.UpdateUser) (*model.User, error) {
	domainUser, err := handler.service.Update(*graphQLUpdateToUser(&input))

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), err
}

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

func (handler *UserGraphqlHandler) QueryById(id string) (*model.User, error) {
	domainUser, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), nil
}

func (handler *UserGraphqlHandler) QueryByUsername(username string) (*model.User, error) {
	domainUser, err := handler.service.GetByUsername(username)

	if err != nil {
		return nil, err
	}

	return userToGraphQL(&domainUser), nil
}

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
