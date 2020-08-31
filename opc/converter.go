package opc

import (
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/gopcua/opcua/ua"
)

func parse(val *ua.Variant) (interface{}, error) {
	switch val.Type() {
	case ua.TypeIDBoolean:
		return val.Bool(), nil
	case ua.TypeIDFloat:
		return val.Float(), nil
	case ua.TypeIDDouble:
		return val.Float(), nil
	case ua.TypeIDInt16, ua.TypeIDInt32, ua.TypeIDInt64:
		return val.Int(), nil
	case ua.TypeIDUint16, ua.TypeIDUint32, ua.TypeIDUint64:
		return val.Uint(), nil
	case ua.TypeIDByteString:
		return val.ByteString(), nil
	default:
		return nil, errors.Errorf("unsupported type")
	}
}

func validateAndTransform(source interface{}, fieldType string) (*ua.Variant, error) {
	var value interface{}
	var ok bool
	var num float64
	if fieldType != Bool {
		num, ok = source.(float64)
	}
	switch fieldType {
	case Bool:
		var b bool
		b, ok = source.(bool)
		value = b
	case Int16:
		i16 := int16(num)
		value = i16
	case UInt16:
		u16 := uint16(num)
		value = u16
	case Int32:
		i32 := int32(num)
		value = i32
	case UInt32:
		u32 := uint32(num)
		value = u32
	case Int64:
		i64 := int64(num)
		value = i64
	case UInt64:
		u64 := uint64(num)
		value = u64
	case Float32:
		f32 := float32(num)
		value = f32
	case Float64:
		value = num
	default:
		return nil, errors.Errorf("unsupported field type [%s]", fieldType)
	}
	if !ok {
		return nil, errors.Errorf("value [%v] not compatible with type [%s] ", source, fieldType)
	}
	res, err := ua.NewVariant(value)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}
