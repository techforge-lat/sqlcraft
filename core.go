package sqlcraft

import (
	"errors"
	"strings"
)

var (
	ErrMissingTableName = errors.New("missing table name in query")
	ErrMissingColumns   = errors.New("missing columns in query")
)

type Query interface {
	defaultOpts() Options
	optionKeys() optionKeys
	sql() string
	Err() error
}

type SQLCraft struct {
	Sql  string
	Args []any
}

func Build(query Query, opts ...Option) (SQLCraft, error) {
	if query.Err() != nil {
		return SQLCraft{}, query.Err()
	}

	// TODO: validate that the option can be use for the `query`

	sql := strings.Builder{}
	sql.WriteString(query.sql())

	opts = append(query.defaultOpts(), opts...)

	args := []any{}
	for _, opt := range opts {
		var option options
		if err := opt(&option); err != nil {
			return SQLCraft{}, err
		}
		if option.key.IsInList(query.optionKeys()) {
			return SQLCraft{}, ErrDuplicatedOption
		}

		sql.WriteString(" ")
		sql.WriteString(option.sql)

		args = append(args, option.args...)
	}

	return SQLCraft{
		Sql:  sql.String(),
		Args: args,
	}, nil
}
