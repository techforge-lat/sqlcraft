package sqlcraft

import (
	"bytes"
	"fmt"
	"strings"
)

type Insert struct {
	query      string
	optionKeys optionKeys
	err        error
}

func NewInsert(tableName string, columns []string) Insert {
	if tableName == "" {
		return Insert{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return Insert{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}
	args := bytes.Buffer{}

	query.WriteString("INSERT INTO ")
	query.WriteString(tableName)
	query.WriteString(" (")

	columnsLength := len(columns) - 1
	for i := range columns {
		args.WriteString(fmt.Sprintf("$%d", i+1))

		if i < columnsLength {
			args.WriteString(", ")
		}
	}

	query.WriteString(strings.Join(columns, ", "))
	query.WriteString(")")
	query.WriteString(" VALUES (")
	query.WriteString(args.String())
	query.WriteString(")")

	return Insert{
		query: query.String(),
	}
}

func (i Insert) sql() string {
	return i.query
}

func WithReturning(columns ...string) Option {
	return func(option *options) error {
		option.sql = fmt.Sprintf("RETURNING %s", strings.Join(columns, ", "))
		option.key = returning

		return nil
	}
}
