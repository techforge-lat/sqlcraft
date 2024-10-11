package sqlcraft

import (
	"strings"
)

type DeleteQuery struct {
	table           string
	returningValues []string
}

func DeleteFrom(table string) DeleteQuery {
	return DeleteQuery{
		table:           table,
		returningValues: []string{},
	}
}

func (u DeleteQuery) Returning(columns ...string) DeleteQuery {
	u.returningValues = columns

	return u
}

func (u DeleteQuery) ToSQL() (Result, error) {
	builder := strings.Builder{}

	builder.WriteString("DELETE FROM ")
	builder.WriteString(u.table)

	if len(u.returningValues) > 0 {
		builder.WriteString(" RETURNING ")
		builder.WriteString(strings.Join(u.returningValues, ", "))
	}

	return Result{
		Sql: builder.String(),
	}, nil
}
