package sqlcraft

import (
	"bytes"
	"strings"
)

type DeleteQuery struct {
	query          string
	defaultOptions SQLClauses
	err            error
}

func Delete(tableName string, defualtOpts ...SQLClause) DeleteQuery {
	if tableName == "" {
		return DeleteQuery{
			err: ErrMissingTableName,
		}
	}

	query := bytes.Buffer{}

	query.WriteString("DELETE FROM ")
	query.WriteString(tableName)

	return DeleteQuery{
		query:          query.String(),
		defaultOptions: defualtOpts,
	}
}

func RawDelete(sql string, defualtOpts ...SQLClause) DeleteQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	defualtOpts = append(defualtOpts, withExcludeWhereKeyword(hasWhere))

	return DeleteQuery{
		query:          sql,
		defaultOptions: defualtOpts,
	}
}

func (d *DeleteQuery) Returning(columns ...string) DeleteQuery {
	d.defaultOptions = append(d.defaultOptions, WithReturning(columns...))
	return *d
}

func (d *DeleteQuery) Where(collection ...FilterItem) DeleteQuery {
	d.defaultOptions = append(d.defaultOptions, WithWhere(collection...))
	return *d
}

func (d *DeleteQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) DeleteQuery {
	d.defaultOptions = append(d.defaultOptions, WithSafeWhere(allowedColumns, collection...))
	return *d
}

func (d DeleteQuery) sql() string {
	return d.query
}

func (d DeleteQuery) defaultOpts() SQLClauses {
	return d.defaultOptions
}

func (d DeleteQuery) paramsCount() uint {
	return 0
}

func (d DeleteQuery) Err() error {
	return d.err
}
