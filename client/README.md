## 数据包组成

> 每个客户端将保持一条tcp长连接与服务端进行通信，数据包格式为`|len (2 bytes)| + |data|`，其中len为2字节的data的总长度，使用大端字节序。data的组成为：
>> `|msgtype (4 bytes)| + |session (4 bytes)| + protobuf_data`
> msgtype和session都使用大端字节序，protobuf_data是protobuf序列化数据。

## msgtype与session

> msgtype是一个无符号32位整数，代表消息的类型，所有的消息类型定义在msgtype.txt，需要用工具生成对应编程语言的常量，数值从1开始逐行递增，例如目前的msgtype.txt内容为：
```
MROLELIST
MCREATEROLE
MLOGIN
MLOGOUT
MECHO
MCHAT
```
> 对应生成的c++代码应该是：
```
const uint32_t MROLELIST = 1;
const uint32_t MCREATEROLE = 2;
...
```
> 通常请求和回复都使用同一个消息类型，有需要的话可以为两者定义不同的消息类型。
session同样是一个无符号32位整数，客户端每发一个消息包就递增1，用于防止简单的抓包重放。

## 连接与登陆细节
* 客户端通过端口号1234与服务器（ip:107.170.209.97)连接，连接成功后发送字符串指令"login"，然后使用D-H密钥交换算法确定通信密钥，后续通信将以该密钥为key以aes加密。
* 密钥确定后发送用户名（测试阶段所有用户名无条件成功），然后读取服务器返回的session值。
* 登陆成功。
> 而对于断线重连来说，忽略了D-H密钥交换过程，连接成功时发送字符串指令"reconnect:用户名"，然后读取服务器返回的随机挑战，使用旧的密钥加密返回给服务端，挑战成功后读取服务端返回的新的session值。

然后进入角色的登陆流程，这个时候的通信都需要使用上述的数据包格式。目前的流程是：
* 请求角色列表（MQRolelist）－－》 角色列表返回（MRRolelist）
* 假若角色列表为空创建角色（MQCreateRole）－－》角色创建返回（MRCreateRole）
* 角色列表不为空直接登陆（MQLogin）－－》登陆返回（MRLogin）

> 协议文件放在proto目录下。
> *到正式场合时，客户端会通过某个服务地址获得服务器列表，并且会先登陆第三方平台获得一个token，然后发送token到服务器，服务器再对token进行验证，通过时会返回用户数据。
