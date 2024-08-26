# syntax=docker/dockerfile:1

# bullseye必须加上，除非更新宿主机系统和docker
# 参考：https://github.com/docker-library/golang/issues/467#issuecomment-1601845758
FROM golang:1.21-bullseye AS build

ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

WORKDIR /app

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Use build cache if go.mod and go.sum are not changed
COPY go.mod go.sum ./
RUN go mod download

# generate a executable file
COPY . ./
RUN go build -o /app/exec

FROM debian:bullseye AS export

RUN apt-get update
RUN apt-get install -y ca-certificates

WORKDIR /app
COPY ./config.yml ./
COPY --from=build /app/exec ./

ENV LANG=C.UTF-8
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD [ "/app/exec" ]
