package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

const orgCollectionName = "organizations"

type OrganizationService struct {
	repository ports.Repository
	config     domain.Config
}

func NewOrgService(repo ports.Repository, config domain.Config) *OrganizationService {
	return &OrganizationService{
		repository: repo,
		config:     config,
	}
}

func (srv *OrganizationService) List(page *int, pageSize *int) ([]domain.Organization, error) {
	results := []domain.Organization{}
	_, pageSizeVal, skip := pagination(page, pageSize, srv.config)
	err := srv.repository.List(orgCollectionName, &results, skip, pageSizeVal)

	return results, err
}

func (srv *OrganizationService) Get(id string) (domain.Organization, error) {
	result := domain.Organization{}
	err := srv.repository.Get(orgCollectionName, id, &result)
	return result, err
}

func (srv *OrganizationService) Create(name string, description string, logo string) (domain.Organization, error) {
	entity := domain.Organization{
		Name:        name,
		Description: description,
		Logo:        logo,
	}

	newId, err := srv.repository.Create(orgCollectionName, &entity)
	entity.Id = newId
	return entity, err
}

func (srv *OrganizationService) Update(entity domain.Organization) (domain.Organization, error) {
	return entity, srv.repository.Update(orgCollectionName, entity.Id, &entity)
}

func (srv *OrganizationService) Delete(id string, hard bool) error {
	return srv.repository.Delete(orgCollectionName, id)
}
