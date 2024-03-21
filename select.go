package sqlcraft

import (
	"bytes"
	"strings"
)

// TODO: select with own sql query just to add options to it

type Select struct {
	query          string
	defaultOptions Options
	optionKeyList  optionKeys
	argsCount      uint
	err            error
}

func NewRawSelect(sql string, defualtOpts ...Option) Select {
	return Select{
		query:          sql,
		defaultOptions: defualtOpts,
	}
}

func NewSelect(tableName string, columns []string, defualtOpts ...Option) Select {
	if tableName == "" {
		return Select{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return Select{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(columns, ", "))
	query.WriteString(" FROM ")
	query.WriteString(tableName)

	return Select{
		query:          query.String(),
		defaultOptions: defualtOpts,
	}
}

func (i Select) sql() string {
	return i.query
}

func (i Select) defaultOpts() Options {
	return i.defaultOptions
}

func (i Select) optionKeys() optionKeys {
	return i.optionKeyList
}

func (i Select) paramsCount() uint {
	return i.argsCount
}

func (i Select) Err() error {
	return i.err
}
