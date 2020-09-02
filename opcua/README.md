# opcua

## 简介
opcua模块基于opcua协议读取与写入数据, 该模块支持配置多个从设备，设置读取周期定时读取数据。
读取数据点可以通过采集设备的nodeid, 变量名称，变量类型定义。
读取的数据会以JSON发送到配置的MQTT broker主题。模块可以通过配置仅从设备采集数据，将数据进行JSON序列化后
发送到指定MQTT主题。
opcua模块可以通过 [baetyl](https://github.com/baetyl/baetyl) 与 [baetyl-cloud](https://github.com/baetyl/baetyl-broker) 部署，
并结合baetyl-broker使用，baetyl会自动配置opcua与 [baetyl-broker](https://github.com/baetyl/baetyl-broker)
之间的双向tls连接。
可以参考 [baetyl文档](https://docs.baetyl.io/zh_CN/latest/) 与[最佳实践](https://docs.baetyl.io/zh_CN/latest/practice/application-deployment-practice.html)。

针对模块的配置可以分为3个部分：1. 连接设备配置 2. 任务配置 3. 数据发送配置：

## 连接设备配置
devices配置项用于配置与模块连接的设备, 支持配置多个设备, 每个设备必须有唯一的id，否则后定义的设备会覆盖先定义的设备连接配置信息。
一个设备基本连接需要包括设备id和设备端点。此外，对设备的配置还包括连接超时时间、安全、认证、证书等选项。
另外，模块也支持自动重连机制，即设备因故障与模块连接断开后，从故障恢复后，模块会自动重连设备并继续读取数据。
   ```yaml
   id: 1 # 设备id
   endpoint: opc.tcp://test.local:53530/OPCUA/Server # 设备端点
   timeout: 20s # 默认为10s
   security:
     policy: None # 可选None、Basic128Rsa15、Basic256、Basic256Sha256、Aes128Sha256RsaOaep、Aes256Sha256RsaPss
     mode: None # 可选Invalid、None、Sign、SignAndEncrypt
   auth:
     username: test # 用户名
     password: test # 密码
   ```

## 任务
任务提供了对一系列数据点和读取周期的定义。在模块配置中可以定义多个任务, 一个任务对应一个设备,
此外，任务还需要配置读取后数据将发送的MQTT broker主题。

   配置详情如下：
   ```yaml
   deviceid: 1 # 任务对应设备id
   interval: 20s # 读取周期
   time: # 任务的时间信息
     name: time # 时间field名，默认为'time'
     type: integer # 可配置为integer和string，默认为integer
     format: '2006-01-02 15:04:05' # 时间格式，需配置为与2006-01-02 15:04:05相同时间的格式，默认为2006-01-02 15:04:05
     precision: s # 可配置为s或ns，即精确到秒或者纳秒，默认为s
   properties:
   - name: var # 读取数据的变量名
     type: float64 # 读取数据的数据类型
     nodeid: ns=1;i=1001 # 读取数据点对应的nodeid
   publish:
     topic: # 发送的主题
     qos: #发送qos
   ``` 
   该配置解释为，从设备id为1的设备上定时读取数据，周期为20s。读取数据的nodeid为ns=1;i=1001，
   采集后的数据会与配置数据类型进行校验, key通过变量名称指定，在这里是var。因此待发送数据（非真实数据）格式为
   ```json
   {
       "deviceid": 1, 
       "time": "2020-05-20 15:04:05",
       "attr": {
           "var": 35.32 
       }
   }
   ```
   该数据会发送至jobs中定义的主题。当运行在baetyl中时，可以不配置主题，默认为<service-name>/<deviceid>。

* nodeid
NodeId的标识符部分唯一地标识名称空间中的节点，需要名称空间加上标识符才能形成完全限定的标识符。

* 数据类型
模块支持的解析数据类型包括有bool、int16、uint16、int32、uint32、int64、uint64、float32、float64、string。在
指定解析项type时，应为以上类型之一。


## 发送
opcua模块目前支持将采集解析后的数据通过MQTT协议发送至MQTT broker。连接broker支持tcp/ws/ssl/wss等方式。通过
指定待发送的MQTT主题，采集或解析后的数据会发送至该主题。
MQTT连接配置:
   ```yaml
  broker:
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

## 采集数据典型配置如下：
```yaml
broker:
  address: tcp://127.0.0.1:1883 # 连接mqtt hub的地址 
  clientid: opcua-1 # 连接mqtt hub时使用的client id，基于baetyl框架运行时可不配置
devices:
  - id: 1 # 设备id
    address: opc.tcp://test.local:53530/OPCUA/Server # 设备端点 
jobs:
  - deviceid: 1 # 读取任务的对应设备
    properties:
      - name: var1 # 变量名
        type: int16 # 数据类型
        nodeid: ns=3;i=1001 # 读取数据的nodeid
      - name: var2 # 变量名
        type: float64 # 数据类型
        nodeid: ns=2;i=1002 # 读取数据的nodeid
    publish:
      topic: test # 读取数据发送的mqtt主题
logger:
  filename: var/log/baetyl/service.log # 日志路径 
  level: info # 日志级别
```

以上述配置读取opcua设备为例，发送到MQTT broker的test主题中数据结构为
```json
{
  "deviceid": 1,
  "attr": {
    "var1": 43,
    "var2": 23.24 
  },
  "time": 234435345,
}
```
其中，deviceid是在配置文件中定义的从设备id, 标识数据是从哪一个从设备读取。attr域即读取的的数据。
最后，time是读取数据时的时间戳，如前文所述，可以在配置文件中对该字段的格式进行配置。

#### 日志路径配置
模块配置文件路径相对可执行文件的相对路径应为etc/baetyl/conf.yml, 模块日志文件相对可执行文件路径为var/log/baetyl/service.log

## 反控
通过opcua设备读取数据是常用功能，但有时也需要通过一定方式控制opcua设备中的数据。前者是通过该模块读取设备
的数据，而后者则是通过该模块向设备写入数据。读取时需要通过配置文件指定所需要读取的数据相关信息以及解析变量，写入
数据至opcua设备主要是通过MQTT发送数据。即配置模块订阅MQTT broker的topic, 接收写入数据的消息（类指令），再将
消息中的数据写入设备中。

需要注意的是，反控使用的是配置中的相关变量，即在配置文件汇总定义了相关变量名与变量类型后，反控才能生效，
还是以上述采集与解析的典型配置为例，在jobs中定义了两个变量var1和var2, 数据类型为int16和float64
需要注意的是，需要在MQTT中配置订阅反控主题，反控主题可以为多个, 代表可以通过多个主题对设备进行控制。

```yaml
broker:
  address: tcp://127.0.0.1:1883 # 连接mqtt hub的地址 
  clientid: opcua-1 # 连接mqtt hub时使用的client id，基于baetyl框架运行时可不配置
  subscriptions:
      - topic: control # 反控向从设备写入数据时订阅的主题
devices:
  - id: 1 # 设备id
    address: opc.tcp://test.local:53530/OPCUA/Server # 设备端点 
jobs:
  - deviceid: 1 # 读取任务的对应设备
    properties:
      - name: var1 # 变量名
        type: int16 # 数据类型
        nodeid: ns=3;i=1001 # 读取数据的nodeid
      - name: var2 # 变量名
        type: float64 # 数据类型
        nodeid: ns=2;i=1002 # 读取数据的nodeid
    publish:
      topic: test # 读取数据发送的mqtt主题
logger:
  filename: var/log/baetyl/service.log # 日志路径 
  level: info # 日志级别
```

通过发送消息至MQTT broker的control主题，可以向从设备写入数据
```json
{
  "deviceid": 1,
  "attr": {
    "var1": 12,
    "var2": 89.87 
  }
}
```
和读取时的数据结构类似，deviceid标识需要反控的从设备id, 与配置文件中定义的deviceid一致。attr域中包含了期望
控制的变量，变量名也与配置文件中jobs中定义的变量一致。需要注意，如果发送配置文件中未定义的变量名与值、数据
类型与配置文件中定义不一致会被忽略。