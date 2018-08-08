FROM alpine:3.6

RUN apk add --update ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY redis_queue_prom_exporter /usr/local/bin/redis_queue_prom_exporter
