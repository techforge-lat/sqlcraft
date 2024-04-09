package sqlcraft

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidFieldName = errors.New("invalid field name")

type AllowedColumns map[string]string

const (
	And = "AND"
	Or  = "OR"
)

const (
	Equal              = "="
	LessThan           = "<"
	GreaterThan        = ">"
	LessThanOrEqual    = "<="
	GreaterThanOrEqual = ">="
	NotEqual           = "<>"
	Is                 = "IS"
	IsNot              = "IS_NOT"
	ILike              = "ILIKE"
	NotILike           = "NOT_ILIKE"
	Like               = "LIKE"
	NotLike            = "NOT_LIKE"
	In                 = "IN"
)

type FilterItem interface {
	GetField() string
	GetOperator() string
	GetParsedOperator() (string, error)
	GetValue() any
	GetChainingKey() string
	GetGroupOpen() string
	HasGroupOpen() bool
	GetGroupClose() string
	HasGroupClose() bool
}

type FilterItems []FilterItem

func WithWhere(collection ...FilterItem) SQLClause {
	return WithSafeWhere(nil, collection...)
}

func WithSafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) SQLClause {
	return func(config *sqlClauseConfig) error {
		if len(collection) == 0 {
			return nil
		}

		builder := bytes.Buffer{}
		if !config.excludeWhereKeyword {
			builder.WriteString(" WHERE ")
		}

		count := int(config.paramCountStartFrom)

		for index, item := range collection {
			columnName := item.GetField()

			if allowedColumns != nil {
				column, ok := allowedColumns[item.GetField()]
				if !ok {
					return fmt.Errorf("field %s not found, %w", item.GetField(), ErrInvalidFieldName)
				}

				columnName = column
			}

			if item.HasGroupOpen() {
				builder.WriteString(" ")
				builder.WriteString(item.GetChainingKey())
				builder.WriteString(item.GetGroupOpen())
				continue
			}

			if item.HasGroupClose() {
				builder.WriteString(item.GetGroupClose())

				builder.WriteString(" ")
				builder.WriteString(item.GetChainingKey())
				builder.WriteString(" ")
				continue
			}

			operator, err := item.GetParsedOperator()
			if err != nil {
				return fmt.Errorf("error parsing operator %s, %w", item.GetOperator(), err)
			}

			builder.WriteString(columnName)
			builder.WriteString(" ")
			builder.WriteString(operator)

			if operator == In {
				in, inArgs := buildIn(item.GetValue, count)
				builder.WriteString(" ")
				builder.WriteString(in)

				count += len(inArgs)
				config.args = append(config.args, inArgs...)
			} else {
				builder.WriteString(" ")
				builder.WriteString("$")
				builder.WriteString(strconv.Itoa(count + 1))
				count++
			}

			if item.GetChainingKey() != "" && len(collection)-1 > index {
				builder.WriteString(" ")
				builder.WriteString(strings.ToUpper(item.GetChainingKey()))
				builder.WriteString(" ")
			}

			if operator == In {
				continue
			}

			config.args = append(config.args, item.GetValue())
		}

		config.expression = strings.TrimSpace(builder.String())
		config.sqlClause = where

		return nil

	}
}

func isChainingKey(value string) bool {
	return strings.EqualFold(value, "AND") || strings.EqualFold(value, "OR")
}
