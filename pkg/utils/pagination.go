package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
	SortBy string
	SortDir string
}

var allowedSortDirs = map[string]bool{
	"ASC": true,
	"DESC": true,
}

func GeneratePagination(pageStr, limitStr, sortBy, sortDir string, allowedColumns map[string]bool) Pagination {
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	if !allowedColumns[sortBy] {
		sortBy = "id"
	}

	sortDir = strings.ToUpper(sortDir)
	if sortDir != "DESC" {
		sortDir = "ASC"
	}

	return Pagination{
		Page: page,
		Limit: limit,
		Offset: offset,
		SortBy: sortBy,
		SortDir: sortDir,
	}
}

func (p Pagination) BuildOrderBy() string {
	return fmt.Sprintf(" ORDER BY %s %s", p.SortBy, p.SortDir)
}
