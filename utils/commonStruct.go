package utils

import "github.com/google/gopacket/layers"

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
	BlackPortList []int  `json:"blackPortList,omitempty"`
	Iface         string `json:"iface,omitempty"`
}

type IfaceStruct struct {
	Iface string `json:"iface,omitempty"`
}

// ProtoRule 协议规则库
type ProtoRule struct {
	Id           int    `json:"id"`            // id
	ProtocolName string `json:"protocol_name"` // 协议名称
	ReqType      string `json:"req_type"`      // 请求类型
	ReqRegx      string `json:"req_regx"`      // 请求正则
	RspType      string `json:"rsp_type"`      // 相应类型
	RspRegx      string `json:"rsp_regx"`      // 响应正则
	StartPort    int    `json:"start_port"`    // 起始端口
	EndPort      int    `json:"end_port"`      // 结束端口
}

// FiveTuple 报文解析后结果
type FiveTuple struct {
	Protocol   layers.IPProtocol // IP上层协议
	SrcAddr    string            // 源IP
	DstAddr    string            // 目的IP
	SrcPort    int               // 源端口
	DstPort    int               // 目的端口
	NextAck    uint32            //
	LayerTcp   *layers.TCP       // TCP
	Payload    []byte            // 报文内容
	PayloadHex []byte            // 报文内容
	Iface      string            // 网卡
}
