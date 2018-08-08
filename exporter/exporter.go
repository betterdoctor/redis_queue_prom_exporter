package exporter

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-redis/redis"
)

var (
	registry        *prometheus.Registry
	up              *prometheus.Desc
	redisQueueCount *prometheus.Desc
)

// Exporter collects and exposes prometheus metrics
type Exporter struct {
	RedisClient *redis.Client
	Queues      string
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri, queues, namespace string) (*Exporter, error) {
	if err := validateRedisURI(uri); err != nil {
		return nil, err
	}
	// we've validated for properly formatted redis URI
	u, _ := url.Parse(uri)
	p := strings.Split(u.Path, "/")
	db, _ := strconv.Atoi(p[len(p)-1])
	client := redis.NewClient(&redis.Options{
		Addr: u.Host,
		DB:   db,
	})

	// create metrics with provided namespace
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"was the last query to redis successful",
		nil, nil,
	)
	redisQueueCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "redis_queue_count"),
		"length of redis queues",
		nil, nil,
	)

	return &Exporter{
		RedisClient: client,
		Queues:      queues,
	}, nil
}

// Describe satifies prometheus Collector interface
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- redisQueueCount
}

// Collect satifies prometheus Collector interface
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	_, err := e.RedisClient.Ping().Result()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
	}

	var size int32
	for _, key := range strings.Split(e.Queues, ",") {
		cmd := e.RedisClient.LLen(key)
		if err := cmd.Err(); err != nil {
			log.Errorf("%v", err)
		}
		size += int32(cmd.Val())
	}
	ch <- prometheus.MustNewConstMetric(redisQueueCount, prometheus.GaugeValue, float64(size))
}

func validateRedisURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("redis_url (%s) failed to parse: %s", uri, err)
	}
	p := strings.Split(u.Path, "/")
	if len(p) != 2 {
		return fmt.Errorf("redis_url (%s) must be in redis://host:port/db format", uri)
	}
	_, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return fmt.Errorf("redis_url (%s) must be in redis://host:port/db format: %s", uri, err)
	}

	return nil
}
