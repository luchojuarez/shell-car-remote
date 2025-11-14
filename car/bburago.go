package car

import (
	"fmt"

	"github.com/shell-car-remote/service"
)

var _ CarFactory = BburagoFactory{}

type BburagoFactory struct{}

func (f BburagoFactory) GetCharacteristicsRepo() CharacteristicsRepo {
	return BburagoCharacteristicsRepo{}
}
func (f BburagoFactory) GetInitialStatus() Status {
	return GetInitialStatus()
}
func (f BburagoFactory) GetCipher() service.Cipher {
	return service.NoopCipher{}
}

func (f BburagoFactory) GetBatteryNotificationHandler(carName string) BatteryNotificationHandler {
	return func(input []byte) {
		fmt.Printf("%s: %d%%\n", carName, input[0])
	}
}

// BBstatus implements Status
var _ Status = BBStatus{}

type BBStatus struct {
	payload *[]byte
}

var _ CharacteristicsRepo = BburagoCharacteristicsRepo{}

type BburagoCharacteristicsRepo struct{}

func (repo BburagoCharacteristicsRepo) GetDriveID() string {
	return "0000fff1-0000-1000-8000-00805f9b34fb"
}
func (repo BburagoCharacteristicsRepo) GetbatteryID() string {
	return "00002a19-0000-1000-8000-00805f9b34fb"
}

func GetInitialStatus() *BBStatus {
	return &BBStatus{
		payload: &[]byte{
			0x01, //[0] always One
			0x00, //[1] forward
			0x00, //[2] backward
			0x00, //[3] left
			0x00, //[4] rigth
			0x00, //[5] ligths
			0x01, //[6] turbo
			0x00, //[7] always zero
		},
	}
}

func (m BBStatus) Throttle() {
	(*m.payload)[1] = One
	(*m.payload)[2] = Zero
}
func (m BBStatus) ThrottleRelease() {
	(*m.payload)[1] = Zero
}

func (m BBStatus) Reverse() {
	(*m.payload)[2] = One
	(*m.payload)[1] = Zero
}
func (m BBStatus) ReverseRelease() {
	(*m.payload)[2] = Zero
}

func (m BBStatus) Rigth() {
	(*m.payload)[4] = One
}
func (m BBStatus) Left() {
	(*m.payload)[3] = One
}

func (m BBStatus) Straight() {
	(*m.payload)[3] = Zero
	(*m.payload)[4] = Zero
}

func (m BBStatus) Turbo() {
	(*m.payload)[6] = One
}

func (m BBStatus) Normal() {
	(*m.payload)[6] = Zero
}

func (m BBStatus) Ligths() {
	(*m.payload)[5] = (^(*m.payload)[5]) & Mask

}

func (m BBStatus) Hexa() string {
	return fmt.Sprintf("%x", (*m.payload))
}

func (m BBStatus) Human() string {
	return fmt.Sprintf("not implemented %x", (*m.payload))
}

func (m BBStatus) Payload() []byte {
	return *(m.payload)
}
