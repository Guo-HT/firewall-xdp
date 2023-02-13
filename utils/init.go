package utils

import (
	"log"
	"xdpEngine/systemConfig"
)

var (
	logger *log.Logger
)

func init() {
	logger = systemConfig.Logger
}
