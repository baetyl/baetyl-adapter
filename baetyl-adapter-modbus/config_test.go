package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	confString := `
slaves:
- id: 1
  address: tcp://127.0.0.1:502
jobs:
- slaveid: 1
  interval: 3s
  kind: binary
  maps:
  - address: 0
    quantity: 1
    function: 3
publish:
  topic: test`
	dir, err := ioutil.TempDir("", "template")
	assert.NoError(t, err)
	fileName := "conf"
	ioutil.WriteFile(filepath.Join(dir, fileName), []byte(confString), 0755)
	var cfg Config
	utils.LoadYAML(filepath.Join(dir, fileName), &cfg)
	cfg2 := Config{
		Slaves: []SlaveConfig{
			{
				ID:          1,
				Address:     "tcp://127.0.0.1:502",
				Timeout:     10 * time.Second,
				IdleTimeout: 1 * time.Minute,
				BaudRate:    19200,
				DataBits:    8,
				StopBits:    1,
				Parity:      "E",
			},
		},
		Jobs: []Job{{
			SlaveId: 1,
			Interval: 3 * time.Second,
			Kind: BinaryKind,
			Time: Time{
				Type:      LongTimeType,
				Format:    "2006-01-02 15:04:05",
				Precision: "second",
			},
			Maps: []MapConfig{{
				Address:  0,
				Quantity: 1,
				Function: 3,
			},},
		},},
		Publish: Publish{
			QOS:   0,
			Topic: "test",
		},
	}
	assert.Equal(t, cfg, cfg2)
}
