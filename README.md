# IP-Port Firewall

## 简介
这个项目是一个部署于Linux终端设备的，可完成基于IP、端口、协议的基础防火墙

## 具体功能
1. 实现基于IP的黑、白名单访问控制
2. 实现基于端口的黑、白名单访问控制
3. 实现基于协议的访问控制
4. 实现用户自定义协议规则
5. 实现多网卡设备的挂载与防护
6. 实现通过Web页面的简单配置

## 环境依赖
- golang 1.18.4
- clang 6.0.0
- linux (kernel>=5.10)
- sqlite3
- libpcap
- libpcap-dev
- net-tools

## 调试工具
- tcpdump
- bpftool (需要编译内核)
- sqlite3
