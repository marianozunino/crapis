package command

import "testing"

func TestParseCommand(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandType
		wantErr bool
	}{
		{
			name: "PING",
			args: args{in: "PING"},
			want: PING,
		},
		{
			name: "GET",
			args: args{in: "GET"},
			want: GET,
		},
		{
			name: "SET",
			args: args{in: "SET"},
			want: SET,
		},
		{
			name: "SETEX",
			args: args{in: "SETEX"},
			want: SETEX,
		},
		{
			name: "DEL",
			args: args{in: "DEL"},
			want: DEL,
		},
		{
			name:    "Invalid Command",
			args:    args{in: "INVALID"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ParseCommand(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCommand() got = %v, want %v", got, tt.want)
			}

		})
	}
}
