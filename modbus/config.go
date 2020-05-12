package main

import (
	"errors"
	"fmt"
	"gopkg.in/validator.v2"
	"time"
)

func init() {
	validator.SetValidationFunc("validjobs", validateJobs)
}

// Config custom configuration of the timer module
type Config struct {
	// slave list
	Slaves []SlaveConfig `yaml:"slaves" json:"slaves"`
	// job list
	Jobs []Job `yaml:"jobs" json:"jobs" validate:"validjobs"`
	// publish topic of collected data
	Publish Publish `yaml:"publish" json:"publish" validate:"nonnil"`
}

type Job struct {
	SlaveId  byte          `yaml:"slaveid" json:"slaveid"`
	Interval time.Duration `yaml:"interval" json:"interval" default:"5s"`
	Encoding string        `yaml:"encoding" json:"encoding" validate:"regexp=^(binary|json)?$" default:"json"`
	Time     Time          `yaml:"time" json:"time" default:"{\"name\":\"time\", \"type\":\"integer\"}"`
	Maps     []MapConfig   `yaml:"maps" json:"maps"`
}

type Field struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

type Time struct {
	Field     `yaml:",inline" json:",inline"`
	Format    string `yaml:"format" json:"format" default:"2006-01-02 15:04:05"`
	Precision string `yaml:"precision" json:"precision" default:"s" validate:"regexp=^(s|ns)?$"`
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
}

// MapConfig map point configuration
type MapConfig struct {
	// Function
	Function byte `yaml:"function" json:"function" validate:"min=1, max=4" validate:"nonzero"`
	// Address
	Address uint16 `yaml:"address" json:"address" validate:"nonzero"`
	// Quantity
	Quantity uint16 `yaml:"quantity" json:"quantity"`
	// parsed attributes
	Field Field `yaml:"field" json:"field"`
}

// Publish publish topic
type Publish struct {
	QOS   uint32 `yaml:"qos" json:"qos" validate:"min=0, max=1"`
	Topic string `yaml:"topic" json:"topic" default:"timer" validate:"nonzero"`
}

func validateJobs(v interface{}, param string) error {
	jobs, ok := v.([]Job)
	if !ok {
		return errors.New("only support job array")
	}
	for _, job := range jobs {
		for _, m := range job.Maps {
			if job.Encoding == JsonEncoding {
				if m.Field.Name == "" || m.Field.Type == "" {
					return fmt.Errorf("field name or type of map %+v shall not be empty when encoding is json", m)
				}
			} else if job.Encoding == BinaryEncoding {
				if m.Quantity == 0 {
					return fmt.Errorf("quantity of map %+v shall not be zero when encoding is binary", m)
				}
			}
			if _, ok := SysName[m.Field.Name]; ok {
				return fmt.Errorf("please use another name, '%s' is reserved by the system", m.Field.Name)
			}
			if _, ok := SysType[m.Field.Type]; !ok {
				return fmt.Errorf("unsupported field type: %s", m.Field.Type)
			}
		}
	}
	return nil
}

func (job *Job) SetDefaults() {
	var ms []MapConfig
	for _, m := range job.Maps {
		if job.Encoding == JsonEncoding {
			populateQuantityIfNeeds(&m)
		}
		ms = append(ms, m)
	}
	job.Maps = ms
}

func populateQuantityIfNeeds(cfg *MapConfig) {
	switch cfg.Field.Type {
	case Bool:
		cfg.Quantity = 1
	case Int16:
		cfg.Quantity = 1
	case UInt16:
		cfg.Quantity = 1
	case Int32:
		cfg.Quantity = 2
	case UInt32:
		cfg.Quantity = 2
	case Int64:
		cfg.Quantity = 4
	case UInt64:
		cfg.Quantity = 4
	case Float32:
		cfg.Quantity = 2
	case Float64:
		cfg.Quantity = 4
	default:
	}
}
