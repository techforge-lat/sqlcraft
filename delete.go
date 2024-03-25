package sqlcraft

import (
	"strings"
)

type DeleteQuery struct {
	query             string
	defaultSQLClauses SQLClauses
	err               error
}

func Delete(tableName string, defualtOpts ...SQLClause) DeleteQuery {
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
		defaultSQLClauses: defualtOpts,
	}
}

func RawDelete(sql string, defualtOpts ...SQLClause) DeleteQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	defualtOpts = append(defualtOpts, withExcludeWhereKeyword(hasWhere))

	return DeleteQuery{
		query:             sql,
		defaultSQLClauses: defualtOpts,
	}
}

func (d *DeleteQuery) Returning(columns ...string) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithReturning(columns...))
	return *d
}

func (d *DeleteQuery) Where(collection ...FilterItem) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithWhere(collection...))
	return *d
}

func (d *DeleteQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) DeleteQuery {
	d.defaultSQLClauses = append(d.defaultSQLClauses, WithSafeWhere(allowedColumns, collection...))
	return *d
}

func (d DeleteQuery) sql() string {
	return d.query
}

func (d DeleteQuery) defaultSQLClouseConfigs() SQLClauses {
	return d.defaultSQLClauses
}

func (d DeleteQuery) paramsCount() uint {
	return 0
}

func (d DeleteQuery) Err() error {
	return d.err
}
