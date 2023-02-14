package xdp

import (
	"context"
	"github.com/dropbox/goebpf"
	"os/exec"
	"sync"
	"xdpEngine/utils"
)

var (
	XDP_FILE            string = "./ebpf_prog/my_xdp.elf"
	EBPF_MAP_WHITE_PORT string = "white_port"
	EBPF_MAP_BLACK_PORT string = "black_port"
	EBPF_MAP_WHITE_IP   string = "white_ip"
	EBPF_MAP_BLACK_IP   string = "black_ip"
	XDP_PROGRAM_NAME    string = "firewall"
	DEFAULT_IFACE       string = "ens33"

	IfaceXdpDict map[string]*IfaceXdpObj // 多网口下存储策略
)

type IfaceXdpObj struct {
	Iface string

	WhitePortMap    goebpf.Map     // port白名单
	BlackPortMap    goebpf.Map     // port黑名单
	WhiteIpMap      goebpf.Map     // ip白名单
	BlackIpMap      goebpf.Map     // ip黑名单
	FirewallProgram goebpf.Program // xdp程序

	WhitePortList []int    // port白名单
	BlackPortList []int    // port黑名单
	WhiteIpList   []string // ip白名单
	BlackIpList   []string // ip黑名单

	Ctx    context.Context    // 上下文
	Cancel context.CancelFunc // 上下文信号

	Lock sync.RWMutex
}

/*
	InitEBpfMap:
	通过clang编译得到的.elf文件，获取内部的map和program，并挂载到指定网卡
*/
func InitEBpfMap() {

	bpf := goebpf.NewDefaultEbpfSystem()
	err := bpf.LoadElf(XDP_FILE)
	if err != nil {
		logger.Fatalf("LoadElf() failed: %v\n", err)
	}
	printBpfInfo(bpf)

	// 获取port白名单map
	mapWhitePort := bpf.GetMapByName(EBPF_MAP_WHITE_PORT)
	if mapWhitePort == nil {
		logger.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_WHITE_PORT)
	} else {
		logger.Printf("Get Map '%s' success\n", EBPF_MAP_WHITE_PORT)
	}
	// 获取port黑名单map
	mapBlackPort := bpf.GetMapByName(EBPF_MAP_BLACK_PORT)
	if mapBlackPort == nil {
		logger.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_BLACK_PORT)
	} else {
		logger.Printf("Get Map '%s' success\n", EBPF_MAP_BLACK_PORT)
	}
	// 获取ip白名单map
	mapWhiteIp := bpf.GetMapByName(EBPF_MAP_WHITE_IP)
	if mapWhiteIp == nil {
		logger.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_WHITE_IP)
	} else {
		logger.Printf("Get Map '%s' success\n", EBPF_MAP_WHITE_IP)
	}
	// 获取ip黑名单ip
	mapBlackIp := bpf.GetMapByName(EBPF_MAP_BLACK_IP)
	if mapBlackIp == nil {
		logger.Fatalf("eBPF map '%s' not found\n", EBPF_MAP_BLACK_IP)
	} else {
		logger.Printf("Get Map '%s' success\n", EBPF_MAP_BLACK_IP)
	}

	// Program name matches function name in xdp.c:
	//      int xdp_dump(struct xdp_md *ctx)
	xdp := bpf.GetProgramByName(XDP_PROGRAM_NAME)
	if xdp == nil {
		logger.Fatalf("Program '%s' not found\n", XDP_PROGRAM_NAME)
	} else {
		logger.Printf("Get Process '%s' success\n", XDP_PROGRAM_NAME)
	}

	// Load XDP program into kernel
	err = xdp.Load()
	if err != nil {
		logger.Fatalf("xdp.Load(): %v\n", err)
	} else {
		logger.Println("xdp.Load() success")
	}

	utils.UpIfaceState(DEFAULT_IFACE)

	// Attach to interface
	err = xdp.Attach(DEFAULT_IFACE)
	if err != nil {
		logger.Fatalf("xdp.Attach(): %v\n", err)
	} else {
		logger.Println("xdp.Attach() success")
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger.Println("绑定完成...")
	//defer func() {
	//	err := xdp.Detach()
	//	if err != nil {
	//		logger.Println("xdp.Detach error:", err)
	//	}else{
	//		logger.Println("xdp.Detach success")
	//	}
	//}()

	IfaceXdpDict[DEFAULT_IFACE] = &IfaceXdpObj{
		Iface:           DEFAULT_IFACE,
		WhitePortMap:    mapWhitePort,
		BlackPortMap:    mapBlackPort,
		WhiteIpMap:      mapWhiteIp,
		BlackIpMap:      mapBlackIp,
		FirewallProgram: xdp,
		WhitePortList:   []int{},
		BlackPortList:   []int{},
		WhiteIpList:     []string{},
		BlackIpList:     []string{},
		Ctx:             ctx,
		Cancel:          cancel,
	}
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
		// ip link set dev ens33 xdp off
		cmd := exec.Command("ip", "link", "set", "dev", iface, "xdp", "off")
		if err := cmd.Start(); err != nil {
			logger.Println("DetachRestXdp in starting xdpEngine error:", err)
		}
	}
}

// DetachIfaceXdp 卸载已挂在的xdp程序
func DetachIfaceXdp() {
	for Iface, value := range IfaceXdpDict {
		logger.Printf("[%s]XDP程序正在卸载...", Iface)
		_ = value.FirewallProgram.Detach()
		value.Cancel()
		logger.Printf("[%s]XDP程序卸载完成", Iface)
	}
}
