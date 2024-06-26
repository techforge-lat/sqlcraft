package sqlcraft

import (
	"errors"
	"testing"
)

func TestNewDelete(t *testing.T) {
	type args struct {
		tableName   string
		defualtOpts []SQLClause
	}
	tests := []struct {
		name    string
		args    args
		want    DeleteQuery
		wantErr bool
	}{
		{
			name: "withour where",
			args: args{
				tableName: "users",
			},
			want: DeleteQuery{
				query: "DELETE FROM users",
			},
			wantErr: false,
		},
		{
			name: "with ONLY",
			args: args{
				tableName: "ONLY users",
			},
			want: DeleteQuery{
				query: "DELETE FROM ONLY users",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Delete(tt.args.tableName, tt.args.defualtOpts...)
			if got.query != tt.want.query {
				t.Errorf("NewInsert() = %v, want %v", got.query, tt.want.query)
			}

			if tt.wantErr && !errors.Is(got.err, tt.want.err) {
				t.Errorf("NewInsert() = %v, want %v", got.err, tt.want.err)
			}
		})
	}
}
