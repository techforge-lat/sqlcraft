package sqlcraft

import (
	"bytes"
	"fmt"
)

type Update struct {
	query          string
	defaultOptions Options
	optionKeyList  optionKeys
	argsCount      uint
	err            error
}

// NewUpdate creates a base UPDATE sql expression with optional default expressions
// When executing the sql, you must first pass the args for the SET expression
// and then pass the args for the used option expresion in the corresponding order
func NewUpdate(tableName string, columns []string, defualtOpts ...Option) Update {
	if tableName == "" {
		return Update{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return Update{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}

	query.WriteString("UPDATE ")
	query.WriteString(tableName)
	query.WriteString(" SET ")

	columnsLength := len(columns) - 1
	for i, column := range columns {
		query.WriteString(column)
		query.WriteString(fmt.Sprintf(" = $%d", i+1))

		if i < columnsLength {
			query.WriteString(", ")
		}
	}

	return Update{
		query:          query.String(),
		defaultOptions: defualtOpts,
		argsCount:      uint(columnsLength),
	}
}

func (u Update) sql() string {
	return u.query
}

func (u Update) defaultOpts() Options {
	return u.defaultOptions
}

func (u Update) optionKeys() optionKeys {
	return u.optionKeyList
}

func (u Update) paramsCount() uint {
	return u.argsCount
}

func (u Update) Err() error {
	return u.err
}
