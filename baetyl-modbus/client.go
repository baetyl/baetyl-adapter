package main

import (
	"fmt"
	"strings"

	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
)

type client struct {
	mbClient   modbus.Client
	cfg        Slave
	tcpHandler *modbus.TCPClientHandler
	rtuHandler *modbus.RTUClientHandler
}

func newClient(cfg Slave) (*client, error) {
	cli := client{
		cfg:     cfg,
	}
	if strings.HasPrefix(cfg.Address, "tcp://") {
		// Modbus TCP
		h := modbus.NewTCPClientHandler(cfg.Address[6:])
		h.SlaveId = cfg.ID
		h.Timeout = cfg.Timeout
		h.IdleTimeout = cfg.IdleTimeout
		err := h.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %s", err.Error())
		}
		cli.tcpHandler = h
		cli.mbClient = modbus.NewClient(h)
	} else {
		// Modbus RTU
		h := modbus.NewRTUClientHandler(cfg.Address)
		h.BaudRate = cfg.BaudRate
		h.DataBits = cfg.DataBits
		h.Parity = cfg.Parity
		h.StopBits = cfg.StopBits
		h.SlaveId = cfg.ID
		h.Timeout = cfg.Timeout
		h.IdleTimeout = cfg.IdleTimeout
		h.RS485 = serial.RS485Config{
			Enabled:            cfg.RS485.Enabled,
			DelayRtsBeforeSend: cfg.RS485.DelayRtsBeforeSend,
			DelayRtsAfterSend:  cfg.RS485.DelayRtsAfterSend,
			RtsHighDuringSend:  cfg.RS485.RtsHighDuringSend,
			RtsHighAfterSend:   cfg.RS485.RtsHighAfterSend,
			RxDuringTx:         cfg.RS485.RxDuringTx,
		}
		err := h.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %s", err.Error())
		}
		cli.rtuHandler = h
		cli.mbClient = modbus.NewClient(h)
	}
	return &cli, nil
}

func (cli *client) read(item ParseItem) (results []byte, err error) {
	address := item.Address
	quantity := item.Quantity
	if reader, ok := readers[item.Function]; ok {
		return reader.read(cli.mbClient, address, quantity)
	} else {
		return nil, fmt.Errorf("function code (%d) not supported", item.Function)
	}
}

func (cli *client) close() error {
	if cli.tcpHandler != nil {
		err := cli.tcpHandler.Close()
		if err != nil {
			return fmt.Errorf("failed to close tcp client: %s", err.Error())
		}
	} else if cli.rtuHandler != nil {
		err := cli.tcpHandler.Close()
		if err != nil {
			return fmt.Errorf("failed to close rtu client: %s", err.Error())
		}
	}
	return nil
}
