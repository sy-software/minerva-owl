package service

import (
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/utils"
)

func pagination(reqPage *int, reqPageSize *int, config domain.Config) (page int, pageSize int, skip int) {
	page = utils.CoalesceInt(reqPage, 1) - 1
	pageSize = utils.CoalesceInt(reqPageSize, config.Pagination.PageSize)

	if pageSize > config.Pagination.MaxPageSize {
		pageSize = config.Pagination.MaxPageSize
	}

	if pageSize < 0 {
		pageSize = config.Pagination.PageSize
	}

	if page < 0 {
		page = 0
	}

	skip = page * pageSize

	return
}
