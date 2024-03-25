package sqlcraft

import (
	"bytes"
	"strings"
)

type SelectQuery struct {
	query             string
	defaultSQLClouses SQLClauses
	argsCount         uint
	err               error
}

func RawSelect(sql string, sqlClouseConfigs ...SQLClause) SelectQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	sqlClouseConfigs = append(sqlClouseConfigs, withExcludeWhereKeyword(hasWhere))

	return SelectQuery{
		query:             sql,
		defaultSQLClouses: sqlClouseConfigs,
	}
}

func Select(tableName string, columns []string, sqlClouseConfigs ...SQLClause) SelectQuery {
	if tableName == "" {
		return SelectQuery{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return SelectQuery{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(columns, ", "))
	query.WriteString(" FROM ")
	query.WriteString(tableName)

	return SelectQuery{
		query:             query.String(),
		defaultSQLClouses: sqlClouseConfigs,
	}
}

func (d SelectQuery) Where(collection ...FilterItem) SelectQuery {
	d.defaultSQLClouses = append(d.defaultSQLClouses, WithWhere(collection...))
	return d
}

func (d SelectQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) SelectQuery {
	d.defaultSQLClouses = append(d.defaultSQLClouses, WithSafeWhere(allowedColumns, collection...))
	return d
}

func (d SelectQuery) GroupBy(columns ...string) SelectQuery {
	d.defaultSQLClouses = append(d.defaultSQLClouses, WithGroupBy(columns...))
	return d
}

func (i SelectQuery) sql() string {
	return i.query
}

func (i SelectQuery) defaultSQLClouseConfigs() SQLClauses {
	return i.defaultSQLClouses
}

func (i SelectQuery) paramsCount() uint {
	return i.argsCount
}

func (i SelectQuery) Err() error {
	return i.err
}
