# 第一阶段构建编译环境
FROM golang:alpine AS builder

WORKDIR /go/src/go-sword
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && go generate && go env && go build -o server .
# 第二阶段构建运行环境
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update && apk add git

WORKDIR /go/src/go-sword

COPY --from=builder /go/src/go-sword ./

EXPOSE 8888

ENTRYPOINT ./server -c config.docker.yaml
