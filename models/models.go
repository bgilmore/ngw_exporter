package models

import "net/url"

type Model interface {
	Path() string
	Query() string
}

func URL(m Model, target string) string {
	addr := url.URL{
		Scheme:   "http",
		Host:     target,
		Path:     m.Path(),
		RawQuery: m.Query(),
	}
	return addr.String()
}

// GatewayInfo is returned from /main_web_app.cgi
type GatewayInfo struct {
	Model string `json:"gwmodel"`

	Device []struct {
		HWVersion     string `json:"HardwareVersion"`
		Product       string `json:"ProductClass"`
		SWVersion     string `json:"SoftwareVersion"`
		Serial        string `json:"SerialNumber"`
		UptimeSeconds int    `json:"UpTime"`
	} `json:"device_app_status"`

	Client []struct {
		Active        int    `json:"Active"`
		AddressSource string `json:"AddressSource"`
		InterfaceType string `json:"InterfaceType"`
	} `json:"device_cfg"`

	// TODO: work out these linkStatus flags.
	LinkLTE []struct {
		Status string `json:"linkStatus"`
	} `json:"link_status_LTE"`
	Link5G []struct {
		Status string `json:"linkStatus"`
	} `json:"link_status_5G"`
}

func (GatewayInfo) Path() string  { return "/main_web_app.cgi" }
func (GatewayInfo) Query() string { return "" }

// RadioStats is returned from /fastmile_radio_status_web_app.cgi
type RadioStats struct {
	CellularStats []struct {
		BytesReceived int `json:"BytesReceived"`
		BytesSent     int `json:"BytesSent"`
	} `json:"cellular_stats"`

	Stats5G []struct {
		Stats BandStats `json:"stat"`
	} `json:"cell_5G_stats_cfg"`

	StatsLTE []struct {
		Stats BandStats `json:"stat"`
	} `json:"cell_LTE_stats_cfg"`
}

func (RadioStats) Path() string  { return "/fastmile_radio_status_web_app.cgi" }
func (RadioStats) Query() string { return "" }

// BandStats is used for RadioStats
type BandStats struct {
	Band         string `json:"Band"`
	CellID       int    `json:"PhysicalCellID"`
	DownlinkLTE  int    `json:"DownlinkEarfcn"`
	Downlink5G   int    `json:"Downlink_NR_ARFCN"`
	RSRP         int    `json:"RSRPCurrent"`
	RSRPStrength int    `json:"RSRPStrengthIndexCurrent"`
	RSRQ         int    `json:"RSRQCurrent"`
	RSSI         int    `json:"RSSICurrent"`
	SNR          int    `json:"SNRCurrent"`
}

// NetworkStats is returned by /lan_status_web_app.cgi?lan
type NetworkStats struct {
	Router struct {
		Enabled int    `json:"Enable"`
		IPAddr  string `json:"IPInterfaceIPAddress"`
		IPMask  string `json:"IPInterfaceSubnetMask"`
	} `json:"lan_ifip"`

	WLAN []struct {
		RadioEnabled    int    `json:"RadioEnabled"`
		NetworkEnabled  int    `json:"Enable"`
		SSID            string `json:"SSID"`
		Channel         int    `json:"Channel"`
		BytesSent       int    `json:"TotalBytesSent"`
		BytesReceived   int    `json:"TotalBytesReceived"`
		PacketsSent     int    `json:"TotalPacketsSent"`
		PacketsReceived int    `json:"TotalPacketsReceived"`
		Associations    int    `json:"TotalAssociations"`
		PowerValue      int    `json:"X_CT_COM_PowerValue"`
		PowerLevel      int    `json:"X_CT_COM_Powerlevel"`
		TransmitPower   int    `json:"TransmitPower"`
		RXErrors        int    `json:"X_ASB_COM_RxErrors"`
		RXDrops         int    `json:"X_ASB_COM_RxDrops"`
		TXErrors        int    `json:"X_ASB_COM_TxErrors"`
		TXDrops         int    `json:"X_ASB_COM_TxDrops"`
	} `json:"wlan_status_glb"`

	LAN []struct {
		PortEnabled int    `json:"Enable"`
		PortState   string `json:"Status"`
		MacAddr     string `json:"MACAddress"`
		MaxBitRate  string `json:"MaxBitRate"`
		DuplexMode  string `json:"DuplexMode"`
		Stats       struct {
			BytesSent                   int `json:"BytesSent"`
			BytesReceived               int `json:"BytesReceived"`
			PacketsSent                 int `json:"PacketsSent"`
			PacketsReceived             int `json:"PacketsReceived"`
			TXErrors                    int `json:"ErrorsSent"`
			RXErrors                    int `json:"ErrorsReceived"`
			UnicastPacketsSent          int `json:"UnicastPacketsSent"`
			UnicastPacketsReceived      int `json:"UnicastPacketsReceived"`
			DiscardPacketsSent          int `json:"DiscardPacketsSent"`
			DiscardPacketsReceived      int `json:"DiscardPacketsReceived"`
			MulticastPacketsSent        int `json:"MulticastPacketsSent"`
			MulticastPacketsReceived    int `json:"MulticastPacketsReceived"`
			BroadcastPacketsSent        int `json:"BroadcastPacketsSent"`
			BroadcastPacketsReceived    int `json:"BroadcastPacketsReceived"`
			UnknownProtoPacketsReceived int `json:"UnknownProtoPacketsReceived"`
		} `json:"stat"`
	} `json:"lan_ether"`
}

func (NetworkStats) Path() string  { return "/lan_status_web_app.cgi" }
func (NetworkStats) Query() string { return "lan" }
