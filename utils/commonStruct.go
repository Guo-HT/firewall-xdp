package utils

import "github.com/google/gopacket/layers"

type WhiteIpStruct struct {
	WhiteIpList []string `json:"whiteIpList"`
	Iface       string   `json:"iface"`
}

type WhitePortStruct struct {
	WhitePortList []int  `json:"whitePortList"`
	Iface         string `json:"iface"`
}

type BlackIpStruct struct {
	BlackIpList []string `json:"blackIpList"`
	Iface       string   `json:"iface"`
}

type BlackPortStruct struct {
	BlackPortList []int  `json:"blackPortList"`
	Iface         string `json:"iface"`
}

type IfaceStruct struct {
	Iface string `json:"iface"`
}

type ProtoId struct {
	Id int `json:"id"`
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
	IsEnable     bool   `json:"is_enable"`     // 是否启用
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

// SessionTuple 会话对
type SessionTuple struct {
	ServerAddr string
	ServerPort int

	ProtoID    int   // 检测到的协议ID
	UpdateTime int64 // 最新一次检测到的时间
	HitReq     bool  // 请求命中
	HitRsp     bool  // 响应命中
}

// ProtoStatusConf 配置单个协议开关
type ProtoStatusConf struct {
	ProtoName string `json:"protoName"`
	Status    bool   `json:"status"`
}

// IpPort ip-port
type IpPort struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
