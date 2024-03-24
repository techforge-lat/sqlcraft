package sqlcraft

import (
	"bytes"
	"strings"
)

type SelectQuery struct {
	query          string
	defaultOptions SQLClauses
	argsCount      uint
	err            error
}

func RawSelect(sql string, defualtOpts ...SQLClause) SelectQuery {
	hasWhere := strings.Contains(strings.ToUpper(sql), strings.ToUpper(string(where)))
	defualtOpts = append(defualtOpts, withExcludeWhereKeyword(hasWhere))

	return SelectQuery{
		query:          sql,
		defaultOptions: defualtOpts,
	}
}

func Select(tableName string, columns []string, defualtOpts ...SQLClause) SelectQuery {
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
		query:          query.String(),
		defaultOptions: defualtOpts,
	}
}

func (i SelectQuery) sql() string {
	return i.query
}

func (i SelectQuery) defaultOpts() SQLClauses {
	return i.defaultOptions
}

func (i SelectQuery) paramsCount() uint {
	return i.argsCount
}

func (i SelectQuery) Err() error {
	return i.err
}
