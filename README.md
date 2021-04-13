# ngw_exporter

A [Prometheus](https://prometheus.io/) exporter for the Nokia FastMile 5G internet gateway, or "**n**okia **g**ate**w**ay **exporter**"

This implementation was thrown together using the rebranded T-Mobile version dubbed "T-Mobile High-Speed Internet Gateway," but the exporter will likely work with units sold by other carriers and potentially with the Nokia Beacon more broadly.

## example usage
```
% ngw_exporter -h
Usage of ./ngw_exporter:
  -listen string
    	Exporter listen address (default ":9099")
  -scrape_timeout duration
    	HTTP timeout for gateway scrapes (default 10s)
  -target string
    	IP address of your Nokia Gateway device

% ngw_exporter -target 192.168.12.1
exporter.go:39: listening for connections at :9099

% curl -s 127.0.0.1:9099/metrics | grep _receive_bytes_total
# HELP ngw_network_lan_receive_bytes_total Number of received bytes
# TYPE ngw_network_lan_receive_bytes_total counter
ngw_network_lan_receive_bytes_total{port="1"} 16368
ngw_network_lan_receive_bytes_total{port="2"} 127406
# HELP ngw_network_wlan_receive_bytes_total Number of received bytes
# TYPE ngw_network_wlan_receive_bytes_total counter
ngw_network_wlan_receive_bytes_total{channel="1",ssid="banana"} 359197
ngw_network_wlan_receive_bytes_total{channel="136",ssid="banana"} 8.0054944e+07
ngw_network_wlan_receive_bytes_total{channel="64",ssid="banana"} 5.9869662e+07
# HELP ngw_radio_receive_bytes_total How many bytes have been received by the cellular radio
# TYPE ngw_radio_receive_bytes_total counter
ngw_radio_receive_bytes_total 1.132949765e+09
```
