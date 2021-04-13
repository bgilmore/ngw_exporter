package collectors

import (
	"log"

	"github.com/bgilmore/ngw_exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

type gatewayInfoCollector struct {
	target string

	uptimeMetric       *prometheus.Desc
	softwareInfoMetric *prometheus.Desc
	hardwareInfoMetric *prometheus.Desc
}

func NewGatewayInfoCollector(target string) *gatewayInfoCollector {
	return &gatewayInfoCollector{
		target: target,

		uptimeMetric: prometheus.NewDesc(
			metricName("gateway", "uptime_seconds"),
			"How many seconds a gateway has been continuously powered on.",
			[]string{"serial"},
			nil,
		),

		softwareInfoMetric: prometheus.NewDesc(
			metricName("gateway", "software_info"),
			"Gateway software information",
			[]string{"serial", "version"},
			nil,
		),

		hardwareInfoMetric: prometheus.NewDesc(
			metricName("gateway", "hardware_info"),
			"Gateway hardware information",
			[]string{"serial", "name", "model", "version"},
			nil,
		),
	}
}

func (c *gatewayInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptimeMetric
	ch <- c.softwareInfoMetric
	ch <- c.hardwareInfoMetric
}

func (c *gatewayInfoCollector) Collect(ch chan<- prometheus.Metric) {
	m := new(models.GatewayInfo)
	if err := scrape(c.target, m); err != nil {
		log.Printf("Failed to collect gateway info from %s: %v", c.target, err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}
	for _, device := range m.Device {
		ch <- prometheus.MustNewConstMetric(c.uptimeMetric, prometheus.CounterValue, float64(device.UptimeSeconds), device.Serial)
		ch <- prometheus.MustNewConstMetric(c.softwareInfoMetric, prometheus.GaugeValue, 1.0, device.Serial, device.SWVersion)
		ch <- prometheus.MustNewConstMetric(c.hardwareInfoMetric, prometheus.GaugeValue, 1.0, device.Serial, m.Model, device.Product, device.HWVersion)
	}
}
