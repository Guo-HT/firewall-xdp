package dpiEngine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// getPacketFromChannel 从指定网卡的缓冲队列中读出数据进行分析
func getPacketFromChannel(iface string) {
	fromIndex := 0
	for {
		fromIndex = fromIndex % xdp.IfaceXdpDict[iface].ChannelListLength // 负载均衡
		select {
		case key := <-xdp.IfaceXdpDict[iface].ProtoPoolChannel[fromIndex]:
			go analyse(key)
		case <-CtrlC:
			logger.Printf("[%s] getPacketFromChannel stop...", iface)
			xdp.DetachIfaceXdp()
			os.Exit(0)
		case <-xdp.IfaceXdpDict[iface].CtxP.Done():
			logger.Printf("[%s] getPacketFromChannel ProtoEngine stop...", iface)
			return
		default:
		}
		fromIndex++
	}
}

// analyse 从协议规则列表中识别报文特征
func analyse(key utils.FiveTuple) {
	for _, rule := range ProtoRuleList {
		if (key.DstPort > rule.StartPort && key.DstPort < rule.EndPort) || (key.SrcPort > rule.StartPort && key.SrcPort < rule.EndPort) {
			// 端口符合要求
			reReq := regexp.MustCompile(rule.ReqRegx)
			resultReq := reReq.Find(key.Payload)
			reRsp := regexp.MustCompile(rule.RspRegx)
			resultRsp := reRsp.Find(key.Payload)
			if len(resultReq) != 0 || len(resultRsp) != 0 {
				fmt.Printf("[!]识别到%s: %s:%d - %s:%d, %v | %v\n", rule.ProtocolName, key.SrcAddr, key.SrcPort, key.DstAddr, key.DstPort, resultReq, resultRsp)
			}
		}
	}
}

// InitProtoRules 初始化协议规则列表
func InitProtoRules() {
	fileContent := readProtoRuleFile()
	ProtoRuleList = parseProtoRules(fileContent)
}

// parseProtoRules 将文件中读出的内容解析为【协议规则列表】
func parseProtoRules(fileContent []byte) (rule []utils.ProtoRule) {
	type protoRuleList struct {
		Rules []utils.ProtoRule `json:"rules"`
	}
	var RuleList protoRuleList
	err := json.Unmarshal(fileContent, &RuleList)
	if err != nil {
		errlog.Println("parseProtoRules error", err)
	}
	return RuleList.Rules
}

// readProtoRuleFile 读取协议规则文件
func readProtoRuleFile() (ruleContent []byte) {
	file, err := os.Open(systemConfig.ProtoRuleFile)
	if err != nil {
		errlog.Println("readProtoRuleFile open file error:", err)
		return
	}
	defer file.Close()
	// 定义接收文件读取的字节数组
	var buf [128]byte
	var content []byte
	for {
		n, err := file.Read(buf[:])
		if err == io.EOF {
			// 读取结束
			break
		}
		if err != nil {
			errlog.Println("readProtoRuleFile read file error", err)
			return
		}
		content = append(content, buf[:n]...)
	}
	//return *(*string)(unsafe.Pointer(&content))
	return content
}
