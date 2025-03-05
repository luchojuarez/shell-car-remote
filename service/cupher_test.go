package service

import (
	"fmt"
	"testing"

	"github.com/shell-car-remote/models"
	"github.com/stretchr/testify/assert"
)

func TestNewAesEcbCipher(t *testing.T) {
	message := models.NewMessage()
	type want struct {
		errMessage string
		output     string
	}
	type params struct {
		key   string
		input []byte
	}
	tests := []struct {
		name   string
		want   want
		params params
	}{
		{
			name: "default values",
			params: params{
				key:   "34522a5b7a6e492c08090a9d8d2a23f8",
				input: message.Payload(),
			},
			want: want{
				output: "85cbe096209a54c892027caed30577f3",
			},
		},
		{
			name: "invalid input",
			params: params{
				key:   "34522a5b7a6e492c08090a9d8d2a23f8",
				input: []byte{0x00, 0x00, 0x00},
			},
			want: want{
				errMessage: "input length must be a multiple of the AES block size (16 bytes)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAesEcbCipher([]byte(tt.params.key))
			result, err := c.Encrypt(tt.params.input)
			if tt.want.errMessage != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.want.errMessage, err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.want.output, fmt.Sprintf("%x", result))
		})
	}
}
