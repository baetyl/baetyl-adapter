module github.com/baetyl/baetyl-adapter/baetyl-adapter-modbus

go 1.13

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20191007211215-3e077fc8667a
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.1-0.20190307181833-2b18fe1d885e
)

require (
	github.com/256dpi/gomqtt v0.12.2
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/baetyl/baetyl v0.0.0-20191024053808-fa151b0276b9
	github.com/frankban/quicktest v1.5.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/tbrandon/mbserver v0.0.0-20170611213546-993e1772cc62
	gotest.tools v2.2.0+incompatible // indirect
)
