package ports

import (
	"fmt"

	"github.com/sy-software/minerva-owl/internal/core/domain"
)

// ErrItemNotFound must be thrown when a operation is tried against a nonexisting item
type ErrItemNotFound struct {
	// The Id of the requested item
	Id string
	// Which domain model this item belongs to
	Model string
}

func (err ErrItemNotFound) Error() string {
	return fmt.Sprintf("Can't find %v with Id: %v", err.Model, err.Id)
}

// OrganizationRepo is the commong interface for repository providers for the Organization model
type OrganizationRepo interface {
	// List returns a single page of items
	List(skip int, limit int) ([]domain.Organization, error)
	// Get returns a single item filter by id
	Get(id string) (domain.Organization, error)
	// Create saves a new item into the repository
	Create(entity domain.Organization) (string, error)
	// Update looks for an existing item and update the values
	Update(entity domain.Organization) error
	// Delete removes the item with the specified id from the repo
	Delete(id string) error
}
