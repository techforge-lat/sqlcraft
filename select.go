package sqlcraft

import (
	"bytes"
	"strings"
)

type SelectQuery struct {
	query      string
	sqlClauses SQLClauses
	argsCount  uint
	err        error
}

func RawSelect(sql string, sqlClauseConfigs ...SQLClause) SelectQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	sqlClauseConfigs = append(sqlClauseConfigs, withExcludeWhereKeyword(hasWhere))

	return SelectQuery{
		query:      sql,
		sqlClauses: sqlClauseConfigs,
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
		query:      query.String(),
		sqlClauses: sqlClauseConfigs,
	}
}

func (s SelectQuery) Where(collection ...FilterItem) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithWhere(collection...))
	return s
}

func (s SelectQuery) SafeWhere(allowedColumns AllowedColumns, collection ...FilterItem) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithSafeWhere(allowedColumns, collection...))
	return s
}

func (s SelectQuery) GroupBy(columns ...string) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithGroupBy(columns...))
	return s
}

func (s SelectQuery) OrderBy(collection ...SortItem) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithOrderBy(collection...))
	return s
}

func (s SelectQuery) SafeOrderBy(allowedColumns AllowedColumns, collection ...SortItem) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithSafeOrderBy(allowedColumns, collection...))
	return s
}

func (s SelectQuery) Limit(limit uint) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithLimit(limit))
	return s
}

func (s SelectQuery) Offset(offset uint) SelectQuery {
	s.sqlClauses = append(s.sqlClauses, WithOffset(offset))
	return s
}

func (s SelectQuery) sql() string {
	return s.query
}

func (s SelectQuery) defaultSQLClauseConfigs() SQLClauses {
	return s.sqlClauses
}

func (s SelectQuery) paramsCount() uint {
	return s.argsCount
}

func (s SelectQuery) Err() error {
	return s.err
}
