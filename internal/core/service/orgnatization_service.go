package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"

	"github.com/google/uuid"
)

type OrganizationService struct {
	repository ports.OrganizationRepo
}

func NewOrgService(repo ports.OrganizationRepo) *OrganizationService {
	return &OrganizationService{
		repository: repo,
	}
}

func (srv *OrganizationService) All() ([]domain.Organization, error) {
	return srv.repository.All()
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

	entity.Id = uuid.New().String()
	return entity, srv.repository.Save(entity)
}

func (srv *OrganizationService) Update(entity domain.Organization) (domain.Organization, error) {
	return entity, srv.repository.Save(entity)
}

func (srv *OrganizationService) Delete(id string, hard bool) error {
	return srv.repository.Delete(id)
}
