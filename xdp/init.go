package xdp

import (
	"log"
	"os"
	"xdpEngine/systemConfig"
)

var (
	logger *log.Logger
	errlog *log.Logger
	CtrlC  chan os.Signal
)

func init() {
	logger = systemConfig.Logger
	errlog = systemConfig.Errlog
	CtrlC = systemConfig.CtrlC
	IfaceXdpDict = make(map[string]*IfaceXdpObj)
	//DetachRestXdp()
}
