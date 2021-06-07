package ports

import "github.com/sy-software/minerva-owl/internal/core/domain"

type OrganizationService interface {
	All() ([]domain.Organization, error)
	Get(id string) (domain.Organization, error)

	Create(name string, Description string, logo string) (domain.Organization, error)
	Update(entity domain.Organization) (domain.Organization, error)

	Delete(id string, hard bool) error
}
