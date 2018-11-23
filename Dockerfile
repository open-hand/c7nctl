FROM golang:1.9.4-alpine3.7 as builder
WORKDIR /go/src/github.com/choerodon/c7n
ADD . .

RUN go build .

FROM vinkdong/helm:2.11.0

COPY --from=builder /go/src/github.com/choerodon/c7n/c7n usr/local/bin/

RUN \
set -ex \
   && apk add --no-cache ca-certificates

WORKDIR /etc/c7n

ADD install.yml /etc/c7n/install.yml

ENTRYPOINT ["/usr/local/bin/c7n"]
CMD ["install","-c","/etc/c7n/install.yml","--no-timeout"]