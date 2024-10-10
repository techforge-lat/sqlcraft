package sqlcraft

import (
	"reflect"
	"testing"
)

func TestInsert_ToSql(t *testing.T) {
	type fields struct {
		table            string
		columns          []string
		returningColumns []string
		values           []any
	}
	tests := []struct {
		name    string
		fields  fields
		want    Result
		wantErr bool
	}{
		{
			name: "standard insert",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{},
				values:           []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			want: Result{
				Sql:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			wantErr: false,
		},
		{
			name: "standard insert with returning",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{"id", "created_at"},
				values:           []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			want: Result{
				Sql:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			wantErr: false,
		},
		{
			name: "standard insert with returning and multiple row values",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{"id", "created_at"},
				values:           []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash", "Brownie", nil, "brownie@gmail.com", "secrethash"},
			},
			want: Result{
				Sql:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) RETURNING id, created_at",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash", "Brownie", nil, "brownie@gmail.com", "secrethash"},
			},
			wantErr: false,
		},
		{
			name: "error empty values",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{"id", "created_at"},
				values:           []any{},
			},
			want:    Result{},
			wantErr: true,
		},
		{
			name: "missmatch values",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{"id", "created_at"},
				values:           []any{"hernan"},
			},
			want:    Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Insert{
				table:            tt.fields.table,
				columns:          tt.fields.columns,
				returningColumns: tt.fields.returningColumns,
				values:           tt.fields.values,
			}
			got, err := i.ToSql()
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert.ToSql() error = \n%v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert.ToSql() = %v, want %v", got, tt.want)
			}
		})
	}
}
