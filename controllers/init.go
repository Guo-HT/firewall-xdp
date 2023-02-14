package controllers

import (
	"log"
	"xdpEngine/systemConfig"
)

var (
	logger *log.Logger
	errlog *log.Logger
)

func init() {
	logger = systemConfig.Logger
	errlog = systemConfig.Errlog

}
