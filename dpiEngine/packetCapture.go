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
			logger.Println("PacketCapture正在退出...")
			xdp.DetachIfaceXdp()
			os.Exit(0)
		default:
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
				targetIndex := rand.Intn(400) % xdp.IfaceXdpDict[iface].ChannelListLength
				select {
				case xdp.IfaceXdpDict[iface].ProtoPoolChannel[targetIndex] <- *key:
				default:
				}
			}
			// 协议功能结束

		}
	}
}

// StatisticLossRate 网卡抓包性能统计
func StatisticLossRate(afpacketHandle *afpacketHandle, iface string) {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			_, afpacketStats, err := afpacketHandle.SocketStats()
			if err != nil {
				errlog.Println("StatisticLossRate error ", err)
			}
			logger.Println("")
			for i := 0; i < xdp.IfaceXdpDict[iface].ChannelListLength; i++ {
				logger.Printf("[%s] Stats {received dropped queue-freeze}: %d - Proto: %d  Drop_Proto: %d", iface, afpacketStats, len(xdp.IfaceXdpDict[iface].ProtoPoolChannel[i]), ProtoDrop)
			}

		case <-CtrlC:
			logger.Println("StatisticLossRate 退出统计...")
			xdp.DetachIfaceXdp()
			os.Exit(0)
		}
	}

}
