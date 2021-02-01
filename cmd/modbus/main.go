package main

import (
	"github.com/baetyl/baetyl-go/v2/dmcontext"

	"github.com/baetyl/baetyl-adapter/v2/modbus"
)

func main() {
	// Running module in baetyl context
	dmcontext.Run(func(ctx dmcontext.Context) error {
		var cfg modbus.Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		modbus, err := modbus.NewModbus(ctx, cfg)
		if err != nil {
			return err
		}
		modbus.Start()
		defer modbus.Close()
		ctx.Wait()
		return nil
	})
}
