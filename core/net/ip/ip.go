package ip

import (
	"net"
	"strings"
)

// Atoi convert ip addr to uint32.
func Atoi(s string) (i uint32) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
	ip = ip.To4()
	if ip == nil {
		return
	}
	i += uint32(ip[0]) << 24
	i += uint32(ip[1]) << 16
	i += uint32(ip[2]) << 8
	i += uint32(ip[3])
	return
}

// Itoa convert uint32 to ip addr.
func Itoa(i uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((i >> 24) & 0xFF)
	ip[1] = byte((i >> 16) & 0xFF)
	ip[2] = byte((i >> 8) & 0xFF)
	ip[3] = byte(i & 0xFF)
	return ip.String()
}

// ExternalIP get external ip.
func ExternalIP() (res []string) {
	inters, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				if ipNet.IP.IsLoopback() || ipNet.IP.IsLinkLocalMulticast() || ipNet.IP.IsLinkLocalUnicast() {
					continue
				}
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					switch true {
					case ip4[0] == 10:
						continue
					case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
						continue
					case ip4[0] == 192 && ip4[1] == 168:
						continue
					default:
						res = append(res, ipNet.IP.String())
					}
				}
			}
		}
	}
	return
}

// InternalIP get internal ip.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !isUp(inter.Flags) {
			continue
		}
		if strings.HasPrefix(inter.Name, "lo") {
			continue
		}
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}

// isUp Interface is up
func isUp(v net.Flags) bool {
	return v&net.FlagUp == net.FlagUp
}
