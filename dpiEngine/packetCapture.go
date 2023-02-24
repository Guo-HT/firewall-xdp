package dpiEngine

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math/rand"
	"os"
	"time"
	"xdpEngine/systemConfig"
	"xdpEngine/xdp"
)

var (
	ProtoDrop int
)

func PacketCapture(iface string) {
	logger.Printf("Starting on interface %q", iface)
	if snaplen <= 0 {
		snaplen = 65535
	}
	szFrame, szBlock, numBlocks, err := afpacketComputeSize(bufferSize, snaplen, os.Getpagesize())
	if err != nil {
		errlog.Fatalln("PacketCapture afpacketComputeSize error", err)
	}
	afpacketHandle, err := newAfpacketHandle(iface, szFrame, szBlock, numBlocks, false, 2000*time.Millisecond)
	if err != nil {
		errlog.Fatalln("PacketCapture newAfpacketHandle error", err)
	}
	//err = afpacketHandle.SetBPFFilter(*filter, *snaplen)
	//if err != nil {
	//	log.Fatal(err)
	//}
	source := gopacket.ZeroCopyPacketDataSource(afpacketHandle)
	defer afpacketHandle.Close()

	if systemConfig.RunMode == "debug" {
		go StatisticLossRate(afpacketHandle, iface)
	}
	for {
	flagContinue:
		select {
		case <-CtrlC:
			// 监听程序退出信息
			logger.Printf("[%s]PacketCapture正在退出...", iface)
			xdp.DetachIfaceXdp()
			os.Exit(0)
		case <-xdp.IfaceXdpDict[iface].Ctx.Done():
			// 监听网口程序退出信息
			logger.Printf("[%s]PacketCapture正在退出...", iface)
			return
		case <-xdp.IfaceXdpDict[iface].CtxP.Done():
			// 监听网口程序退出信息
			logger.Printf("[%s]PacketCapture正在退出...", iface)
			return
		default:
			// 无异常信息，开始抓包
			data, _, err := source.ZeroCopyReadPacketData()
			if err != nil {
				errlog.Println("PacketCapture ZeroCopyReadPacketData error", err)
				goto flagContinue
			}
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
			key := decodePacket(iface, packet)

			if key == nil || len(key.Payload) == 0 {
				//logger.Println("nil or 0")
				continue
			}

			// 如果协议分析开关打开
			if xdp.IfaceXdpDict[iface].ProtoSwitch {
				rand.Seed(time.Now().UnixNano())
				randInt := rand.Intn(400)
				loopCount := 0 // 重试计数器
			LoopP:
				for {
					targetIndex := randInt % xdp.IfaceXdpDict[iface].ChannelListLength
					select {
					case xdp.IfaceXdpDict[iface].ProtoPoolChannel[targetIndex] <- *key:
						break LoopP
					default:
						loopCount++ // 重试次数+1
						randInt++   // 随机数+1
					}
					if loopCount > xdp.IfaceXdpDict[iface].ChannelListLength {
						ProtoDrop++
						break LoopP
					}
				}
			}
			// 协议功能结束
		}
	}
}

// StatisticLossRate 网卡抓包性能统计
func StatisticLossRate(afpacketHandle *afpacketHandle, iface string) {
	ticker := time.Tick(3 * time.Second)
	for {
		select {
		case <-ticker:
			_, afpacketStats, err := afpacketHandle.SocketStats()
			if err != nil {
				errlog.Println("StatisticLossRate error ", err)
			}
			logger.Println("")
			logger.Printf("[%s] 会话流表:\n%+v", iface, xdp.IfaceXdpDict[iface].SessionFlow)
			for i := 0; i < xdp.IfaceXdpDict[iface].ChannelListLength; i++ {
				logger.Printf("[%s] Stats {received dropped queue-freeze}: %d - Proto: %d  Drop_Proto: %d", iface, afpacketStats, len(xdp.IfaceXdpDict[iface].ProtoPoolChannel[i]), ProtoDrop)
			}

		case <-CtrlC:
			logger.Println("StatisticLossRate 退出统计...")
			xdp.DetachIfaceXdp()
			os.Exit(0)
		case <-xdp.IfaceXdpDict[iface].Ctx.Done():
			// 监听网口程序退出信息
			logger.Printf("[%s]StatisticLossRate 退出统计...", iface)
			return
		case <-xdp.IfaceXdpDict[iface].CtxP.Done():
			// 监听网口程序退出信息
			logger.Printf("[%s]StatisticLossRate 退出统计...", iface)
			return
		}
	}

}
