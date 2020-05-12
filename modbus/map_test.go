package modbus

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/baetyl/baetyl-go/log"
	"github.com/stretchr/testify/assert"
)

func TestMapRead(t *testing.T) {
	server := MbSlave{}
	server.StartTCPSlave()
	slaveCfg := SlaveConfig{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client := NewClient(slaveCfg)
	client.Connect()
	slave := NewSlave(slaveCfg, client)
	log := log.With(log.Any("modbus", "map_test"))

	cfg1 := MapConfig{
		Function: 2,
		Address:  0,
		Quantity: 1,
	}
	m := NewMap(cfg1, slave, log)
	results, err := m.read()
	assert.NoError(t, err)
	expected1 := []byte{1}
	assert.Equal(t, results, expected1)

	cfg2 := MapConfig{
		Function: 4,
		Address:  0,
		Quantity: 2,
	}
	m = NewMap(cfg2, slave, log)

	results, err = m.read()
	assert.NoError(t, err)
	expected2 := []byte{144, 144}
	assert.Equal(t, results, expected2)

	client.WriteSingleCoil(0, 0xFF00)
	cfg3 := MapConfig{
		Function: 1,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(cfg3, slave, log)
	results, err = m.read()
	assert.NoError(t, err)
	expected3 := []byte{1}
	assert.Equal(t, results, expected3)

	client.WriteSingleRegister(0, 65535)
	cfg4 := MapConfig{
		Function: 3,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(cfg4, slave, log)
	results, err = m.read()
	assert.NoError(t, err)
	expected4 := []byte{255, 255}
	assert.Equal(t, results, expected4)

	cfg5 := MapConfig{
		Function: 3,
		Address:  0,
		Quantity: 0,
	}
	m = NewMap(cfg5, slave, log)
	results, err = m.read()
	assert.Error(t, err)
	server.Stop()
	client.Close()

	cfg6 := MapConfig{
		Function: 3,
		Address:  0,
		Quantity: 1,
	}
	m = NewMap(cfg6, slave, log)
	_, err = m.read()
	assert.Error(t, err)
}

func TestMapCollect(t *testing.T) {
	server := MbSlave{}
	server.StartTCPSlave()

	cfg := SlaveConfig{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client := NewClient(cfg)
	client.Connect()
	defer client.Close()
	slave := NewSlave(cfg, client)

	mapConfig := MapConfig{
		Function: 2,
		Address:  0,
		Quantity: 1,
	}

	log := log.With(log.Any("modbus", "test"))
	ma := NewMap(mapConfig, slave, log)

	results, err := ma.Collect()
	assert.NoError(t, err)
	expected := []byte{1}
	// address is 0(uint16) corresponding to []byte{0, 0}
	assert.Equal(t, results[:2], []byte{0, 0})
	// quantity is 1(uint16) corresponding to []byte{0, 1}
	assert.Equal(t, results[2:4], []byte{0, 1})
	// result should be 1 (slave have already set it to 1)
	assert.Equal(t, results[4:], expected)

	// invalid quantity
	mapConfig2 := MapConfig{
		Function: 2,
		Address:  0,
		Quantity: 0,
	}
	ma2 := NewMap(mapConfig2, slave, log)
	_, err = ma2.Collect()
	assert.Error(t, err)
	server.Stop()
}

func TestParse(t *testing.T) {
	m := NewMap(MapConfig{}, NewSlave(SlaveConfig{}, NewClient(SlaveConfig{})), log.With(log.Any("modbus", "test")))
	cfgs := []MapConfig{
		{Field: Field{Type: Bool}},
		{Field: Field{Type: Int16}},
		{Field: Field{Type: UInt16}},
		{Field: Field{Type: Int32}},
		{Field: Field{Type: UInt32}},
		{Field: Field{Type: Int64}},
		{Field: Field{Type: UInt64}},
		{Field: Field{Type: Float32}},
		{Field: Field{Type: Float64}},
		{Field: Field{Type: "string"}},
		{Function: Coil},
		{Function: Coil, Field: Field{Type: Int16}},
	}
	source := [][]byte{
		convertToByte(true),
		convertToByte(int16(1)),
		convertToByte(uint16(2)),
		convertToByte(int32(3)),
		convertToByte(uint32(4)),
		convertToByte(int64(5)),
		convertToByte(uint64(6)),
		convertToByte(float32(7)),
		convertToByte(float64(8)),
		convertToByte(""),
		{0, 1},
		convertToByte(false),
	}
	expected := []interface{}{
		true,
		int16(1),
		uint16(2),
		int32(3),
		uint32(4),
		int64(5),
		uint64(6),
		float32(7),
		float64(8),
		nil,
		nil,
		nil,
	}
	for i, src := range source {
		m.cfg = cfgs[i]
		res, err := m.Parse(src)
		assert.Equal(t, res, expected[i])
		if i <= 8 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func convertToByte(v interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}
