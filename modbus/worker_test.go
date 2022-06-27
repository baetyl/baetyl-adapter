package modbus

import (
	"testing"
)

func TestExecute(t *testing.T) {
	//server := MbSlave{}
	//server.StartTCPSlave()
	//mockCtl := gomock.NewController(t)
	//ms := mock.NewMockSender(mockCtl)
	//
	//slaveCfg := SlaveConfig{
	//	Id:      1,
	//	Mode:    string(ModeTcp),
	//	Address: "tcp://127.0.0.1:50200",
	//}
	//client, err := NewClient(slaveCfg)
	//assert.NoError(t, err)
	//client.Connect()
	//slave := NewSlave(slaveCfg, client)
	//log := log.With(log.Any("modbus", "worker_test"))
	//job := Job{
	//	SlaveID:  1,
	//	Maps: []MapConfig{{
	//		Function: 4,
	//		Address:  0,
	//		Quantity: 2,
	//		Field: Field{
	//			Name: "a",
	//			Type: Int16,
	//		},
	//	}},
	//}
	//defaults.Set(job)
	//w := NewWorker(job, slave, ms, log)
	//ms.EXPECT().Send(gomock.Any()).Return(nil)
	//err = w.Execute()
	//assert.NoError(t, err)
	//
	//job.Encoding = JsonEncoding
	//w = NewWorker(job, slave, ms, log)
	//ms.EXPECT().Send(gomock.Any()).Return(nil)
	//err = w.Execute()
	//assert.NoError(t, err)
	//
	//ms.EXPECT().Send(gomock.Any()).Return(errors.New("send error"))
	//err = w.Execute()
	//assert.Error(t, err)
	//server.Stop()
}
