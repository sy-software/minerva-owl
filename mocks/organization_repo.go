package mocks

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

type OrgInMemoryRepo struct {
	DummyData []domain.Organization
}

func (memRepo *OrgInMemoryRepo) All() ([]domain.Organization, error) {
	return memRepo.DummyData, nil
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

func (memRepo *OrgInMemoryRepo) Save(entity domain.Organization) error {
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
