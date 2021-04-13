package collectors

import (
	"log"

	"github.com/bgilmore/ngw_exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

type radioStatsCollector struct {
	target string

	receiveBytesMetric  *prometheus.Desc
	transmitBytesMetric *prometheus.Desc

	cellIDMetric          *prometheus.Desc
	downlinkChannelMetric *prometheus.Desc
	rsrpMetric            *prometheus.Desc
	rsrqMetric            *prometheus.Desc
	rssiMetric            *prometheus.Desc
	snrMetric             *prometheus.Desc
}

func NewRadioStatsCollector(target string) *radioStatsCollector {
	return &radioStatsCollector{
		target: target,

		receiveBytesMetric: prometheus.NewDesc(
			metricName("radio", "receive_bytes_total"),
			"How many bytes have been received by the cellular radio",
			nil, nil,
		),

		transmitBytesMetric: prometheus.NewDesc(
			metricName("radio", "transmit_bytes_total"),
			"How many bytes have been transmitted by the cellular radio",
			nil, nil,
		),

		cellIDMetric: prometheus.NewDesc(
			metricName("radio", "cell_id"),
			"Physical cell identifier",
			[]string{"type", "band"}, nil,
		),

		downlinkChannelMetric: prometheus.NewDesc(
			metricName("radio", "downlink_channel"),
			"Cellular channel number",
			[]string{"type", "band"}, nil,
		),

		rsrpMetric: prometheus.NewDesc(
			metricName("radio", "rsrp"),
			"Cellular reference signal RX power",
			[]string{"type", "band"}, nil,
		),

		rsrqMetric: prometheus.NewDesc(
			metricName("radio", "rsrq"),
			"Cellular reference signal RX quality",
			[]string{"type", "band"}, nil,
		),

		rssiMetric: prometheus.NewDesc(
			metricName("radio", "rssi"),
			"Cellular received signal strength indicator",
			[]string{"type", "band"}, nil,
		),

		snrMetric: prometheus.NewDesc(
			metricName("radio", "snr"),
			"Cellular signal to noise ratio",
			[]string{"type", "band"}, nil,
		),
	}
}

func (c *radioStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.receiveBytesMetric
	ch <- c.transmitBytesMetric
	ch <- c.cellIDMetric
	ch <- c.downlinkChannelMetric
	ch <- c.rsrpMetric
	ch <- c.rsrqMetric
	ch <- c.rssiMetric
	ch <- c.snrMetric
}

func (c *radioStatsCollector) Collect(ch chan<- prometheus.Metric) {
	m := new(models.RadioStats)
	if err := scrape(c.target, m); err != nil {
		log.Printf("Failed to collect radio stats from %s: %v", c.target, err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	var receiveTotal, transmitTotal int
	for _, entry := range m.CellularStats {
		receiveTotal += entry.BytesReceived
		transmitTotal += entry.BytesSent
	}
	ch <- prometheus.MustNewConstMetric(c.receiveBytesMetric, prometheus.CounterValue, float64(receiveTotal))
	ch <- prometheus.MustNewConstMetric(c.transmitBytesMetric, prometheus.CounterValue, float64(transmitTotal))

	for _, entry := range m.StatsLTE {
		c.collectBand(ch, "LTE", entry.Stats)
	}

	for _, entry := range m.Stats5G {
		c.collectBand(ch, "5G", entry.Stats)
	}
}

func (c *radioStatsCollector) collectBand(ch chan<- prometheus.Metric, band string, m models.BandStats) {
	labels := []string{band, m.Band}
	switch band {
	case "LTE":
		ch <- prometheus.MustNewConstMetric(c.downlinkChannelMetric, prometheus.GaugeValue, float64(m.DownlinkLTE), labels...)
	case "5G":
		ch <- prometheus.MustNewConstMetric(c.downlinkChannelMetric, prometheus.GaugeValue, float64(m.Downlink5G), labels...)
	}
	ch <- prometheus.MustNewConstMetric(c.cellIDMetric, prometheus.GaugeValue, float64(m.CellID), labels...)
	ch <- prometheus.MustNewConstMetric(c.rsrpMetric, prometheus.GaugeValue, float64(m.RSRP), labels...)
	ch <- prometheus.MustNewConstMetric(c.rsrqMetric, prometheus.GaugeValue, float64(m.RSRQ), labels...)
	ch <- prometheus.MustNewConstMetric(c.rssiMetric, prometheus.GaugeValue, float64(m.RSSI), labels...)
	ch <- prometheus.MustNewConstMetric(c.snrMetric, prometheus.GaugeValue, float64(m.SNR), labels...)
}
