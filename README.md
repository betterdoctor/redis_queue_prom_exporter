redis queue prometheus exporter
===============================

## Install

`go get -u github.com/betterdoctor/redis_queue_prom_exporter`

## Run

```
redis_queue_prom_exporter -logtostderr=true -redis-uri redis://some.host:6379/1 \
  -queues important_queue,awesome_queue -namespace my_app
```

This will export the following Prometheus metric `my_app_redis_queue_count`

