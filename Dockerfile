FROM golang:1.13.8 as builder

WORKDIR /go/src/github.com/choerodon/c7nctl
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor

FROM vinkdong/helm:2.11.0

COPY --from=builder /go/src/github.com/choerodon/c7nctl/c7nctl usr/local/bin/

RUN set -ex \
   && apk add --no-cache ca-certificates

WORKDIR /etc/c7nctl

COPY install.yml /etc/c7nctl/install.yml

ENTRYPOINT ["/usr/local/bin/c7nctl"]
CMD ["install", "-c", "/etc/c7nctl/install.yml", "--no-timeout", "--skip-input"]