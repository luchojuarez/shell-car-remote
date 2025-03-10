package scanner

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

var bleConnectionParams = bluetooth.ConnectionParams{
	Timeout:           bluetooth.NewDuration(500 * time.Microsecond),
	ConnectionTimeout: bluetooth.NewDuration(50 * time.Microsecond),
	MinInterval:       bluetooth.NewDuration(9 * time.Microsecond),
	MaxInterval:       bluetooth.NewDuration(100 * time.Microsecond),
}

func NewBLE() (*BLE, error) {
	// Initialize the Bluetooth adapter
	adapter := bluetooth.DefaultAdapter

	// Enable the Bluetooth adapter
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("bluetooth adapter failed to enable: %w", err)
	}
	return &BLE{
		adapter:      adapter,
		foundDevices: []*Result{},
	}, nil
}

func (ble *BLE) Scan(ctx context.Context) chan *Result {
	ch := make(chan *Result)
	go func(ctx context.Context, b *BLE, ch chan *Result) {
		if err := b.scan(ctx, ch); err != nil {
			panic(err)
		}
	}(ctx, ble, ch)
	return ch
}
func (ble *BLE) scan(ctx context.Context, newDevices chan *Result) error {
	go func(ctx context.Context, newDevices chan *Result, bleScanner *BLE) {
		err := bleScanner.adapter.Scan(func(adapter *bluetooth.Adapter, resutl bluetooth.ScanResult) {
			if strings.HasPrefix(resutl.LocalName(), "QCAR") {
				if bleScanner.Device(resutl.LocalName()) != nil {
					log.Printf("already found device %s", resutl.LocalName())
				} else {
					device, err := adapter.Connect(resutl.Address, bleConnectionParams)

					if err != nil {
						log.Printf("Failed to connect to device: %s (%s)\n", resutl.LocalName(), resutl.Address.String())
						return
					}

					r := &Result{
						name:       fmt.Sprintf("CAR %d", len(ble.foundDevices)+1),
						scanResult: resutl,
						device:     &device,
						paired:     false,
					}

					ble.foundDevices = append(ble.foundDevices, r)
					newDevices <- r
				}
			}
		})
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(ctx, newDevices, ble)

	for {
		select {
		case <-ctx.Done(): // Listen for cancel signal
			if err := ble.adapter.StopScan(); err != nil {
				return fmt.Errorf("bluetooth adapter failed to stop scan: %w", err)
			}
			time.Sleep(100 * time.Millisecond)
			return nil
		}
	}
}

func (ble *BLE) Reconnect(address bluetooth.Address) (bluetooth.Device, error) {
	return ble.adapter.Connect(address, bleConnectionParams)
}
