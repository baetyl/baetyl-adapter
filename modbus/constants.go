package modbus

const (
	BinaryEncoding = "binary"
	JsonEncoding = "json"
	SecondPrecision = "s"
	NanoPrecision = "ns"
	IntegerTime = "integer"
	StringTime = "string"
	Coil = 1
	DiscreteInput = 2
	HoldingRegister = 3
	InputRegister = 4
	SlaveId = "slaveid"
	SysTime = "time"

	Bool = "bool"
	Int16 = "int16"
	UInt16 = "uint16"
	Int32 = "int32"
	UInt32 = "uint32"
	Int64 = "int64"
	UInt64 = "uint64"
	Float32 = "float32"
	Float64 = "float64"
)

var SysType =  map[string]struct{}{
	Bool: {},
	Int16: {},
	UInt16: {},
	Int32: {},
	UInt32: {},
	Int64: {},
	UInt64: {},
	Float32: {},
	Float64: {},
}
var SysName =  map[string]struct{}{
	SysTime: {},
	SlaveId: {},
}
