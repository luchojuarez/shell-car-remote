package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	carPkg "github.com/shell-car-remote/car"
	"github.com/shell-car-remote/input"
	"github.com/shell-car-remote/service/scanner"
)

func main() {
	bleScanner, err := scanner.NewBLE()
	if err != nil {
		panic(err)
	}

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

		car, err := carPkg.NewCar(BLEcar, ch, bleScanner)

		if err != nil {
			panic(fmt.Sprintf("error building car '%s'", err.Error()))
		}

		BLEcar.Paired = true //mark as paired.
		car.Start()
	}

	// check for unpaired cars.
	BLECars, err = bleScanner.UnpairedDevices()
	if err != nil {
		panic(err)
	}
	if len(BLECars) != 0 {
		fmt.Println("setting keyboard input")
		keyboard := input.NewKeyboardInput()
		ch := keyboard.Listen()
		BLE := BLECars[0]

		car, err := carPkg.NewCar(BLE, ch, bleScanner)

		if err != nil {
			panic(fmt.Sprintf("error building car '%s'", err.Error()))
		}
		car.Start()
		BLE.Paired = true
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
