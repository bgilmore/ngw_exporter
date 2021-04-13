package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bgilmore/ngw_exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listen = flag.String("listen", ":9099", "Exporter listen address")
	target = flag.String("target", "", "IP address of your Nokia Gateway device")

	// The default timeout is a little too generous for our upstream Prometheus' taste.
	scrapeTimeout = flag.Duration("scrape_timeout", 10*time.Second, "HTTP timeout for gateway scrapes")
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if *target == "" {
		log.Fatalf("the required -target flag is currently unset, exiting.")
	}

	http.DefaultClient.Timeout = *scrapeTimeout

	prometheus.MustRegister(collectors.NewGatewayInfoCollector(*target))
	prometheus.MustRegister(collectors.NewRadioStatsCollector(*target))
	prometheus.MustRegister(collectors.NewNetworkStatsCollector(*target))
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listening for connections at %v", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("ListenAndServe at %v failed: %v", *listen, err)
	}
}
