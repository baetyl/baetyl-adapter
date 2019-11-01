# baetyl-adapter-modbus

baetyl-adapter-modbus模块基于modbus协议采集数据, 该模块支持配置多个从设备（简称为slave）,针对设备设置采集周期定时采集数据。
连接方式可选TCP或RTU模式。采集数据点可以通过采集设备的id, 数据起始地址，数据单元长度，功能码定义。
采集后的数据会以二进制流的形式发送到配置的mqtt hub主题，典型配置如下：

```yaml
hub:
  address: tcp://127.0.0.1:1883
  username: test
  password: test
  clientid: modbus-1
slaves:
  - id: 1
    address: tcp://127.0.0.1:502
    interval: 3s
maps:
  - slaveid: 1
    address: 0
    quantity: 1
    function: 3
publish:
 topic: test
logger:
  path: var/log/baetyl/service.log
  level: "debug"
```

1. hub定义模块连接mqtt hub的配置信息，包括hub的地址，用户名，密码和连接使用的client id。
2. slaves用于配置与模块连接的slave, 支持配置多个slave, 每个slave必须有唯一的id，否则后定义的slave会覆盖先定义的slave连接配置信息。slave连接可以通过TCP或RTU模式连接。此外，模块会以统一的时间间隔采集同一设备数据。
   1. TCP模式：地址以tcp://开头时，会默认为使用tcp连接，TCP配置详情如下:
   ```yaml
   id: 1 #设备id
   address: tcp://127.0.0.1:502 #设备地址
   timeout: 10s #缺省超时时间为10s
   idletimeout: 1m #空闲tcp连接保留时间,缺省1分钟
   interval: 5s #设备采集时间间隔，缺省为5s
   ```
   2. RTU模式：配置详情如下：
   ```yaml
   id: 2 #设备id
   address: /dev/ttyUSB4 #设备地址
   timeout: 10s #缺省超时时间为10s
   idletimeout: 1m #空闲tcp连接保持时间,缺省1分钟
   interval: 5s #设备采集时间间隔，缺省为5s
   baudrate: 19200 #波特率 缺省为19200
   databits: 8 #数据位，可选(5,6,7,8)缺省为8 
   stopbits: 1 #停止位，可选(1,2)缺省为1
   parity: E #奇偶校验类型，可选N(无，对应stopbits应配置为2)，E(奇校验)，O(偶校验)缺省为E
   ```
3. 采集点配置:每个采集点需要定义采集设备的id, 数据起始地址，采集数据长度以及modbus功能码，具体功能码根据modbus协议解释为
    ```
   1: 线圈状态
   2: 离线输入
   3: 保持寄存器
   4: 输入寄存器
   ```
   典型配置为例：
   ```yaml
   slaveid: 1 #设备id
   address: 0 #起始地址
   quantity: 2 #数据单元数量
   function: 3 #功能码
   ```
   该配置解释为，从slave id 为1的从站上定时采集数据，采集起始地址为0, 由于功能码为3，对应采集保持寄存器的数据。
   quantity定义为2，单个寄存器的数据长度为16bit(2字节)，表示从起始地址开始采集4字节的数据，其它采集类型配置与此类似。
   每个点采集后的数据会以二进制字节流的形式发送，具体格式为:
   1字节（采集设备slave id）+ 2字节（采集起始地址）+ 2字节（采集数据单元数量）+ 4字节（时间戳）+ 采集数据
4. 采集后的数据会发送至一个统一的主题，可以在配置文件中指定：
   ```yaml
   publish:
     topic: test #发送主题
     qos: 1 # 消息QOS
   ```
   即配置发送主题为test, 发送消息QOS=1

模块配置文件路径应为etc/service.yml, 模块日志文件路径为var/log/service.log