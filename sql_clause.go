package sqlcraft

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var ErrDuplicatedOption = errors.New("cannot have duplicated sql clauses")

type sqlClause string

const (
	returning sqlClause = "RETURNING"
	where     sqlClause = "WHERE"
	order     sqlClause = "ORDER"
	limit     sqlClause = "LIMIT"
	offset    sqlClause = "OFFSET"
	from      sqlClause = "FROM"
)

type sqlClauses []sqlClause

type sqlClauseConfig struct {
	// expression the final sql expression that will be joined with the main SQL query.
	expression string

	// args store the values that must be used for the parameters in the SQL expression.
	args []any

	// sqlClause indicates what clause this config belongs to.
	sqlClause sqlClause

	// excludeWhereKeyword is only used by the WHERE clause
	// to indicate that the WHERE keyword has been already set
	// in the final SQL query.
	excludeWhereKeyword bool

	// paramCountStartFrom is used for those SQL clauses
	// where you need to pass some data and to know the parameter you
	// have to start from, every SQL clause uses this value as the start point.
	// NOTE: this field is meant to be set by the core.Build function.
	paramCountStartFrom uint
}

// SQLClause constraint of modification used to define the data you want to get or manipulate.
type SQLClause func(config *sqlClauseConfig) error

// SQLClauses list of SQLClause.
type SQLClauses []SQLClause

// WithReturning creates a RETURNING SQL clause.
// If `columns` param is empty, the `*` will be used insetead.
func WithReturning(columns ...string) SQLClause {
	return func(option *sqlClauseConfig) error {
		if len(columns) == 0 {
			columns = append(columns, "*")
		}

		option.expression = fmt.Sprintf("RETURNING %s", strings.Join(columns, ", "))
		option.sqlClause = returning

		return nil
	}
}

// WithFrom creates a FROM SQL clause
// if `expression` param is empty, an error will be returned
// NOTE: It is not intended to be used with a subquery.
func WithFrom(expression string) SQLClause {
	return func(option *sqlClauseConfig) error {
		if expression == "" {
			return errors.New("expression in FROM option cannot be empty")
		}

		option.expression = fmt.Sprintf("FROM %s", expression)
		option.sqlClause = from

		return nil
	}
}

type SortItem interface {
	GetField() string
	GetOrder() string
}

func WithSort(allowedColumns AllowedColumns, items ...SortItem) SQLClause {
	return func(option *sqlClauseConfig) error {
		if len(items) == 0 {
			return errors.New("sort items cannot be empty in ORDER BY option")
		}

		builder := bytes.Buffer{}
		builder.WriteString(" ORDER BY ")
		for _, item := range items {
			builder.WriteString(item.GetField())
			builder.WriteString(" ")
			builder.WriteString(item.GetOrder())
			builder.WriteString(", ")
		}

		// removes the last `, `
		builder.Truncate(builder.Len() - 2)
		option.expression = builder.String()
		option.sqlClause = order

		return nil
	}
}

func WithLimit(value uint) SQLClause {
	return func(config *sqlClauseConfig) error {
		if value == 0 {
			return nil
		}

		limitParamNumber := config.paramCountStartFrom + 1

		config.expression = fmt.Sprintf("LIMIT $%d", limitParamNumber)
		config.args = append(config.args, value)
		config.sqlClause = limit

		return nil
	}
}

func WithOffset(value uint) SQLClause {
	return func(config *sqlClauseConfig) error {
		if value == 0 {
			return nil
		}

		offsetParamNumber := config.paramCountStartFrom + 1

		config.expression = fmt.Sprintf("OFFSET $%d", offsetParamNumber)
		config.args = append(config.args, value)
		config.sqlClause = limit

		return nil
	}
}

func WithParamCount(count int) SQLClause {
	return func(option *sqlClauseConfig) error {
		option.paramCountStartFrom = uint(count)
		return nil
	}
}

func withExcludeWhereKeyword(exclude bool) SQLClause {
	return func(option *sqlClauseConfig) error {
		option.excludeWhereKeyword = exclude
		return nil
	}
}
