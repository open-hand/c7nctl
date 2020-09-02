FROM nginx:1.19.2-alpine

RUN adduser -h /home/client -s /bin/sh -u 1001 -D client && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk upgrade && \
    apk add --no-cache mysql-client postgresql-client curl

USER 1001
