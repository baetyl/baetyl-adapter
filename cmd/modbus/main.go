package main

import (
	"fmt"
	"strconv"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/utils"

	"github.com/baetyl/baetyl-adapter/v2/modbus"
)

func main() {
	// Running module in baetyl context
	dm.Run(func(ctx dm.Context) error {
		cfg, err := genConfig(ctx)
		if err != nil {
			return err
		}
		mod, err := modbus.NewModbus(ctx, *cfg)
		if err != nil {
			return err
		}
		mod.Start()
		defer mod.Close()
		ctx.Wait()
		return nil
	})
}

func genConfig(ctx dm.Context) (*modbus.Config, error) {
	cfg := &modbus.Config{}
	devProps := ctx.GetPropertiesConfig()
	var slaves []modbus.SlaveConfig
	var jobs []modbus.Job
	for name, acc := range ctx.GetAccessConfig() {
		if acc.Modbus == nil {
			continue
		}
		slave := modbus.SlaveConfig{
			Device:      name,
			Id:          acc.Modbus.Id,
			Timeout:     acc.Modbus.Timeout,
			IdleTimeout: acc.Modbus.IdleTimeout,
		}
		if tcp := acc.Modbus.Tcp; tcp != nil {
			slave.Mode = string(modbus.ModeTcp)
			slave.Address = fmt.Sprintf("%s:%d", tcp.Address, tcp.Port)
		} else if rtu := acc.Modbus.Rtu; rtu != nil {
			slave.Mode = string(modbus.ModeRtu)
			slave.Address = rtu.Port
			slave.BaudRate = rtu.BaudRate
			slave.DataBits = rtu.DataBit
			slave.StopBits = rtu.StopBit
			slave.Parity = rtu.Parity
		}
		slaves = append(slaves, slave)
		var jobMaps []modbus.MapConfig
		for _, prop := range devProps[name] {
			if visitor := prop.Visitor.Modbus; visitor != nil {
				address, _ := strconv.ParseUint(visitor.Address[2:], 10, 16)
				m := modbus.MapConfig{
					Name:         prop.Name,
					Type:         visitor.Type,
					Function:     visitor.Function,
					Address:      uint16(address),
					Quantity:     visitor.Quantity,
					SwapRegister: visitor.SwapRegister,
					SwapByte:     visitor.SwapByte,
				}
				jobMaps = append(jobMaps, m)
			}
		}
		job := modbus.Job{
			SlaveID:  acc.Modbus.Id,
			Interval: acc.Modbus.Interval,
			Maps:     jobMaps,
		}
		jobs = append(jobs, job)
	}
	cfg.Jobs = jobs
	cfg.Slaves = slaves
	if err := utils.SetDefaults(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
