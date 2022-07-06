package dmp

import (
	"encoding/json"
	"strconv"

	"github.com/baetyl/baetyl-go/v2/dmcontext"
	"github.com/baetyl/baetyl-go/v2/errors"
)

const (
	BIEKey   = "bie"
	DMPKey   = "dmp"
	Method   = "thing.event.post"
	Version  = "1.0"
	BindName = "MAIN"
)

type DMP struct {
	ReqId     string                 `yaml:"reqId,omitempty" json:"reqId,omitempty"`
	Method    string                 `yaml:"method,omitempty" json:"method,omitempty"`
	Version   string                 `yaml:"version,omitempty" json:"version,omitempty"`
	Timestamp int64                  `yaml:"timestamp,omitempty" json:"timestamp,omitempty"`
	BindName  string                 `yaml:"bindName,omitempty" json:"bindName,omitempty"`
	Events    map[string]interface{} `yaml:"events,omitempty" json:"events,omitempty"`
}

func GetMappingName(id string, template *dmcontext.AccessTemplate) (string, error) {
	var name string

	for _, deviceProperty := range template.Properties {
		if id == deviceProperty.Id {
			name = deviceProperty.Name
			break
		}
	}
	if name == "" {
		return "", errors.New("unknown property id")
	}
	return name, nil
}

func GetConfigIdByModelName(name string, template *dmcontext.AccessTemplate) (string, error) {
	for _, modelMapping := range template.Mappings {
		if modelMapping.Attribute == name {
			ids, err := dmcontext.ParseExpression(modelMapping.Expression)
			if err != nil {
				return "", err
			}
			if len(ids) > 0 {
				return ids[0][1:], nil
			}
		}
	}
	return "", errors.New("config id not exist")
}

func GetPropValueByModelName(name string, val interface{}, template *dmcontext.AccessTemplate) (float64, error) {
	for _, modelMapping := range template.Mappings {
		if modelMapping.Attribute == name {
			value, err := parseValueToFloat64(val)
			if err != nil {
				return 0.0, err
			}
			if modelMapping.Type == dmcontext.MappingValue {
				return value, nil
			}
			propVal, err := dmcontext.SolveExpression(modelMapping.Expression, value)
			if err != nil {
				return 0.0, err
			}
			return propVal, nil
		}
	}
	return 0.0, errors.New("prop value not exist")
}

func parseValueToFloat64(v interface{}) (float64, error) {
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case int16:
		return float64(v.(int16)), nil
	case int32:
		return float64(v.(int32)), nil
	case int64:
		return float64(v.(int64)), nil
	case float32:
		s := strconv.FormatFloat(float64(v.(float32)), 'e', -1, 32)
		return strconv.ParseFloat(s, 64)
	case float64:
		return v.(float64), nil
	default:
		return 0, errors.New("unsupported value type")
	}
}

func ParsePropertyValue(tpy string, val interface{}) (interface{}, error) {
	// it is json.Number (string actually) when val is number
	switch tpy {
	case dmcontext.TypeInt16:
		num, _ := val.(json.Number)
		i, err := strconv.ParseInt(num.String(), 10, 16)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return int16(i), nil
	case dmcontext.TypeInt32:
		num, _ := val.(json.Number)
		i, err := strconv.ParseInt(num.String(), 10, 32)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return int32(i), nil
	case dmcontext.TypeInt64:
		num, _ := val.(json.Number)
		i, err := strconv.ParseInt(num.String(), 10, 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return i, nil
	case dmcontext.TypeFloat32:
		num, _ := val.(json.Number)
		f, err := strconv.ParseFloat(num.String(), 32)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return float32(f), nil
	case dmcontext.TypeFloat64:
		num, _ := val.(json.Number)
		f, err := strconv.ParseFloat(num.String(), 64)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return f, nil
	case dmcontext.TypeBool, dmcontext.TypeString:
		return val, nil
	default:
		return nil, errors.Trace(dmcontext.ErrTypeNotSupported)
	}
}
