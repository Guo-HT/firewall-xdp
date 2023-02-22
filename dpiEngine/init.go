package dpiEngine

import (
	"log"
	"os"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
)

var (
	logger        *log.Logger
	errlog        *log.Logger
	CtrlC         chan os.Signal
	ProtoRuleList []utils.ProtoRule // 协议规则列表
)

func init() {
	logger = systemConfig.Logger
	errlog = systemConfig.Errlog
	CtrlC = systemConfig.CtrlC
	InitProtoRules() // 初始化协议规则列表
	//DetachRestXdp()
}
