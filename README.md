## Build

配置go开发环境，参见：
http://pkg.golang.org/doc/code.html
安装依赖包：
```
go get github.com/golang/protobuf
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
安装mysql及生成数据库：
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
