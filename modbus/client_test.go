package main

import (
	"context"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"github.com/tbrandon/mbserver"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	slave := MbSlave{}
	slave.StartTCPSlave()

	cfg := SlaveConfig{}
	client := NewClient(cfg)
	err := client.Connect()
	assert.Error(t, err)

	cfg = SlaveConfig{
		ID:      1,
		Address: "tcp://127.0.0.1:50200",
	}
	client = NewClient(cfg)
	err = client.Connect()
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)

	cfg = SlaveConfig{
		ID:      2,
		Address: "tcp://127.0.0.1:50201",
	}
	client = NewClient(cfg)
	err = client.Connect()
	assert.Error(t, err)
	slave.Stop()
}

type MbSlave struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	started chan bool
}

func (m *MbSlave) StartTCPSlave() {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.started = make(chan bool)
	m.wg.Add(1)
	go m.startTCP()
	<-m.started
}

func (m *MbSlave) startTCP() error {
	server := mbserver.NewServer()
	err := server.ListenTCP("0.0.0.0:50200")
	if err != nil {
		return err
	}

	server.RegisterFunctionHandler(2,
		func(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
			data := frame.GetData()
			register := int(binary.BigEndian.Uint16(data[0:2]))
			numRegs := int(binary.BigEndian.Uint16(data[2:4]))
			endRegister := register + numRegs
			if endRegister > 65535 {
				return []byte{}, &mbserver.IllegalDataAddress
			}
			dataSize := numRegs / 8
			if (numRegs % 8) != 0 {
				dataSize++
			}
			data = make([]byte, 1+dataSize)
			data[0] = byte(dataSize)
			for i := range s.DiscreteInputs[register:endRegister] {
				shift := uint(i) % 8
				data[1+i/8] |= byte(1 << shift)
			}
			return data, &mbserver.Success
		})

	server.RegisterFunctionHandler(4,
		func(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
			data := frame.GetData()
			register := int(binary.BigEndian.Uint16(data[0:2]))
			numRegs := int(binary.BigEndian.Uint16(data[2:4]))
			endRegister := register + numRegs
			if endRegister > 65535 {
				return []byte{}, &mbserver.IllegalDataAddress
			}
			data = make([]byte, 1+numRegs)
			data[0] = byte(numRegs)
			for i := range s.InputRegisters[register:endRegister] {
				data[1+i] = 144
			}
			return data, &mbserver.Success
		})
	m.started <- true
	for {
		select {
		case <-m.ctx.Done():
			m.wg.Done()
			server.Close()
			return nil
		}
	}
}

func (m *MbSlave) Stop() {
	m.cancel()
	m.wg.Wait()
	close(m.started)
}
