package repositories

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

type OrgRepo struct {
	cassandra *Cassandra
}

func NewOrgRepo(cassandra *Cassandra) (*OrgRepo, error) {
	err := cassandra.session.Query("CREATE TABLE IF NOT EXISTS minerva.organizations (id text, name text, description text, logo text, PRIMARY KEY (id));").Exec()
	if err != nil {
		return nil, err
	}

	return &OrgRepo{
		cassandra: cassandra,
	}, nil
}

func (memRepo *OrgRepo) All() ([]domain.Organization, error) {
	return []domain.Organization{}, nil
}
func (memRepo *OrgRepo) Get(id string) (domain.Organization, error) {
	return domain.Organization{}, ports.ErrItemNotFound{
		Id:    id,
		Model: "Organization",
	}
}

func (memRepo *OrgRepo) Save(entity domain.Organization) error {
	return nil
}

func (memRepo *OrgRepo) Delete(id string) error {
	return nil
}
