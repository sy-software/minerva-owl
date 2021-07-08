package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/utils"
)

type OrganizationService struct {
	repository ports.OrganizationRepo
	config     domain.Config
}

func NewOrgService(repo ports.OrganizationRepo, config domain.Config) *OrganizationService {
	return &OrganizationService{
		repository: repo,
		config:     config,
	}
}

func (srv *OrganizationService) List(page *int, pageSize *int) ([]domain.Organization, error) {
	pageVal := utils.CoalesceInt(page, 1) - 1
	pageSizeVal := utils.CoalesceInt(pageSize, srv.config.Pagination.PageSize)

	if pageSizeVal > srv.config.Pagination.MaxPageSize || pageSizeVal <= 0 {
		pageSizeVal = srv.config.Pagination.PageSize
	}

	if pageVal < 0 {
		pageVal = 0
	}

	return srv.repository.List(pageVal*pageSizeVal, pageSizeVal)
}

func (srv *OrganizationService) Get(id string) (domain.Organization, error) {
	return srv.repository.Get(id)
}

func (srv *OrganizationService) Create(name string, description string, logo string) (domain.Organization, error) {
	entity := domain.Organization{
		Name:        name,
		Description: description,
		Logo:        logo,
	}

	newId, err := srv.repository.Create(entity)
	entity.Id = newId
	return entity, err
}

func (srv *OrganizationService) Update(entity domain.Organization) (domain.Organization, error) {
	return entity, srv.repository.Update(entity)
}

func (srv *OrganizationService) Delete(id string, hard bool) error {
	return srv.repository.Delete(id)
}
