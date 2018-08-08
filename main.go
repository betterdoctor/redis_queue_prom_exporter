package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/deepthawtz/redis_queue_prom_exporter/exporter"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

// Version gets set in build
var Version string

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("port", ":9107", "Address to listen on to serve metrics")
		metricsPath   = flag.String("path", "/metrics", "Path under which to expose metrics.")
		redisURI      = flag.String("redis-uri", "", "Redis URI to check length of queues (e.g., redis://some.host:6379/1)")
		queues        = flag.String("queues", "", "Comma-separated list of queues to check length of")
		namespace     = flag.String("namespace", "", "Namespace to segment metrics with")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println("version:", Version)
		fmt.Println(version.Print("redis_queue_prom_exporter"))
		os.Exit(0)
	}

	if *redisURI == "" {
		fmt.Println("-redis-uri is required")
		os.Exit(1)
	}

	if *queues == "" {
		fmt.Println("-queues is required")
		os.Exit(1)
	}

	if *namespace == "" {
		fmt.Println("-namespace is required")
		os.Exit(1)
	}

	glog.Infoln("Starting consul_exporter:", version.Info())
	glog.Infoln("Build context:", version.BuildContext())
	exporter, err := exporter.NewExporter(*redisURI, *queues, *namespace)
	if err != nil {
		glog.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Redis Queue Prometheus Exporter</title></head>
		<body>
		<h1>Redis Queue Prometheus Exporter</h1>
		<p>redis: ` + *redisURI + `</p>` +
			`<p>queues: ` + *queues + `</p>` +
			`<p><a href='` + *metricsPath + `'>Metrics</a></p>
		</body></html>`))
	})

	glog.Infoln("Listening on", *listenAddress)
	glog.Flush()
	glog.Fatal(http.ListenAndServe(*listenAddress, nil))
}
