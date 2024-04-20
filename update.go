package sqlcraft

import (
	"bytes"
	"fmt"
	"strings"
)

type UpdateQuery struct {
	query      string
	sqlClauses SQLClauses
	argsCount  uint
	err        error
}

func RawUpdate(sql string, sqlClauseConfigs ...SQLClause) UpdateQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	sqlClauseConfigs = append(sqlClauseConfigs, withExcludeWhereKeyword(hasWhere))

	return UpdateQuery{
		query:      sql,
		sqlClauses: sqlClauseConfigs,
	}
}

// Update creates a base UPDATE sql expression with optional default expressions
// When executing the sql, you must first pass the args for the SET expression
// and then pass the args for the used option expression in the corresponding order
func Update(tableName string, columns []string, sqlClauseConfigs ...SQLClause) UpdateQuery {
	if tableName == "" {
		return UpdateQuery{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return UpdateQuery{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}

	query.WriteString("UPDATE ")
	query.WriteString(tableName)
	query.WriteString(" SET ")

	columnsLength := len(columns) - 1
	for i, column := range columns {
		query.WriteString(column)
		query.WriteString(fmt.Sprintf(" = $%d", i+1))

		if i < columnsLength {
			query.WriteString(", ")
		}
	}

	return UpdateQuery{
		query:      query.String(),
		sqlClauses: sqlClauseConfigs,
		argsCount:  uint(columnsLength) + 1,
	}
}

func (u UpdateQuery) Returning(columns ...string) UpdateQuery {
	u.sqlClauses = append(u.sqlClauses, WithReturning(columns...))
	return u
}

func (u UpdateQuery) Where(items ...FilterItem) UpdateQuery {
	u.sqlClauses = append(u.sqlClauses, WithWhere(items...))
	return u
}

func (u UpdateQuery) SafeWhere(allowedColumns AllowedColumns, items ...FilterItem) UpdateQuery {
	u.sqlClauses = append(u.sqlClauses, WithSafeWhere(allowedColumns, items...))
	return u
}

func (u UpdateQuery) sql() string {
	return u.query
}

func (u UpdateQuery) defaultSQLClauseConfigs() SQLClauses {
	return u.sqlClauses
}

func (u UpdateQuery) paramsCount() uint {
	return u.argsCount
}

func (u UpdateQuery) Err() error {
	return u.err
}
