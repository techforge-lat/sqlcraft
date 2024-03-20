package sqlcraft

import (
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	type args struct {
		query Query
		opts  []Option
	}
	tests := []struct {
		name    string
		args    args
		want    SQLCraft
		wantErr bool
	}{
		{
			name: "insert with returning option",
			args: args{
				query: NewInsert("users", []string{"id", "name", "email", "password"}),
				opts:  []Option{WithReturning("id", "created_at")},
			},
			want: SQLCraft{
				Sql:  "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				Args: []any{},
			},
			wantErr: false,
		},
		{
			name: "insert with default options",
			args: args{
				query: NewInsert("users", []string{"id", "name", "email", "password"}, WithReturning("id", "created_at")),
			},
			want: SQLCraft{
				Sql:  "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
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