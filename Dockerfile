FROM golang:1.19 as app-builder
MAINTAINER bewakes bibek.pandey@togglecorp.com

WORKDIR /go/app
COPY . /go/app

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get . && \
    go build -o /go/app/coding-challenge /go/app/main.go


FROM ubuntu:18.04 as app

RUN apt update -y && apt install openssh-client -y

RUN mkdir -p $HOME/.ssh/

WORKDIR /go/app
COPY . /go/app
COPY --from=app-builder /go/app/coding-challenge /go/app/coding-challenge

ENTRYPOINT ["/go/app/scripts/server_entrypoint.sh"]
