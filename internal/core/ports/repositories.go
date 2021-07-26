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

// Filter is used to privide an abstraction for repository specific filters
// Each repository have the responsibility to parse the filters into the right
// query representation (E.G.: SQL, CQL, etc.)
type Filter struct {
	Name  string
	Value interface{}
}

type Repository interface {
	// List returns a single page of items
	List(collection string, results interface{}, skip int, limit int, filters ...Filter) error
	// Get returns a single item filter by id
	Get(collection string, id string, result interface{}) error
	// Get returns a single item filtered with the provided filters
	GetOne(collection string, result interface{}, filter ...Filter) error
	// Create saves a new item into the repository and returns the assigned Id
	Create(collection string, entity interface{}) (string, error)
	// Update looks for an existing item and update the values omiting the fields in omit
	Update(collection string, id string, entity interface{}, omit ...string) error
	// Delete removes the item with the specified id from the repo
	Delete(collection string, id string) error
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
