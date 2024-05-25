# syntax=docker/dockerfile:1

# bullseye必须加上，除非更新宿主机系统和docker
# 参考：https://github.com/docker-library/golang/issues/467#issuecomment-1601845758
FROM golang:1.21-bullseye AS build

WORKDIR /app

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
# RUN go mod download

COPY . ./
# generate a executable file
RUN go build -o /app/exec

EXPOSE 8080

CMD [ "/app/exec" ]