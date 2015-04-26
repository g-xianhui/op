## Build

安装mysql及protobuf

配置go开发环境，参见：
http://pkg.golang.org/doc/code.html

安装依赖包：
```
go get github.com/golang/protobuf/{proto,protoc-gen-go}
go get github.com/go-sql-driver/mysql
go get github.com/bitly/go-simplejson
```
安装游戏框架：
```
go get github.com/g-xianhui/op
cd $GOPATH/src/github.com/g-xianhui/op/server
go build
cd $GOPATH/src/github.com/g-xianhui/op/client
go build
```
生成数据库：
```
cd $GOPATH/src/github.com/g-xianhui/op/database
mysql <game.sql
```

## Test

现在server及client目录下分别生成了server及client可执行文件，在不同终端分别执行：
```
./server
./client
```
client可执行的指令有请求角色列表(rolelist)，创建角色(createrole），登录(login)，登出(logout)，echo以及chat，具体用法可参考cmd.go。

## Usage
在正式的场合应该将server进程变为精灵进程，由于go语言本身不支持daemonize，所以需要daemonize之类的工具的支持。作为例子，安装daemonize之后可使用control.sh脚本启动及关闭服务端进程。

默认的服务端配置在config.json，正式场合切记修改。

## Description

这是使用go语言实现的一个Actor模型的网游服务器框架，利益于go语言对多核编程的良好支持，理论上可榨干cpu的性能。

1.核心
服务器由'服务'与'Agent'两种实体组成，两者都是消息驱动，消息处理运行在自身独立的goroutine中以保证时序。

服务可使用锁的方式实现，以这种方式实现的服务不需要独立的goroutine，这方面的例子可参考agentcenter。

与多数服务不同，Agent除了接收进程内部的消息之外还需要接收来自客户端的消息，所以它实际上有两个goroutine，一个是上面所说的消息处理，另一个则用于接收网络数据包，具体可见netio.go。

框架核心解决了两部分内容：

一个是客户端与自身Agent之间的通信。应用层只需要为每种消息需要定义一个类型，然后注册该类型的回调函数即可。以框架默认通信协议为例（参考2），通过：
```
registerHandler(pb.MECHO, &pb.MQEcho{}, echo)
```
即注册了类型为pb.MECHO的消息处理函数为echo，echo会得到一个proto.Message类型的参数，这个参数就是客户端传过来的协议包。

而在发送消息到客户端则只需要调用:
```
replyMsg(agent, pb.MECHO, rep)
```
作为客户端与Agent之间通信的例子可参考echo.go。

框架的解决的另一个内容是Agent与服务之间的内部通信（当然还包括Agent对Agent，服务对服务）。目前Agent与服务的处理方式并不一样，Agent接收的内源消息类型定义为InnserMsg，这个结构体包括了消息类型(cmd)、消息参数(ud)以及消息返回(reply)，考虑到内源消息类型并不多，cmd使用的是string类型，而消息参数及返回都是interface{}类型，需要每个消息处理函数自行转换到相应的实际类型。目前提供了两种方式给Agent发送内源消息，分别是sendInnerMsg和call，前者只发送消息不期待返回，而后者会阻塞直至Agent返回数据。

而对于服务的话，我希望服务能向外部提供更好的调用接口，外部调用服务的接口时可以无视'消息'这一概念，例如:
```
broadcast(bc, msg)
```
而不是:
```
sendInnerMsg(bc, "push", msg)
```
所以目前实际上框架没有在服务对内源消息的处理上提供支持，需要服务的编写人员自行按照go的范式（锁或goroutine）去实现。梳理内源消息的处理也是接下来的工作。

2.通信协议
首先是客户端与服务端的通信协议。每个客户端将保持一条tcp长连接与服务端进行通信，数据包格式为|len (2 bytes)| + |data|，其中len为2字节的包总长度，使用大端字节序，data为自定义的协议数据。框架默认使用protobuf作为data的编码协议，可改动netmsg.go以改用其他协议（但注意现有的示例模块将不能工作）。

其次是内源消息。由于所有goroutine处于同一进程内，所以定义了一个通用的InnerMsg结构，传递InnerMsg的指针就可以了。

3.代码规范
文件名: 消息处理入口(xxxhandler.go)，消息处理逻辑主体(xxx.go)，当逻辑简单时可以直接写在xxxhandler.go。例如taskhandler.go和task.go。数据加载保存(xxxdb.go)，服务(xxx.go)。

类型及变量名使用驼峰式命名，类型名称首字母大写，变量名首字母小写。所有文件使用utf8编码，使用gofmt进行模式化。

4.作为游戏
作为游戏，框架已处理了登陆退出流程，玩家数据会定时以及在关服时定时保存，玩家各个模块的加载与保存需要自行实现，并在roledb.go里的load与save添加代码。另外实现了一个简陋的广播服务。
