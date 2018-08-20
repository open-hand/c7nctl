FROM golang:1.9.4-alpine3.7 as builder
WORKDIR /go/src/github.com/choerodon/c7n
ADD . .

RUN \
  apk update &&\
  apk add git

RUN go build .

FROM alpine

COPY --from=builder /go/src/github.com/choerodon/c7n/c7n /c7n

RUN \
set -ex \
   && apk add --no-cache ca-certificates

CMD ["/c7n"]