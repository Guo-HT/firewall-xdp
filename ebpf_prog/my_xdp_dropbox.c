//#include "headers/bpf.h"
#include "headers/bpf_helpers_dropbox.h"
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#define MAX_SIZE 10240

// 可读的端口数据 [0001 0110 0000 0000](5632) -> [0000 0000 0001 0110](22)
#define GET_PORT(x) ( (x & 0xff00)>>8 | ((x & 0x00ff)<<8) )
// 可读的IP数据
#define GET_IP(x) (((x&0xff000000)>>24)|((x&0x00ff0000)>>8)|((x&0x0000ff00)<<8)|((x&0x000000ff)<<24))
// 可读的MAC地址
#define GET_MAC(x) ()
// 封装bpf_trace_printk函数
#define bpfprint(fmt, ...) ({ char fmt_char[] = fmt;  bpf_trace_printk(fmt_char, sizeof(fmt_char),  ##__VA_ARGS__); })


// port 白名单
BPF_MAP_DEF(white_port) = {
	.map_type    = BPF_MAP_TYPE_HASH,
	.key_size    = sizeof(__u32),
	.value_size  = sizeof(__u32),
	.max_entries = MAX_SIZE,
};
BPF_MAP_ADD(white_port);

// port 黑名单
BPF_MAP_DEF(black_port) = {
	.map_type    = BPF_MAP_TYPE_HASH,
	.key_size    = sizeof(__u32),
	.value_size  = sizeof(__u32),
	.max_entries = MAX_SIZE,
};
BPF_MAP_ADD(black_port);

// ip 白名单
BPF_MAP_DEF(white_ip) = {
	.map_type    = BPF_MAP_TYPE_LPM_TRIE,
	.key_size    = sizeof(__u64),
	.value_size  = sizeof(__u32),
	.max_entries = MAX_SIZE,
	.map_flags   = 1,
};
BPF_MAP_ADD(white_ip);

// ip 黑名单
BPF_MAP_DEF(black_ip) = {
	.map_type    = BPF_MAP_TYPE_LPM_TRIE,
	.key_size    = sizeof(__u64),
	.value_size  = sizeof(__u32),
	.max_entries = MAX_SIZE,
	.map_flags   = 1,
};
BPF_MAP_ADD(black_ip);

// 功能开关map
BPF_MAP_DEF(function_switch) = {
	.map_type    = BPF_MAP_TYPE_HASH,
	.key_size    = sizeof(__u32),
	.value_size  = sizeof(__u32),
	.max_entries = MAX_SIZE,
};
BPF_MAP_ADD(function_switch);

SEC("xdp")
int firewall(struct xdp_md *ctx)
{
    int ipsize = 0;
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;
    struct ethhdr *eth = data;
    struct iphdr *ip;
    struct tcphdr *tcp;
    struct udphdr *udp;
    int is_tcp = 0;
    int is_udp = 0;
    int src_port = 0;
    int dst_port = 0;

    struct {
        __u32 prefixlen;
        __u32 saddr;
    } key;
    key.prefixlen = 32;

    if ((void *)eth + sizeof(*eth) <= data_end)
    {
        ip = data + sizeof(*eth);
        if ((void *)ip + sizeof(*ip) <= data_end)
        {
            // bpfprint("ip: %d, %d", ip->saddr, ip->daddr);
            if (ip->protocol == IPPROTO_UDP)
            {
                udp = (void *)ip + sizeof(*ip);
                if ((void *)udp + sizeof(*udp) <= data_end)
                {
                    // char sure_udp[] = "sure_udp, port: %u";
                    // bpf_trace_printk(sure_udp, sizeof(sure_udp), GET_PORT(udp->dest));
                    is_udp = 1;
                    is_tcp = 0;
                    goto process;
                }
            }

            if (ip->protocol == IPPROTO_TCP)
            {
                tcp = (void *)ip + sizeof(*ip);
                if ((void *)tcp + sizeof(*tcp) <= data_end)
                {
                    // char sure_tcp[] = "sure_tcp, port: %u";
                    // bpf_trace_printk(sure_tcp, sizeof(sure_tcp), GET_PORT(tcp->dest));
                    is_udp = 0;
                    is_tcp = 1;
                    goto process;
                }
            }
        }
    }

    process:

    if(is_udp || is_tcp){
        if(is_udp){
            // bpfprint("src: %d, dst: %d", GET_PORT(udp->source), GET_PORT(udp->dest));
            src_port = GET_PORT(udp->source);
            dst_port = GET_PORT(udp->dest);
        }
        if(is_tcp){
            // bpfprint("src: %d, dst: %d", GET_PORT(tcp->source), GET_PORT(tcp->dest));
            src_port = GET_PORT(tcp->source);
            dst_port = GET_PORT(tcp->dest);
        }
        bpfprint("[ IP ] src: %u, dst: %u", ip->saddr, ip->daddr);
        bpfprint("[Port] src: %u, dst: %u", src_port, dst_port);

        // Port 白名单  限制目的端口
        int *lookup_port_white = bpf_map_lookup_elem(&white_port, &dst_port);
        if(lookup_port_white){
            bpfprint("[!] Hitted! port White...");
            return XDP_PASS;
        }

        // IP 白名单
        key.saddr = ip->saddr;  // 限制源ip
        int *lookup_ip_white = bpf_map_lookup_elem(&white_ip, &key);
        if(lookup_ip_white){
            bpfprint("[!] Hitted! ip White...");
            return XDP_PASS;
        }

        // Port 黑名单  限制目的端口
        int *lookup_port_black = bpf_map_lookup_elem(&black_port, &dst_port);
        if(lookup_port_black){
            bpfprint("[!] Hitted! port Black...");
            return XDP_DROP;
        }

        // IP 黑名单
        key.saddr = ip->saddr;  // 限制源ip
        int *lookup_ip_black = bpf_map_lookup_elem(&black_ip, &key);
        if(lookup_ip_black){
            bpfprint("[!] Hitted! ip Black...");
            return XDP_DROP;
        }

    }


    return XDP_PASS;
}


char _license[] SEC("license") = "GPL";

// 挂载xdp程序：
// 		ip link set dev ens33 xdp obj my_xdp.o sec xdp
// 卸载xdp程序：
// 		ip link set dev ens33 xdp off
