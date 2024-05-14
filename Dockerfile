# syntax=docker/dockerfile:1

FROM golang:1.21 AS build

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