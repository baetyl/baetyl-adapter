package main

import (
	"github.com/baetyl/baetyl-adapter/opc"

	"github.com/baetyl/baetyl-go/v2/context"
)

func main() {
	// Running module in baetyl context
	context.Run(func(ctx context.Context) error {
		var cfg opc.Config
		// load custom config
		if err := ctx.LoadCustomConfig(&cfg); err != nil {
			return err
		}
		opc, err := opc.NewOpc(ctx, cfg)
		if err != nil {
			return err
		}
		defer opc.Close()
		ctx.Wait()
		return nil
	})
}
