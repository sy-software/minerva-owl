package service

import (
	"regexp"
	"testing"

	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/mocks"
)

func TestOrganizationIsCreated(t *testing.T) {
	repo := mocks.OrgInMemoryRepo{
		DummyData: []domain.Organization{},
	}

	service := OrganizationService{
		repository: &repo,
	}

	created, err := service.Create("name", "description", "logo")

	if err != nil {
		t.Errorf("Item should be created without errors: %v", err)
	}

	match, err := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89AB][0-9a-f]{3}-[0-9a-f]{12}$", created.Id)
	if !match || err != nil {
		t.Errorf("ID is not V4 UUID got: %v", created.Id)
	}

	if created.Name != "name" {
		t.Errorf("Name was not assigned expected: name got: %v", created.Name)
	}

	if created.Description != "description" {
		t.Errorf("Description was not assigned expected: description got: %v", created.Description)
	}

	if created.Logo != "logo" {
		t.Errorf("Logo was not assigned expected: logo got: %v", created.Logo)
	}
}
