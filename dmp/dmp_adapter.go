package dmp

import (
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
