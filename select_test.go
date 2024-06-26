package sqlcraft

import (
	"reflect"
	"testing"
)

func TestNewSelect(t *testing.T) {
	type args struct {
		tableName   string
		columns     []string
		defualtOpts []SQLClause
	}
	tests := []struct {
		name string
		args args
		want SelectQuery
	}{
		{
			name: "select without where",
			args: args{
				tableName: "users",
				columns:   []string{"id", "name", "email", "password"},
			},
			want: SelectQuery{
				query: "SELECT id, name, email, password FROM users",
			},
		},
		{
			name: "select one column",
			args: args{
				tableName: "users",
				columns:   []string{"id"},
			},
			want: SelectQuery{
				query: "SELECT id FROM users",
			},
		},
		{
			name: "missing table",
			args: args{
				tableName: "",
				columns:   []string{"id"},
			},
			want: SelectQuery{
				query: "",
				err:   ErrMissingTableName,
			},
		},
		{
			name: "missing columns",
			args: args{
				tableName: "users",
				columns:   []string{},
			},
			want: SelectQuery{
				query: "",
				err:   ErrMissingColumns,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Select(tt.args.tableName, tt.args.columns, tt.args.defualtOpts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}
