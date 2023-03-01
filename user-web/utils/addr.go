package utils

import (
	"fmt"
	"net"
	"strings"
)

// GetFreePort用于自动获取可用端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// getIP() 获取当前校园网内的ip
func GeIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ip = ipnet.IP.String(); strings.HasPrefix(ip, "219.228") {
					return
				}
				if ip = ipnet.IP.String(); strings.HasPrefix(ip, "172.20") {
					return
				}
				if ip = ipnet.IP.String(); strings.HasPrefix(ip, "172.30") {
					return
				}
			}
		}
	}
	return
}
