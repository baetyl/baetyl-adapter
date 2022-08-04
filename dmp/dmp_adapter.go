package dmp

import (
    "strconv"

    "github.com/baetyl/baetyl-go/v2/dmcontext"
    "github.com/baetyl/baetyl-go/v2/errors"
)

const (
    BLinkKey      = "blink"
    Method        = "thing.property.post"
    Version       = "1.0"
)

type DMP struct {
    ReqId      string                 `yaml:"reqId,omitempty" json:"reqId,omitempty"`
    Method     string                 `yaml:"method,omitempty" json:"method,omitempty"`
    Version    string                 `yaml:"version,omitempty" json:"version,omitempty"`
    Timestamp  int64                  `yaml:"timestamp,omitempty" json:"timestamp,omitempty"`
    Properties map[string]interface{} `yaml:"properties,omitempty" json:"properties,omitempty"`
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

func GetPropValueByModelName(name string, val interface{}, template *dmcontext.AccessTemplate) (interface{}, error) {
    for _, modelMapping := range template.Mappings {
        if modelMapping.Attribute == name {
            if modelMapping.Type == dmcontext.MappingValue {
                return val, nil
            }
            value, err := parseValueToFloat64(val)
            if err != nil {
                return nil, err
            }
            propVal, err := dmcontext.SolveExpression(modelMapping.Expression, value)
            if err != nil {
                return nil, err
            }
            return propVal, nil
        }
    }
    return nil, errors.New("prop value not exist")
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

func ParsePropertyValue(tpy string, val float64) (interface{}, error) {
    switch tpy {
    case dmcontext.TypeInt16:
        return int16(val), nil
    case dmcontext.TypeInt32:
        return int32(val), nil
    case dmcontext.TypeInt64:
        return int64(val), nil
    case dmcontext.TypeFloat32:
        return float32(val), nil
    case dmcontext.TypeFloat64:
        return val, nil
    default:
        return nil, errors.Trace(dmcontext.ErrTypeNotSupported)
    }
}
