package main

import (
	"encoding/binary"
	"time"

	"github.com/baetyl/baetyl/utils"
	"github.com/tbrandon/mbserver"
)

type MbSlave struct {
	tomb utils.Tomb
}

func (m *MbSlave) StartTCPSlave() error {
	err := m.tomb.Go(m.startTCP)
	time.Sleep(10 * time.Millisecond)
	return err
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

	defer server.Close()

	for {
		select {
		case <-m.tomb.Dying():
			return nil
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (m *MbSlave) Stop() {
	m.tomb.Kill(nil)
	m.tomb.Wait()
}
