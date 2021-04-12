package opcua

import (
	"context"

	dm "github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Device struct {
	info        *dm.DeviceInfo
	opcuaClient *opcua.Client
	cfg         DeviceConfig
}

func NewDevice(info *dm.DeviceInfo, cfg DeviceConfig) (*Device, error) {
	opts := []opcua.Option{
		opcua.RequestTimeout(cfg.Timeout),
		opcua.SecurityPolicy(cfg.Security.Policy),
		opcua.SecurityModeString(cfg.Security.Mode),
	}
	if cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		opts = append(opts,
			opcua.AuthUsername(cfg.Auth.Username, cfg.Auth.Password),
		)
	} else {
		cli := opcua.NewClient(cfg.Endpoint)
		var ctx, cancel = context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		if err := cli.Dial(ctx); err != nil {
			return nil, err
		}
		defer cli.Close()
		var res, err = cli.GetEndpoints()
		if err != nil {
			return nil, err
		}
		endpoints := res.Endpoints
		var ep = opcua.SelectEndpoint(endpoints, cfg.Security.Policy, ua.MessageSecurityModeFromString(cfg.Security.Mode))
		opts = append(opts,
			opcua.AuthAnonymous(),
			opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
		)
	}
	// TODO add certificate options

	// optimize timeout
	var ctx, cancel = context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	var client = opcua.NewClient(cfg.Endpoint, opts...)
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Trace(err)
	}
	return &Device{info: info, cfg: cfg, opcuaClient: client}, nil
}
