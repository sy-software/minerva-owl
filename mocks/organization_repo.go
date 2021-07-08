package mocks

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

type OrgInMemoryRepo struct {
	DummyData []domain.Organization
}

func (memRepo *OrgInMemoryRepo) List(skip int, limit int) ([]domain.Organization, error) {
	if skip >= len(memRepo.DummyData) {
		return []domain.Organization{}, nil
	}

	available := len(memRepo.DummyData) - skip
	capLimit := limit
	if limit > available {
		capLimit = available
	}
	// log.Debug().Msgf("MemRepo: Skip = %d, Limit = %d, Cap Limit = %d, available = %d", skip, limit, capLimit, available)
	return memRepo.DummyData[skip : skip+capLimit], nil
}
func (memRepo *OrgInMemoryRepo) Get(id string) (domain.Organization, error) {
	for _, item := range memRepo.DummyData {
		if item.Id == id {
			return item, nil
		}
	}

	return domain.Organization{}, ports.ErrItemNotFound{
		Id:    id,
		Model: "Organization",
	}
}

func (memRepo *OrgInMemoryRepo) Create(entity domain.Organization) (string, error) {
	entity.Id = uuid.New().String()
	log.Debug().Msgf("creating an id: %q", entity.Id)
	return entity.Id, memRepo.save(entity)
}

func (memRepo *OrgInMemoryRepo) Update(entity domain.Organization) error {
	return memRepo.save(entity)
}

func (memRepo *OrgInMemoryRepo) save(entity domain.Organization) error {
	for index, item := range memRepo.DummyData {
		if item.Id == entity.Id {
			memRepo.DummyData[index] = entity
			return nil
		}
	}

	memRepo.DummyData = append(memRepo.DummyData, entity)
	return nil
}

func (memRepo *OrgInMemoryRepo) Delete(id string) error {
	newData := []domain.Organization{}

	for _, item := range memRepo.DummyData {
		if item.Id != id {
			newData = append(newData, item)
		}
	}

	memRepo.DummyData = newData
	return nil
}
