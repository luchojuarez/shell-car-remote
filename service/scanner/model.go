package scanner

import (
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

type Result struct {
	Name       string
	ScanResult bluetooth.ScanResult
	Device     *bluetooth.Device
	Paired     bool
	CarBrand   int
}

func (r Result) String() string {
	status := "unpaired"
	if r.Paired {
		status = "paired"
	}
	return fmt.Sprintf("%s (%s/%s) %s", r.Name, r.ScanResult.LocalName(), r.Device.Address.String(), status)
}

type BLE struct {
	adapter      *bluetooth.Adapter
	foundDevices []*Result
	mu           sync.Mutex
}

func (ble *BLE) UnpairedDevices() ([]*Result, error) {
	ble.mu.Lock()
	defer ble.mu.Unlock()
	var result []*Result
	for _, d := range ble.foundDevices {
		if !d.Paired {
			result = append(result, d)
		}
	}
	return result, nil
}
func (ble *BLE) Device(name string) *Result {
	for _, d := range ble.foundDevices {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (r *Result) Devices() *bluetooth.Device {
	return r.Device
}
