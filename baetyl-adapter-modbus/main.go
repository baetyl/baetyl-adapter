package main

import (
	"github.com/baetyl/baetyl-go/context"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		if err := validate(cfg); err != nil {
			return err
		}
		setDefault(&cfg, ctx)
		modbus, err := NewModbus(ctx, cfg)
		if err != nil {
			return err
		}
		defer modbus.Close()
		ctx.Wait()
		return nil
	})
}
