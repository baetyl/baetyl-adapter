package opcua

import (
	"time"
)

type Config struct {
	Devices []DeviceConfig `yaml:"devices" json:"devices"`
	Jobs    []Job          `yaml:"jobs" json:"jobs"`
}

type DeviceConfig struct {
	Device      string        `yaml:"device" json:"device"`
	Endpoint    string        `yaml:"endpoint" json:"endpoint"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout" default:"10s"`
	Security    Security      `yaml:"security" json:"security"`
	Auth        Auth          `yaml:"auth" json:"auth"`
	Certificate Certificate   `yaml:"certificate" json:"certificate"`
}

type Security struct {
	Policy string `yaml:"policy" json:"policy"`
	Mode   string `yaml:"mode" json:"mode"`
}

type Auth struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type Certificate struct {
	Cert string `yaml:"certFile" json:"certFile"`
	Key  string `yaml:"keyFile" json:"keyFile"`
}

type Job struct {
	Device     string        `yaml:"device" json:"device"`
	Time       Time          `yaml:"time" json:"time" default:"{\"name\":\"time\", \"type\":\"integer\"}"`
	Interval   time.Duration `yaml:"interval" json:"interval" default:"20s"`
	Properties []Property    `yaml:"properties" json:"properties"`
	Publish    Publish       `yaml:"publish" json:"publish"`
}

type Time struct {
	Name      string `yaml:"name" json:"name"`
	Type      string `yaml:"type" json:"type"`
	Format    string `yaml:"format" json:"format" default:"2006-01-02 15:04:05"`
	Precision string `yaml:"precision" json:"precision" default:"s" validate:"regexp=^(s|ns)?$"`
}

type Property struct {
	Name   string `yaml:"name" json:"name"`
	Type   string `yaml:"type" json:"type"`
	NodeID string `yaml:"nodeid" json:"nodeid"`
}

type Publish struct {
	QOS   uint32 `yaml:"qos" json:"qos" validate:"min=0, max=1"`
	Topic string `yaml:"topic" json:"topic" default:"timer" validate:"nonzero"`
}
