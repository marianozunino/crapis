package command

import (
	"reflect"
	"testing"

	"github.com/marianozunino/crapis/internal/resp"
	"github.com/marianozunino/crapis/internal/store"
)

func Test_get(t *testing.T) {
	type args struct {
		args []resp.Value
	}
	tests := []struct {
		name string
		args args
		want resp.Value
		run  func(*store.Store)
	}{
		{
			name: "No Args",
			args: args{args: []resp.Value{}},
			want: resp.NewError("wrong number of arguments for 'get' command"),
		},
		{
			name: "Get",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewString("test"),
			run: func(s *store.Store) {
				s.StoreValue("test", "test")
			},
		},
		{
			name: "Get Nil",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewNull(),
		},
		{
			name: "Invalid Amount of Args",
			args: args{args: []resp.Value{
				resp.NewBulk(stringPtr("test")),
				resp.NewBulk(stringPtr("test")),
			}},
			want: resp.NewError("wrong number of arguments for 'get' command"),
		},
	}
	for _, tt := range tests {

		e := executor{
			handlers: make(map[CommandType]CommandHandler),
			db:       store.NewStore(),
		}

		e.handlers[GET] = e.get

		t.Run(tt.name, func(t *testing.T) {
			if tt.run != nil {
				tt.run(e.db)
			}
			if got := e.Execute(GET, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}
