package dpiEngine

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// GetPacketFromChannel 从指定网卡的缓冲队列中读出数据进行分析
func GetPacketFromChannel(iface string) {
	fromIndex := 0
	for {
		fromIndex = fromIndex % xdp.IfaceXdpDict[iface].ChannelListLength // 负载均衡
		select {
		case key := <-xdp.IfaceXdpDict[iface].ProtoPoolChannel[fromIndex]:
			go analyse(key, iface)
		case <-CtrlC:
			logger.Printf("[%s] GetPacketFromChannel stop...", iface)
			xdp.DetachIfaceXdp()
			os.Exit(0)
		case <-xdp.IfaceXdpDict[iface].CtxP.Done():
			logger.Printf("[%s] GetPacketFromChannel ProtoEngine stop...", iface)
			return
		default:
		}
		fromIndex++
	}
}

// analyse 从协议规则列表中识别报文特征
func analyse(key utils.FiveTuple, iface string) {
	for _, rule := range ProtoRuleList {
		// 遍历协议规则列表
		if rule.IsEnable == false {
			continue // 如果当前协议未启用，不检测
		}
		if (key.DstPort > rule.StartPort && key.DstPort < rule.EndPort) || (key.SrcPort > rule.StartPort && key.SrcPort < rule.EndPort) {
			// 端口符合要求
			reReq := regexp.MustCompile(rule.ReqRegx)
			resultReq := reReq.Find(key.Payload)
			reRsp := regexp.MustCompile(rule.RspRegx)
			resultRsp := reRsp.Find(key.Payload)
			if len(resultReq) != 0 {
				// 请求
				fmt.Printf("[!]识别到%s请求: %s:%d - %s:%d\n", rule.ProtocolName, key.SrcAddr, key.SrcPort, key.DstAddr, key.DstPort)
				target := key.DstAddr + "_" + strconv.Itoa(key.DstPort)
				if _, ok := xdp.IfaceXdpDict[iface].SessionFlow[target]; ok {
					// 如果会话流表中已经存在
					if xdp.IfaceXdpDict[iface].SessionFlow[target].ProtoID == rule.Id {
						// 相同协议
						xdp.IfaceXdpDict[iface].SessionFlow[target].HitReq = true // 标记请求命中
						if xdp.IfaceXdpDict[iface].SessionFlow[target].HitReq && xdp.IfaceXdpDict[iface].SessionFlow[target].HitRsp {
							// 如果请求和响应都命中，更新最近一次命中时间，并下发策略
							xdp.IfaceXdpDict[iface].SessionFlow[target].UpdateTime = time.Now().UnixNano()
							fmt.Println("开始阻断：", target)
						}
					}
				} else {
					// 如果会话流表中不存在，则加入，并标注请求命中
					xdp.IfaceXdpDict[iface].SessionFlow[target] = &utils.SessionTuple{
						ServerAddr: key.DstAddr,
						ServerPort: key.DstPort,
						ProtoID:    rule.Id,
						HitReq:     true,
						UpdateTime: 9999999999, // 10位
					}
				}
			} else if len(resultRsp) != 0 {
				// 响应
				fmt.Printf("[!]识别到%s响应: %s:%d - %s:%d\n", rule.ProtocolName, key.SrcAddr, key.SrcPort, key.DstAddr, key.DstPort)
				target := key.SrcAddr + "_" + strconv.Itoa(key.SrcPort)
				if _, ok := xdp.IfaceXdpDict[iface].SessionFlow[target]; ok {
					// 如果会话流表中已经存在
					if xdp.IfaceXdpDict[iface].SessionFlow[target].ProtoID == rule.Id {
						// 相同协议
						xdp.IfaceXdpDict[iface].SessionFlow[target].HitRsp = true // 标记响应命中
						if xdp.IfaceXdpDict[iface].SessionFlow[target].HitRsp && xdp.IfaceXdpDict[iface].SessionFlow[target].HitReq {
							// 如果响应和请求都命中，更新最近一次命中时间，并下发策略
							xdp.IfaceXdpDict[iface].SessionFlow[target].UpdateTime = time.Now().UnixNano()
							fmt.Println("开始阻断：", target)
						}
					}
				} else {
					// 如果会话流表中不存在，则加入，并标注请求命中
					xdp.IfaceXdpDict[iface].SessionFlow[target] = &utils.SessionTuple{
						ServerAddr: key.SrcAddr,
						ServerPort: key.SrcPort,
						ProtoID:    rule.Id,
						HitRsp:     true,
						UpdateTime: 9999999999, // 10位
					}
				}
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

// WriteProtoRuleFile 将协议规则写入文件
func WriteProtoRuleFile() (err error) {
	type protoRuleList struct {
		Rules []utils.ProtoRule `json:"rules"`
	}
	logger.Println("更新协议规则文件...")
	ruleList := protoRuleList{
		Rules: ProtoRuleList,
	}
	// 序列化
	result, err := json.Marshal(ruleList)
	if err != nil {
		errlog.Println("WriteProtoRuleFile json.Marshal error: ", err.Error())
		return errors.New(err.Error())
	}
	// 保存文件
	ruleFile, err := os.OpenFile(systemConfig.ProtoRuleFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		errlog.Println("WriteProtoRuleFile write File err: " + err.Error())
		return err
	}
	defer ruleFile.Close()
	// offset
	writer := bufio.NewWriter(ruleFile)
	_, err = writer.WriteString(string(result))
	_ = writer.Flush() // 刷新缓冲区，强制写出
	return nil
}
