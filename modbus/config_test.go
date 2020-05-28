package modbus

import (
	"github.com/baetyl/baetyl-go/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "modbus")
	assert.NoError(t, err)
	fileName := "conf"
	var cfg Config
	confString := `
slaves:
- id: 1
  address: tcp://127.0.0.1:502
jobs:
- slaveid: 1
  interval: 3s
  encoding: json
  maps:
  - function: 1
    address: 1
    field:
      name: a
      type: bool
  - function: 3
    address: 2
    field:
      name: b
      type: int16
  - function: 3
    address: 3
    field:
      name: c
      type: uint16
  - function: 3
    address: 4
    field:
      name: d
      type: int32
  - function: 3
    address: 5
    field:
      name: e
      type: uint32
  - function: 3
    address: 6
    field:
      name: f
      type: int64
  - function: 3
    address: 7
    field:
      name: g
      type: uint64
  - function: 3
    address: 8
    field:
      name: h
      type: float32
  - function: 3
    address: 9
    field:
      name: i
      type: float64
publish:
  topic: test`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	cfg2 := Config{
		Slaves: []SlaveConfig{{
			ID:          1,
			Address:     "tcp://127.0.0.1:502",
			Timeout:     10 * time.Second,
			IdleTimeout: 1 * time.Minute,
			BaudRate:    19200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "E",
		}},
		Jobs: []Job{{
			SlaveId:  1,
			Interval: 3 * time.Second,
			Encoding: JsonEncoding,
			Time: Time{
				Field: Field{
					Name: SysTime,
					Type: IntegerTime,
				},
				Format:    "2006-01-02 15:04:05",
				Precision: "s",
			},
			Maps: []MapConfig{
				{
					Address:  1,
					Quantity: 1,
					Function: 1,
					Field:    Field{Name: "a", Type: Bool},
				},
				{
					Address:  2,
					Quantity: 1,
					Function: 3,
					Field:    Field{Name: "b", Type: Int16},
				},
				{
					Address:  3,
					Quantity: 1,
					Function: 3,
					Field:    Field{Name: "c", Type: UInt16},
				},
				{
					Address:  4,
					Quantity: 2,
					Function: 3,
					Field:    Field{Name: "d", Type: Int32},
				},
				{
					Address:  5,
					Quantity: 2,
					Function: 3,
					Field:    Field{Name: "e", Type: UInt32},
				},
				{
					Address:  6,
					Quantity: 4,
					Function: 3,
					Field:    Field{Name: "f", Type: Int64},
				},
				{
					Address:  7,
					Quantity: 4,
					Function: 3,
					Field:    Field{Name: "g", Type: UInt64},
				},
				{
					Address:  8,
					Quantity: 2,
					Function: 3,
					Field:    Field{Name: "h", Type: Float32},
				},
				{
					Address:  9,
					Quantity: 4,
					Function: 3,
					Field:    Field{Name: "i", Type: Float64},
				},
			},
		}},
		Publish: Publish{
			QOS:   0,
			Topic: "test",
		},
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
    field:
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
    field:
      name: a
      type: string
`
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	err = utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	assert.Error(t, err)
}
