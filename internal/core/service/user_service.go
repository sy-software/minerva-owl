package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/utils"
)

const userCollectionName = "users"

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
	pageVal := utils.CoalesceInt(page, 1) - 1
	pageSizeVal := utils.CoalesceInt(pageSize, srv.config.Pagination.PageSize)

	if pageSizeVal > srv.config.Pagination.MaxPageSize || pageSizeVal <= 0 {
		pageSizeVal = srv.config.Pagination.PageSize
	}

	if pageVal < 0 {
		pageVal = 0
	}

	results := []domain.User{}
	err := srv.repository.List(userCollectionName, &results, pageVal*pageSizeVal, pageSizeVal)

	return results, err
}

func (srv *UserService) ListByRole(role string, page *int, pageSize *int) ([]domain.User, error) {
	results := []domain.User{}
	return results, nil
}

func (srv *UserService) Get(id string) (domain.User, error) {
	result := domain.User{}
	err := srv.repository.Get(userCollectionName, id, &result)
	return result, err
}

func (srv *UserService) GetByUsername(id string) (domain.User, error) {
	result := domain.User{}
	return result, nil
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
