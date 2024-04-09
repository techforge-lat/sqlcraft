package sqlcraft

import (
	"bytes"
	"strings"
)

type SelectQuery struct {
	query             string
	defaultSQLClauses SQLClauses
	argsCount         uint
	err               error
}

func RawSelect(sql string, sqlClauseConfigs ...SQLClause) SelectQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	sqlClauseConfigs = append(sqlClauseConfigs, withExcludeWhereKeyword(hasWhere))

	return SelectQuery{
		query:             sql,
		defaultSQLClauses: sqlClauseConfigs,
	}
}

func Select(tableName string, columns []string, sqlClauseConfigs ...SQLClause) SelectQuery {
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
		defaultSQLClauses: sqlClauseConfigs,
	}
}

func (s SelectQuery) Where(collection ...FilterItem) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithWhere(collection...))
	return s
}

func (s SelectQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithSafeWhere(allowedColumns, collection...))
	return s
}

func (s SelectQuery) GroupBy(columns ...string) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithGroupBy(columns...))
	return s
}

func (s SelectQuery) OrderBy(collection ...SortItem) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithOrderBy(collection...))
	return s
}

func (s SelectQuery) SafeOrderBy(allowedColumns AllowedColumns, collection ...SortItem) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithSafeOrderBy(allowedColumns, collection...))
	return s
}

func (s SelectQuery) Limit(limit uint) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithLimit(limit))
	return s
}

func (s SelectQuery) Offset(offset uint) SelectQuery {
	s.defaultSQLClauses = append(s.defaultSQLClauses, WithOffset(offset))
	return s
}

func (s SelectQuery) sql() string {
	return s.query
}

func (s SelectQuery) defaultSQLClauseConfigs() SQLClauses {
	return s.defaultSQLClauses
}

func (s SelectQuery) paramsCount() uint {
	return s.argsCount
}

func (s SelectQuery) Err() error {
	return s.err
}
