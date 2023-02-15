package xdp

import (
	"errors"
	"os"
	"xdpEngine/utils"
)

func ListenExit() {
	logger.Println("开始监听系统退出信号...")
	for {
		select {
		case <-CtrlC:
			DetachIfaceXdp()
			os.Exit(0)
		}
	}
}

// InsertWhitePortMap 向Port白名单Map中插入数据
func InsertWhitePortMap(portList []int, iface string) (err error) {
	defer func() {

	}()
	logger.Printf("[%s]正在导入[ %d ]个Port白名单", iface, len(portList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		for _, port := range portList {
			err := IfaceXdpDict[iface].WhitePortMap.Upsert(port, 0)
			if err != nil {
				errlog.Println("插入Port白名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// GetAllWhitePortMap 从Map中获取Port白名单
func GetAllWhitePortMap(iface string) (portList []int, err error) {
	defer func() {

	}()
	if _, ok := IfaceXdpDict[iface]; ok {
		// 存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		nextKey, err := IfaceXdpDict[iface].WhitePortMap.GetNextKeyInt("") // 获取第一个
		if err == nil {
			portList = utils.AppendPortListDeduplicate(portList, []int{nextKey})
			for {
				nextKey, err = IfaceXdpDict[iface].WhitePortMap.GetNextKeyInt(nextKey)
				if err != nil {
					break
				}
				portList = utils.AppendPortListDeduplicate(portList, []int{nextKey})
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// DeleteWhitePortMap 删除map中指定的数据
func DeleteWhitePortMap(portList []int, iface string) (err error) {
	defer func() {

	}()
	logger.Printf("[%s]正在删除[ %d ]个Port白名单", iface, len(portList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		for _, port := range portList {
			// 先查找，如果查不到，则不删除
			_, err := IfaceXdpDict[iface].WhitePortMap.Lookup(port)
			if err != nil {
				errlog.Println("white_port lookup error:", err.Error())
				continue
			}
			// 如果可以查找匹配目标，则删除目标
			err = IfaceXdpDict[iface].WhitePortMap.Delete(port)
			if err != nil {
				errlog.Println("删除Port白名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// InsertBlackPortMap 向Port黑名单Map中插入数据
func InsertBlackPortMap(portList []int, iface string) (err error) {
	defer func() {

	}()
	logger.Printf("[%s]正在导入[ %d ]个Port黑名单", iface, len(portList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		for _, port := range portList {
			err := IfaceXdpDict[iface].BlackPortMap.Upsert(port, 0)
			if err != nil {
				errlog.Println("插入Port黑名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// GetAllBlackPortMap 从Map中获取Port黑名单
func GetAllBlackPortMap(iface string) (portList []int, err error) {
	defer func() {

	}()
	if _, ok := IfaceXdpDict[iface]; ok {
		// 存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		nextKey, err := IfaceXdpDict[iface].BlackPortMap.GetNextKeyInt("") // 获取第一个
		if err == nil {
			portList = utils.AppendPortListDeduplicate(portList, []int{nextKey})
			for {
				nextKey, err = IfaceXdpDict[iface].BlackPortMap.GetNextKeyInt(nextKey)
				if err != nil {
					break
				}
				portList = utils.AppendPortListDeduplicate(portList, []int{nextKey})
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// DeleteBlackPortMap 删除map中指定的数据
func DeleteBlackPortMap(portList []int, iface string) (err error) {
	defer func() {

	}()
	logger.Printf("[%s]正在删除[ %d ]个Port黑名单", iface, len(portList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		for _, port := range portList {
			// 先查找，如果查不到，则不删除
			_, err := IfaceXdpDict[iface].BlackPortMap.Lookup(port)
			if err != nil {
				errlog.Println("black_port lookup error:", err.Error())
				continue
			}
			// 如果可以查找匹配目标，则删除目标
			err = IfaceXdpDict[iface].BlackPortMap.Delete(port)
			if err != nil {
				errlog.Println("删除Port黑名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}
