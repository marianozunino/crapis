package command

import (
	"reflect"
	"testing"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

func Test_setex(t *testing.T) {
	type args struct {
		args []resp.Value
	}
	tests := []struct {
		name string
		args args
		want resp.Value
	}{
		{
			name: "Setex No Args",
			args: args{args: []resp.Value{}},
			want: resp.NewError("wrong number of arguments for 'setex' command"),
		},
		{
			name: "Setex Invalid TTL",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("value is not an integer or out of range"),
		},
		{
			name: "Invalid Amount of Args (2)",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("wrong number of arguments for 'setex' command"),
		},
		{
			name: "Invalid Amount of Args (4)",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("2")),
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("wrong number of arguments for 'setex' command"),
		},
		{
			name: "Setex",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("2")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.Ok(),
		},
		{
			name: "Invalid TTL (zero)",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("0")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("invalid expire time in 'setex' command"),
		},
		{
			name: "Invalid TTL (negative)",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("-1")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("invalid expire time in 'setex' command"),
		},
	}
	for _, tt := range tests {
		e := executor{
			handlers: make(map[CommandType]CommandHandler),
			db:       store.NewStore(),
		}
		e.handlers[SETEX] = e.setex
		t.Run(tt.name, func(t *testing.T) {
			if got := e.Execute(SETEX, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setex() = %v, want %v", got, tt.want)
			}
		})
	}
}
