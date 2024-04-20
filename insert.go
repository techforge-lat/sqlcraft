package sqlcraft

import (
	"bytes"
	"fmt"
	"strings"
)

type InsertQuery struct {
	query      string
	sqlClauses SQLClauses
	argsCount  uint
	err        error
}

func Insert(tableName string, columns []string, defaultOpts ...SQLClause) InsertQuery {
	if tableName == "" {
		return InsertQuery{
			err: ErrMissingTableName,
		}
	}

	if len(columns) == 0 {
		return InsertQuery{
			err: ErrMissingColumns,
		}
	}

	query := bytes.Buffer{}
	args := bytes.Buffer{}

	query.WriteString("INSERT INTO ")
	query.WriteString(tableName)
	query.WriteString(" (")

	columnsLength := len(columns) - 1
	for i := range columns {
		args.WriteString(fmt.Sprintf("$%d", i+1))

		if i < columnsLength {
			args.WriteString(", ")
		}
	}
	query.WriteString(strings.Join(columns, ", "))
	query.WriteString(")")
	query.WriteString(" VALUES (")
	query.WriteString(args.String())
	query.WriteString(")")

	return InsertQuery{
		query:      query.String(),
		sqlClauses: defaultOpts,
		argsCount:  uint(columnsLength) + 1,
	}
}

func (i InsertQuery) Returning(columns ...string) InsertQuery {
	i.sqlClauses = append(i.sqlClauses, WithReturning(columns...))
	return i
}

func (i InsertQuery) sql() string {
	return i.query
}

func (i InsertQuery) defaultSQLClauseConfigs() SQLClauses {
	return i.sqlClauses
}

func (i InsertQuery) paramsCount() uint {
	return i.argsCount
}

func (i InsertQuery) Err() error {
	return i.err
}
