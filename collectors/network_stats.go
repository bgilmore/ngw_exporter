package collectors

import (
	"log"
	"strconv"

	"github.com/bgilmore/ngw_exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

type networkStatsCollector struct {
	target string

	// LAN
	lanPortUpMetric          *prometheus.Desc
	lanTransmitBytesMetric   *prometheus.Desc
	lanReceiveBytesMetric    *prometheus.Desc
	lanTransmitPacketsMetric *prometheus.Desc
	lanReceivePacketsMetric  *prometheus.Desc
	lanTransmitErrorsMetric  *prometheus.Desc
	lanReceiveErrorsMetric   *prometheus.Desc

	// WLAN
	wlanAssociationsMetric    *prometheus.Desc
	wlanTransmitBytesMetric   *prometheus.Desc
	wlanReceiveBytesMetric    *prometheus.Desc
	wlanTransmitPacketsMetric *prometheus.Desc
	wlanReceivePacketsMetric  *prometheus.Desc
	wlanTransmitErrorsMetric  *prometheus.Desc
	wlanReceiveErrorsMetric   *prometheus.Desc
	wlanTransmitDropsMetric   *prometheus.Desc
	wlanReceiveDropsMetric    *prometheus.Desc
}

func NewNetworkStatsCollector(target string) *networkStatsCollector {
	return &networkStatsCollector{
		target: target,

		lanPortUpMetric: prometheus.NewDesc(
			metricName("network", "lan_port_link_up"),
			"Whether this LAN port has a link",
			[]string{"port"}, nil,
		),

		lanTransmitBytesMetric: prometheus.NewDesc(
			metricName("network", "lan_transmit_bytes_total"),
			"Number of transmitted bytes",
			[]string{"port"}, nil,
		),

		lanReceiveBytesMetric: prometheus.NewDesc(
			metricName("network", "lan_receive_bytes_total"),
			"Number of received bytes",
			[]string{"port"}, nil,
		),

		lanTransmitPacketsMetric: prometheus.NewDesc(
			metricName("network", "lan_transmit_packets_total"),
			"Number of transmitted packets",
			[]string{"port"}, nil,
		),

		lanReceivePacketsMetric: prometheus.NewDesc(
			metricName("network", "lan_receive_packets_total"),
			"Number of received packets",
			[]string{"port"}, nil,
		),

		lanTransmitErrorsMetric: prometheus.NewDesc(
			metricName("network", "lan_transmit_errors_total"),
			"Number of TX errors",
			[]string{"port"}, nil,
		),

		lanReceiveErrorsMetric: prometheus.NewDesc(
			metricName("network", "lan_receive_errors_total"),
			"Number of RX errors",
			[]string{"port"}, nil,
		),

		wlanAssociationsMetric: prometheus.NewDesc(
			metricName("network", "wlan_associations"),
			"Number of actively associated stations",
			[]string{"ssid", "channel"}, nil,
		),

		wlanTransmitBytesMetric: prometheus.NewDesc(
			metricName("network", "wlan_transmit_bytes_total"),
			"Number of transmitted bytes",
			[]string{"ssid", "channel"}, nil,
		),

		wlanReceiveBytesMetric: prometheus.NewDesc(
			metricName("network", "wlan_receive_bytes_total"),
			"Number of received bytes",
			[]string{"ssid", "channel"}, nil,
		),

		wlanTransmitPacketsMetric: prometheus.NewDesc(
			metricName("network", "wlan_transmit_packets_total"),
			"Number of transmitted packets",
			[]string{"ssid", "channel"}, nil,
		),

		wlanReceivePacketsMetric: prometheus.NewDesc(
			metricName("network", "wlan_receive_packets_total"),
			"Number of received packets",
			[]string{"ssid", "channel"}, nil,
		),

		wlanTransmitErrorsMetric: prometheus.NewDesc(
			metricName("network", "wlan_transmit_errors_total"),
			"Number of TX errors",
			[]string{"ssid", "channel"}, nil,
		),

		wlanReceiveErrorsMetric: prometheus.NewDesc(
			metricName("network", "wlan_receive_errors_total"),
			"Number of RX errors",
			[]string{"ssid", "channel"}, nil,
		),

		wlanTransmitDropsMetric: prometheus.NewDesc(
			metricName("network", "wlan_transmit_drops_total"),
			"Number of TX drops",
			[]string{"ssid", "channel"}, nil,
		),

		wlanReceiveDropsMetric: prometheus.NewDesc(
			metricName("network", "wlan_receive_drops_total"),
			"Number of RX drops",
			[]string{"ssid", "channel"}, nil,
		),
	}
}

func (c *networkStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.lanPortUpMetric
	ch <- c.lanTransmitBytesMetric
	ch <- c.lanReceiveBytesMetric
	ch <- c.lanTransmitPacketsMetric
	ch <- c.lanReceivePacketsMetric
	ch <- c.lanTransmitErrorsMetric
	ch <- c.lanReceiveErrorsMetric
	ch <- c.wlanAssociationsMetric
	ch <- c.wlanTransmitBytesMetric
	ch <- c.wlanReceiveBytesMetric
	ch <- c.wlanTransmitPacketsMetric
	ch <- c.wlanReceivePacketsMetric
	ch <- c.wlanTransmitErrorsMetric
	ch <- c.wlanReceiveErrorsMetric
	ch <- c.wlanTransmitDropsMetric
	ch <- c.wlanReceiveDropsMetric
}

func (c *networkStatsCollector) Collect(ch chan<- prometheus.Metric) {
	m := new(models.NetworkStats)
	if err := scrape(c.target, m); err != nil {
		log.Printf("Failed to collect network stats from %s: %v", c.target, err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	for i, lan := range m.LAN {
		if lan.PortEnabled == 0 {
			continue
		}

		labels := []string{strconv.Itoa(i + 1)}
		var portUp float64
		if lan.PortState == "Up" {
			portUp = 1
		}
		ch <- prometheus.MustNewConstMetric(c.lanPortUpMetric, prometheus.GaugeValue, portUp, labels...)

		ch <- prometheus.MustNewConstMetric(c.lanTransmitBytesMetric, prometheus.CounterValue, float64(lan.Stats.BytesSent), labels...)
		ch <- prometheus.MustNewConstMetric(c.lanReceiveBytesMetric, prometheus.CounterValue, float64(lan.Stats.BytesReceived), labels...)
		ch <- prometheus.MustNewConstMetric(c.lanTransmitPacketsMetric, prometheus.CounterValue, float64(lan.Stats.PacketsSent), labels...)
		ch <- prometheus.MustNewConstMetric(c.lanReceivePacketsMetric, prometheus.CounterValue, float64(lan.Stats.PacketsReceived), labels...)
		ch <- prometheus.MustNewConstMetric(c.lanTransmitErrorsMetric, prometheus.CounterValue, float64(lan.Stats.TXErrors), labels...)
		ch <- prometheus.MustNewConstMetric(c.lanReceiveErrorsMetric, prometheus.CounterValue, float64(lan.Stats.RXErrors), labels...)
	}

	for _, wlan := range m.WLAN {
		if wlan.RadioEnabled == 0 || wlan.NetworkEnabled == 0 {
			continue
		}

		labels := []string{wlan.SSID, strconv.Itoa(wlan.Channel)}
		ch <- prometheus.MustNewConstMetric(c.wlanAssociationsMetric, prometheus.GaugeValue, float64(wlan.Associations), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanTransmitBytesMetric, prometheus.CounterValue, float64(wlan.BytesSent), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanReceiveBytesMetric, prometheus.CounterValue, float64(wlan.BytesReceived), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanTransmitPacketsMetric, prometheus.CounterValue, float64(wlan.PacketsSent), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanReceivePacketsMetric, prometheus.CounterValue, float64(wlan.PacketsReceived), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanTransmitErrorsMetric, prometheus.CounterValue, float64(wlan.TXErrors), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanReceiveErrorsMetric, prometheus.CounterValue, float64(wlan.RXErrors), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanTransmitDropsMetric, prometheus.CounterValue, float64(wlan.TXDrops), labels...)
		ch <- prometheus.MustNewConstMetric(c.wlanReceiveDropsMetric, prometheus.CounterValue, float64(wlan.RXDrops), labels...)
	}
}
