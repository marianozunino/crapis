package resp

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "NewError",
			args: args{message: "test"},
			want: NewError("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewString(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "NewString",
			args: args{message: "test"},
			want: NewString("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewString(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBulk(t *testing.T) {
	type args struct {
		message *string
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "NewBulk",
			args: args{message: stringPtr("test")},
			want: NewBulk(stringPtr("test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBulk(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBulk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInteger(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "NewInteger",
			args: args{number: 1},
			want: NewInteger(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInteger(tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewArray(t *testing.T) {
	type args struct {
		values []Value
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "NewArray",
			args: args{values: []Value{NewString("test")}},
			want: NewArray(NewString("test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArray(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNull(t *testing.T) {
	tests := []struct {
		name string
		want Value
	}{
		{
			name: "NewNull",
			want: NewNull(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNull(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOk(t *testing.T) {
	tests := []struct {
		name string
		want Value
	}{
		{
			name: "Ok",
			want: Ok(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ok(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ok() = %v, want %v", got, tt.want)
			}
		})
	}
}
