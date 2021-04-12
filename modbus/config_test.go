package modbus

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "modbus")
	assert.NoError(t, err)
	fileName := "conf"
	var cfg Config
	confString := `
slaves:
- id: 1
  mode: tcp
  address: 127.0.0.1:502
jobs:
- slaveid: 1
  interval: 3s
  maps:
  - function: 1
    address: 1
    name: a
    type: bool
  - function: 3
    address: 2
    name: b
    type: int16
  - function: 3
    address: 3
    name: c
    type: uint16
  - function: 3
    address: 4
    name: d
    type: int32
  - function: 3
    address: 5
    name: e
    type: uint32
  - function: 3
    address: 6
    name: f
    type: int64
  - function: 3
    address: 7
    name: g
    type: uint64
  - function: 3
    address: 8
    name: h
    type: float32
  - function: 3
    address: 9
    name: i
    type: float64
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	cfg2 := Config{
		Slaves: []SlaveConfig{{
			Id:          1,
			Mode:        "tcp",
			Address:     "127.0.0.1:502",
			Timeout:     10 * time.Second,
			IdleTimeout: 1 * time.Minute,
			BaudRate:    19200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "E",
		}},
		Jobs: []Job{{
			SlaveID:  1,
			Interval: 3 * time.Second,
			Maps: []MapConfig{
				{
					Address:  1,
					Quantity: 1,
					Function: 1,
					Name:     "a",
					Type:     Bool,
				},
				{
					Address:  2,
					Quantity: 1,
					Function: 3,
					Name:     "b",
					Type:     Int16,
				},
				{
					Address:  3,
					Quantity: 1,
					Function: 3,
					Name:     "c",
					Type:     UInt16,
				},
				{
					Address:  4,
					Quantity: 2,
					Function: 3,
					Name:     "d",
					Type:     Int32,
				},
				{
					Address:  5,
					Quantity: 2,
					Function: 3,
					Name:     "e",
					Type:     UInt32,
				},
				{
					Address:  6,
					Quantity: 4,
					Function: 3,
					Name:     "f",
					Type:     Int64,
				},
				{
					Address:  7,
					Quantity: 4,
					Function: 3,
					Name:     "g",
					Type:     UInt64,
				},
				{
					Address:  8,
					Quantity: 2,
					Function: 3,
					Name:     "h",
					Type:     Float32,
				},
				{
					Address:  9,
					Quantity: 4,
					Function: 3,
					Name:     "i",
					Type:     Float64,
				},
			},
		}},
	}
	assert.Equal(t, cfg, cfg2)

	// encoding is json and field is empty
	confString = `
jobs:
- slaveid: 1
  encoding: json
  maps:
  - address: 1
    function: 1
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	err = utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	assert.Error(t, err)

	// encoding is binary and quantity not configured
	confString = `
jobs:
- slaveid: 1
  encoding: binary
  maps:
  - address: 2
    function: 3
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	err = utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	assert.Error(t, err)

	//field name is system time name
	confString = `
jobs:
- slaveid: 1
  encoding: json
  maps:
  - address: 1
    function: 3
    name: time
    type: int16
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	err = utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	assert.Error(t, err)

	//field type is not required type
	confString = `
jobs:
- slaveid: 1
  interval: 3s
  encoding: json
  maps:
  - address: 1
    function: 1
    name: a
    type: string
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	err = utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	assert.Error(t, err)
}
