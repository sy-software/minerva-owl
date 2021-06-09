package service

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/mocks"
)

func TestOrganizationIsCreated(t *testing.T) {
	expected := domain.Organization{
		Name:        "name",
		Description: "description",
		Logo:        "logo",
	}

	repo := mocks.OrgInMemoryRepo{
		DummyData: []domain.Organization{},
	}

	service := OrganizationService{
		repository: &repo,
	}

	created, err := service.Create(expected.Name, expected.Description, expected.Logo)

	if err != nil {
		t.Errorf("Item should be created without errors: %v", err)
	}

	match, err := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", created.Id)
	if !match || err != nil {
		t.Errorf("ID is not V4 UUID got: %q with error: %v", created.Id, err)
	}

	if created.Name != expected.Name {
		t.Errorf(
			"Name was not assigned expected: %q got: %q",
			expected.Name,
			created.Name,
		)
	}

	if created.Description != expected.Description {
		t.Errorf(
			"Description was not assigned expected: %q got: %q",
			expected.Description,
			created.Description,
		)
	}

	if created.Logo != expected.Logo {
		t.Errorf(
			"Logo was not assigned expected: %q got: %q",
			expected.Logo,
			created.Logo,
		)
	}
}

func TestOrganizationIsRead(t *testing.T) {
	expected := []domain.Organization{
		{
			Id:          "1",
			Name:        "name 1",
			Description: "description 1",
			Logo:        "logo 1",
		},
		{
			Id:          "2",
			Name:        "name 2",
			Description: "description 2",
			Logo:        "logo 2",
		},
	}

	repo := mocks.OrgInMemoryRepo{
		DummyData: expected,
	}

	service := OrganizationService{
		repository: &repo,
	}

	t.Run("Get All Organizations", func(t *testing.T) {
		got, err := service.All()

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != len(expected) {
			t.Errorf("Expected %d elements got %d", len(expected), len(got))
		}

		if !cmp.Equal(got[0], expected[0]) {
			t.Errorf("Expected first item id to be: %v got: %v", expected[0], got[0])
		}
	})

	t.Run("Get an Organization by id", func(t *testing.T) {
		for _, expect := range expected {
			got, err := service.Get(expect.Id)
			if err != nil {
				t.Errorf(
					"Got error while getting organization by id: %v, Error: %v",
					expect.Id,
					err,
				)
			}

			if !cmp.Equal(got, expect) {
				t.Errorf("Expected item to be: %v got: %v", expect, got)
			}
		}
	})

	t.Run("Get an Organization by a not existing Id", func(t *testing.T) {
		expectedError := "Can't find Organization with Id: not_exists_id"
		_, err := service.Get("not_exists_id")
		if err == nil {
			t.Errorf("Expected 'ErrItemNotFound' error got: %v", err)
		}

		if err.Error() != expectedError {
			t.Errorf("Expected error: %q got: %q", expectedError, err)
		}
	})
}

func TestOrganizationIsUpdated(t *testing.T) {
	base := []domain.Organization{
		{
			Id:          "1",
			Name:        "name 1",
			Description: "description 1",
			Logo:        "logo 1",
		},
		{
			Id:          "2",
			Name:        "name 2",
			Description: "description 2",
			Logo:        "logo 2",
		},
	}

	expected := []domain.Organization{
		{
			Id:          "1",
			Name:        "name updated",
			Description: "description 1",
			Logo:        "logo 1",
		},
		{
			Id:          "2",
			Name:        "name 2",
			Description: "description 2",
			Logo:        "logo 2",
		},
	}

	repo := mocks.OrgInMemoryRepo{
		DummyData: base,
	}

	service := OrganizationService{
		repository: &repo,
	}

	got, err := service.Update(expected[0])

	if err != nil {
		t.Errorf("Got error while getting updating organizations: %v", err)
	}

	if !cmp.Equal(got, expected[0]) {
		t.Errorf("Expected item to be: %v got: %v", expected[0], got)
	}

	for _, expect := range expected {
		got, _ := service.Get(expect.Id)

		if !cmp.Equal(got, expect) {
			t.Errorf("Expected item to be: %v got: %v", expect, got)
		}
	}
}

func TestOrganizationIsDeleted(t *testing.T) {
	base := []domain.Organization{
		{
			Id:          "1",
			Name:        "name",
			Description: "description",
			Logo:        "logo",
		},
		{
			Id:          "2",
			Name:        "name",
			Description: "description",
			Logo:        "logo",
		},
	}

	expected := []domain.Organization{
		{
			Id:          "2",
			Name:        "name",
			Description: "description",
			Logo:        "logo",
		},
	}

	repo := mocks.OrgInMemoryRepo{
		DummyData: base,
	}

	service := OrganizationService{
		repository: &repo,
	}

	err := service.Delete("1", false)

	if err != nil {
		t.Errorf("Got error while getting updating organizations: %v", err)
	}

	all, err := service.All()

	if len(all) != len(expected) {
		t.Errorf("Excted to have %d items got %d", len(expected), len(all))
	}

	for _, expect := range expected {
		got, _ := service.Get(expect.Id)

		if !cmp.Equal(got, expect) {
			t.Errorf("Expected item to be: %v got: %v", expect, got)
		}
	}
}
