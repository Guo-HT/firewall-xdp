package utils

import (
	"net"
	"os/exec"
)

func GetLocalIP() (ipList []string) {
	return
}

func GetIfaceList() (netcardList []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}

	for _, i := range interfaces {
		netcardList = append(netcardList, i.Name)
	}
	return netcardList, nil
}

// UpIfaceState 通过"ifconfig"命令，启用指定网卡
func UpIfaceState(Iface string) {
	defer func() {
		if e := recover(); e != nil {
			logger.Println("UpIfaceState error:	", e)
		}
	}()
	cmd := exec.Command("ifconfig", Iface, "up")
	if err := cmd.Start(); err != nil {
		logger.Println("UpIfaceState in starting xdpEngine error:", err)
	}
}
