# modbus

## 简介
modbus模块基于modbus协议采集解析数据, 该模块支持配置多个从设备（简称为slave），设置采集周期定时采集数据。
连接方式可选TCP或RTU模式。采集数据点可以通过采集设备的id, 数据起始地址，数据单元长度，功能码，字段信息定义。
采集后的数据会以二进制流的形式发送到配置的mqtt hub主题。模块可以通过配置仅从设备采集数据，将数据以二进制流
发送到指定mqtt主题。也可以配置将采集后的数据进行解析得到所需的数据，以JSON格式发送到指定的mqtt主题。针对模块
的配置可以分为3个部分：1. slave配置 2. 任务配置 3. 数据发送配置

## slave配置
slaves用于配置与模块连接的slave, 支持配置多个slave, 每个slave必须有唯一的id，否则后定义的slave会覆盖先定义的slave连接配置信息。
slave连接可以通过TCP或RTU模式连接。此外，模块也支持自动重连机制，即slave因故障与模块连接断开后，从故障恢复后，模块
会自动重连设备并继续采集数据
   * TCP模式：地址以tcp://开头时，会默认为使用tcp连接，TCP配置详情如下:
   ```yaml
   id: 1 # 设备id
   address: tcp://127.0.0.1:502 # 设备地址
   timeout: 10s # 超时时间，默认10s
   idletimeout: 1m # 空闲tcp连接保留时间，默认1分钟
   ```
   * RTU模式：配置详情如下：
   ```yaml
   id: 1 # 设备id
   address: /dev/ttyUSB4 # 设备地址   timeout: 10s # 超时时间，默认10s
   idletimeout: 1m # 空闲tcp连接保持时间，默认1分钟
   baudrate: 19200 # 波特率 默认为19200
   databits: 8 # 数据位，可选(5,6,7,8) 默认为8 
   stopbits: 1 # 停止位，可选(1,2) 默认为1
   parity: E # 奇偶校验类型，可选N(无，对应stopbits应配置为2)，E(奇校验)，O(偶校验) 默认为E
   ```

## 任务
任务提供了对一系列数据点和采集周期，仅采集或采集并解析的定义。在模块配置中可以定义多个任务, 一个任务对应一个slave,
在任务定义中可以配置采集周期，指定任务为采集或是采集并解析。

   配置详情如下：
   ```yaml
   slaveid: 1 # 任务采集的设备id
   interval: 20s # 采集周期
   encoding: json # 编码，可配置为binary或json，binary即采集后以二进制流发送，json即采集并解析以JSON发送 
   time: # 任务采集的时间信息, encoding为binary时仅支持配置precision
     name: time # 时间field名，默认为'time'
     type: integer # 可配置为integer和string，默认为integer
     format: '2006-01-02 15:04:05' # 时间格式，需配置为与2006-01-02 15:04:05相同时间的格式，默认为2006-01-02 15:04:05
     precision: s # 可配置为s或ns，即精确到秒或者纳秒，默认为s
   maps:
   - function: 3 # 功能码 可配置为（1，2，3，4），下文有详细说明
     address: 40011 # 起始地址
     quantity: 4 # 采集数量，encoding为json时无需配置，会自动推断
     field: # encoding为json时必须配置field
       name: temperature # 解析field名
       type: float64 # 解析数据类型
   ``` 
   该配置解释为，从slave id 为1的slave上定时采集数据，采集周期为20s。采集起始地址为40011, 由于功能码为3，对应采集保持寄存器的数据。
   quantity定义为4，单个寄存器的数据长度为16bit(2字节)，表示从起始地址开始采集8字节的数据。
   采集后的数据会按照big endian解析为float64类型的value, key通过field.name指定，在这里是temperature。因此待发送数据（非真实数据）格式为
   ```json
   {
       "slaveid": 1, 
       "time": "2020-05-20 15:04:05",
       "temperature": 35.32
   }
   ```

* 解析类型
模块支持的解析数据类型包括有bool、int16、uint16、int32、uint32、int64、uint64、float32、float64。在
指定解析项type时，应为以上类型之一。解析时使用big endian方式

* 功能码
​Modbus​可​访问​的​数据​存储​在​四​个​数据​库​或​地址​范围​的​其中​一个： 线圈​状态、​离散​量​输入、​保持​寄存器​和​输入​寄存器。
其中线圈转态和离散量输入的数据以bit为单位,解析后的数据仅支持bool类型。线圈状态对应功能码1，离散量输入对应
功能码2。保持寄存器和输入寄存器的数据以双字节（16bit）为单位，解析数据支持前文所有数据类型。保持寄存机对应
功能码3，输入寄存器对应功能码4

* 采集数量
任务的encoding指定为json时，map必须对field进行配置，指定以JSON解析时的name和数据类型, map中的quantity无需
配置。因为各种数据类型对应的quantity是固定的。例如当配置field.type为int32时，即4字节，而保持寄存器或输入寄存器
的单位是16bit（2字节），因此quantity必然为2（32/16）。当配置field.type为float64时，即8字节，quantity必然为
4(64/16)。此外，当任务encoding指定为binary时，即仅进行采集，不对采集后数据进行解析。map对field配置是无效的，
且quantity必须进行配置


* 仅采集数据
将任务的encoding配置为json（默认值），模块会将采集后的数据进行解析并以JSON发送，将encoding配置为binary，模块
会将采集后的数据直接发送（不进行解析）
   具体格式为:
   ```
   |----|----|----|----|---------|----|---------|  
   |    ts   | id |a+l |  data   |a+l |  data   |
   ```

   8字节（时间戳）+ 4字节（采集设备id）+ 2字节（采集起始地址）+ 2字节（采集数量）+ 采集数据 + ...

## 发送
modbus模块目前支持将采集解析后的数据通过mqtt协议发送至mqtt hub。连接mqtt支持tcp/ws/ssl/wss等方式。通过
指定待发送的mqtt主题，采集或解析后的数据会发送至该主题。
mqtt配置:
   ```yaml
  mqtt:
    clientid: Client 连接 Hub 的 Client ID。cleansession 为 false 则不允许为空
    address: [必须] Client 连接Hub的地址
    username: Client 连接Hub的用户名
    password: 如果采用账号密码，必须填 Client 连接Hub的密码，否者不用填写
    ca: 如果采用证书双向认证，必须填 Client 连接Hub的CA证书路径
    key: 如果采用证书双向认证，必须填 Client 连接Hub的客户端私钥路径
    cert: 如果采用证书双向认证，必须填 Client 连接Hub的客户端公钥路径
    timeout: 默认值：30s，Client 连接 Hub 的超时时间
    maxReconnectInterval: 默认值：3m，Client 连接 Hub 的重连最大间隔时间，从500微秒翻倍增加到最大值
    keepalive: 默认值：30s，Client 连接Hub的保持连接时间
    cleansession: 默认值：false，Client 连接 Hub 的是否保持 Session
    disableAutoAck: 默认值：false，禁用自动ack
    subscriptions: 订阅hub的主题列表
    maxCacheMessages: 默认值：10，Client 发送消息给 Hub 的内存队列大小，异常退出会导致消息丢失，恢复后 QoS 为1的消息依赖 Hub 重发    
   ```

发送主题配置:
   ```yaml
    publish:
      topic: 发送的主题
      qos: 发送qos 
   ```   

## 采集数据典型配置如下：
```yaml
mqtt:
  address: tcp://127.0.0.1:1883 # 连接mqtt hub的地址 
  clientid: modbus-1 # 连接mqtt hub时使用的client id，基于baetyl框架运行时可不配置
slaves:
  - id: 1 # slave id
    address: tcp://127.0.0.1:502 # 基于tcp连接slave时的地址
jobs:
  - slaveid: 1 # 采集任务的对应设备
    encoding: binary # 指定仅采集数据，数据以二进制流发送
    maps:
      - function: 1 # 功能码，对应线圈状态
        address: 32 # 起始地址
        quantity: 1 # 采集数量，线圈状态对应1bit，encoding为binary时，quantity不可缺失
      - function: 3 # 功能码，对应保持寄存器
        address: 40011 # 起始地址
        quantity: 1 # 采集数量，保持寄存器对应16bit，encoding为binary时，quantity不可缺失
publish:
 topic: test # 采集数据发送的mqtt主题
logger:
  filename: var/log/baetyl/service.log # 日志路径 
  level: info # 日志级别
```

## 采集并解析数据典型配置如下：
```yaml
mqtt:
  address: tcp://127.0.0.1:1883 # 连接mqtt hub的地址 
  clientid: modbus-1 # 连接mqtt hub时使用的client id，基于baetyl框架运行时可不配置
slaves:
  - id: 1 # slave id
    address: tcp://127.0.0.1:502 # 基于tcp连接slave时的地址
jobs:
  - slaveid: 1 # 采集任务的对应设备
    encoding: json # 指定采集解析数据，数据以JSON发送，默认为json
    maps:
      - function: 2 # 功能码，对应离散量输入
        address: 47 # 起始地址
        quantity: 1 # 采集数量，离散量输入对应1bit, 当解析数据时模块可自动确定
        field:
          name: switch # 解析后数据field名
          type: bool # 针对float32数据类型解析数据
      - function: 4 # 功能码，对应输入寄存器
        address: 30027 # 起始地址
        quantity: 1 # 采集数量，输入寄存器对应16bit
        field:
          name: humidity # 解析后数据field名
          type: float32 # 针对float32数据类型解析数据
publish:
 topic: test # 采集数据发送的mqtt主题
logger:
  filename: var/log/baetyl/service.log # 日志路径 
  level: info # 日志级别
```

#### 日志路径配置
模块配置文件路径应为etc/service.yml, 模块日志文件路径为var/log/baetyl/service.log