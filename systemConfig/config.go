package systemConfig

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
)

var (
	Logger      *log.Logger
	Errlog      *log.Logger
	CtrlC       chan os.Signal
	LogFileName string = "xdpEngine.log"
)

func init() {
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
	signal.Notify(CtrlC, os.Interrupt, os.Kill)
}

// PrintBanner 输出程序Banner
func PrintBanner() {
	fmt.Println(`===================================================================================================`)
	fmt.Println(``)
	fmt.Println(`__    __   _____      _____      _______    ___      _      ____      __    ___      _    _______  `)
	fmt.Println(`\ \  / /  |  __  \   |  __  \   |  _____|  |   \    | |   /  ___ \   |  |  |   \    | |  |  _____| `)
	fmt.Println(` \ \/ /   | |  \  |  | |  \  |  | |        | |\ \   | |  |  /   \ |  |  |  | |\ \   | |  | |       `)
	fmt.Println(`  \  /    | |   | |  | |__/  |  | |_____   | | \ \  | |  | |   ___   |  |  | | \ \  | |  | |_____  `)
	fmt.Println(`  /  \    | |   | |  |  ___ /   |  _____|  | |  \ \ | |  | |  |__ |  |  |  | |  \ \ | |  |  _____| `)
	fmt.Println(` / /\ \   | |__/  |  | |        | |_____   | |   \ \| |  |  \___/ |  |  |  | |   \ \| |  | |_____  `)
	fmt.Println(`/_/  \_\  |______/   |_|        |_______|  |_|    \___|   \______/   |__|  |_|    \___|  |_______| `)
	fmt.Println(``)
	fmt.Println(`===================================================================================================`)
	fmt.Println(`[-] An XDP-based firewall`)
	fmt.Println(``)
}
