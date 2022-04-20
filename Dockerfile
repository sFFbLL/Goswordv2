FROM golang:alpine

WORKDIR /go/src/go-sword

COPY . .

EXPOSE 8888

RUN chmod 777 ./server

ENTRYPOINT ./server -c config.docker.yaml
