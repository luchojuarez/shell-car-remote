package service

import (
	"fmt"
	"github.com/shell-car-remote/service/scanner"
	"time"

	"github.com/shell-car-remote/input"
	"github.com/shell-car-remote/models"
	"tinygo.org/x/bluetooth"
)

var characteristicDrive = "d44bc439-abfd-45a2-b575-925416129600"

type QCAR struct {
	cipher              AesEcbCipher
	ble                 *bluetooth.Device
	driveCharacteristic *bluetooth.DeviceCharacteristic
	controller          *chan input.Command
	carStatus           *models.Message
	//TODO add battery characteristic
}

func (car *QCAR) Disconnect() error {
	e := car.ble.Disconnect()
	time.Sleep(100 * time.Millisecond)
	return e
}
func (car *QCAR) Reconnect() error {
	driveCharacteristic, err := GetDriveCharacteristic(car.ble)
	if err != nil {
		return err
	}
	car.driveCharacteristic = driveCharacteristic

	return nil
}

func GetDriveCharacteristic(ble *bluetooth.Device) (*bluetooth.DeviceCharacteristic, error) {
	characteristic, err := scanner.GetCharacteristicByUUID(*ble, characteristicDrive)
	if err != nil {
		return nil, fmt.Errorf("connect error: %s", err.Error())
	}
	return characteristic, nil
}

// NewQCar initializes and returns a new QCAR instance.
func NewQCar(
	cipher AesEcbCipher,
	device *bluetooth.Device,
	controller *chan input.Command,

) (*QCAR, error) {
	driveCharacteristic, err := GetDriveCharacteristic(device)
	if err != nil {
		return nil, err
	}
	idle := models.NewMessage()

	// Create and return the QCAR instance
	instance := &QCAR{
		cipher:              cipher,
		ble:                 device,
		driveCharacteristic: driveCharacteristic,
		controller:          controller,
		carStatus:           &idle, // new idle status
	}
	instance.ListenController()

	return instance, nil
}

func (car *QCAR) ListenController() {
	go func(ch *chan input.Command, output *models.Message) {

		listenController(ch, output)
	}(car.controller, car.carStatus)
}

func listenController(inputChannel *chan input.Command, output *models.Message) {
	for {
		select {

		case event, ok := <-*inputChannel:
			fmt.Printf("Received: %+v\n", event)
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
					output.ThrottleRelease()
				} else {
					output.Throttle()
				}
			case input.Backward:
				if event.Value == input.Release {
					output.ReverseRelease()
				} else {
					output.Reverse()
				}
			case input.Right:
				if event.Value == input.Release {
					output.Straight()
				} else {
					output.Rigth()
				}
			case input.Left:
				if event.Value == input.Release {
					output.Straight()
				} else {
					output.Left()
				}
			case input.Headlights:
				if event.Value == input.Press {
					output.Ligths()
				}
			case input.Turbo:
				if event.Value == input.Press {
					output.Turbo()
				} else {
					output.Normal()
				}
			}
		}
	}
}

func (car *QCAR) StartTransmission() {
	for {
		go func(car *QCAR) {
			err := car.SendMessage(*car.carStatus)
			if err != nil {
				fmt.Printf("error sending to car: %s\n", err.Error())
			}
		}(car)

		time.Sleep(10 * time.Millisecond)
	}
}

func (car *QCAR) SendMessage(message models.Message) error {
	initTime := time.Now()
	encriptedMsg, err := car.cipher.Encrypt(message.Payload())
	if err != nil {
		return fmt.Errorf("error while cipher message '%+v'", err)
	}

	_, err = car.driveCharacteristic.WriteWithoutResponse(encriptedMsg)
	if err != nil {
		switch err.Error() {
		case "Not connected":
			fmt.Printf("reconnect... ")
			err := car.Reconnect()
			if err != nil {
				return fmt.Errorf("reconnect error '%s': %s", car.ble.Address, err.Error())
			}
			fmt.Printf("reconnected\n")
		case "In Progress":
			time.Sleep(50 * time.Microsecond)
		default:
			fmt.Printf("unkow send error %s\n", err.Error())

		}
	}

	// Flag this param as 'show statics'.
	if false {
		fmt.Printf("%s sent, took %dms\n", message.Human(), time.Now().Sub(initTime).Milliseconds())
	}

	return nil
}
