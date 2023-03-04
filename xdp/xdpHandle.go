package xdp

import (
	"errors"
	"github.com/dropbox/goebpf"
	"os"
	"strconv"
	"strings"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
)

func ListenExit() {
	logger.Println("[!] 开始监听系统退出信号...")
	logger.Println("[!] Ctrl+C 退出引擎")
	for {
		select {
		case <-CtrlC:
			StopAllXdpEngine()
			os.Exit(0)
		}
	}
}

// ******************** Port 操作 **************************
// *********************** 白 ******************************

// InsertWhitePortMap 向Port白名单Map中插入数据
func InsertWhitePortMap(portList []int, iface string) (err error) {
	defer func() {

	}()
	portList = utils.AppendPortListDeduplicate(portList, []int{systemConfig.ServerPort}) // 将引擎服务端口默认加白

	logger.Printf("[%s]正在导入[ %d ]个Port白名单: %v", iface, len(portList), portList)
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
			if port == systemConfig.ServerPort {
				continue // 跳过服务Port，避免阻断后无法操作
			}
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

// *********************** 黑 ******************************

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

// ********************* IP 操作 ***************************
// *********************** 白 ******************************

// InsertWhiteIpMap 向IP白名单MAP中插入数据
func InsertWhiteIpMap(ipList []string, iface string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("InsertWhiteIpMap error, ", err)
		}
	}()
	logger.Printf("[%s]正在导入[ %d ]个IP白名单", iface, len(ipList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()
		defer IfaceXdpDict[iface].Lock.RUnlock()
		// ********************** 删除所有 *************************
		logger.Println("正在预删除所有IP白名单")
		for {
			nextKey, err := IfaceXdpDict[iface].WhiteIpMap.GetNextKey("")
			if err != nil {
				//errlog.Println("InsertWhiteIpMap GetNextKey error:", err)
				break
			}
			logger.Println("删除一个,", utils.IpFormat(nextKey))
			err = IfaceXdpDict[iface].WhiteIpMap.Delete(nextKey)
			if err != nil {
				errlog.Println("InsertWhiteIpMap preDelete error:", err)
				break
			}
		}
		// ************************ 写入 ***************************
		logger.Println("正在插入所有IP白名单")
		for _, ip := range ipList {
			logger.Println("正在导入IP白名单: ", goebpf.CreateLPMtrieKey(ip))
			err := IfaceXdpDict[iface].WhiteIpMap.Insert(goebpf.CreateLPMtrieKey(ip), 0)
			if err != nil {
				errlog.Println("插入IP白名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// GetAllWhiteIpMap 从Map中获取IP白名单
func GetAllWhiteIpMap(iface string) (ipList []string, err error) {
	defer func() {

	}()
	if _, ok := IfaceXdpDict[iface]; ok {
		// 存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		nextKey, err := IfaceXdpDict[iface].WhiteIpMap.GetNextKey("")
		for {
			if err != nil {
				//errlog.Println("GetAllWhiteIpMap GetNextKey error:", err)
				break
			}
			ipString := utils.IpFormat(nextKey)
			logger.Println("获取白名单IP,", ipString)
			ipList = append(ipList, ipString)
			nextKey, err = IfaceXdpDict[iface].WhiteIpMap.GetNextKey(nextKey)
			if err != nil {
				errlog.Println("GetAllWhiteIpMap GetNextKey error:", err)
				break
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// DeleteWhiteIpMap 删除IP白名单Map中指定的数据
func DeleteWhiteIpMap(ipList []string, iface string) (err error) {
	defer func() {

	}()
	logger.Printf("[%s]正在删除[ %d ]个IP白名单", iface, len(ipList))
	if _, ok := IfaceXdpDict[iface]; ok {
		//iface存在
		IfaceXdpDict[iface].Lock.RLock()
		defer IfaceXdpDict[iface].Lock.RUnlock()

		for _, ip := range ipList {
			//先查找，查不到跳过
			_, err := IfaceXdpDict[iface].WhiteIpMap.Lookup(goebpf.CreateLPMtrieKey(ip))
			if err != nil {
				errlog.Println("white_ip lookup error:", err.Error())
				continue
			}
			// 如果可以查找到匹配，则删除
			err = IfaceXdpDict[iface].WhiteIpMap.Delete(goebpf.CreateLPMtrieKey(ip))
			if err != nil {
				errlog.Println("删除IP白名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// *********************** 黑 ******************************

// InsertBlackIpMap 向IP黑名单MAP中插入数据
func InsertBlackIpMap(ipList []string, iface string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("InsertBlackIpMap error, ", err)
		}
	}()
	logger.Printf("[%s]正在导入[ %d ]个IP黑名单", iface, len(ipList))
	if _, ok := IfaceXdpDict[iface]; ok {
		// iface存在
		IfaceXdpDict[iface].Lock.RLock()
		defer IfaceXdpDict[iface].Lock.RUnlock()
		// ********************** 删除所有 *************************
		logger.Println("正在预删除所有IP黑名单")
		for {
			nextKey, err := IfaceXdpDict[iface].BlackIpMap.GetNextKey("")
			if err != nil {
				//errlog.Println("InsertBlackIpMap GetNextKey error:", err)
				break
			}
			logger.Println("删除一个,", utils.IpFormat(nextKey))
			err = IfaceXdpDict[iface].BlackIpMap.Delete(nextKey)
			if err != nil {
				errlog.Println("InsertBlackIpMap preDelete error:", err)
				break
			}
		}
		// ************************ 写入 ***************************
		logger.Println("正在插入所有IP黑名单")
		for _, ip := range ipList {
			logger.Println("正在导入IP黑名单: ", goebpf.CreateLPMtrieKey(ip))
			err := IfaceXdpDict[iface].BlackIpMap.Insert(goebpf.CreateLPMtrieKey(ip), 0)
			if err != nil {
				errlog.Println("插入IP黑名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// GetAllBlackIpMap 从Map中获取IP黑名单
func GetAllBlackIpMap(iface string) (ipList []string, err error) {
	defer func() {

	}()
	if _, ok := IfaceXdpDict[iface]; ok {
		// 存在
		IfaceXdpDict[iface].Lock.RLock()         // 读锁
		defer IfaceXdpDict[iface].Lock.RUnlock() // 解锁

		nextKey, err := IfaceXdpDict[iface].BlackIpMap.GetNextKey("")
		for {
			if err != nil {
				//errlog.Println("GetAllBlackIpMap GetNextKey error:", err)
				break
			}
			ipString := utils.IpFormat(nextKey)
			logger.Println("获取黑名单IP,", ipString)
			ipList = append(ipList, ipString)
			nextKey, err = IfaceXdpDict[iface].BlackIpMap.GetNextKey(nextKey)
			if err != nil {
				errlog.Println("GetAllBlackIpMap GetNextKey error:", err)
				break
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// DeleteBlackIpMap 删除IP黑名单Map中指定的数据
func DeleteBlackIpMap(ipList []string, iface string) (err error) {
	defer func() {

	}()

	logger.Printf("[%s]正在删除[ %d ]个IP黑名单", iface, len(ipList))
	if _, ok := IfaceXdpDict[iface]; ok {
		//iface存在
		IfaceXdpDict[iface].Lock.RLock()
		defer IfaceXdpDict[iface].Lock.RUnlock()

		for _, ip := range ipList {
			//先查找，查不到跳过
			_, err := IfaceXdpDict[iface].BlackIpMap.Lookup(goebpf.CreateLPMtrieKey(ip))
			if err != nil {
				errlog.Println("black_ip lookup error:", err.Error())
				continue
			}
			// 如果可以查找到匹配，则删除
			err = IfaceXdpDict[iface].BlackIpMap.Delete(goebpf.CreateLPMtrieKey(ip))
			if err != nil {
				errlog.Println("删除IP黑名单失败, ", err.Error())
			}
		}
	} else {
		err = errors.New("该网卡无效")
	}
	return
}

// *********************** 协议 *****************************

// UpdateProtoIpPortMap 将协议分析结果以IP-PORT形式写入对应MAP
func UpdateProtoIpPortMap() (err error) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("UpdateProtoIpPortMap error, ", err)
		}
	}()
	for iface, xdpObj := range IfaceXdpDict {
		IfaceXdpDict[iface].Lock.RLock()
		logger.Printf("[%s]正在更新协议阻断Map: %v", iface, xdpObj.SessionFlow)
		// ********************** 删除所有 *************************
		for {
			// [192 168 113 1 87 52 0 0]
			nextKey, err := IfaceXdpDict[iface].ProtoDetectMap.GetNextKey("")
			if err != nil {
				//errlog.Println("UpdateProtoIpPortMap GetNextKey error:", err)
				break
			}
			logger.Println("删除一个,", utils.IpPortFormat(nextKey))
			err = IfaceXdpDict[iface].ProtoDetectMap.Delete(nextKey)
			if err != nil {
				errlog.Println("UpdateProtoIpPortMap preDelete error:", err)
				break
			}
		}
		// ************************ 写入 ***************************
		for _, session := range xdpObj.SessionFlow {
			logger.Printf("正在导入协议阻断数据: %+v", utils.IpPort2Byte(session.ServerAddr, session.ServerPort))
			err := IfaceXdpDict[iface].ProtoDetectMap.Insert(utils.IpPort2Byte(session.ServerAddr, session.ServerPort), 0)
			if err != nil {
				errlog.Println("插入协议阻断数据, ", err.Error())
			}
		}

		IfaceXdpDict[iface].Lock.RUnlock()
	}

	return
}

// GetAllProtoIpPortMap 查询MAP中所有协议分析结果
func GetAllProtoIpPortMap() (ipPortList []utils.IpPort) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("GetAllProtoIpPortMap error, ", err)
		}
	}()
	for iface, xdpObj := range IfaceXdpDict {
		IfaceXdpDict[iface].Lock.RLock()

		nextKey, err := xdpObj.ProtoDetectMap.GetNextKey("")
		for {
			if err != nil {
				//errlog.Println("GetAllProtoIpPortMap GetNextKey error:", err)
				break
			}
			ipPort := utils.IpPortFormat(nextKey)
			ipPortArray := strings.Split(ipPort, ":")
			port, err := strconv.Atoi(ipPortArray[1])
			if err != nil {
				errlog.Println("GetAllProtoIpPortMap Atoi error: ", err.Error())
			}
			ipPortList = append(ipPortList, utils.IpPort{
				IP:   ipPortArray[0],
				Port: port,
			})
			nextKey, err = xdpObj.ProtoDetectMap.GetNextKey(nextKey)
			if err != nil {
				errlog.Println("GetAllProtoIpPortMap GetNextKey error:", err)
				break
			}
		}
		IfaceXdpDict[iface].Lock.RUnlock()
	}
	return
}

// ********************* 功能开关 ****************************

// SetFunctionSwitch 配置功能开关, ["proto", ...]
func SetFunctionSwitch(funcName string, mode string) (err error) {
	defer func() {

	}()

	var switchFlag uint32 // map中对应的key标识
	var modeFlag uint32   // 开关状态标识
	switch funcName {
	case "proto":
		switchFlag = systemConfig.Func2flag[funcName]
		if mode == "start" {
			modeFlag = 1 // 协议开->1
		} else {
			modeFlag = 0 // 协议关->2
		}
	default:
		// 没有匹配的功能
		switchFlag = 0
	}
	if switchFlag == 0 {
		errlog.Printf("SetFunctionSwitch error: 不存在该功能")
		return errors.New("不存在该功能")
	}
	for iface, value := range IfaceXdpDict {
		// 遍历所有网口
		err := value.FunctionSwitchMap.Upsert(switchFlag, modeFlag)
		if err != nil {
			errlog.Printf("[%s]SetFunctionSwitch Upsert error: %v", iface, err.Error())
			return errors.New("[" + iface + "]SetFunctionSwitch Upsert error: " + err.Error())
		}
	}
	return
}

// GetFunctionSwitch 从Map中获取分析功能开关状态
func GetFunctionSwitch(funcName string, iface string) (status int, err error) {
	switchFlag := systemConfig.Func2flag[funcName] // map中对应的key标识
	switch funcName {
	case "proto":
		status, err := IfaceXdpDict[iface].FunctionSwitchMap.LookupInt(switchFlag)
		if err != nil {
			errlog.Printf("[%s]获取%s开关状态失败, 默认返回0, %v", iface, funcName, err)
			return 0, nil
		}
		return status, nil
	default:
		err = errors.New("暂无此功能")
		return
	}
}
