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
	defaultOpts() []Option
	optionKeys() optionKeys
	sql() string
}

type SQLCraft struct {
	Sql  string
	Args []any
}

func Build(query Query, opts ...Option) (SQLCraft, error) {
	sql := strings.Builder{}
	sql.WriteString(query.sql())

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
