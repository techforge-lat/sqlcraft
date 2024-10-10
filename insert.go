package sqlcraft

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrEmptyValues     = errors.New("empty values in insert")
	ErrMissMatchValues = errors.New("miss match values for given columns")
)

type Insert struct {
	table            string
	columns          []string
	returningColumns []string
	values           []any
}

func InsertInto(tableName string) Insert {
	return Insert{
		table:  tableName,
		values: make([]any, 0),
	}
}

func (i Insert) Columns(columns []string) Insert {
	i.columns = columns

	return i
}

func (i Insert) Values(args ...any) Insert {
	i.values = append(i.values, args)

	return i
}

func (i Insert) Returning(columns []string) Insert {
	i.returningColumns = columns

	return i
}

func (i Insert) ToSql() (Result, error) {
	if len(i.values) == 0 {
		return Result{}, ErrEmptyValues
	}

	if len(i.values)%len(i.columns) != 0 {
		return Result{}, ErrMissMatchValues
	}

	builder := strings.Builder{}

	builder.WriteString("INSERT INTO ")
	builder.WriteString(i.table)
	builder.WriteString(" (")
	builder.WriteString(strings.Join(i.columns, ", "))
	builder.WriteString(") ")

	builder.WriteString("VALUES ")

	valueRowCount := 0
	for index := range i.values {
		valueRowCount += 1

		if valueRowCount == 1 && index > 0 {
			builder.WriteString(", ")
		}

		if valueRowCount == 1 {
			builder.WriteString("(")
		}

		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(index + 1))

		if valueRowCount == len(i.columns) {
			builder.WriteString(")")
			valueRowCount = 0
			continue
		}

		builder.WriteString(", ")
	}

	if len(i.returningColumns) > 0 {
		builder.WriteString(" RETURNING ")
		builder.WriteString(strings.Join(i.returningColumns, ", "))
	}

	return Result{
		Sql:  builder.String(),
		Args: i.values,
	}, nil
}
