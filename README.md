# http3传输文件
http3标准尚未出来，但是quic协议已出来很长时间。http3简单来说 = http2 + quic协议，quic底层基于UDP协议，能够解决队头阻塞和提升传输效率。  

## 项目使用
1. 安装golang1.14  
 [下载地址](https://gomirrors.org/)，选择合适的平台进行下载，并添加go二进制文件路径到PATH环境变量中  

2. 安装依赖  
项目使用quic-go包，需要下载该依赖，下载方式如下：
```shell
go get -u github.com/lucas-clemente/quic-go
```  

下载依赖需要连接github，一般比较慢，可在.bashrc中或者.zshrc中添加PROXY，方法如下：
```shell
export GOPROXY=https://mirrors.aliyun.com/goproxy/
保存
source .bashrc
重新执行go get操作
```

3. 编译  
服务端和客户端均在各自的项目主目录下，使用以下简单命令编译
```shell
go build
```  
编译之后生成二进制可执行文件，包含了所有依赖

4. 启动服务端
```shell
./http-file-server
```  

5. 启动客户端
```shell
./http-file-client go1.14.3.darwin-amd64.tar.gz
```  
其中go1.14.3.darwin-amd64.tar.gz为当前目录下的文件




