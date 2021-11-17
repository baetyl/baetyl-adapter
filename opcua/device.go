package opcua

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

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
	var ep = opcua.SelectEndpoint(res.Endpoints, cfg.Security.Policy, ua.MessageSecurityModeFromString(cfg.Security.Mode))

	if cfg.Auth != nil && cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		opts = append(opts,
			opcua.AuthUsername(cfg.Auth.Username, cfg.Auth.Password),
			opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeUserName),
		)
	} else if cfg.Certificate != nil{
		cert, err := decodeCert([]byte(cfg.Certificate.Cert))
		if err != nil {
			return nil, err
		}
		key, err := decodeKey([]byte(cfg.Certificate.Key))
		if err != nil {
			return nil, err
		}
		opts = append(opts, opcua.AuthCertificate(cert),
			opcua.PrivateKey(key),
			opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeCertificate))
	} else {
		opts = append(opts,
			opcua.AuthAnonymous(),
			opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
		)
	}

	// optimize timeout
	ctx, cancel = context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	var client = opcua.NewClient(cfg.Endpoint, opts...)
	if err = client.Connect(ctx); err != nil {
		return nil, errors.Trace(err)
	}
	return &Device{info: info, cfg: cfg, opcuaClient: client}, nil
}

func decodeCert(certPEM []byte) ([]byte, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.Errorf("failed to decode cert")
	}
	return block.Bytes, nil
}

func decodeKey(keyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.Errorf("failed to decode key")
	}
	var pk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Errorf("failed to parse key")
	}
	return pk, nil
}

