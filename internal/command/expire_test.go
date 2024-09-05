package command

import (
	"reflect"
	"testing"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

func Test_expire(t *testing.T) {
	type args struct {
		args []resp.Value
	}
	tests := []struct {
		name string
		args args
		want resp.Value
	}{
		{
			name: "Del No Args",
			args: args{args: []resp.Value{}},
			want: resp.NewError("wrong number of arguments for 'del' command"),
		},
		{
			name: "Del",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("del_key1")),
			}},
			want: resp.NewInteger(0),
		},
		{
			name: "Multiple Keys",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("del_key1")),
				resp.NewBulk(stringPtr("del_key2")),
			}},
			want: resp.NewInteger(0),
		},
	}
	for _, tt := range tests {
		e := executor{
			handlers: make(map[CommandType]CommandHandler),
			db:       store.NewStore(),
		}
		e.handlers[DEL] = e.del
		t.Run(tt.name, func(t *testing.T) {
			if got := e.Execute(DEL, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("del() = %v, want %v", got, tt.want)
			}
		})
	}
}
