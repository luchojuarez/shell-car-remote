package car

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/shell-car-remote/service"
)

var _ CarFactory = BrandFactory{}

type BrandFactory struct{}

func (f BrandFactory) GetCharacteristicsRepo() CharacteristicsRepo {
	return BrandCharacteristicsRepo{}
}
func (f BrandFactory) GetInitialStatus() Status {
	return NewBrandMessage()
}
func (f BrandFactory) GetCipher() service.Cipher {
	hexKey, err := hex.DecodeString("34522a5b7a6e492c08090a9d8d2a23f8")
	if err != nil {
		panic(err)
	}
	cipher := service.NewAesEcbCipher(hexKey)
	return cipher
}

func (f BrandFactory) GetBatteryNotificationHandler(carName string) BatteryNotificationHandler {
	return func(input []byte) {
		batteryPercentage, _ := strconv.ParseUint(fmt.Sprintf("%x", input[4]), 16, 16)
		fmt.Printf("%s battery level %d%% \n", carName, batteryPercentage)
	}
}

var _ Status = BrandStatus{}

type BrandStatus struct {
	payload *[]byte
}

var _ CharacteristicsRepo = BrandCharacteristicsRepo{}

type BrandCharacteristicsRepo struct{}

func (repo BrandCharacteristicsRepo) GetDriveID() string {
	return "d44bc439-abfd-45a2-b575-925416129600"
}
func (repo BrandCharacteristicsRepo) GetbatteryID() string {
	return "d44bc439-abfd-45a2-b575-925416129601"
}

const (
	SpeedFast   byte = 0x64
	SpeedNormal byte = 0x50
)

func NewBrandMessage() BrandStatus {
	return BrandStatus{
		payload: &[]byte{
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

func (m BrandStatus) Hexa() string {
	return fmt.Sprintf("%x", *(m.payload))
}

func (m BrandStatus) Human() string {
	move := "none"
	if (*m.payload)[5] != Zero {
		move = "back"
	}
	if (*m.payload)[4] != Zero {
		move = "forward"
	}

	speed := "normal"
	if (*m.payload)[8] == SpeedFast {
		speed = "turbo"
	}

	steeringwheel := "none"
	if (*m.payload)[7] != Zero {
		steeringwheel = "rigth"
	}
	if (*m.payload)[6] != Zero {
		steeringwheel = "left"
	}

	return fmt.Sprintf("[move:%s, steeringwheel:%s, speed:%s]", move, steeringwheel, speed)
}

func (m BrandStatus) Payload() []byte {
	return *(m.payload)
}

func (m BrandStatus) Throttle() {
	(*m.payload)[4] = One
	(*m.payload)[5] = Zero
}
func (m BrandStatus) ThrottleRelease() {
	(*m.payload)[4] = Zero
}

func (m BrandStatus) Reverse() {
	(*m.payload)[5] = One
	(*m.payload)[4] = Zero
}
func (m BrandStatus) ReverseRelease() {
	(*m.payload)[5] = Zero
}

func (m BrandStatus) Rigth() {
	(*m.payload)[7] = One
}
func (m BrandStatus) Left() {
	(*m.payload)[6] = One
}

func (m BrandStatus) Straight() {
	(*m.payload)[6] = Zero
	(*m.payload)[7] = Zero
}

func (m BrandStatus) Turbo() {
	(*m.payload)[9] = 0x64
}

func (m BrandStatus) Normal() {
	(*m.payload)[9] = 0x50
}

func (m BrandStatus) Ligths() {
	(*m.payload)[8] = (^(*m.payload)[8]) & Mask

}
