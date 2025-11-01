package car

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shell-car-remote/input"
)

func (thisCar *Car) ListenController() {

	go func(ch *chan input.Command, output *Status) {

		for {
			select {

			case event, ok := <-*ch:
				//fmt.Printf("Received: %+v\n", event)
				if event.Value == input.Hold {
					break
				}
				if !ok {
					// Channel closed, exit the loop
					fmt.Println("Channel closed.")
					return
				}
				switch event.Key {
				case input.Forward:
					if event.Value == input.Release {
						(*output).ThrottleRelease()
					} else {
						(*output).Throttle()
					}
				case input.Backward:
					if event.Value == input.Release {
						(*output).ReverseRelease()
					} else {
						(*output).Reverse()
					}
				case input.Right:
					if event.Value == input.Release {
						(*output).Straight()
					} else {
						(*output).Rigth()
					}
				case input.Left:
					if event.Value == input.Release {
						(*output).Straight()
					} else {
						(*output).Left()
					}
				case input.Headlights:
					if event.Value == input.Press {
						(*output).Ligths()
					}
				case input.Turbo:
					if event.Value == input.Press {
						(*output).Turbo()
					} else {
						(*output).Normal()
					}
				}
			}
		}
	}(thisCar.controller, thisCar.carStatus)
}

func (thisCar *Car) StartTransmission() {
	for {
		go func(car *Car) {
			err := car.SendMessage(thisCar.carStatus)
			if err != nil {
				fmt.Printf("error sending to car: %s\n", err.Error())
			}
		}(thisCar)

		time.Sleep(10 * time.Millisecond)
	}
}

func (thisCar *Car) SendMessage(status *Status) error {
	initTime := time.Now()
	s := *status

	encryptedMsg, err := thisCar.cipher.Encrypt(s.Payload())
	if err != nil {
		return fmt.Errorf("error while cipher message '%+v'", err)
	}

	_, err = thisCar.driveCharacteristic.WriteWithoutResponse(encryptedMsg)
	if err != nil {
		switch err.Error() {
		case "Not connected":
			err, ok := thisCar.Reconnect()

			if err != nil {
				return fmt.Errorf("reconnect error '%s': %s", thisCar.ble.Address, err.Error())
			}
			if !ok {
				fmt.Printf("reconnecting...")
				return nil
			}
		case "In Progress":
			time.Sleep(50 * time.Microsecond)
		default:
			fmt.Printf("unkow send error %s\n", err.Error())

		}
	}

	// Flag this param as 'show statics'.
	if false {
		fmt.Printf("%s sent, took %dms\n", (*status).Human(), time.Since(initTime).Milliseconds())
	}

	return nil
}

func (car *Car) Disconnect() error {
	e := car.ble.Disconnect()
	time.Sleep(100 * time.Millisecond)
	return e
}
func (car *Car) Reconnect() (error, bool) {
	device, err := car.adapter.Reconnect(car.ble.Address)
	if err != nil {
		if strings.Contains(err.Error(), "Operation already in progress") {
			return nil, false
		}
		panic(fmt.Sprintf("error '%s'", err.Error()))
	}
	car.ble = &device

	return nil, true
}

func (car *Car) EnableBatteryNotification() {
	e := car.batteryCharacteristic.EnableNotifications(
		func(buf []byte) {
			message, err := car.cipher.Decrypt(buf)
			if err != nil {
				fmt.Printf("cant decrypt battery status message [%x]: %s\n", buf, err.Error())
			} else {
				car.batteryHandler(message)
			}
		},
	)
	if e != nil {
		panic(errors.New("enable battery notification error: " + e.Error()))
	}
}
