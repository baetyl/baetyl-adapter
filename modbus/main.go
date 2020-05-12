package main

import (
	"github.com/baetyl/baetyl-go/context"
	"github.com/baetyl/baetyl-go/mqtt"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		mqttCfg := ctx.ServiceConfig().MQTT
		if mqttCfg.MaxCacheMessages < len(cfg.Jobs)*2 {
			mqttCfg.MaxCacheMessages = len(cfg.Jobs) * 2
		}
		if mqttCfg.ClientID == "" {
			mqttCfg.ClientID = ctx.ServiceName()
		}
		option, err := mqttCfg.ToClientOptions(nil)
		if err != nil {
			return err
		}
		sender := mqttSender{publish: cfg.Publish, Client: mqtt.NewClient(*option)}
		defer sender.Close()
		modbus, err := NewModbus(ctx, cfg, sender)
		if err != nil {
			return err
		}
		defer modbus.Close()
		ctx.Wait()
		return nil
	})
}
