package resp

import (
	"reflect"
	"testing"
)

func Test_ping(t *testing.T) {
	type args struct {
		args []Value
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "Ping No Args",
			args: args{args: []Value{}},
			want: Value{kind: STRING, strVal: "PONG"},
		},
		{
			name: "Ping With Args",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: STRING, strVal: "test"},
		},
		{
			name: "Ping Invalid Amount of Args",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'ping' command"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ping(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set(t *testing.T) {
	type args struct {
		args []Value
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "Set No Args",
			args: args{args: []Value{}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'set' command"},
		},
		{
			name: "Set With Args",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: STRING, strVal: "OK"},
		},
		{
			name: "Set Invalid Amount of Args",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'set' command"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setex(t *testing.T) {
	type args struct {
		args []Value
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "Setex No Args",
			args: args{args: []Value{}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'setex' command"},
		},
		{
			name: "Setex Invalid TTL",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "value is not an integer or out of range"},
		},
		{
			name: "Invalid Amount of Args (2)",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'setex' command"},
		},
		{
			name: "Invalid Amount of Args (4)",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("2")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'setex' command"},
		},
		{
			name: "Setex",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("2")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: STRING, strVal: "OK"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setex(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		args []Value
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "Get No Args",
			args: args{args: []Value{}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'get' command"},
		},
		{
			name: "Get",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: STRING, strVal: "test"},
		},
		{
			name: "Get Invalid Amount of Args",
			args: args{args: []Value{
				Value{kind: BULK, bulkVal: stringPtr("test")},
				Value{kind: BULK, bulkVal: stringPtr("test")},
			}},
			want: Value{kind: ERROR, strVal: "wrong number of arguments for 'get' command"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}
