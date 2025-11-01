package car

import "github.com/shell-car-remote/service"

type CarBrand string

const (
	//first iteration of this cars, with removable battery.
	Branbase CarBrand = "Brandbase"

	//second iteration, battery inside.
	Bburago CarBrand = "Bburago"
)

const (
	Zero byte = 0x00
	One  byte = 0x01
	Mask byte = One
)

type CarFactory interface {
	GetCharacteristicsRepo() CharacteristicsRepo
	GetInitialStatus() Status
	GetCipher() service.Cipher
	GetBatteryNotificationHandler(carName string) BatteryNotificationHandler
}
type BatteryNotificationHandler func(input []byte)

type CharacteristicsRepo interface {
	GetDriveID() string
	GetbatteryID() string
}

type Status interface {
	Hexa() string
	Payload() []byte
	Human() string
	Throttle()
	ThrottleRelease()
	Reverse()
	ReverseRelease()
	Rigth()
	Left()
	Straight()
	Turbo()
	Normal()
	Ligths()
}
