package sqlcraft

import (
	"bytes"
)

type Delete struct {
	query          string
	defaultOptions Options
	optionKeyList  optionKeys
	err            error
}

func NewDelete(tableName string, defualtOpts ...Option) Delete {
	if tableName == "" {
		return Delete{
			err: ErrMissingTableName,
		}
	}

	query := bytes.Buffer{}

	query.WriteString("DELETE FROM ")
	query.WriteString(tableName)

	return Delete{
		query:          query.String(),
		defaultOptions: defualtOpts,
	}
}

func (d Delete) sql() string {
	return d.query
}

func (d Delete) defaultOpts() Options {
	return d.defaultOptions
}

func (d Delete) optionKeys() optionKeys {
	return d.optionKeyList
}

func (d Delete) Err() error {
	return d.err
}
