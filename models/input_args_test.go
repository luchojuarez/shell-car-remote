package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInputArgs(t *testing.T) {
	type args struct {
		args []string
	}
	type want struct {
		output InputArgs
		err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty parameters should return empty input",
			args: args{
				[]string{"programname"},
			},
			want: want{
				output: InputArgs{},
				err:    nil,
			},
		},
		{
			name: "player parameter should return no error",
			args: args{
				[]string{"programname", "-p", "1"},
			},
			want: want{
				output: InputArgs{
					Players: 1,
				},
				err: nil,
			},
		},
		{
			name: "empty parameters should return empty input",
			args: args{
				[]string{"programname", "-p", "NaN"},
			},
			want: want{
				output: InputArgs{
					Players: 1,
				},
				err: errors.New("players must be a int strconv.Atoi: parsing \"NaN\": invalid syntax"),
			},
		},
		{
			name: "invalid number of params must return error",
			args: args{
				[]string{"programname", "-p"},
			},
			want: want{
				output: InputArgs{
					Players: 1,
				},
				err: InputError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewInputArgs(tt.args.args)
			if tt.want.err != nil {
				assert.Error(t, gotErr)
				assert.Equal(t, tt.want.err.Error(), gotErr.Error())
				return
			}
			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want.output, got)
		})
	}
}
