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

func WithWhere(allowedColumns AllowedColumns, items ...FilterItem) Option {
	return func(option *options) error {
		if len(items) == 0 {
			return errors.New("filter items cannot be empty in WHERE option")
		}

		builder := bytes.Buffer{}
		if !option.excludeWhereKeyword {
			builder.WriteString(" WHERE ")
		}

		args := []any{}

		count := int(option.paramCountStartFrom)
		for index, item := range items {
			columnName, ok := allowedColumns[item.GetField()]
			if !ok {
				return fmt.Errorf("field %s not found, %w", item.GetField(), ErrInvalidFieldName)
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

			op, err := item.GetParsedOperator()
			if err != nil {
				return fmt.Errorf("error parsing operator %s, %w", item.GetOperator(), err)
			}

			builder.WriteString(columnName)
			builder.WriteString(" ")
			builder.WriteString(op)

			if op == In {
				in, inArgs := buildIn(item.GetValue, count)
				builder.WriteString(" ")
				builder.WriteString(in)

				count += len(inArgs)
				args = append(args, inArgs...)
			} else {
				builder.WriteString(" ")
				builder.WriteString("$")
				builder.WriteString(strconv.Itoa(count + 1))
				count++
			}

			if item.GetChainingKey() != "" && len(items)-1 > index {
				builder.WriteString(" ")
				builder.WriteString(strings.ToUpper(item.GetChainingKey()))
				builder.WriteString(" ")
			}

			if op == In {
				continue
			}

			args = append(args, item.GetValue)
		}

		option.sql = strings.TrimSpace(builder.String())
		option.args = args
		option.key = where

		return nil

	}
}

func isChainingKey(value string) bool {
	return strings.EqualFold(value, "AND") || strings.EqualFold(value, "OR")
}
