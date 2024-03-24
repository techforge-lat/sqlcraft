package sqlcraft

import (
	"errors"
	"testing"
)

func TestNewInsert(t *testing.T) {
	type args struct {
		tableName string
		columns   []string
	}
	tests := []struct {
		name    string
		args    args
		want    InsertQuery
		wantErr bool
	}{
		{
			name: "multiple columns",
			args: args{
				tableName: "users",
				columns:   []string{"id", "name", "email", "password"},
			},
			want: InsertQuery{
				query: "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
			},
		},
		{
			name: "one column",
			args: args{
				tableName: "users",
				columns:   []string{"id"},
			},

			want: InsertQuery{
				query: "INSERT INTO users (id) VALUES ($1)",
			},
		},
		{
			name: "missing columns",
			args: args{
				tableName: "users",
				columns:   nil,
			},
			want: InsertQuery{
				err: ErrMissingColumns,
			},
			wantErr: true,
		},
		{
			name: "missing table name",
			args: args{
				tableName: "",
				columns:   nil,
			},
			want: InsertQuery{
				err: ErrMissingTableName,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Insert(tt.args.tableName, tt.args.columns)
			if got.query != tt.want.query {
				t.Errorf("NewInsert() = %v, want %v", got.query, tt.want.query)
			}

			if tt.wantErr && !errors.Is(got.err, tt.want.err) {
				t.Errorf("NewInsert() = %v, want %v", got.err, tt.want.err)
			}
		})
	}
}
