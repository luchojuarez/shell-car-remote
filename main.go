package main

import (
	"context"
	"fmt"
	"github.com/shell-car-remote/input"
	"github.com/shell-car-remote/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	container := NewContainer()
	cipher := container.GetCipher()
	bleScanner := container.GetBLE()

	ctx := context.Background()

	ctxTO, callback := context.WithTimeout(ctx, 5*time.Second)
	defer callback()
	foundCars := bleScanner.Scan(ctxTO)
	loop := true
	for loop {
		select {
		case device := <-foundCars:
			log.Printf("%+v\n", device)
		case <-ctxTO.Done(): // Listen for cancel signal
			loop = false
		}
	}
	inputs, err := input.GetAllDSControllers()
	if err != nil {
		panic(err)
	}
	if len(inputs) == 0 {
		log.Println("no DS4 found")
	}

	BLECars, err := bleScanner.UnpairedDevices()
	if err != nil {
		panic(err)
	}
	if len(BLECars) == 0 {
		fmt.Println("no BLE devices found")
		return
	}

	// pair controllers with cars.
	for i, controller := range inputs {
		BLEcar := BLECars[i]
		if BLEcar == nil {
			continue
		}

		ds := input.NewDS4Input(controller)
		ch := ds.Listen()

		car, err := service.NewQCar(*cipher, BLEcar.Devices(), ch, bleScanner)

		if err != nil {
			panic(fmt.Sprintf("error building car '%s'", err.Error()))
		}

		BLEcar.Paired() //mark as paired.

		go car.StartTransmission()

	}

	// check for unpaired cars.
	BLECars, err = bleScanner.UnpairedDevices()
	if err != nil {
		panic(err)
	}
	if len(BLECars) != 0 {
		keyboard := input.NewKeyboardInput()
		ch := keyboard.Listen()
		BLE := BLECars[0]

		car, err := service.NewQCar(*cipher, BLE.Devices(), ch, bleScanner)
		if err != nil {
			panic(fmt.Sprintf("error building car '%s'", err.Error()))
		}
		go car.StartTransmission()
		BLE.Paired()
	}

	BLECars, err = bleScanner.UnpairedDevices()
	if err != nil {
		panic(err)
	}
	if len(BLECars) == 0 {
		log.Println("LET'S RACE")
	} else {
		log.Printf("some cars cant be paired 😭: %+v", BLECars)
	}
	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopChan // wait for SIGINT
	log.Printf("\nShutting down...\n")
}
