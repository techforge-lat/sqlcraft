package sqlcraft

import (
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	type args struct {
		query Query
		opts  []SQLClause
	}
	tests := []struct {
		name    string
		args    args
		want    SQLQuery
		wantErr bool
	}{
		{
			name: "insert with returning option",
			args: args{
				query: Insert("users", []string{"id", "name", "email", "password"}),
				opts:  []SQLClause{WithReturning("id", "created_at")},
			},
			want: SQLQuery{
				Sql:  "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "insert with default options",
			args: args{
				query: Insert("users", []string{"id", "name", "email", "password"}, WithReturning("id", "created_at")),
			},
			want: SQLQuery{
				Sql:  "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "update with default returnind option",
			args: args{
				query: Update("users", []string{"id", "name", "email", "password"}, WithReturning("id", "created_at")),
			},
			want: SQLQuery{
				Sql:  "UPDATE users SET id = $1, name = $2, email = $3, password = $4 RETURNING id, created_at",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "update with default from option",
			args: args{
				query: Update("users", []string{"id", "name", "email", "password"}, WithFrom("roles")),
			},
			want: SQLQuery{
				Sql:  "UPDATE users SET id = $1, name = $2, email = $3, password = $4 FROM roles",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "update with default from option with table alias",
			args: args{
				query: Update("users", []string{"id", "name", "email", "password"}, WithFrom("roles r")),
			},
			want: SQLQuery{
				Sql:  "UPDATE users SET id = $1, name = $2, email = $3, password = $4 FROM roles r",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "delete with RETURNING",
			args: args{
				query: Delete("users", WithReturning("id", "created_at")),
			},
			want: SQLQuery{
				Sql:  "DELETE FROM users RETURNING id, created_at",
				Args: []any{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Build(tt.args.query, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
