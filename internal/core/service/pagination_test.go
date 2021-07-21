package service

import (
	"testing"

	"github.com/sy-software/minerva-owl/internal/core/domain"
)

func TestPaginationNormalization(t *testing.T) {
	config := domain.DefaultConfig()

	t.Run("Test nil page size and page return default values", func(t *testing.T) {
		gotPage, gotSize, skip := pagination(nil, nil, config)

		if gotPage != 0 {
			t.Errorf("Expect first page to be: 0 got: %d", gotPage)
		}

		if gotSize != config.Pagination.PageSize {
			t.Errorf("Expect page size to be: %d got: %d", config.Pagination.PageSize, gotSize)
		}

		if skip != 0 {
			t.Errorf("Expected 0 skip got: %d", skip)
		}
	})

	t.Run("Test an invalid page number returns default page", func(t *testing.T) {
		page := -1
		gotPage, _, _ := pagination(&page, nil, config)

		if gotPage != 0 {
			t.Errorf("Expect first page to be: 0 got: %d", gotPage)
		}

	})

	t.Run("Test an invalid page size number returns default page size", func(t *testing.T) {
		size := -1
		_, gotSize, _ := pagination(nil, &size, config)

		if gotSize != config.Pagination.PageSize {
			t.Errorf("Expect page size to be: %d got: %d", config.Pagination.PageSize, gotSize)
		}

	})

	t.Run("Test a page size greater than max returns max page size", func(t *testing.T) {
		size := config.Pagination.MaxPageSize + 1
		_, gotSize, _ := pagination(nil, &size, config)

		if gotSize != config.Pagination.MaxPageSize {
			t.Errorf("Expect page size to be: %d got: %d", config.Pagination.MaxPageSize, gotSize)
		}

	})

	t.Run("Test the right skip value is returned ", func(t *testing.T) {
		size := 7
		page := 7
		_, _, skip := pagination(&page, &size, config)

		if skip != size*(page-1) {
			t.Errorf("Expect skip to be: %d got: %d", size*(page-1), skip)
		}

	})
}
