package sqlcraft

import (
	"errors"
	"testing"
)

func TestNewUpdate(t *testing.T) {
	type args struct {
		tableName   string
		columns     []string
		defualtOpts []SQLClause
	}
	tests := []struct {
		name    string
		args    args
		want    UpdateQuery
		wantErr bool
	}{
		{
			name: "multiple columns",
			args: args{
				tableName: "users",
				columns:   []string{"id", "name", "email", "password"},
			},
			want: UpdateQuery{
				query: "UPDATE users SET id = $1, name = $2, email = $3, password = $4",
			},
			wantErr: false,
		},
		{
			name: "one column",
			args: args{
				tableName: "users",
				columns:   []string{"id"},
			},
			want: UpdateQuery{
				query: "UPDATE users SET id = $1",
			},
			wantErr: false,
		},
		{
			name: "missing columns",
			args: args{
				tableName: "users",
				columns:   nil,
			},
			want: UpdateQuery{
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
			want: UpdateQuery{
				err: ErrMissingTableName,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Update(tt.args.tableName, tt.args.columns, tt.args.defualtOpts...)
			if got.query != tt.want.query {
				t.Errorf("NewInsert() = %v, want %v", got.query, tt.want.query)
			}

			if tt.wantErr && !errors.Is(got.err, tt.want.err) {
				t.Errorf("NewInsert() = %v, want %v", got.err, tt.want.err)
			}
		})
	}
}
