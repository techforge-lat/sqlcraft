package sqlcraft

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/techforge-lat/dafi/v2"
)

var ErrInvalidOperator = errors.New("invalid dafi operator")

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
	dafi.IsNot:          "IS NOT",
	dafi.In:             "IN",
	dafi.NotIn:          "NOT IN",
}

func Where(filters dafi.Filters) (Result, error) {
	if filters.IsZero() {
		return Result{}, nil
	}

	builder := strings.Builder{}
	args := []any{}

	builder.WriteString("WHERE ")
	for i, filter := range filters {
		if filter.IsGroupOpen {
			if filter.GroupOpenQty == 0 {
				filter.GroupOpenQty = 1
			}

			builder.WriteString(strings.Repeat("(", filter.GroupOpenQty))
		}

		operator, ok := psqlOperatorByDafiOperator[filter.Operator]
		if !ok {
			return Result{}, errors.Join(fmt.Errorf("operator %q not found", filter.Operator), ErrInvalidOperator)
		}

		if filter.Operator == dafi.In || filter.Operator == dafi.NotIn {
			inResult := In(filter.Value, len(args)+1)
			if inResult.Sql == "" {
				continue
			}

			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)

			builder.WriteString(" ")
			builder.WriteString(inResult.Sql)

			args = append(args, inResult.Args...)
		} else {
			builder.WriteString(string(filter.Field))
			builder.WriteString(" ")
			builder.WriteString(operator)
			builder.WriteString(" $")
			builder.WriteString(strconv.Itoa(i + 1))

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

		builder.WriteString(" ")
		builder.WriteString(string(filter.ChainingKey))
		builder.WriteString(" ")
	}

	return Result{
		Sql:  strings.TrimSpace(builder.String()),
		Args: args,
	}, nil
}
