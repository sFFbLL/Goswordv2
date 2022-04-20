FROM alpine:latest

WORKDIR /go/src/go-sword

COPY . .

EXPOSE 8888

#ENTRYPOINT ./server -c config.docker.yaml

ENTRYPOINT pwd