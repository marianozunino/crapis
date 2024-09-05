package command

import (
	"reflect"
	"testing"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

func Test_ping(t *testing.T) {
	type args struct {
		args []resp.Value
	}
	tests := []struct {
		name string
		args args
		want resp.Value
	}{
		{
			name: "Ping No Args",
			args: args{args: []resp.Value{}},
			want: resp.NewString("PONG"),
		},
		{
			name: "Ping With Args",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewString("test"),
		},
		{
			name: "Ping Invalid Amount of Args",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("wrong number of arguments for 'ping' command"),
		},
	}
	for _, tt := range tests {
		e := executor{
			handlers: make(map[CommandType]CommandHandler),
			db:       store.NewStore(),
		}
		e.handlers[PING] = e.ping

		t.Run(tt.name, func(t *testing.T) {

			if got := e.Execute(PING, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
