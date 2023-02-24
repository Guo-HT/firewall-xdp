package systemConfig

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var (
	Logger         *log.Logger                      // 正常日志
	Errlog         *log.Logger                      // 错误日志
	CtrlC          chan os.Signal                   // 退出信号
	RunMode        string                           // 运行模式
	LogFileName    string         = "xdpEngine.log" // 日志文件
	ServerPort     int            = 1888            // 服务监听端口
	ServerPortStr                 = strconv.Itoa(ServerPort)
	ProtoRuleFile                 = "systemConfig/rule.json" // 协议规则文件
	DefaultIface   string         = "ens33"                  // 默认监听网口
	DefaultChanNum int            = runtime.NumCPU()         // 缓冲池中channel数量为CPU核心数
)

func init() {
	PrintBanner()
	LogInit()
	ListenExit()
}

// LogInit 初始化日志相关配置
func LogInit() {

	file, err := os.OpenFile(LogFileName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	Logger = log.New(io.MultiWriter(file, os.Stdout), "[INFO] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Errlog = log.New(io.MultiWriter(file, os.Stdout), "[ERROR] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

// ListenExit 配置退出信号监听
func ListenExit() {
	CtrlC = make(chan os.Signal, 1)
	signal.Notify(CtrlC, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

// PrintBanner 输出程序Banner
func PrintBanner() {
	//fmt.Println(`===================================================================================================`)
	fmt.Println(``)
	fmt.Println(`          __     ____          _          `)
	fmt.Println(` __ _____/ /__  / __/__  ___ _(_)__  ___  `)
	fmt.Println(` \ \ / _  / _ \/ _// _ \/ _ ·/ / _ \/ -_) `)
	fmt.Println(`/_\_\\_,_/ .__/___/_//_/\_, /_/_//_/\__/  `)
	fmt.Println(`        /_/            /___/              `)
	//fmt.Println(``)
	fmt.Println(`===================================================================================================`)
	fmt.Println(`[-] An XDP-based firewall`)
	fmt.Println(`                     ———— Powered by GuoHT`)
	fmt.Println(``)
}

// SayBye 退出时输出
func SayBye() {
	fmt.Println(`.______   ____    ____  _______   `)
	fmt.Println(`|   _  \  \   \  /   / |   ____|  `)
	fmt.Println(`|  |_)  |  \   \/   /  |  |__     `)
	fmt.Println(`|   _  <    \_    _/   |   __|    `)
	fmt.Println(`|  |_)  |     |  |     |  |____   `)
	fmt.Println(`|______/      |__|     |_______|  `)
}

/*

          __     ____          _
 __ _____/ /__  / __/__  ___ _(_)__  ___
 \ \ / _  / _ \/ _// _ \/ _ ·/ / _ \/ -_)
/_\_\\_,_/ .__/___/_//_/\_, /_/_//_/\__/
        /_/            /___/

.______   ____    ____  _______
|   _  \  \   \  /   / |   ____|
|  |_)  |  \   \/   /  |  |__
|   _  <    \_    _/   |   __|
|  |_)  |     |  |     |  |____
|______/      |__|     |_______|

*/
