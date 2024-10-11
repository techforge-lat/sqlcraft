package sqlcraft

import (
	"strconv"
	"strings"

	"github.com/techforge-lat/dafi/v2"
)

type SelectQuery struct {
	table                  string
	columns                []string
	requiredColumns        map[string]struct{}
	sqlColumnByDomainField map[string]string

	filters    dafi.Filters
	sorts      dafi.Sorts
	pagination dafi.Pagination
}

func Select(columns ...string) SelectQuery {
	return SelectQuery{
		table:           "",
		columns:         columns,
		requiredColumns: make(map[string]struct{}),
	}
}

func (s SelectQuery) From(table string) SelectQuery {
	s.table = table

	return s
}

func (s SelectQuery) Where(filters ...dafi.Filter) SelectQuery {
	s.filters = filters

	return s
}

func (s SelectQuery) OrderBy(sorts ...dafi.Sort) SelectQuery {
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

// RequiredColumns allows you to select just some of the columns provided in the Select func
func (s SelectQuery) RequiredColumns(columns ...string) SelectQuery {
	for _, col := range columns {
		s.requiredColumns[col] = struct{}{}
	}

	return s
}

func (s SelectQuery) SqlColumnByDomainField(sqlColumnByDomainField map[string]string) SelectQuery {
	s.sqlColumnByDomainField = sqlColumnByDomainField

	return s
}

func (s SelectQuery) ToSQL() (Result, error) {
	if len(s.columns) == 0 {
		return Result{}, ErrEmptyColumns
	}

	if len(s.sqlColumnByDomainField) > 0 {
		requiredCols := make(map[string]struct{})
		for k := range s.requiredColumns {
			requiredSqlColumn, ok := s.sqlColumnByDomainField[k]
			if !ok {
				return Result{}, ErrInvalidFieldName
			}

			requiredCols[requiredSqlColumn] = struct{}{}
		}

		s.requiredColumns = requiredCols
	}

	builder := strings.Builder{}

	builder.WriteString("SELECT ")

	if len(s.requiredColumns) == 0 {
		builder.WriteString(strings.Join(s.columns, ", "))
	} else {
		for i, col := range s.columns {
			if _, ok := s.requiredColumns[col]; ok {
				builder.WriteString(col)
			} else {
				builder.WriteString("null")
			}

			if i < len(s.columns)-1 {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(" FROM ")
	builder.WriteString(s.table)

	args := []any{}
	if len(s.filters) > 0 {
		whereResult, err := WhereSafe(s.sqlColumnByDomainField, s.filters...)
		if err != nil {
			return Result{}, err
		}
		args = append(args, whereResult.Args...)

		builder.WriteString(" ")
		builder.WriteString(whereResult.Sql)
	}

	if len(s.sorts) > 0 {
		sortSql := BuildOrderBy(s.sorts)

		builder.WriteString(" ")
		builder.WriteString(sortSql)
	}

	paginationSql := BuildPagination(s.pagination)
	builder.WriteString(paginationSql)

	return Result{
		Sql:  builder.String(),
		Args: args,
	}, nil
}

func BuildOrderBy(sorts dafi.Sorts) string {
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
	builder.WriteString(" LIMIT ")
	builder.WriteString(strconv.Itoa(int(pagination.PageSize)))

	if pagination.HasPageNumber() {
		builder.WriteString(" OFFSET ")
		builder.WriteString(strconv.Itoa(int(pagination.PageSize * (pagination.PageNumber - 1))))
	}

	return builder.String()
}
