package scanner

import (
	"fmt"
	"sync"
	"tinygo.org/x/bluetooth"
)

type Result struct {
	name       string
	scanResult bluetooth.ScanResult
	device     *bluetooth.Device
	paired     bool
}

func (r Result) String() string {
	status := "unpaired"
	if r.paired {
		status = "paired"
	}
	return fmt.Sprintf("%s (%s/%s) %s", r.name, r.scanResult.LocalName(), r.device.Address.String(), status)
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
		if !d.paired {
			result = append(result, d)
		}
	}
	return result, nil
}
func (ble *BLE) Device(name string) *Result {
	ble.mu.Lock()
	defer ble.mu.Unlock()
	for _, d := range ble.foundDevices {
		if d.name == name {
			return d
		}
	}
	return nil
}

func (r *Result) Devices() *bluetooth.Device {
	return r.device
}

func (r *Result) Paired() {
	r.paired = true
}
