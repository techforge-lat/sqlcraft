package sqlcraft

import (
	"strings"
)

type DeleteQuery struct {
	query      string
	sqlClauses SQLClauses
	err        error
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
		query:      query.String(),
		sqlClauses: defaultOpts,
	}
}

func RawDelete(sql string, defaultOpts ...SQLClause) DeleteQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	defaultOpts = append(defaultOpts, withExcludeWhereKeyword(hasWhere))

	return DeleteQuery{
		query:      sql,
		sqlClauses: defaultOpts,
	}
}

func (d DeleteQuery) Returning(columns ...string) DeleteQuery {
	d.sqlClauses = append(d.sqlClauses, WithReturning(columns...))
	return d
}

func (d DeleteQuery) Where(collection ...FilterItem) DeleteQuery {
	d.sqlClauses = append(d.sqlClauses, WithWhere(collection...))
	return d
}

func (d DeleteQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) DeleteQuery {
	d.sqlClauses = append(d.sqlClauses, WithSafeWhere(allowedColumns, collection...))
	return d
}

func (d DeleteQuery) sql() string {
	return d.query
}

func (d DeleteQuery) defaultSQLClauseConfigs() SQLClauses {
	return d.sqlClauses
}

func (d DeleteQuery) paramsCount() uint {
	return 0
}

func (d DeleteQuery) Err() error {
	return d.err
}
