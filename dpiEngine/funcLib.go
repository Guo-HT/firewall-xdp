package dpiEngine

import (
	"encoding/hex"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// StartProtoEngine 开启所有引擎的协议检测
func StartProtoEngine() {
	for iface, xdpObj := range xdp.IfaceXdpDict {
		// 循环开启各个网口的协议分析功能
		if xdpObj.ProtoSwitch {
			logger.Printf("正在开启[%s]的分析功能...", iface)
			go GetPacketFromChannel(iface)
			go PacketCapture(iface)
		}
	}
}

// decodePacket 解析报文
func decodePacket(Iface string, pkt gopacket.Packet) (key *utils.FiveTuple) {
	ipv4, ok := pkt.NetworkLayer().(*layers.IPv4)
	if !ok {
		return // Ignore packets that aren't IPv4
	}

	if ipv4.FragOffset != 0 || (ipv4.Flags&layers.IPv4MoreFragments) != 0 {
		return // Ignore fragmented packets.
	}
	// 仅分析TCP报文

	switch t := pkt.TransportLayer().(type) {
	case *layers.TCP:
		var stream utils.FiveTuple
		LenPayload := len(t.Payload)
		HexPayloadBytes := make([]byte, hex.EncodedLen(LenPayload))
		stream.SrcAddr = ipv4.SrcIP.String()
		stream.DstAddr = ipv4.DstIP.String()
		stream.SrcPort = int(t.SrcPort)
		stream.DstPort = int(t.DstPort)
		stream.NextAck = uint32(t.Seq) + uint32(LenPayload) //响应包的ACK值
		stream.Payload = t.Payload
		stream.PayloadHex = HexPayloadBytes
		stream.LayerTcp = t
		stream.Iface = Iface
		return &stream
	}
	return
}

// GetStartingProto 获取所有开启分析的协议名
func GetStartingProto() (rules []string) {
	for _, value := range ProtoRuleList {
		if value.IsEnable {
			rules = append(rules, value.ProtocolName)
		}
	}
	return
}
