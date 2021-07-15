package ports

import "github.com/sy-software/minerva-owl/internal/core/domain"

// OrganizationService is a common interface for a service provider for organization entity
type OrganizationService interface {
	// List returns a single page of items
	List(page *int, pageSize *int) ([]domain.Organization, error)
	// Get returns a single item filter by id
	Get(id string) (domain.Organization, error)
	// Create saves a new organization item into the repository
	Create(name string, Description string, logo string) (domain.Organization, error)
	// Update looks for an existing item and update the values
	Update(entity domain.Organization) (domain.Organization, error)
	// Delete removes the item with the specified id from the repo.
	//
	// If the hard parameter is false the value is only soft deleted
	// and can be later restored.
	Delete(id string, hard bool) error
}

// AreaService is a common interface for a service provider for Area entity
type AreaService interface {
	// List returns a single page of items
	List(page *int, pageSize *int) ([]domain.Area, error)
	// List returns a single page of items filtered by Organization Id
	ListByOrg(org string, page *int, pageSize *int) ([]domain.Area, error)
	// Get returns a single item filter by id
	Get(id string) (domain.Area, error)
	// Create saves a new organization item into the repository
	Create(name string, Description string, logo string) (domain.Area, error)
	// Update looks for an existing item and update the values
	Update(entity domain.Area) (domain.Area, error)
	// Delete removes the item with the specified id from the repo.
	//
	// If the hard parameter is false the value is only soft deleted
	// and can be later restored.
	Delete(id string, hard bool) error
}
