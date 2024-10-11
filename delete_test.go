package sqlcraft

import (
	"reflect"
	"testing"
)

func TestDeleteQuery_ToSQL(t *testing.T) {
	tests := []struct {
		name    string
		query   DeleteQuery
		want    Result
		wantErr bool
	}{
		{
			name:  "simple delete",
			query: DeleteFrom("users"),
			want: Result{
				Sql: "DELETE FROM users",
			},
			wantErr: false,
		},
		{
			name:  "delete with returning",
			query: DeleteFrom("users").Returning("id"),
			want: Result{
				Sql: "DELETE FROM users RETURNING id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.ToSQL()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteQuery.ToSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteQuery.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
