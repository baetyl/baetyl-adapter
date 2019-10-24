package main

import (
	"fmt"
	"github.com/goburrow/modbus"
)

type MbReader interface {
	read(cli modbus.Client, address, quantity uint16)(results []byte, err error)
}
var readers map[byte]MbReader

type CoilReader struct{}

func (r *CoilReader)read(cli modbus.Client, address, quantity uint16)(results []byte, err error) {
	results, err = cli.ReadCoils(address, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to read Coils: %s", err.Error())
	}
	return
}

type DiscreteInputsReader struct{}

func (r *DiscreteInputsReader)read(cli modbus.Client, address, quantity uint16)(results []byte, err error) {
	results, err = cli.ReadDiscreteInputs(address, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to read Coils: %s", err.Error())
	}
	return
}

type HoldingRegistersReader struct{}

func (r *HoldingRegistersReader)read(cli modbus.Client, address, quantity uint16)(results []byte, err error) {
	results, err = cli.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to read Coils: %s", err.Error())
	}
	return
}

type InputRegistersReader struct{}

func (r *InputRegistersReader)read(cli modbus.Client, address, quantity uint16)(results []byte, err error) {
	results, err = cli.ReadInputRegisters(address, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to read Coils: %s", err.Error())
	}
	return
}

func init() {
	readers = map[byte]MbReader{}
	readers[1] = &CoilReader{}
	readers[2] = &DiscreteInputsReader{}
	readers[3] = &HoldingRegistersReader{}
	readers[4] = &InputRegistersReader{}
}
