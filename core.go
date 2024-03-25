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
	Err() error
	defaultSQLClouseConfigs() SQLClauses
	paramsCount() uint
	sql() string
}

type SQLQuery struct {
	Sql  string
	Args []any
}

func Build(query Query, opts ...SQLClause) (SQLQuery, error) {
	if query.Err() != nil {
		return SQLQuery{}, query.Err()
	}

	// TODO: validate that the option can be use for the `query`
	// TODO: don't allow option duplicates

	sql := strings.Builder{}
	sql.WriteString(query.sql())

	opts = append(query.defaultSQLClouseConfigs(), opts...)

	args := []any{}
	excludeWhereKeyword := false
	for _, opt := range opts {
		var option sqlClauseConfig
		option.excludeWhereKeyword = excludeWhereKeyword
		option.paramCountStartFrom = uint(len(args)) + query.paramsCount()

		if err := opt(&option); err != nil {
			return SQLQuery{}, err
		}

		// in case the WHERE option was used as a default
		// and then the client still sends extra filters
		if option.sqlClause == where {
			excludeWhereKeyword = true
		}

		sql.WriteString(" ")
		sql.WriteString(option.expression)

		args = append(args, option.args...)
	}

	return SQLQuery{
		Sql:  sql.String(),
		Args: args,
	}, nil
}
