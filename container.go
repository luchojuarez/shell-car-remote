package main

import (
	"encoding/hex"
	"github.com/shell-car-remote/service"
	"github.com/shell-car-remote/service/scanner"
)

type Container struct {
	cipher *service.AesEcbCipher
	ble    *scanner.BLE
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) GetCipher() *service.AesEcbCipher {
	if c.cipher == nil {
		hexKey, err := hex.DecodeString("34522a5b7a6e492c08090a9d8d2a23f8")
		if err != nil {
			panic(err)
		}
		cipher := service.NewAesEcbCipher(hexKey)
		c.cipher = &cipher
	}
	return c.cipher
}

func (c *Container) GetBLE() *scanner.BLE {
	if c.ble == nil {
		ble, err := scanner.NewBLE()
		if err != nil {
			panic(err)
		}
		c.ble = ble
	}
	return c.ble
}
