package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/baetyl/baetyl/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	confString := `
slaves:
  - id: 1
    address: tcp://127.0.0.1:502
    interval: 3s
tables:
  - slaveid: 1
    address: 0
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
		Slaves: []SlaveItem{
			{
				ID:          1,
				Address:     "tcp://127.0.0.1:502",
				Interval:    3 * time.Second,
				Timeout:     10 * time.Second,
				IdleTimeout: 1 * time.Minute,
				BaudRate:    19200,
				DataBits:    8,
				StopBits:    1,
				Parity:      "E",
			},
		},
		Tables: []MapItem{
			{
				SlaveID:  1,
				Address:  0,
				Quantity: 1,
				Function: 3,
			},
		},
		Publish: &Publish{
			QOS:   0,
			Topic: "test",
		},
	}
	assert.Equal(t, cfg, cfg2)
}
