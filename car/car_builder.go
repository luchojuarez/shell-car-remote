package car

import (
	"strings"

	"github.com/shell-car-remote/input"
	"github.com/shell-car-remote/service"
	"github.com/shell-car-remote/service/scanner"
	"tinygo.org/x/bluetooth"
)

type Car struct {
	ble                   *bluetooth.Device
	driveCharacteristic   *bluetooth.DeviceCharacteristic
	batteryCharacteristic *bluetooth.DeviceCharacteristic
	adapter               *scanner.BLE
	controller            *chan input.Command
	carStatus             *Status
	cipher                service.Cipher
	batteryHandler        BatteryNotificationHandler
}

func NewCar(
	scannerResult *scanner.Result,
	controller *chan input.Command,
	adapter *scanner.BLE,

) (*Car, error) {
	deviceName := scannerResult.ScanResult.LocalName()

	//choose car factory
	var carFactory CarFactory = BburagoFactory{}
	if strings.HasPrefix(deviceName, "QCAR") {
		carFactory = BrandFactory{}
	}

	driveCharacteristic, err := scanner.GetCharacteristicByUUID(*scannerResult.Devices(), carFactory.GetCharacteristicsRepo().GetDriveID())
	if err != nil {
		panic(err)
	}
	batteryCharacteristic, err := scanner.GetCharacteristicByUUID(*scannerResult.Devices(), carFactory.GetCharacteristicsRepo().GetbatteryID())
	if err != nil {
		panic(err)
	}

	initialStatus := carFactory.GetInitialStatus()

	car := Car{
		ble:                   scannerResult.Device,
		driveCharacteristic:   driveCharacteristic,
		batteryCharacteristic: batteryCharacteristic,
		controller:            controller,
		carStatus:             &initialStatus,
		cipher:                carFactory.GetCipher(),
		batteryHandler:        carFactory.GetBatteryNotificationHandler(deviceName),
	}

	return &car, nil

}

func (c *Car) Start() {
	go c.StartTransmission()
	go c.ListenController()
	go c.EnableBatteryNotification()
}
