FROM golang:1.14-alpine AS builder

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on
WORKDIR /go/src/github.com/choerodon/c7nctl

COPY . .

RUN CGO_ENABLED=0 go install -ldflags '-s -w' ./cmd/c7n-tool/



FROM alpine:3.12

COPY --from=builder /go/bin/c7n-tool /c7n-tool

RUN adduser -h /home/client -s /bin/sh -u 1001 -D client && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk upgrade && \
    apk add --no-cache mysql-client postgresql-client curl wget

USER 1001

ENTRYPOINT [ "/c7n-tool" ]
