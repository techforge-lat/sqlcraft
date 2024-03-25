package sqlcraft

import (
	"bytes"
	"fmt"
	"strings"
)

type UpdateQuery struct {
	query             string
	defaultSQLClauses SQLClauses
	argsCount         uint
	err               error
}

func RawUpdate(sql string, sqlClouseConfigs ...SQLClause) UpdateQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	sqlClouseConfigs = append(sqlClouseConfigs, withExcludeWhereKeyword(hasWhere))

	return UpdateQuery{
		query:             sql,
		defaultSQLClauses: sqlClouseConfigs,
	}
}

// Update creates a base UPDATE sql expression with optional default expressions
// When executing the sql, you must first pass the args for the SET expression
// and then pass the args for the used option expresion in the corresponding order
func Update(tableName string, columns []string, sqlClouseConfigs ...SQLClause) UpdateQuery {
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
		query:             query.String(),
		defaultSQLClauses: sqlClouseConfigs,
		argsCount:         uint(columnsLength),
	}
}

func (u *UpdateQuery) Returning(columns ...string) UpdateQuery {
	u.defaultSQLClauses = append(u.defaultSQLClauses, WithReturning(columns...))
	return *u
}

func (u *UpdateQuery) Where(items ...FilterItem) UpdateQuery {
	u.defaultSQLClauses = append(u.defaultSQLClauses, WithWhere(items...))
	return *u
}

func (u *UpdateQuery) SafeWhere(allowedColumns AllowedColumns, items ...FilterItem) UpdateQuery {
	u.defaultSQLClauses = append(u.defaultSQLClauses, WithSafeWhere(allowedColumns, items...))
	return *u
}

func (u UpdateQuery) sql() string {
	return u.query
}

func (u UpdateQuery) defaultSQLClouseConfigs() SQLClauses {
	return u.defaultSQLClauses
}

func (u UpdateQuery) paramsCount() uint {
	return u.argsCount
}

func (u UpdateQuery) Err() error {
	return u.err
}
