# nsq example

## 部署

### docker部署

拉取镜像
```shell
docker pull nsqio/nsq
```

`nsqio/nsq`镜像中包含三个组件: nsqlookupd, nsqd, nsqadmin。每一个组件都可以通过指定组件名的方式去启动, 启动命令的形式如下:

#### nsqlookupd组件

```shell
docker run --name lookupd -p 4160:4160 -p 4161:4161 -d nsqio/nsq /nsqlookupd
```

#### nsqd组件

首先通过`ifconfig`命令查看主机ip, 以本机ip`192.168.71.57`为例

```shell
docker run --name nsqd -p 4150:4150 -p 4151:4151 -d nsqio/nsq /nsqd --broadcast-address=192.168.71.57 --lookupd-tcp-address=192.168.71.57:4160
```

#### nsqadmin组件

```shell
docker run -d --name nsqadmin -p 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=192.168.71.57:4161
```

访问 http://localhost:4161 就可以查看nsq系统详情