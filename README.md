Golang开发Oracle应用程序

### 1.前言

出于项目需要，并行运维组的监控系统要支持Oracle、DB2等商业数据库。

而Golang语言发展的正是如火如荼，所以Telegraf的采集端使用Golang开发。而Golang访问Oracle数据库的驱动目前只有go-oci8

> https://github.com/wendal/go-oci8

所以最终采集端使用Golang+go-oci8进行开发。

### 2.准备开发环境

#### 2.1.操作系统选择

Oracle客户端开发环境支持全系的Windows、Mac和Linux系统。所以可以选择一个自己习惯的系统作为自己的开发环境。

本文档以CentOS6 64位系统为例。

#### 2.1.1.下载Oracle客户端开发包

Oracle客户端开发包的下载首页为：

> http://www.oracle.com/technetwork/cn/database/features/instant-client/index-097480.html

此页面提供了不同的操作系统开发包，本文档选择"Instant Client for Linux x86-64"，进入后先选择“Accept License Agreement”统一条款。下面的链接就可用了。

由于Oracle 11Gr2是现在占比非常高的版本，各种的bug也最少。所以，本文档选择“Version 11.2.0.4.0 ”这个版本，需要下载这个版本中的两个文件

>  oracle-instantclient11.2-basic-11.2.0.4.0-1.x86_64.rpm 

>  oracle-instantclient11.2-devel-11.2.0.4.0-1.x86_64.rpm 

以上两个文件basic为库文件，devel为头文件和库文件。

#### 2.1.2.安装Oracle客户端开发包

本文档的系统为CentOS6，下载好两个rpm包后用yun安装

> yum install -y oracle-instantclient11.2-basic-11.2.0.4.0-1.x86_64.rpm oracle-instantclient11.2-devel-11.2.0.4.0-1.x86_64.rpm

安装后会在如下目录中生成一系列文件：

1. /usr/include/oracle/11.2/client64开发用的头文件。
1. /usr/lib/oracle/11.2/client64/lib开发库文件。

#### 2.1.3.安装Golang

可以参考go官方的安装方法安装此软件，也可以使用系统默认带的版本安装。

#### 2.1.4.创建一个系统用户

此用户作为开发用户，如果不创建用户也可以使用root用户

> \# useradd go
> \# passwd go

#### 2.1.5.在go用户的家目录新建工作目录

> $ mkdir /home/go/go
> $ cd go
> $ mkdir bin pkg src


#### 2.1.6.为新建的go用户配置环境变量

修改.bashrc文件，追加如下内容

> export GOPATH="/home/go/go"

> export GOBIN=$GOPATH/bin

> export CGO_ENABLED=1

> export ORACLE_HOME=/usr/lib/oracle/11.2/client64

> export PATH=$PATH:$GOBIN:$ORACLE_HOME/bin

> export export GO_OCI8_CONNECT_STRING="system/oracle@172.18.7.201:1521/DB11G"

:star2: 环境变量GO_OCI8_CONNECT_STRING是目标数据库的连接信息 :star2:

> export LD_LIBRARY_PATH=/usr/lib/oracle/11.2/client64/lib

#### 2.1.7.配置pkgconfig配置文件

创建/usr/share/pkgconfig/oci8.pc文件，此文件内容为：

> prefix=/usr

> includedir=${prefix}/include/oracle/11.2/client64

> libdir=${prefix}/lib/oracle/11.2/client64/lib

> 

> Name: oci8

> Description: Oracle Instant Client

> Version: 11.2

> Cflags: -I${includedir}

> Libs: -L${libdir} -lclntsh -locci



#### 2.1.7.使环境变量生效

> $ source .bashrc

### 3.第一个程序

#### 3.1.下载go-oci8包

> $ go get github.com/wendal/go-oci8

#### 3.2.创建应用程序

在/home/go/go/src/目录中新建一个目录

> $ mkdir /home/go/go/src/pdefcon

新建一个目录

> $ cp /home/go/go/src/github.com/mattn/go-oci8/\_example/lastinsertid/main.go /home/go/go/src/pdefcon/main.go

拷贝实例文件到src目录下的pdefcon目录中

> $ go install pdefcon

编译并且安装可执行文件

> $ pdefcon

执行编译好的二进制文件。

### 4.需要注意的地方

go env中的CGO_ENABLED必须为1,这样才能调用Oracle的动态库文件。




