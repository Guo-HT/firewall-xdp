package systemConfig

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	signal.Notify(CtrlC, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

// PrintBanner 输出程序Banner
func PrintBanner() {
	fmt.Println(`===================================================================================================`)
	fmt.Println(``)
	fmt.Println(`          __     ____          _`)
	fmt.Println(` __ _____/ /__  / __/__  ___ _(_)__  ___`)
	fmt.Println(` \ \ / _  / _ \/ _// _ \/ _ ·/ / _ \/ -_)`)
	fmt.Println(`/_\_\\_,_/ .__/___/_//_/\_, /_/_//_/\__/`)
	fmt.Println(`        /_/            /___/`)
	fmt.Println(``)
	fmt.Println(`===================================================================================================`)
	fmt.Println(`[-] An XDP-based firewall`)
	fmt.Println(``)
}

/*

          __     ____          _
 __ _____/ /__  / __/__  ___ _(_)__  ___
 \ \ / _  / _ \/ _// _ \/ _ ·/ / _ \/ -_)
/_\_\\_,_/ .__/___/_//_/\_, /_/_//_/\__/
        /_/            /___/

*/
