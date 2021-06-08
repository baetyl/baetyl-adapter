package main

import (
	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-adapter/v2/opcua"
)

func main() {
	// Running module in baetyl context
	dm.Run(func(ctx dm.Context) error {
		cfg, err := genConfig(ctx)
		if err != nil {
			return err
		}
		o, err := opcua.NewOpcua(ctx, *cfg)
		if err != nil {
			return err
		}
		defer o.Close()
		ctx.Wait()
		return nil
	})
}

func genConfig(ctx dm.Context) (*opcua.Config, error) {
	cfg := &opcua.Config{}
	devProps := ctx.GetPropertiesConfig()
	var devices []opcua.DeviceConfig
	var jobs []opcua.Job
	for name, acc := range ctx.GetAccessConfig() {
		if acc.Opcua == nil {
			continue
		}
		dev := opcua.DeviceConfig{Device: name}
		if err := copier.Copy(&dev, acc.Opcua); err != nil {
			return nil, err
		}
		devices = append(devices, dev)
		var jobProps []opcua.Property
		for _, prop := range devProps[name] {
			if visitor := prop.Visitor.Opcua; visitor != nil {
				jobProps = append(jobProps, opcua.Property{
					Name:   prop.Name,
					Type:   visitor.Type,
					NodeID: visitor.NodeID,
				})
			}
		}
		job := opcua.Job{
			Device:     name,
			Interval:   acc.Opcua.Interval,
			Properties: jobProps,
		}
		jobs = append(jobs, job)
	}
	cfg.Devices = devices
	cfg.Jobs = jobs
	if err := utils.SetDefaults(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
