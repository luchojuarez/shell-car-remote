package models

import (
	"fmt"
)

type Message struct {
	payload []byte
}

const (
	SpeedFast   byte = 0x64
	SpeedNormal byte = 0x50

	Zero byte = 0x00
	One  byte = 0x01
	Mask byte = One
)

func NewMessage() Message {
	return Message{
		payload: []byte{
			0x00,        // 0: Unknown purpose
			0x43,        // 1: C
			0x54,        // 2: T
			0x4c,        // 3: L
			0x00,        // 4: fordward
			0x00,        // 5: backward
			0x00,        // 6: left
			0x00,        // 7: rigth
			0x01,        // 8: ligths
			SpeedNormal, // 9: speed

			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //padding
		},
	}
}
func MessageFromPayload(payload []byte) Message {
	return Message{
		payload: payload,
	}
}

func (m *Message) Hexa() string {
	return fmt.Sprintf("%x", m.payload)
}

func (m *Message) Human() string {
	move := "none"
	if m.payload[5] != Zero {
		move = "back"
	}
	if m.payload[4] != Zero {
		move = "forward"
	}

	speed := "normal"
	if m.payload[8] == SpeedFast {
		speed = "turbo"
	}

	steeringwheel := "none"
	if m.payload[7] != Zero {
		steeringwheel = "rigth"
	}
	if m.payload[6] != Zero {
		steeringwheel = "left"
	}

	return fmt.Sprintf("[move:%s, steeringwheel:%s, speed:%s]", move, steeringwheel, speed)
}

func (m *Message) Payload() []byte {
	return m.payload
}

func (m *Message) Throttle() {
	m.payload[4] = One
	m.payload[5] = Zero
}
func (m *Message) ThrottleRelease() {
	m.payload[4] = Zero
}

func (m *Message) Reverse() {
	m.payload[5] = One
	m.payload[4] = Zero
}
func (m *Message) ReverseRelease() {
	m.payload[5] = Zero
}

func (m *Message) Rigth() {
	m.payload[7] = One
	//m.payload[6] = Zero
}
func (m *Message) Left() {
	m.payload[6] = One
	//m.payload[7] = Zero
}

func (m *Message) Straight() {
	m.payload[6] = Zero
	m.payload[7] = Zero
}

func (m *Message) Turbo() {
	m.payload[9] = 0x64
}

func (m *Message) Normal() {
	m.payload[9] = 0x50
}

func (m *Message) Ligths() {
	m.payload[8] = (^m.payload[8]) & Mask

}
