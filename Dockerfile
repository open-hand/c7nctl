FROM golang:1.9.4-alpine3.7 as builder
WORKDIR /go/src/github.com/choerodon/c7nctl
ADD . .

RUN go build .

FROM vinkdong/helm:2.11.0

COPY --from=builder /go/src/github.com/choerodon/c7nctl/c7nctl usr/local/bin/

RUN \
set -ex \
   && apk add --no-cache ca-certificates

WORKDIR /etc/c7nclt

ADD install.yml /etc/c7nctl/install.yml

ENTRYPOINT ["/usr/local/bin/c7nctl"]
CMD ["install","-c","/etc/c7nctl/install.yml","--no-timeout","--skip-input"]