package sqlcraft

import (
	"strconv"
	"strings"

	"github.com/techforge-lat/dafi/v2"
)

type SelectQuery struct {
	table   string
	columns []string
	values  []any

	filters    dafi.Filters
	sorts      dafi.Sorts
	pagination dafi.Pagination

	offset uint
}

func Select(columns ...string) SelectQuery {
	return SelectQuery{
		table:   "",
		columns: columns,
		values:  []any{},
	}
}

func (s SelectQuery) From(table string) SelectQuery {
	s.table = table

	return s
}

func (s SelectQuery) Where(filters dafi.Filters) SelectQuery {
	s.filters = filters

	return s
}

func (s SelectQuery) OrderBy(sorts dafi.Sorts) SelectQuery {
	s.sorts = sorts

	return s
}

func (s SelectQuery) Limit(limit uint) SelectQuery {
	s.pagination.PageSize = limit

	return s
}

func (s SelectQuery) Page(page uint) SelectQuery {
	s.pagination.PageNumber = page

	return s
}

func (s SelectQuery) Offset(offset uint) SelectQuery {
	s.offset = offset

	return s
}

func OrderBy(sorts dafi.Sorts) string {
	if sorts.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("ORDER BY ")
	for i, sort := range sorts {
		builder.WriteString(string(sort.Field))

		if sort.Type != dafi.None {
			builder.WriteString(" ")
			builder.WriteString(strings.ToUpper(string(sort.Type)))
		}

		if i < len(sorts)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func BuildPagination(pagination dafi.Pagination) string {
	if pagination.HasPageSize() && !pagination.HasPageNumber() {
		pagination.PageNumber = 1
	}

	if pagination.IsZero() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("LIMIT ")
	builder.WriteString(strconv.Itoa(int(pagination.PageSize)))

	if pagination.HasPageNumber() {
		builder.WriteString(" OFFSET ")
		builder.WriteString(strconv.Itoa(int(pagination.PageSize * (pagination.PageNumber - 1))))
	}

	return builder.String()
}
