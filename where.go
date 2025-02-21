package sqlcraft

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/techforge-lat/dafi/v2"
)

var (
	ErrInvalidOperator  = errors.New("invalid dafi operator")
	ErrInvalidFieldName = errors.New("invalid field name")
)

var psqlOperatorByDafiOperator = map[dafi.FilterOperator]string{
	dafi.Equal:          "=",
	dafi.NotEqual:       "<>",
	dafi.Greater:        ">",
	dafi.GreaterOrEqual: ">=",
	dafi.Less:           "<",
	dafi.LessOrEqual:    "<=",
	dafi.Contains:       "ILIKE",
	dafi.NotContains:    "NOT ILIKE",
	dafi.Is:             "IS",
	dafi.IsNull:         "IS NULL",
	dafi.IsNot:          "IS NOT",
	dafi.IsNotNull:      "IS NOT NULL",
	dafi.In:             "IN",
	dafi.NotIn:          "NOT IN",
	dafi.Default:        "",
}

// WhereSafe maps domain field names to sql column names,
// if a filter with an unknow domain field name is found it will return an error
func WhereSafe(initialArgCount int, sqlColumnByDomainField map[string]string, filters ...dafi.Filter) (Result, error) {
	if len(sqlColumnByDomainField) > 0 {
		for i, filter := range filters {
			sqlColumnName, ok := sqlColumnByDomainField[string(filter.Field)]
			if !ok {
				return Result{}, ErrInvalidFieldName
			}

			filters[i].Field = dafi.FilterField(sqlColumnName)
		}
	}

	return Where(initialArgCount, filters...)
}

// Where returns a WHERE sql sentence and if an invalid operator is found, it will return an error
func Where(initialArgCount int, filters ...dafi.Filter) (Result, error) {
	if len(filters) == 0 {
		return Result{}, nil
	}

	builder := strings.Builder{}
	args := []any{}

	builder.WriteString(" WHERE ")
	for i, filter := range filters {
		if filter.IsGroupOpen {
			if filter.GroupOpenQty == 0 {
				filter.GroupOpenQty = 1
			}

			builder.WriteString(strings.Repeat("(", filter.GroupOpenQty))
		}

		if filter.Operator == "" {
			filter.Operator = dafi.Equal
		}

		operator, ok := psqlOperatorByDafiOperator[filter.Operator]
		if !ok {
			return Result{}, errors.Join(fmt.Errorf("operator %q not found", filter.Operator), ErrInvalidOperator)
		}

		switch filter.Operator {
		case dafi.In, dafi.NotIn:
			inResult := In(filter.Value, len(args)+1+initialArgCount)
			if inResult.Sql == "" {
				continue
			}

			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)

			builder.WriteString(" ")
			builder.WriteString(inResult.Sql)

			args = append(args, inResult.Args...)
		case dafi.Is, dafi.IsNot, dafi.IsNull, dafi.IsNotNull:
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)
		case dafi.Default:
			builder.WriteString(string(filter.Field))
		default:

			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)
			builder.WriteString(" $")
			builder.WriteString(strconv.Itoa(i + 1 + initialArgCount))

			args = append(args, filter.Value)
		}

		if i < len(filters)-1 && filter.ChainingKey == "" {
			filter.ChainingKey = dafi.And
		}

		if filter.IsGroupClose {
			if filter.GroupCloseQty == 0 {
				filter.GroupCloseQty = 1
			}

			builder.WriteString(strings.Repeat(")", filter.GroupCloseQty))
		}

		if i < len(filters)-1 {
			builder.WriteString(" ")
			builder.WriteString(string(filter.ChainingKey))
			builder.WriteString(" ")
		}
	}

	return Result{
		Sql:  strings.TrimRight(builder.String(), " "),
		Args: args,
	}, nil
}
