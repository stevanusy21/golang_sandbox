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

func GeneratePaginationData(pageStr, limitStr, sortBy, sortDir string, allowedSortColumns map[string]bool) Pagination {
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	if !allowedSortColumns[sortBy] {
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

type PaginationQueryBuilder struct {
	query strings.Builder
	args []any
	paramIndex int
}

func NewPaginationQueryBuilder(baseQuery string) *PaginationQueryBuilder {
	q := &PaginationQueryBuilder{
		paramIndex: 1,
	}

	q.query.WriteString(baseQuery)

	return q
}

func (q *PaginationQueryBuilder) Where(condition string, column string, operator string, value any) {
	q.query.WriteString(fmt.Sprintf(" %s %s %s $%d", condition, column, operator, q.paramIndex))
	q.args = append(q.args, value)
	q.paramIndex++
}

func (q *PaginationQueryBuilder) WhereAny(condition string, column string, values any) {
	q.query.WriteString(fmt.Sprintf(" %s %s = ANY($%d)", condition, column, q.paramIndex))
	q.args = append(q.args, values)
	q.paramIndex++
}

func (q *PaginationQueryBuilder) WhereNull(condition string, column string) {
	q.query.WriteString(fmt.Sprintf(" %s %s IS NULL", condition, column))
}

func (q *PaginationQueryBuilder) OrderBy(sortBy string, sortDir string) {
	q.query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir))
}

func (q *PaginationQueryBuilder) LimitOffset(limit, offset int) {
	q.query.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", q.paramIndex, q.paramIndex+1))
	q.args = append(q.args, limit, offset)
	q.paramIndex += 2
}

func (q *PaginationQueryBuilder) Build() (string, []any) {
	return q.query.String(), q.args
}

type PaginationResponse struct {
	Data       any `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

func NewPaginationResponse(data any, total, page, limit int) PaginationResponse {
	totalPages := 0
	if limit > 0 {
		totalPages = total / limit
		if total%limit != 0 {
			totalPages++
		}
	}

	return PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}