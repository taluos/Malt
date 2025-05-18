// this file is a copy of the https://github.com/go-kratos/kratos/blob/main/internal/host/host.go

package host

import (
	"fmt"
	"net"
	"strconv"
)

// ExtractHostPort from a given address
func ExtractHostPort(address string) (host string, port uint64, err error) {
	var ports string
	// 从address中提取host和port
	host, ports, err = net.SplitHostPort(address)
	if err != nil {
		return
	}
	// 将port转换为uint64
	port, err = strconv.ParseUint(ports, 10, 16)
	if err != nil {
		return
	}
	return
}

// isValidIP checks if the IP address is valid.
func isValidIP(address string) bool {
	ip := net.ParseIP(address)                                     // 解析ip地址
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast() // 判断是否是全局单播地址，并且不是接口本地多播地址
}

// Port return a real port.
func Port(listener net.Listener) (int, bool) {
	// 如果listener是tcp监听器，则返回监听器的地址
	if address, ok := listener.Addr().(*net.TCPAddr); ok {
		return address.Port, true
	}
	return 0, false
}

// Extract returns a private Host and port.
func Extract(hostPort string, listener net.Listener) (string, error) {
	address, port, err := net.SplitHostPort(hostPort)
	if err != nil && listener == nil {
		return "", err
	}

	if listener != nil {
		if p, ok := Port(listener); ok {
			port = strconv.Itoa(p)
		} else {
			return "", fmt.Errorf("failed to get extrat port: %v", listener.Addr())
		}
	}

	if len(address) > 0 && (address != "0.0.0.0" && address != "[::]" && address != "::") {
		return net.JoinHostPort(address, port), nil
	}

	Interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	lowest := int(^uint(0) >> 1)
	var result net.IP

	for _, iface := range Interfaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index < lowest || result == nil {
			lowest = iface.Index
		} else if result != nil {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPNet:
				ip = addr.IP
			case *net.IPAddr:
				ip = addr.IP
			default:
				continue
			}
			if isValidIP(ip.String()) {
				result = ip
			}
		}
	}
	// if result is not nil, return the result
	if result != nil {
		return net.JoinHostPort(result.String(), port), nil
	}
	// if result is nil, return an error
	return "", nil

}
