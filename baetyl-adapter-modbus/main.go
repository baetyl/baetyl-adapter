package main

import (
	"github.com/baetyl/baetyl-go/context"
	"github.com/baetyl/baetyl-go/mqtt"
	uuid "github.com/satori/go.uuid"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg Config
		// load custom config
		err := ctx.LoadCustomConfig(&cfg)
		if err != nil {
			return err
		}

		mqttConfig := ctx.ServiceConfig().MQTT
		if mqttConfig.MaxCacheMessages < len(cfg.Jobs) * 2 {
			mqttConfig.MaxCacheMessages = len(cfg.Jobs) * 2
		}
		if mqttConfig.ClientID == "" {
			mqttConfig.ClientID = uuid.NewV4().String()
		}
		option, err := mqttConfig.ToClientOptions(nil)
		if err != nil {
			return err
		}
		mqtt := mqtt.NewClient(*option)
		defer mqtt.Close()

		log := ctx.Log()
		modbus := NewModbus(cfg, mqtt, log)
		modbus.Start(ctx)
		return modbus.Close()
	})
}
