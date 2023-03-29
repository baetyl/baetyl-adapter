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
	var slaves []modbus.SlaveConfig
	var jobs []modbus.Job

	for _, deviceInfo := range ctx.GetAllDevices() {
		accessConfig := deviceInfo.AccessConfig
		if accessConfig.Modbus == nil {
			continue
		}
		slave := modbus.SlaveConfig{
			Device:      deviceInfo.Name,
			Id:          accessConfig.Modbus.Id,
			Timeout:     accessConfig.Modbus.Timeout,
			IdleTimeout: accessConfig.Modbus.IdleTimeout,
		}
		if tcp := accessConfig.Modbus.Tcp; tcp != nil {
			slave.Mode = string(modbus.ModeTcp)
			slave.Address = fmt.Sprintf("%s:%d", tcp.Address, tcp.Port)
		} else if rtu := accessConfig.Modbus.Rtu; rtu != nil {
			slave.Mode = string(modbus.ModeRtu)
			slave.Address = rtu.Port
			slave.BaudRate = rtu.BaudRate
			slave.DataBits = rtu.DataBit
			slave.StopBits = rtu.StopBit
			slave.Parity = rtu.Parity
		}
		slaves = append(slaves, slave)

		var jobMaps []modbus.MapConfig
		deviceTemplate, _ := ctx.GetAccessTemplates(&deviceInfo)
		if deviceTemplate != nil && deviceTemplate.Properties != nil && len(deviceTemplate.Properties) > 0 {
			for _, prop := range deviceTemplate.Properties {
				if visitor := prop.Visitor.Modbus; visitor != nil {
					address, _ := strconv.ParseUint(visitor.Address[2:], 16, 16)
					m := modbus.MapConfig{
						Id:           prop.Id,
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
		}
		job := modbus.Job{
			SlaveID:  accessConfig.Modbus.Id,
			Interval: accessConfig.Modbus.Interval,
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
