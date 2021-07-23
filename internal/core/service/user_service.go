package service

import (
	"errors"
	"fmt"

	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/utils"
)

const userCollectionName = domain.USER_COL_NAME

type UserService struct {
	repository ports.Repository
	config     domain.Config
}

func NewUserService(repo ports.Repository, config domain.Config) *UserService {
	return &UserService{
		repository: repo,
		config:     config,
	}
}

func (srv *UserService) List(page *int, pageSize *int) ([]domain.User, error) {
	_, pageSizeVal, skip := pagination(page, pageSize, srv.config)

	results := []domain.User{}
	err := srv.repository.List(userCollectionName, &results, skip, pageSizeVal)

	return results, err
}

func (srv *UserService) ListByRole(role string, page *int, pageSize *int) ([]domain.User, error) {
	_, pageSizeVal, skip := pagination(page, pageSize, srv.config)

	results := []domain.User{}
	err := srv.repository.List(userCollectionName, &results, skip, pageSizeVal, ports.Filter{
		Name:  "role",
		Value: role,
	})
	return results, err
}

func (srv *UserService) Get(id string) (domain.User, error) {
	result := domain.User{}
	err := srv.repository.Get(userCollectionName, id, &result)
	return result, err
}

func (srv *UserService) GetByUsername(username string) (domain.User, error) {
	result := domain.User{}
	err := srv.repository.GetOne(userCollectionName, &result, ports.Filter{
		Name:  "username",
		Value: username,
	})
	return result, err
}

func (srv *UserService) Create(
	name string,
	username string,
	picture string,
	role string,
	provider string,
	tokenID string,
	status string,
) (domain.User, error) {

	current, err := srv.GetByUsername(username)

	if err == nil && current.Username == username {
		return domain.User{}, errors.New(fmt.Sprintf("duplicated Username: %s", username))
	}

	if _, ok := err.(ports.ErrItemNotFound); err != nil && !ok {
		return domain.User{}, err
	}

	encryptedToken, err := utils.AES256Encrypt(srv.config.Keys.Auth, tokenID)

	if err != nil {
		return domain.User{}, err
	}
	now := utils.UnixNow()
	entity := domain.User{
		Name:       name,
		Username:   username,
		Picture:    picture,
		Role:       role,
		Provider:   provider,
		TokenID:    encryptedToken,
		Status:     status,
		CreateDate: now,
		UpdateDate: now,
	}

	newId, err := srv.repository.Create(userCollectionName, &entity)
	entity.Id = newId
	return entity, err
}

func (srv *UserService) Update(entity domain.User) (domain.User, error) {
	entity.UpdateDate = utils.UnixUTCNow()

	current, err := srv.Get(entity.Id)

	if err != nil {
		return entity, err
	}

	// TODO: Perform this operation without a Get
	if current.TokenID != entity.TokenID {
		encryptedToken, err := utils.AES256Encrypt(srv.config.Keys.Auth, entity.TokenID)

		if err != nil {
			return entity, err
		}

		entity.TokenID = encryptedToken
	}

	return entity, srv.repository.Update(userCollectionName, entity.Id, &entity, "createDate")
}

func (srv *UserService) Delete(id string, hard bool) error {
	return srv.repository.Delete(userCollectionName, id)
}
