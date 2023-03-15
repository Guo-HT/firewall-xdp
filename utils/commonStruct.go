package utils

import (
	"github.com/google/gopacket/layers"
	"time"
)

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

// PortHitCount Port及命中次数
type PortHitCount struct {
	Port int `json:"port"`
	Hit  int `json:"hit"`
}

// IpHitCount Port及命中次数
type IpHitCount struct {
	IP  string `json:"ip"`
	Hit int    `json:"hit"`
}

// UserLoginForm 登录时提交数据
type UserLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangeUserPassword 修改密码时提交数据
type ChangeUserPassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type NetcardInfo struct {
	NetcardName  string   `json:"netcard_name"`
	IP           []string `json:"ip"`
	MAC          string   `json:"mac"`
	Flags        string   `json:"flags"`
	EngineStatus bool     `json:"engine_status"`
}

type UserInfo struct {
	ID       int       `json:"id"`        // ID
	UserName string    `json:"username"`  // 用户名
	Email    string    `json:"email"`     // 邮箱
	Role     int       `json:"role"`      // 用户角色[0-admin; 1-操作员; 2-审计员]
	CreateAt time.Time `json:"create_at"` // 创建时间
}

type DelUserCheck struct {
	Password       string `json:"password"`
	TargetUserName string `json:"target_user_name"`
}

// UserAdd 新增用户需要的信息
type UserAdd struct {
	CurUserPassword string `json:"cur_user_password"` // 当前操作用户的密码
	UserName        string `json:"user_name"`         // 用户名
	Password        string `json:"password"`          // 密码
	Email           string `json:"email"`             // 邮箱
	Role            int    `json:"role"`              // 用户角色[0-admin; 1-操作员; 2-访客]
}
