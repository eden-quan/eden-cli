ARG GOPATH=/myworkspace/golang
ARG WORK_DIR=/myworkspace/golang/src/workspace
ARG GOPROXY=https://goproxy.io,direct

FROM golang:1.21 AS builder

ARG WORK_DIR
ARG GOPATH
ARG GOPROXY

ENV GO111MODULE=on
ENV GOPATH=$GOPATH
ENV GOPROXY=$GOPROXY
ENV CGO_ENABLED=0

WORKDIR $WORK_DIR

COPY . .

#RUN ls -al
RUN go mod download && go mod tidy
RUN go build -o ./bin/main ./app/grant-service/cmd/grant-service/...

FROM alpine:latest as release

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add -U tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone

ARG WORK_DIR

WORKDIR $WORK_DIR

COPY --from=builder $WORK_DIR/bin/main $WORK_DIR/main

EXPOSE 8000
EXPOSE 9000

ENTRYPOINT ./main -conf=./configs
