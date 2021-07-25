package handlers

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/mocks"
)

func TestOrgToGraphQLModel(t *testing.T) {
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

	got := orgToGraphQLModel(&test)

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

func TestOrgCreateOperation(t *testing.T) {
	t.Run("Create an Organization without logo", func(t *testing.T) {
		data := map[string][]map[string]interface{}{
			"organizations": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

		match, err := regexp.MatchString(mocks.ID_REGEX, got.ID)
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

		if len(repo.Data["organizations"]) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.Data["organizations"]))
		}
	})

	t.Run("Create an Organization with logo", func(t *testing.T) {
		data := map[string][]map[string]interface{}{
			"organizations": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

		match, err := regexp.MatchString(mocks.ID_REGEX, got.ID)
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

		if len(repo.Data["organizations"]) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.Data["organizations"]))
		}
	})
}

func TestOrgQueryOperations(t *testing.T) {
	t.Run("Query all items", func(t *testing.T) {
		base := []map[string]interface{}{
			{
				"id":          "myid1",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
			{
				"id":          "myid2",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
			{
				"id":          "myid3",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
		}
		data := map[string][]map[string]interface{}{
			"organizations": base,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(nil, nil)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(got) != len(repo.Data["organizations"]) {
			t.Errorf("Expected: %d elements got: %d", len(repo.Data["organizations"]), len(got))
		}
	})

	t.Run("Query items by Id", func(t *testing.T) {
		base := []map[string]interface{}{
			{
				"id":          "myid1",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
			{
				"id":          "myid2",
				"name":        "originalName2",
				"description": "originalDescription2",
				"logo":        "",
			},
			{
				"id":          "myid3",
				"name":        "originalName3",
				"description": "originalDescription3",
				"logo":        "",
			},
		}
		data := map[string][]map[string]interface{}{
			"organizations": base,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		logo := ""
		expected := model.Organization{
			ID:          "myid2",
			Name:        "originalName2",
			Description: "originalDescription2",
			Logo:        &logo,
		}
		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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
		data := map[string][]map[string]interface{}{
			"organizations": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

func TestOrgQueryPagination(t *testing.T) {
	dummydata := make([]domain.Organization, 20)
	dummyDict := make([]map[string]interface{}, 20)

	for i := 0; i < 20; i++ {
		str := strconv.Itoa(i)
		dummydata[i] = domain.Organization{
			Id:          str,
			Name:        "name " + str,
			Description: "description " + str,
			Logo:        "logo " + str,
		}

		dummyDict[i] = map[string]interface{}{
			"id":          str,
			"name":        "name " + str,
			"description": "description " + str,
			"logo":        "logo " + str,
		}
	}

	data := map[string][]map[string]interface{}{
		"organizations": dummyDict,
	}

	repo := mocks.MemRepo{
		Data: data,
	}

	t.Run("Get organizations with page size", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5

		orgService := service.NewOrgService(&repo, domain.Config{
			Pagination: domain.Pagination{
				PageSize:    pageSize,
				MaxPageSize: maxPageSize,
			},
		})

		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(nil, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != pageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		if got[0].ID != dummydata[0].Id {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[0], got[0])
		}
	})

	t.Run("Get Organizations second page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5
		page := 2

		orgService := service.NewOrgService(&repo, domain.Config{
			Pagination: domain.Pagination{
				PageSize:    pageSize,
				MaxPageSize: maxPageSize,
			},
		})

		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != pageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		startIndex := (page - 1) * pageSize
		if got[0].ID != dummydata[startIndex].Id {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[startIndex], got[0])
		}
	})

	t.Run("Get Organizations last page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 9
		page := 3

		orgService := service.NewOrgService(&repo, domain.Config{
			Pagination: domain.Pagination{
				PageSize:    pageSize,
				MaxPageSize: maxPageSize,
			},
		})

		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		expectedSize := len(dummydata) - ((page - 1) * pageSize)
		if len(got) != expectedSize {
			t.Errorf("Expected %d elements got %d", expectedSize, len(got))
		}

		startIndex := (page - 1) * pageSize
		if got[0].ID != dummydata[startIndex].Id {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[startIndex], got[0])
		}
	})

	t.Run("Get non existing page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5
		page := 100

		orgService := service.NewOrgService(&repo, domain.Config{
			Pagination: domain.Pagination{
				PageSize:    pageSize,
				MaxPageSize: maxPageSize,
			},
		})

		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}
	})

	t.Run("Pass invalid values for pageSize and page", func(t *testing.T) {
		pageSize := -5
		page := -2

		expectedPageSize := 5

		orgService := service.NewOrgService(&repo, domain.Config{
			Pagination: domain.Pagination{
				PageSize:    expectedPageSize,
				MaxPageSize: 10,
			},
		})

		handlerInstance := NewOrgGraphqlHandler(*orgService)

		got, err := handlerInstance.Query(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != expectedPageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		if got[0].ID != dummydata[0].Id {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[0], got[0])
		}
	})
}

func TestOrgUpdateOperation(t *testing.T) {
	t.Run("Partial Update", func(t *testing.T) {
		base := []map[string]interface{}{
			{
				"id":          "myid",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
		}
		data := map[string][]map[string]interface{}{
			"organizations": base,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

		if len(repo.Data["organizations"]) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.Data["organizations"]))
		}
	})

	t.Run("Complete Update", func(t *testing.T) {
		base := []map[string]interface{}{
			{
				"id":          "myid",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
		}
		data := map[string][]map[string]interface{}{
			"organizations": base,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

		if len(repo.Data["organizations"]) != 1 {
			t.Errorf("Expected repository to have 1 element got: %d", len(repo.Data["organizations"]))
		}
	})

	t.Run("Update a non-existing id", func(t *testing.T) {
		data := map[string][]map[string]interface{}{
			"organizations": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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

func TestOrgDeleteOperaton(t *testing.T) {
	t.Run("Delete an item", func(t *testing.T) {
		base := []map[string]interface{}{
			{
				"id":          "myid",
				"name":        "originalName",
				"description": "originalDescription",
				"logo":        "",
			},
		}
		data := map[string][]map[string]interface{}{
			"organizations": base,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
		handlerInstance := NewOrgGraphqlHandler(*orgService)

		_, err := handlerInstance.Delete("myid")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(repo.Data["organizations"]) != 0 {
			t.Errorf("Expected 0 items in repo got: %d", len(repo.Data["organizations"]))
		}
	})

	t.Run("Delete a non-existing id", func(t *testing.T) {
		data := map[string][]map[string]interface{}{
			"organizations": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		orgService := service.NewOrgService(&repo, domain.DefaultConfig())
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
