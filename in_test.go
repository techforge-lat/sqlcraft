package sqlcraft

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_buildIn(t *testing.T) {
	type args struct {
		value any
		index int
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []any
	}{
		{
			name: "string",
			args: args{
				value: []string{"facebook", "apple", "microsoft"},
				index: 0,
			},
			want:     "($1, $2, $3)",
			wantArgs: []any{"facebook", "apple", "microsoft"},
		},
		{
			name: "uuid",
			args: args{
				value: uuid.UUIDs{uuid.MustParse("4b320e5f-642d-4e64-ba1d-b99fa17ff353"), uuid.MustParse("6ba4bb63-86b2-41c1-90ac-3d6324de8bc0")},
				index: 0,
			},
			want:     "($1, $2)",
			wantArgs: []any{uuid.MustParse("4b320e5f-642d-4e64-ba1d-b99fa17ff353"), uuid.MustParse("6ba4bb63-86b2-41c1-90ac-3d6324de8bc0")},
		},
		{
			name: "int",
			args: args{
				value: []int{1, 2, 3},
				index: 0,
			},
			want:     "($1, $2, $3)",
			wantArgs: []any{1, 2, 3},
		},
		{
			name: "uint",
			args: args{
				value: []uint{1, 2, 3},
				index: 0,
			},
			want:     "($1, $2, $3)",
			wantArgs: []any{uint(1), uint(2), uint(3)},
		},
		{
			name: "float",
			args: args{
				value: []float64{1.1, 2.2, 3.3},
				index: 0,
			},
			want:     "($1, $2, $3)",
			wantArgs: []any{1.1, 2.2, 3.3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildIn(tt.args.value, tt.args.index)
			if got != tt.want {
				t.Errorf("buildIn() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantArgs) {
				t.Errorf("buildIn() got1 = %v, want %v", got1, tt.wantArgs)
			}
		})
	}
}
