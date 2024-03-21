package sqlcraft

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var ErrDuplicatedOption = errors.New("cannot have duplicated options")

type optionKey string

func (o optionKey) IsInList(list optionKeys) bool {
	for _, item := range list {
		if item == o {
			return true
		}
	}

	return false
}

const (
	returning  optionKey = "RETURNING"
	where      optionKey = "WHERE"
	order      optionKey = "ORDER"
	limit      optionKey = "LIMIT"
	offset     optionKey = "OFFSET"
	pagination optionKey = "PAGINATION"
	from       optionKey = "FROM"
)

type optionKeys []optionKey

type options struct {
	sql   string
	args  []any
	key   optionKey
	order uint
	err   error

	excludeWhereKeyword bool
	paramCountStartFrom uint
}

type Option func(option *options) error

type Options []Option

// WithReturning creates a RETURNING sql expression
// if `columns` param is empty, the `*` will be used insetead
func WithReturning(columns ...string) Option {
	return func(option *options) error {
		if len(columns) == 0 {
			columns = append(columns, "*")
		}

		option.sql = fmt.Sprintf("RETURNING %s", strings.Join(columns, ", "))
		option.key = returning

		return nil
	}
}

// WithFrom creates a FROM sql expression
// if `expression` param is empty, an error will be returned
// NOTE: It is not intended to be used with a subquery
func WithFrom(expression string) Option {
	return func(option *options) error {
		if expression == "" {
			return errors.New("expression in FROM option cannot be empty")
		}

		option.sql = fmt.Sprintf("FROM %s", expression)
		option.key = from

		return nil
	}
}

type Sorter interface {
	Field() string
	Order() string
}

func WithSort(allowedColumns AllowedColumns, items ...Sorter) Option {
	return func(option *options) error {
		if len(items) == 0 {
			return errors.New("sort items cannot be empty in ORDER BY option")
		}

		builder := bytes.Buffer{}
		builder.WriteString(" ORDER BY ")
		for _, v := range items {
			builder.WriteString(v.Field())
			builder.WriteString(" ")
			builder.WriteString(v.Order())
			builder.WriteString(", ")
		}

		// removes the last `, `
		builder.Truncate(builder.Len() - 2)
		option.sql = builder.String()
		option.key = order

		return nil
	}
}

func WithPagination(limit, offset uint) Option {
	return func(option *options) error {
		if limit == 0 || offset == 0 {
			return nil
		}

		limitParamNumber := option.paramCountStartFrom + 1
		offsetParamNumber := limitParamNumber + 1

		option.sql = fmt.Sprintf("LIMIT $%d OFFSET $%d", limitParamNumber, offsetParamNumber)
		option.args = append(option.args, limit, offset)
		option.key = pagination

		return nil
	}
}
