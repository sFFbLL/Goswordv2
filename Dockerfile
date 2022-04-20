FROM alpine:latest

WORKDIR /go/src/go-sword

COPY . .

EXPOSE 8888

#ENTRYPOINT ./server -c config.docker.yaml

ENTRYPOINT ls

ENTRYPOINT /go/src/go-sword/server -c config.docker.yaml