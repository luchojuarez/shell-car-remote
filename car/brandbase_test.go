package car

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	type want struct{}
	tests := []struct {
		name    string
		want    want
		wantErr bool
	}{
		{
			name: "default values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := NewBrandMessage()
			assert.Equal(t, len(message.Payload()), 16)

			assert.Equal(t, message.Human(), "[move:none, steeringwheel:none, speed:normal]")
			assert.Equal(t, message.Hexa(), "0043544c000000000150000000000000")
		})
	}
}
