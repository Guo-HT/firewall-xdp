package utils

type WhiteIpStruct struct {
	WhiteIpList []string `json:"whiteIpList,omitempty"`
	Iface       string   `json:"iface,omitempty"`
}

type WhitePortStruct struct {
	WhitePortList []int  `json:"whitePortList,omitempty"`
	Iface         string `json:"iface,omitempty"`
}

type BlackIpStruct struct {
	BlackIpList []string `json:"blackIpList,omitempty"`
	Iface       string   `json:"iface,omitempty"`
}

type BlackPortStruct struct {
	BlackPortList []string `json:"blackPortList,omitempty"`
	Iface         string   `json:"iface,omitempty"`
}

type IfaceStruct struct {
	Iface string `json:"iface,omitempty"`
}
