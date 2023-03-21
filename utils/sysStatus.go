package utils

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

const durationTime = 500

// GetCpuPercent 获取CPU占用
func GetCpuPercent(chanCpuLoad chan float64) {
	loads, err := cpu.Percent(time.Second, false)
	if err != nil {
		errlog.Println("GetCpuPercent error", err.Error())
		chanCpuLoad <- float64Decimal2(float64(0))
		return
	}
	chanCpuLoad <- float64Decimal2(loads[0])
	return
}

// GetMemPercent 获取内存占用
func GetMemPercent() (memPercent float64) {
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		errlog.Println("GetMemPercent error", err.Error())
		memPercent = 0
		return
	}
	memPercent = float64Decimal2(virtualMem.UsedPercent)
	return
}

// GetDiskPercent 获取磁盘占用
func GetDiskPercent() (diskPercent float64) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		errlog.Println("GetDiskPercent error", err.Error())
		diskPercent = 0
		return
	}
	diskPercent = float64Decimal2(diskUsage.UsedPercent)
	return
}

// GetCpuTemperature 获取CPU温度
func GetCpuTemperature() (temperature float64) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		errlog.Println("GetCpuTemperature error:", err.Error())
		return 40
	}
	if len(temps) == 0 {
		return 40
	}
	for index, temp := range temps {
		temperature = (temperature*float64(index) + temp.Temperature) / (float64(index) + 1)
	}
	return
}

// GetAllNetcardIOSpeed 获取所有网卡的IO速率
func GetAllNetcardIOSpeed(ioSpeed chan []IOSpeed) {
	var netSpeed []IOSpeed
	preIO, err := NetIOCount()
	if err != nil {
		errlog.Println("NetIOCount pre error:", err.Error())
		ioSpeed <- []IOSpeed{}
		return
	}
	time.Sleep(time.Millisecond * durationTime) // 0.5s
	afterIO, err := NetIOCount()
	if err != nil {
		errlog.Println("NetIOCount after error:", err.Error())
		ioSpeed <- []IOSpeed{}
		return
	}
	for _, pre := range preIO {
		for _, after := range afterIO {
			if pre.Name == after.Name {
				sentSpeed := calcNetSpeed(pre.BytesSent, after.BytesSent)
				recvSpeed := calcNetSpeed(pre.BytesRecv, after.BytesRecv)
				eachNetSpeed := IOSpeed{
					pre.Name,
					sentSpeed,
					recvSpeed,
				}
				netSpeed = append(netSpeed, eachNetSpeed)
			}
		}
	}
	ioSpeed <- netSpeed
	return
}

// NetIOCount 获取网卡IO数据量
func NetIOCount() (netStat []IOStat, err error) {
	defer func() {

	}()
	ioStat, err := net.IOCounters(true)
	if err != nil {
		errlog.Println("NetIOCount net.IOCounters error:", err.Error())
		err = errors.New(err.Error())
		return
	}
	for _, eachNet := range ioStat {
		eachNetstat := IOStat{
			eachNet.Name,
			eachNet.BytesSent,
			eachNet.BytesRecv,
		}
		netStat = append(netStat, eachNetstat)
	}
	return netStat, nil
}

func calcNetSpeed(pre, after uint64) (netSpeed uint64) {
	netSpeed = ((after - pre) * 1000) / durationTime
	return
}
