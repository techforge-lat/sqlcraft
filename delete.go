package sqlcraft

import (
	"strings"
)

type DeleteQuery struct {
	query             string
	defaultSQLClauses SQLClauses
	err               error
}

func Delete(tableName string, defaultOpts ...SQLClause) DeleteQuery {
	if tableName == "" {
		return DeleteQuery{
			err: ErrMissingTableName,
		}
	}

	query := strings.Builder{}

	query.WriteString("DELETE FROM ")
	query.WriteString(tableName)

	return DeleteQuery{
		query:             query.String(),
		defaultSQLClauses: defaultOpts,
	}
}

func RawDelete(sql string, defaultOpts ...SQLClause) DeleteQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	defaultOpts = append(defaultOpts, withExcludeWhereKeyword(hasWhere))

	return DeleteQuery{
		query:             sql,
		defaultSQLClauses: defaultOpts,
	}
}

func (d DeleteQuery) Returning(columns ...string) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithReturning(columns...))
	return d
}

func (d DeleteQuery) Where(collection ...FilterItem) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithWhere(collection...))
	return d
}

func (d DeleteQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithSafeWhere(allowedColumns, collection...))
	return d
}

func (d DeleteQuery) sql() string {
	return d.query
}

func (d DeleteQuery) defaultSQLClauseConfigs() SQLClauses {
	return d.defaultSQLClauses
}

func (d DeleteQuery) paramsCount() uint {
	return 0
}

func (d DeleteQuery) Err() error {
	return d.err
}
