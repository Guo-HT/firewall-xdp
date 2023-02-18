package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// IsIpRightSingle 判断IP是否合法 true-合法 false-不合法
func IsIpRightSingle(ip string) (result bool) {
	address := net.ParseIP(ip)
	_, addressNet, _ := net.ParseCIDR(ip)
	if address != nil || addressNet != nil {
		// 正确
		return true
	} else {
		// 异常
		return false
	}
}

func IsIpListRight(ipList []string) bool {
	for _, ip := range ipList {
		isRight := IsIpRightSingle(ip)
		if !isRight {
			return false
		}
	}
	return true
}

func IpFormat(ipByte []byte) (ip string) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("IpFormat error:", err)
		}
	}()
	var ev struct {
		Mask int32
		Ip   int32
	}
	buf := bytes.NewBuffer(ipByte)
	if err := binary.Read(buf, binary.LittleEndian, &ev); err != nil {
		errlog.Println("IpFormat error: ", err)
	} else {
		return fmt.Sprintf("%s/%d", InetNtoA(ev.Ip), ev.Mask)
	}
	return
}

func InetNtoA(ip int32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
}
