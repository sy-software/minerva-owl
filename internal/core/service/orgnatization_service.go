package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"

	"github.com/google/uuid"
)

type OrganizationService struct {
	repository ports.OrganizationRepo
}

func (srv *OrganizationService) All() ([]domain.Organization, error) {
	return []domain.Organization{}, nil
}

func (srv *OrganizationService) Get(id string) (domain.Organization, error) {
	return domain.Organization{}, nil
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
	return domain.Organization{}, nil
}

func (srv *OrganizationService) Delete(id string, hard bool) error {
	return nil
}
