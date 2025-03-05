package input

import (
	"fmt"
	"log"

	evdev "github.com/gvalkov/golang-evdev"
)

const (
	// Map Keyboard
	forward    = evdev.KEY_UP
	backward   = evdev.KEY_DOWN
	right      = evdev.KEY_RIGHT
	left       = evdev.KEY_LEFT
	turbo      = evdev.KEY_T
	headlights = evdev.KEY_L
)

func NewKeyboardInput() KeyboardInput {
	ch := make(chan Command, 10)
	return KeyboardInput{
		pressedKeys: &ch,
	}
}

type KeyboardInput struct {
	pressedKeys *chan Command
}

func (k KeyboardInput) Listen() *chan Command {
	go func(k KeyboardInput) {
		listen(k)
	}(k)
	return k.pressedKeys
}

func listen(k KeyboardInput) error {
	// Find and open the keyboard device
	devices, err := evdev.ListInputDevices()
	if err != nil {
		log.Fatalf("Failed to list input devices: %v", err)
	}

	var device *evdev.InputDevice
	for _, d := range devices {
		//TODO implement best way to detect rigth keyboard
		if d.Name == "Dell Computer Corp Dell Universal Receiver" {
			device = d
			break
		}
	}

	if device == nil {
		return fmt.Errorf("No keyboard device found")
	}

	fmt.Printf("Listening for events on: %s (%s)\n", device.Name, device.Fn)

	// Read events from the keyboard
	for {
		events, err := device.Read()
		if err != nil {
			return fmt.Errorf("Error reading events: %v", err)
		}

		for _, event := range events {
			if event.Type == evdev.EV_KEY {
				var value ValueCommand = Press
				switch event.Value {
				case int32(evdev.KeyDown):
					value = Press
				case int32(evdev.KeyUp):
					value = Release
				case int32(evdev.KeyHold):
					value = Hold
				}

				switch event.Code {
				case forward:
					*k.pressedKeys <- Command{Key: Forward, Value: value}
				case backward:
					*k.pressedKeys <- Command{Key: Backward, Value: value}
				case right:
					*k.pressedKeys <- Command{Key: Right, Value: value}
				case left:
					*k.pressedKeys <- Command{Key: Left, Value: value}
				case headlights:
					*k.pressedKeys <- Command{Key: Headlights, Value: value}
				case turbo:
					*k.pressedKeys <- Command{Key: Turbo, Value: value}
				}
			}

		}
	}
}
