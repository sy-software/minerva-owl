package service

import (
	"regexp"
	"strconv"
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
		config:     domain.DefaultConfig(),
	}

	created, err := service.Create(expected.Name, expected.Description, expected.Logo)

	if err != nil {
		t.Errorf("Item should be created without errors: %v", err)
	}

	match, err := regexp.MatchString(ID_REGEX, created.Id)
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
		config: domain.Config{
			Pagination: domain.Pagination{
				PageSize:    10,
				MaxPageSize: 100,
			},
		},
	}

	t.Run("Get a list of Organizations", func(t *testing.T) {
		got, err := service.List(nil, nil)

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

func TestPagination(t *testing.T) {
	dummydata := make([]domain.Organization, 20)

	for i := 0; i < 20; i++ {
		str := strconv.Itoa(i)
		dummydata[i] = domain.Organization{
			Id:          str,
			Name:        "name " + str,
			Description: "description " + str,
			Logo:        "logo " + str,
		}
	}

	repo := mocks.OrgInMemoryRepo{
		DummyData: dummydata,
	}

	t.Run("Get Organizations with Page Size", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5

		service := OrganizationService{
			repository: &repo,
			config: domain.Config{
				Pagination: domain.Pagination{
					PageSize:    pageSize,
					MaxPageSize: maxPageSize,
				},
			},
		}

		got, err := service.List(nil, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != pageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		if !cmp.Equal(got[0], dummydata[0]) {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[0], got[0])
		}
	})

	t.Run("Get Organizations second page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5
		page := 2

		service := OrganizationService{
			repository: &repo,
			config: domain.Config{
				Pagination: domain.Pagination{
					PageSize:    pageSize,
					MaxPageSize: maxPageSize,
				},
			},
		}

		got, err := service.List(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != pageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		startIndex := (page - 1) * pageSize
		if !cmp.Equal(got[0], dummydata[startIndex]) {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[startIndex], got[0])
		}
	})

	t.Run("Get Organizations last page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 9
		page := 3

		service := OrganizationService{
			repository: &repo,
			config: domain.Config{
				Pagination: domain.Pagination{
					PageSize:    pageSize,
					MaxPageSize: maxPageSize,
				},
			},
		}

		got, err := service.List(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		expectedSize := len(dummydata) - ((page - 1) * pageSize)
		if len(got) != expectedSize {
			t.Errorf("Expected %d elements got %d", expectedSize, len(got))
		}

		startIndex := (page - 1) * pageSize
		if !cmp.Equal(got[0], dummydata[startIndex]) {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[startIndex], got[0])
		}
	})

	t.Run("Get non existing page", func(t *testing.T) {
		maxPageSize := 10
		pageSize := 5
		page := 100

		service := OrganizationService{
			repository: &repo,
			config: domain.Config{
				Pagination: domain.Pagination{
					PageSize:    pageSize,
					MaxPageSize: maxPageSize,
				},
			},
		}

		got, err := service.List(&page, &pageSize)

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

		service := OrganizationService{
			repository: &repo,
			config: domain.Config{
				Pagination: domain.Pagination{
					PageSize:    expectedPageSize,
					MaxPageSize: 10,
				},
			},
		}

		got, err := service.List(&page, &pageSize)

		if err != nil {
			t.Errorf("Got error while getting all organizations: %v", err)
		}

		if len(got) != expectedPageSize {
			t.Errorf("Expected %d elements got %d", pageSize, len(got))
		}

		if !cmp.Equal(got[0], dummydata[0]) {
			t.Errorf("Expected first item id to be: %v got: %v", dummydata[0], got[0])
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
		config:     domain.DefaultConfig(),
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
		config:     domain.DefaultConfig(),
	}

	err := service.Delete("1", false)

	if err != nil {
		t.Errorf("Got error while getting updating organizations: %v", err)
	}

	all, err := service.List(nil, nil)

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
