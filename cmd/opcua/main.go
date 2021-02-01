package main

import (
	"github.com/baetyl/baetyl-go/v2/context"

	"github.com/baetyl/baetyl-adapter/v2/opcua"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg opcua.Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		o, err := opcua.NewOpcua(ctx, cfg)
		if err != nil {
			return err
		}
		defer o.Close()
		ctx.Wait()
		return nil
	})
}
