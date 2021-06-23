package handlers

import (
	"regexp"
	"testing"

	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/mocks"
)

func TestNilCoalescing(t *testing.T) {
	defaultValue := "default"
	got := nilCoalescing(nil, defaultValue)

	if defaultValue != got {
		t.Errorf("Expected: %q Got: %q", defaultValue, got)
	}

	expected := "expected"
	got = nilCoalescing(&expected, defaultValue)

	if got != expected {
		t.Errorf("Expected: %q Got: %q", defaultValue, got)
	}
}

func TestDomainToGraphQLModel(t *testing.T) {
	logo := "myLogo"
	expected := model.Organization{
		ID:          "myId",
		Name:        "myName",
		Description: "Description",
		Logo:        &logo,
	}

	test := domain.Organization{
		Id:          expected.ID,
		Name:        expected.Name,
		Description: expected.Description,
		Logo:        *expected.Logo,
	}

	got := domainToGraphQLModel(&test)

	if got.ID != expected.ID {
		t.Errorf("Expected ID to be: %q got: %q", expected.ID, got.ID)
	}

	if got.Name != expected.Name {
		t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
	}

	if got.Description != expected.Description {
		t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
	}

	if *got.Logo != *expected.Logo {
		t.Errorf("Expected.Logo to be: %q got: %q", *expected.Logo, *got.Logo)
	}
}

func TestCreateOperation(t *testing.T) {
	t.Run("Create an Organization without logo", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		expected := model.Organization{
			ID:          "...",
			Name:        "myName",
			Description: "Description",
			Logo:        nil,
		}
		got, err := handlerInstance.Create(expected.Name, expected.Description, expected.Logo)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		match, err := regexp.MatchString(service.ID_REGEX, got.ID)
		if !match || err != nil {
			t.Errorf("ID is not V4 UUID got: %q with error: %v", got.ID, err)
		}

		if got.Name != expected.Name {
			t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
		}

		if got.Description != expected.Description {
			t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
		}

		if *got.Logo != "" {
			t.Errorf("Expected.Logo to be \"\" got: %q", *got.Logo)
		}

		if len(repo.DummyData) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.DummyData))
		}
	})

	t.Run("Create an Organization with logo", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		logo := "logo"
		expected := model.Organization{
			ID:          "...",
			Name:        "myName",
			Description: "Description",
			Logo:        &logo,
		}
		got, err := handlerInstance.Create(expected.Name, expected.Description, expected.Logo)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		match, err := regexp.MatchString(service.ID_REGEX, got.ID)
		if !match || err != nil {
			t.Errorf("ID is not V4 UUID got: %q with error: %v", got.ID, err)
		}

		if got.Name != expected.Name {
			t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
		}

		if got.Description != expected.Description {
			t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
		}

		if *got.Logo != *expected.Logo {
			t.Errorf("Expected.Logo to be: %q got: %q", *expected.Logo, *got.Logo)
		}

		if len(repo.DummyData) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.DummyData))
		}
	})
}

func TestQueryOperations(t *testing.T) {
	t.Run("Query all items", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{
				{
					Id:          "myid1",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
				{
					Id:          "myid2",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
				{
					Id:          "myid3",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
			},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(got) != len(repo.DummyData) {
			t.Errorf("Expected: %d elements got: %d", len(repo.DummyData), len(got))
		}
	})

	t.Run("Query items by Id", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{
				{
					Id:          "myid1",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
				{
					Id:          "myid2",
					Name:        "originalName2",
					Description: "originalDescription2",
					Logo:        "",
				},
				{
					Id:          "myid3",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
			},
		}

		logo := ""
		expected := model.Organization{
			ID:          "myid2",
			Name:        "originalName2",
			Description: "originalDescription2",
			Logo:        &logo,
		}
		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.QueryById(expected.ID)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.ID != expected.ID {
			t.Errorf("Expected ID to be: %q got: %q", expected.ID, got.ID)
		}

		if got.Name != expected.Name {
			t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
		}

		if got.Description != expected.Description {
			t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
		}
	})

	t.Run("Query a non-existing id", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		_, err := handlerInstance.QueryById("myid")

		if err == nil {
			t.Errorf("Expected error got nil")
		}

		_, ok := err.(ports.ErrItemNotFound)
		if !ok {
			t.Errorf("Expected error of type ErrItemNotFound got: %T", err)
		}
	})
}

func TestUpdateOperation(t *testing.T) {
	t.Run("Partial Update", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{
				{
					Id:          "myid",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
			},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		logo := "logo"

		expected := model.Organization{
			ID:          "myid",
			Name:        "newName",
			Description: "originalDescription",
			Logo:        &logo,
		}
		got, err := handlerInstance.Update(expected.ID, &expected.Name, nil, &logo)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.Name != expected.Name {
			t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
		}

		if got.Description != expected.Description {
			t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
		}

		if *got.Logo != *expected.Logo {
			t.Errorf("Expected.Logo to be: %q got: %q", *expected.Logo, *got.Logo)
		}

		if len(repo.DummyData) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.DummyData))
		}
	})

	t.Run("Complete Update", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{
				{
					Id:          "myid",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
			},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		logo := "logo"

		expected := model.Organization{
			ID:          "myid",
			Name:        "newName",
			Description: "newDescription",
			Logo:        &logo,
		}
		got, err := handlerInstance.Update(expected.ID, &expected.Name, &expected.Description, &logo)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.Name != expected.Name {
			t.Errorf("Expected Name to be: %q got: %q", expected.Name, got.Name)
		}

		if got.Description != expected.Description {
			t.Errorf("Expected Description to be: %q got: %q", expected.Description, got.Description)
		}

		if *got.Logo != *expected.Logo {
			t.Errorf("Expected.Logo to be: %q got: %q", *expected.Logo, *got.Logo)
		}

		if len(repo.DummyData) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.DummyData))
		}
	})

	t.Run("Update a non-existing id", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		logo := "logo"

		expected := model.Organization{
			ID:          "myid",
			Name:        "newName",
			Description: "originalDescription",
			Logo:        &logo,
		}
		_, err := handlerInstance.Update(expected.ID, &expected.Name, nil, &logo)

		if err == nil {
			t.Errorf("Expected error got nil")
		}

		_, ok := err.(ports.ErrItemNotFound)
		if !ok {
			t.Errorf("Expected error of type ErrItemNotFound got: %T", err)
		}
	})
}

func TestDeleteOperaton(t *testing.T) {
	t.Run("Delete an item", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{
				{
					Id:          "myid",
					Name:        "originalName",
					Description: "originalDescription",
					Logo:        "",
				},
			},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		_, err := handlerInstance.Delete("myid")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(repo.DummyData) != 0 {
			t.Errorf("Expected 0 items in repo got: %d", len(repo.DummyData))
		}
	})

	t.Run("Delete a non-existing id", func(t *testing.T) {
		repo := &mocks.OrgInMemoryRepo{
			DummyData: []domain.Organization{},
		}

		orgService := service.NewOrgService(repo)
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		_, err := handlerInstance.Delete("id")

		if err == nil {
			t.Errorf("Expected error got nil")
		}

		_, ok := err.(ports.ErrItemNotFound)
		if !ok {
			t.Errorf("Expected error of type ErrItemNotFound got: %T", err)
		}
	})
}
