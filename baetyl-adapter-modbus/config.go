package main

import "time"

// Config custom configuration of the timer module
type Config struct {
	// slave.go list
	Slaves []SlaveConfig `yaml:"slaves" json:"slaves"`
	// map list
	Maps []MapConfig `yaml:"maps" json:"maps"`
	// publish topic of collected data
	Publish Publish `yaml:"publish" json:"publish" validate:"nonnil"`
}

// SlaveConfig modbus slave device configuration
type SlaveConfig struct {
	ID byte `yaml:"id" json:"id"`
	// Address Device path (/dev/ttyS0)
	Address string `yaml:"address" json:"address" default:"/dev/ttyS0"`
	// Timeout Read (Write) timeout.
	Timeout time.Duration `yaml:"timeout" json:"timeout" default:"10s"`
	// IdleTimeout Idle timeout to close the connection
	IdleTimeout time.Duration `yaml:"idletimeout" json:"idletimeout" default:"1m"`
	//// RTU only
	// BaudRate (default 19200)
	BaudRate int `yaml:"baudrate" json:"baudrate" default:"19200"`
	// DataBits: 5, 6, 7 or 8 (default 8)
	DataBits int `yaml:"databits" json:"databits" default:"8" validate:"min=5, max=8"`
	// StopBits: 1 or 2 (default 1)
	StopBits int `yaml:"stopbits" json:"stopbits" default:"1" validate:"min=1, max=2"`
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string `yaml:"parity" json:"parity" default:"E" validate:"regexp=^(E|N|O)?$"`
	// RS485 Configuration related to RS485
	RS485 struct {
		// Enabled Enable RS485 support
		Enabled bool `yaml:"enabled" json:"enabled"`
		// DelayRtsBeforeSend Delay RTS prior to send
		DelayRtsBeforeSend time.Duration `yaml:"delay_rts_before_send" json:"delay_rts_before_send"`
		// DelayRtsAfterSend Delay RTS after send
		DelayRtsAfterSend time.Duration `yaml:"delay_rts_after_send" json:"delay_rts_after_send"`
		// RtsHighDuringSend Set RTS high during send
		RtsHighDuringSend bool `yaml:"rts_high_during_send" json:"rts_high_during_send"`
		// RtsHighAfterSend Set RTS high after send
		RtsHighAfterSend bool `yaml:"rts_high_after_send" json:"rts_high_after_send"`
		// RxDuringTx Rx during Tx
		RxDuringTx bool `yaml:"rx_during_tx" json:"rx_during_tx"`
	} `yaml:"rs485" json:"rs485"`
	Interval time.Duration `yaml:"interval" json:"interval" default:"5s" validate:"nonzero"`
}

// MapConfig map point configuration
type MapConfig struct {
	// Slave Id
	SlaveID byte `yaml:"slaveid" json:"slaveid"`
	// Function
	Function byte `yaml:"function" json:"function" validate:"min=1, max=4"`
	// Address
	Address uint16 `yaml:"address" json:"address"`
	// Quantity
	Quantity uint16 `yaml:"quantity" json:"quantity"`
}

// Publish publish topic
type Publish struct {
	QOS   uint32 `yaml:"qos" json:"qos" validate:"min=0, max=1"`
	Topic string `yaml:"topic" json:"topic" default:"timer" validate:"nonzero"`
}
