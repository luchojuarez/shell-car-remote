//go:build linux

package input

import (
	"context"
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/mrasband/ps4"
	"log"
)

const (
// Map Keys
)

func NewDS4Input(inputController *ps4.Input) DS4Input {
	ch := make(chan Command, 10)
	return DS4Input{
		pressedKeys: &ch,
		input:       inputController,
	}
}

type DS4Input struct {
	pressedKeys *chan Command
	input       *ps4.Input
}

func (ds DS4Input) Listen() *chan Command {
	go func(ds DS4Input) {
		listenDS4(context.Background(), ds)
	}(ds)
	return ds.pressedKeys
}

func ScanControllers() chan *ps4.Input {
	ch := make(chan *ps4.Input)
	go func(ch chan *ps4.Input) {
		controllers, err := GetAllDSControllers()
		if err != nil {
			panic(err)
		}
		for _, controller := range controllers {
			ch <- controller
		}
		ch <- nil
	}(ch)
	return ch
}

func GetAllDSControllers() (out []*ps4.Input, e error) {
	inputs, err := ps4.Discover()
	if err != nil {
		if err.Error() == "unable to find any controller inputs, is it paired and on?" {
			return nil, nil
		}
		e = err
		return
	}

	for _, input := range inputs {
		if input.Type == ps4.Controller {
			out = append(out, input)
		}
	}
	return
}

func listenDS4(ctx context.Context, ds DS4Input) error {
	events, _ := ps4.Watch(ctx, ds.input)
	for e := range events {
		//map analog inputs
		absEvent, ok := e.(*ps4.AbsEvent)
		if ok {
			switch absEvent.Button {
			case ps4.R2:
				if absEvent.Value == 0 {
					*ds.pressedKeys <- Command{Key: Forward, Value: Release}
					*ds.pressedKeys <- Command{Key: Turbo, Value: Release}
					break
				}
				if absEvent.Value >= 127 {
					*ds.pressedKeys <- Command{Key: Turbo, Value: Press}
				} else {
					*ds.pressedKeys <- Command{Key: Turbo, Value: Release}
				}
				*ds.pressedKeys <- Command{Key: Forward, Value: Press}
			case ps4.L2:
				if absEvent.Value == 0 {
					*ds.pressedKeys <- Command{Key: Backward, Value: Release}
					*ds.pressedKeys <- Command{Key: Turbo, Value: Release}
					break
				}
				if absEvent.Value >= 127 {
					*ds.pressedKeys <- Command{Key: Turbo, Value: Press}
				} else {
					*ds.pressedKeys <- Command{Key: Turbo, Value: Release}
				}
				*ds.pressedKeys <- Command{Key: Backward, Value: Press}

			// direction
			case ps4.DPadX:
				switch absEvent.Value {
				case 1:
					*ds.pressedKeys <- Command{Key: Right, Value: Press}
				case -1:
					*ds.pressedKeys <- Command{Key: Left, Value: Press}
				default:
					*ds.pressedKeys <- Command{Key: Right, Value: Release}
					*ds.pressedKeys <- Command{Key: Left, Value: Release}
				}
			case ps4.LeftStickX:
				// +- 7 is to prevent false positive readings on analogs
				if absEvent.Value >= 120 && absEvent.Value <= 135 {
					break
				}
				if absEvent.Value <= 85 {
					*ds.pressedKeys <- Command{Key: Left, Value: Press}
					break
				}
				if absEvent.Value >= 170 {
					*ds.pressedKeys <- Command{Key: Right, Value: Press}
					break
				}
				*ds.pressedKeys <- Command{Key: Left, Value: Release}

				//default:
				//	fmt.Printf("%+v\n", absEvent)
			}

		} else {
			fmt.Printf("%T, %+v\n", e, e)
		}

		//map analog inputs
		keyevent, ok := e.(*ps4.KeyEvent)
		if ok {
			switch keyevent.Button {
			//Headlights.
			case ps4.Playstation:
				if keyevent.State == ps4.KeyDown {
					*ds.pressedKeys <- Command{Key: Headlights, Value: Press}
				} else {
					*ds.pressedKeys <- Command{Key: Headlights, Value: Release}
				}
			}
		}

		//fmt.Printf("%T\n", e)

	}

	return nil
}

func listen0DS4(ds DS4Input) error {
	// Find and open the keyboard device
	devices, err := evdev.ListInputDevices("/dev/input/event*")
	if err != nil {
		log.Fatalf("Failed to list input devices: %v", err)
	}

	var device *evdev.InputDevice
	for _, d := range devices {
		//TODO implement best way to detect rigth keyboard
		if d.Name == "Sony Computer Entertainment Wireless Controller Touchpad" {
			device = d
			break
		}
		fmt.Println(d)
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
				fmt.Printf("llego events queue: %s...\n", event.String())
				//var value ValueCommand = Press
				/*
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
					}*/
			}

		}
	}
}
