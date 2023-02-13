package xdp

import (
	"log"
	"os"
	"xdpEngine/systemConfig"
)

var (
	logger *log.Logger
	CtrlC  chan os.Signal
)

func init() {
	logger = systemConfig.Logger
	CtrlC = systemConfig.CtrlC
	//DetachRestXdp()
}
