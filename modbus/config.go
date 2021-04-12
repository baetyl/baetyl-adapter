package modbus

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/validator.v2"
)

func init() {
	validator.SetValidationFunc("validjobs", validateJobs)
}

// Config custom configuration of the timer module
type Config struct {
	// Slaves slave list
	Slaves []SlaveConfig `yaml:"slaves" json:"slaves"`
	// Jobs job list
	Jobs []Job `yaml:"jobs" json:"jobs" validate:"validjobs"`
}

type Job struct {
	// SlaveID slave id defined in slaves
	SlaveID byte `yaml:"slaveid" json:"slaveid"`
	// Interval the interval between task execution
	Interval time.Duration `yaml:"interval" json:"interval" default:"5s"`
	// Maps definition of data points
	Maps []MapConfig `yaml:"maps" json:"maps"`
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
	Device string `yaml:"device" json:"device"`
	// ID slave id
	Id byte `yaml:"id" json:"id"`
	// Mode mode of connecting
	Mode string `yaml:"mode" json:"mode" default:"rtu" validate:"regexp=^(tcp|rtu)?$"`
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
}

// MapConfig map point configuration
type MapConfig struct {
	// Name name of map config
	Name string `yaml:"name" json:"name"`
	// Type type of map type
	Type string `yaml:"type" json:"type"`
	// Function
	Function byte `yaml:"function" json:"function" validate:"min=1, max=4" validate:"nonzero"`
	// Address
	Address uint16 `yaml:"address" json:"address"`
	// Quantity
	Quantity uint16 `yaml:"quantity" json:"quantity"`
	// SwapByte whether swap byte, meaning using big endian or little endian
	SwapByte bool `yaml:"swapByte" json:"swapByte"`
	// SwapRegister whether swap high and low register
	SwapRegister bool `yaml:"swapRegister" json:"swapRegister"`
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
			if _, ok := SysName[m.Name]; ok {
				return fmt.Errorf("please use another name, '%s' is reserved by the system", m.Name)
			}
			if _, ok := SysType[m.Type]; !ok {
				return fmt.Errorf("unsupported field type: %s", m.Type)
			}
		}
	}
	return nil
}

func (job *Job) SetDefaults() {
	var ms []MapConfig
	for _, m := range job.Maps {
		PopulateQuantityIfNeeds(&m)
		ms = append(ms, m)
	}
	job.Maps = ms
}

func PopulateQuantityIfNeeds(cfg *MapConfig) {
	switch cfg.Type {
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
