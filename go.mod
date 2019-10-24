module github.com/baetyl/baetyl-adapter

go 1.13

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20191007211215-3e077fc8667a
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.1-0.20190307181833-2b18fe1d885e
)

require (
	github.com/256dpi/gomqtt v0.12.2
	github.com/baetyl/baetyl v0.0.0-20191023143945-4d673de16a40
	github.com/baidubce/bce-sdk-go v0.0.0-20191012060435-0868fe1d4ceb
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/magiconair/properties v1.8.0
	github.com/pkg/errors v0.8.1
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135
)
