#Tcp Reverse Proxy#
====

一个TCP反向代理，监听指定的IP和Port，并将接收到的所有包代理转发至后端Server。

##使用样例
-----
`trp -l 0.0.0.0:1935 -b 192.168.1.1:1935`

监听在本机1935端口，将所有请求代理制192.168.1.1的1935端口

`trp -l 0.0.0.0:1935 -b 192.168.1.1:1935，192.168.1.2:1935`

监听在本机1935端口，将所有请求代理制192.168.1.1和192.168.1.2的1935端口，多后端采用随机调度方式进行负载均衡

`trp -l 0.0.0.0:1935 -b 192.168.1.1:1935 1>access.log &`

后台执行，并输出访问日志

##其它说明
----
启动后运行日志会输出在stderr，访问日志会输出在stdout

访问日志定时每10秒钟输出一条

日志格式为（空格分隔）

`date time sessionid client backend upstream downstream`

##构建方法
----
###安装GO语言
GO语言可从源代码编译，下载地址为

[go1.3.3.src.tar.gz](http://www.golangtc.com/static/go/go1.3.3.src.tar.gz)

下载完成后，执行以下指令编译

`tar xzvf go1.3.3.src.tar.gz`

`cd go/src`

`./make.bash`

编译好的go指令可在`go/bin`目录下找到，将该目录加入PATH中即可执行
###配置环境变量
执行`go version`确认go语言安装正常

设置环境变量GOPATH，指向项目根目录

`export GOPATH=/path/to/trp`
###构建
执行

`go install`

在项目bin目录下可找到trp可执行文件，将该可执行文件放到要安装的目录即可执行。





