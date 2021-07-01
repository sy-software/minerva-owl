package repositories

import (
	"fmt"

	"github.com/scylladb/gocqlx/table"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

var tableName = "minerva.organization"

// metadata specifies table name and columns it must be in sync with schema.
var orgMetadata = table.Metadata{
	Name:    tableName,
	Columns: []string{"id", "name", "description", "logo"},
	PartKey: []string{"id"},
	SortKey: []string{"name"},
}

// orgTable allows for simple CRUD operations based on orgMetadata.
var orgTable = table.New(orgMetadata)

type OrgRepo struct {
	cassandra *Cassandra
}

func NewOrgRepo(cassandra *Cassandra) (*OrgRepo, error) {
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (id text, name text, description text, logo text, PRIMARY KEY (id));",
		tableName,
	)
	err := cassandra.session.ExecStmt(query)
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
	q := memRepo.cassandra.session.Query(orgTable.Insert()).BindStruct(entity)
	return q.ExecRelease()
}

func (memRepo *OrgRepo) Delete(id string) error {
	return nil
}
