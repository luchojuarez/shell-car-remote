package service

import (
	"encoding/hex"
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

func TestAesEcbCipher_Decrypt(t *testing.T) {
	key, _ := hex.DecodeString("34522a5b7a6e492c08090a9d8d2a23f8")
	type params struct {
		input string
	}
	type want struct {
		errMessage string
		output     string
	}
	tests := []struct {
		name    string
		params  params
		want    want
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "default values",
			params: params{
				input: "727a06d7af7e399d01e7d3a559e23843",
			},
			want: want{
				output: "1C" + //28
					"56" + // V
					"42" + // B
					"54" + // T
					"46" + // battery percentage
					"0000000000000000000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewAesEcbCipher(key)
			input, _ := hex.DecodeString(tt.params.input)
			result, err := c.Decrypt(input)
			if tt.want.errMessage != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.want.errMessage, err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.want.output, fmt.Sprintf("%X", result))
		})
	}
}

/*
message battery level 'aae2deb8232cc26dd4dfd41055875941'
message battery level 'd8abc6291e3f137d05cb6930d9bf08fb'
727a06d7af7e399d01e7d3a559e23843
*/
