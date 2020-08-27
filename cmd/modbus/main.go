package main

import (
	"github.com/baetyl/baetyl-adapter/modbus"
	"github.com/baetyl/baetyl-go/v2/context"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg modbus.Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		modbus, err := modbus.NewModbus(ctx, cfg)
		if err != nil {
			return err
		}
		defer modbus.Close()
		ctx.Wait()
		return nil
	})
}
