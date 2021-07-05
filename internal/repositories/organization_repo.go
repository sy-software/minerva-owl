package repositories

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sy-software/minerva-owl/internal/core/domain"
)

var tableName = "minerva.organizations"

// metadata specifies table name and columns it must be in sync with schema.
var orgMetadata = table.Metadata{
	Name:    tableName,
	Columns: []string{"id", "name", "description", "logo"},
	PartKey: []string{"id"},
}

// orgTable allows for simple CRUD operations based on orgMetadata.
var orgTable = table.New(orgMetadata)

type OrgRepo struct {
	cassandra *Cassandra
	config    *domain.Config
}

// NewOrgRepo creates an instances of the Organization repo with Cassandra DB
func NewOrgRepo(cassandra *Cassandra, config *domain.Config) (*OrgRepo, error) {
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
		config:    config,
	}, nil
}

func (repo *OrgRepo) List(skip int, limit int) ([]domain.Organization, error) {
	var orgs []domain.Organization
	stmt, names, err := gocqlx.CompileNamedQueryString(
		fmt.Sprintf("SELECT * FROM %s LIMIT :limit", tableName),
	)

	if err != nil {
		return []domain.Organization{}, err
	}

	log.Debug().Msgf("Statement: %v Names: %v", stmt, names)
	q := repo.cassandra.session.
		Query(stmt, names).BindMap(qb.M{
		"limit": repo.config.Pagination.PageSize,
	})

	log.Debug().Msgf("Query: %v", q)

	if err := q.SelectRelease(&orgs); err != nil {
		log.Debug().Err(err).Msgf("Error in query %v", q)
		return []domain.Organization{}, err
	}

	log.Debug().Msgf("Quering values: %d", len(orgs))
	return orgs, nil
}

func (repo *OrgRepo) Get(id string) (domain.Organization, error) {
	org := domain.Organization{}

	stmt, names := orgTable.Select()

	log.Debug().Msgf("Statement: %v Names: %v", stmt, names)
	q := repo.cassandra.session.Query(stmt, names).BindMap(qb.M{
		"id": id,
	})
	if err := q.GetRelease(&org); err != nil {
		return org, err
	}

	return org, nil
}

func (repo *OrgRepo) Save(entity domain.Organization) error {
	log.Debug().Msgf("Updating: %v", entity)
	q := repo.cassandra.session.Query(orgTable.Insert()).BindStruct(entity)
	return q.ExecRelease()
}

func (repo *OrgRepo) Delete(id string) error {
	q := repo.cassandra.session.Query(orgTable.Delete()).BindMap(qb.M{
		"id": id,
	})

	err := q.ExecRelease()

	log.Debug().Err(err).Msg("Delete error: ")

	return err
}
