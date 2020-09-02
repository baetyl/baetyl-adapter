package opcua

import (
	"encoding/json"

	"github.com/256dpi/gomqtt/packet"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/mqtt"
	"github.com/gopcua/opcua/ua"
)

type observer struct {
	devices map[byte]*Device
	log     *log.Logger
}

type CtrData struct {
	DeviceID   byte                   `yaml:"deviceid" json:"deviceid"`
	Attributes map[string]interface{} `yaml:"attr" json:"attr"`
}

func NewObserver(devices map[byte]*Device, log *log.Logger) mqtt.Observer {
	return &observer{
		devices: devices,
		log:     log,
	}
}

func (o *observer) OnPublish(pkt *packet.Publish) error {
	var ctrData CtrData
	err := json.Unmarshal(pkt.Message.Payload, &ctrData)
	if err != nil {
		return errors.Trace(err)
	}
	var device *Device
	device, ok := o.devices[ctrData.DeviceID]
	if !ok {
		o.log.Error("device to write data not exist", log.Any("id", ctrData.DeviceID))
		return errors.Errorf("device to write data not exist")
	}
	if err := o.Write(device, ctrData.Attributes); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (o *observer) Write(device *Device, attr map[string]interface{}) error {
	config, ok := configRecoder[device.cfg.ID]
	if !ok {
		o.log.Error("map config of device id not exist", log.Any("id", device.cfg.ID))
		return errors.Errorf("map config of device id [%d] not exist", device.cfg.ID)
	}
	for key, val := range attr {
		cfg, ok := config[key]
		if !ok {
			o.log.Warn("ignore key whose property config not exist", log.Any("key", key))
			continue
		}
		value, err := value2Variant(val, cfg.Type)
		if err != nil {
			o.log.Warn("ignore illegal data type of val", log.Any("value", val), log.Any("type", cfg.Type))
			continue
		}

		id, err := ua.ParseNodeID(cfg.NodeID)
		if err != nil {
			return errors.Trace(err)
		}
		var req = &ua.WriteRequest{
			NodesToWrite: []*ua.WriteValue{
				{
					NodeID:      id,
					AttributeID: ua.AttributeIDValue,
					Value: &ua.DataValue{
						EncodingMask: ua.DataValueValue,
						Value:        value,
					},
				},
			},
		}
		_, err = device.opcuaClient.Write(req)
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	}
	return nil
}

func (o *observer) OnPuback(pkt *packet.Puback) error {
	return nil
}

func (o *observer) OnError(err error) {
	o.log.Error("receive mqtt message error", log.Error(err))
}
