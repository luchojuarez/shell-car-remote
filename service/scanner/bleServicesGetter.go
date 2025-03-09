package scanner

import (
	"fmt"
	"tinygo.org/x/bluetooth"
)

// Search in all services provided by device a service with matches with UUID
func GetCharacteristicByUUID(device bluetooth.Device, characteristicUUID string) (*bluetooth.DeviceCharacteristic, error) {
	// Discover services on the device
	services, err := device.DiscoverServices(nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to discover services: %v", err)
	}

	// Search for the target characteristic UUID
	var targetChar *bluetooth.DeviceCharacteristic
	for _, service := range services {
		characteristics, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to discover characteristics: %v", err)
		}

		for _, char := range characteristics {
			if char.UUID().String() == characteristicUUID {
				targetChar = &char
				break
			}
		}

		if targetChar != nil {
			break
		}
	}

	return targetChar, nil
}
