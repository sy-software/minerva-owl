package ports

import (
	"fmt"

	"github.com/sy-software/minerva-owl/internal/core/domain"
)

type ErrItemNotFound struct {
	Id    string
	Model string
}

func (err ErrItemNotFound) Error() string {
	return fmt.Sprintf("Can't find %v with Id: %v", err.Model, err.Id)
}

type OrganizationRepo interface {
	All() ([]domain.Organization, error)
	Get(id string) (domain.Organization, error)

	Save(entity domain.Organization) error

	Delete(id string) error
}
