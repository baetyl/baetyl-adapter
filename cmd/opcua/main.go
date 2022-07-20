package main

import (
    "fmt"
    "strconv"

    "github.com/baetyl/baetyl-adapter/v2/opcua"
    dm "github.com/baetyl/baetyl-go/v2/dmcontext"
    "github.com/baetyl/baetyl-go/v2/utils"
    "github.com/jinzhu/copier"
)

func main() {
    // Running module in baetyl context
    dm.Run(func(ctx dm.Context) error {
        cfg, err := genConfig(ctx)
        if err != nil {
            return err
        }
        o, err := opcua.NewOpcua(ctx, *cfg)
        if err != nil {
            return err
        }
        defer o.Close()
        ctx.Wait()
        return nil
    })
}

func genConfig(ctx dm.Context) (*opcua.Config, error) {
    cfg := &opcua.Config{}
    var devices []opcua.DeviceConfig
    var jobs []opcua.Job

    for _, deviceInfo := range ctx.GetAllDevices() {
        accessConfig := deviceInfo.AccessConfig
        if accessConfig.Opcua == nil {
            continue
        }
        device := opcua.DeviceConfig{
            Device: deviceInfo.Name,
        }
        if err := copier.Copy(&device, accessConfig.Opcua); err != nil {
            return nil, err
        }
        devices = append(devices, device)

        var jobProps []opcua.Property
        deviceTemplate, _ := ctx.GetAccessTemplates(&deviceInfo)
        if deviceTemplate != nil && deviceTemplate.Properties != nil && len(deviceTemplate.Properties) > 0 {
            for _, prop := range deviceTemplate.Properties {
                if visitor := prop.Visitor.Opcua; visitor != nil {
                    var nodeId string
                    ns := deviceInfo.AccessConfig.Opcua.NsOffset + visitor.NsBase
                    switch visitor.IdType {
                    case dm.OpcuaIdTypeI:
                        idBase, err := strconv.Atoi(visitor.IdBase)
                        if err != nil {
                            continue
                        }
                        nodeId = fmt.Sprintf("ns=%d;i=%d", ns, deviceInfo.AccessConfig.Opcua.IdOffset+idBase)
                    case dm.OpcuaIdTypeS, dm.OpcuaIdTypeG, dm.OpcuaIdTypeB:
                        nodeId = fmt.Sprintf("ns=%d;s=%s", ns, visitor.IdBase)
                    default:
                        continue
                    }
                    jobProps = append(jobProps, opcua.Property{
                        Name:   prop.Name,
                        Type:   visitor.Type,
                        NodeID: nodeId,
                    })
                }
            }
        }
        job := opcua.Job{
            Device:     deviceInfo.Name,
            Interval:   accessConfig.Opcua.Interval,
            Properties: jobProps,
        }
        jobs = append(jobs, job)
    }
    cfg.Devices = devices
    cfg.Jobs = jobs
    if err := utils.SetDefaults(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
