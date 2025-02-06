FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/go_community
WORKDIR /go_community

# 下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 go_community_app
RUN go build -o go_community_app .

###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

# 复制 wait-for.sh 脚本到镜像中，并设置执行权限
COPY ./wait-for.sh /
# 从builder镜像中把静态文件拷贝到当前目录
COPY ./templates /templates
COPY ./static /static
# 从builder镜像中把配置文件拷贝到当前目录
COPY ./conf /conf

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /go_community/go_community_app /

# RUN set -eux; \
# apt-get update; \
# apt-get install -y \
# --no-install-recommends \
# netcat; \
# chmod 755 wait-for.sh

# 设置镜像的源列表为阿里云的镜像源
#RUN echo "deb http://mirrors.aliyun.com/debian/ stretch main non-free contrib" > /etc/apt/sources.list \
#    echo "deb-src http://mirrors.aliyun.com/debian/ stretch main non-free contrib" >> /etc/apt/sources.list \
#    echo "deb http://mirrors.aliyun.com/debian-security stretch/updates main" >> /etc/apt/sources.list \
# 更新镜像源并安装 netcat
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends netcat; \
    chmod 755 wait-for.sh

# 声明服务端口
EXPOSE 8088

# 需要运行的命令（注释掉这一句，因为需要等 MySQL 启动之后再启动我们的Web程序）
# ENTRYPOINT ["/go_community_app", "conf/config.yaml"]