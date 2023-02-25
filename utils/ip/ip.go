// Package ip ip地址相关函数
package ip

import (
	"encoding/binary"
	"fmt"
	"net"
	"runtime"
)

// NetInterfaceIP 通过指定的网卡名称获取对应的ip地址
func NetInterfaceIP(name string) (string, error) {
	itf, err := net.InterfaceByName(name) //here your interface
	if err != nil {
		return "", err
	}
	item, err := itf.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil {
				return v.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not found")
}

// LocalIP 获取系统默认网卡名称获取对应的ip地址
func LocalIP() string {
	name := "eth1"
	if runtime.GOOS == "darwin" {
		name = "en0"
	}
	ip, err := NetInterfaceIP(name)
	if err != nil {
		return ""
	}
	return ip
}

// Uint32ToIPv4 uint32 转 ipv4
func Uint32ToIPv4(intIP uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// IPv4ToString IP 转字符串
func IPv4ToString(ipInt uint32) string {
	ip := Uint32ToIPv4(ipInt)
	return ip.String()
}

// IPv4ToU32 ip 转 uint32
func IPv4ToU32(ip string) uint32 {
	ips := net.ParseIP(ip)
	if len(ips) == 16 {
		return binary.BigEndian.Uint32(ips[12:16])
	} else if len(ips) == 4 {
		return binary.BigEndian.Uint32(ips)
	}
	return 0
}
