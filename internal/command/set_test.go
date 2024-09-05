package command

import (
	"reflect"
	"testing"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

func Test_set(t *testing.T) {
	type args struct {
		args []resp.Value
	}
	tests := []struct {
		name string
		args args
		want resp.Value
	}{
		{
			name: "Set No Args",
			args: args{args: []resp.Value{}},
			want: resp.NewError("wrong number of arguments for 'set' command"),
		},
		{
			name: "Set With Args",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.Ok(),
		},
		{
			name: "Set Invalid Amount of Args",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("wrong number of arguments for 'set' command"),
		},
	}
	for _, tt := range tests {
		e := executor{
			handlers: make(map[CommandType]CommandHandler),
			db:       store.NewStore(),
		}
		e.handlers[SET] = e.set
		t.Run(tt.name, func(t *testing.T) {

			if got := e.Execute(SET, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("set() = %v, want %v", got, tt.want)
			}
		})
	}
}
