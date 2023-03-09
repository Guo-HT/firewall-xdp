package xdp

import (
	"context"
	"github.com/dropbox/goebpf"
	"os/exec"
	"sync"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
)

var (
	XDP_FILE                 string = "./ebpf_prog/my_xdp.elf"
	EBPF_MAP_WHITE_PORT      string = "white_port"
	EBPF_MAP_BLACK_PORT      string = "black_port"
	EBPF_MAP_WHITE_IP        string = "white_ip"
	EBPF_MAP_BLACK_IP        string = "black_ip"
	EBPF_MAP_PROTO           string = "proto_detect"
	EBPF_MAP_FUNCTION_SWITCH string = "function_switch"
	XDP_PROGRAM_NAME         string = "firewall"

	IfaceXdpDict map[string]*IfaceXdpObj // 多网口下存储策略
)

type IfaceXdpObj struct {
	Iface string

	WhitePortMap      goebpf.Map     // port白名单
	BlackPortMap      goebpf.Map     // port黑名单
	WhiteIpMap        goebpf.Map     // ip白名单
	BlackIpMap        goebpf.Map     // ip黑名单
	ProtoDetectMap    goebpf.Map     // 协议ip-port
	FunctionSwitchMap goebpf.Map     // ip黑名单
	FirewallProgram   goebpf.Program // xdp程序

	WhitePortList []int    // port白名单
	BlackPortList []int    // port黑名单
	WhiteIpList   []string // ip白名单
	BlackIpList   []string // ip黑名单

	/*
		SessionFlow:
		会话流表，其中
		key: 会话中响应IP、port二元组（服务端IP-PORT）;
		eg: 会话中，192.168.113.128发送请求，192.168.113.1回复响应，
			则key: "192.168.113.1_80"
	*/
	SessionFlow map[string]*utils.SessionTuple

	ChannelListLength int                    // 缓存中通道数量
	ProtoSwitch       bool                   // 协议分析开关
	ProtoPoolChannel  []chan utils.FiveTuple // 协议缓冲队列

	CtxP    context.Context    // 协议上下文
	CancelP context.CancelFunc // 协议上下文信号

	Ctx    context.Context    // 网口上下文
	Cancel context.CancelFunc // 网口上下文信号

	Lock sync.RWMutex
}

/*
	InitEBpfMap:
	通过clang编译得到的.elf文件，获取内部的map和program，并挂载到指定网卡
*/
func InitEBpfMap(iface string) {
	isExist := utils.IsIfaceExist(iface)
	if !isExist {
		errlog.Println("该网卡不存在")
		return
	}

	bpf := goebpf.NewDefaultEbpfSystem()
	err := bpf.LoadElf(XDP_FILE)
	if err != nil {
		errlog.Fatalf("LoadElf() failed: %v\n", err)
	}
	printBpfInfo(bpf)

	// 获取port白名单map
	mapWhitePort := bpf.GetMapByName(EBPF_MAP_WHITE_PORT)
	if mapWhitePort == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_WHITE_PORT)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_WHITE_PORT)
	}
	// 获取port黑名单map
	mapBlackPort := bpf.GetMapByName(EBPF_MAP_BLACK_PORT)
	if mapBlackPort == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_BLACK_PORT)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_BLACK_PORT)
	}
	// 获取ip白名单map
	mapWhiteIp := bpf.GetMapByName(EBPF_MAP_WHITE_IP)
	if mapWhiteIp == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_WHITE_IP)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_WHITE_IP)
	}
	// 获取ip黑名单map
	mapBlackIp := bpf.GetMapByName(EBPF_MAP_BLACK_IP)
	if mapBlackIp == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_BLACK_IP)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_BLACK_IP)
	}
	// 获取协议黑名单map
	mapProtoDetect := bpf.GetMapByName(EBPF_MAP_PROTO)
	if mapProtoDetect == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_PROTO)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_PROTO)
	}
	// 获取功能开关map
	mapFunctionSwitch := bpf.GetMapByName(EBPF_MAP_FUNCTION_SWITCH)
	if mapFunctionSwitch == nil {
		errlog.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_FUNCTION_SWITCH)
	} else {
		//logger.Printf("Get Map '%s' success\n", EBPF_MAP_FUNCTION_SWITCH)
	}

	// Program name matches function name in xdp.c:
	//      int xdp_dump(struct xdp_md *ctx)
	xdp := bpf.GetProgramByName(XDP_PROGRAM_NAME)
	if xdp == nil {
		logger.Fatalf("Program '%s' not found\n", XDP_PROGRAM_NAME)
	} else {
		//logger.Printf("Get Process '%s' success\n", XDP_PROGRAM_NAME)
	}

	// Load XDP program into kernel
	err = xdp.Load()
	if err != nil {
		errlog.Fatalf("xdp.Load(): %v\n", err)
	} else {
		//logger.Println("xdp.Load() success")
	}

	utils.UpIfaceState(iface)

	// Attach to interface
	err = xdp.Attach(iface)
	if err != nil {
		errlog.Fatalf("xdp.Attach(): %v\n", err)
	} else {
		//logger.Println("xdp.Attach() success")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctxP, cancelP := context.WithCancel(context.Background())

	logger.Println("绑定完成...")
	//defer func() {
	//	err := xdp.Detach()
	//	if err != nil {
	//		logger.Println("xdp.Detach error:", err)
	//	}else{
	//		logger.Println("xdp.Detach success")
	//	}
	//}()

	poolList := make([]chan utils.FiveTuple, systemConfig.DefaultChanNum)
	for i := 0; i < systemConfig.DefaultChanNum; i++ {
		poolList[i] = make(chan utils.FiveTuple, systemConfig.DefaultChanLength)
	}
	sessionFlow := make(map[string]*utils.SessionTuple)
	_ = SetFunctionSwitch("proto", "stop")

	thisXdpObj := IfaceXdpObj{
		Iface:             iface,
		WhitePortMap:      mapWhitePort,
		BlackPortMap:      mapBlackPort,
		WhiteIpMap:        mapWhiteIp,
		BlackIpMap:        mapBlackIp,
		ProtoDetectMap:    mapProtoDetect,
		FunctionSwitchMap: mapFunctionSwitch,
		FirewallProgram:   xdp,

		WhitePortList: []int{systemConfig.ServerPort}, // 引擎服务端口默认加白
		BlackPortList: []int{},
		WhiteIpList:   []string{},
		BlackIpList:   []string{},

		SessionFlow:      sessionFlow,
		ProtoSwitch:      systemConfig.ProtoEngineStatus,
		ProtoPoolChannel: poolList,

		ChannelListLength: systemConfig.DefaultChanNum,

		CtxP:    ctxP,
		CancelP: cancelP,

		Ctx:    ctx,
		Cancel: cancel,
	}
	IfaceXdpDict[iface] = &thisXdpObj
	_ = InsertWhitePortMap([]int{systemConfig.ServerPort}, iface) // 引擎服务端口加白
}

// printBpfInfo 输出当前.elf文件下的map和program（来源：dropbox样例）
func printBpfInfo(bpf goebpf.System) {
	logger.Println("Maps:")
	for _, item := range bpf.GetMaps() {
		m := item.(*goebpf.EbpfMap)
		logger.Printf("\t%s: %v, Fd %v\n", m.Name, m.Type, m.GetFd())
	}
	logger.Println("Programs:")
	for _, prog := range bpf.GetPrograms() {
		logger.Printf("\t%s: %v, size %d, license \"%s\"\n",
			prog.GetName(), prog.GetType(), prog.GetSize(), prog.GetLicense(),
		)

	}
	logger.Println()
}

// DetachRestXdp 卸载所有网卡上残留的xdp
func DetachRestXdp() {
	interfaces, err := utils.GetIfaceList()
	if err != nil {
		logger.Println("GetIfaceList error: ", err)
	}
	for _, iface := range interfaces {
		DetachXdp(iface)
	}
}

// DetachXdp 卸载指定网卡上的XDP
func DetachXdp(iface string) {
	// ip link set dev ens33 xdp off
	cmd := exec.Command("ip", "link", "set", "dev", iface, "xdp", "off")
	if err := cmd.Start(); err != nil {
		logger.Println("DetachRestXdp in starting xdpEngine error:", err)
	}
}

// StopAllXdpEngine 卸载已挂载的xdp程序
func StopAllXdpEngine() {
	for Iface, value := range IfaceXdpDict {
		StopXdpEngine(Iface, value)
	}
}

// StopXdpEngine 关闭指定网卡的xdp程序
func StopXdpEngine(iface string, value *IfaceXdpObj) {
	logger.Printf("[%s]XDP程序正在卸载...", iface)
	_ = value.WhitePortMap.Close()
	_ = value.BlackPortMap.Close()
	_ = value.WhiteIpMap.Close()
	_ = value.BlackIpMap.Close()
	_ = value.FunctionSwitchMap.Close()
	_ = value.ProtoDetectMap.Close()
	_ = value.FirewallProgram.Detach()
	value.CancelP()
	for _, channel := range value.ProtoPoolChannel {
		close(channel)
	}
	value.Cancel()
	value.SessionFlow = make(map[string]*utils.SessionTuple)
	delete(IfaceXdpDict, iface)
	DetachXdp(iface)
	logger.Printf("[%s]XDP程序卸载完成...", iface)
}
