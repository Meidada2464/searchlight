/**
 * Package nic
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/17 23:14
 */

package nic

import (
	"errors"
	"math"
	"net"
)

func IsIpLegal(ip string) bool {
	parseIP := net.ParseIP(ip)
	if parseIP.To16() != nil && parseIP.To4() != nil {
		return true
	}
	return false
}

// GetSameAvailableIp obtain the ip address of the same type
func GetSameAvailableIp(ip string) (string, error) {
	var (
		ipType   string
		ipv4List []string
		ipv6List []string
	)
	if IsIpLegal(ip) {
		if net.ParseIP(ip).To4() != nil {
			ipType = "ipv4"
		}
		if net.ParseIP(ip).To16() != nil {
			ipType = "ipv6"
		}
	}

	ipv4s, ipv6s := GetAllNiCs()

	// Gets the available ip, and returns empty if no ip of the same type is available
	switch ipType {
	case "ipv4":
		if ok := ipv4s[ip]; ok != "" {
			return ip, nil
		}
		if len(ipv4s) == 0 {
			break
		}
		pointer := math.Round(float64(len(ipv4s)) - 1)
		for key, _ := range ipv4s {
			ipv4List = append(ipv4List, key)
		}
		return ipv4List[int(pointer)], nil
	case "ipv6":
		if ok := ipv6s[ip]; ok != "" {
			return ip, nil
		}
		if len(ipv6s) == 0 {
			break
		}
		pointer := math.Round(float64(len(ipv4s)) - 1)
		for key, _ := range ipv4s {
			ipv6List = append(ipv6List, key)
		}
		return ipv6List[int(pointer)], nil
	default:
		// default use ipv4 nic
		pointer := math.Round(float64(len(ipv4s)) - 1)
		for key, _ := range ipv4s {
			ipv4List = append(ipv4List, key)
		}
		return ipv4List[int(pointer)], nil
	}

	return "", errors.New("no network adapter is available")
}

func GetAllNiCs() (ipv4s, ipv6s map[string]string) {
	// get all interface info
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil
	}

	var ipv4nics = make(map[string]string)
	var ipv6nics = make(map[string]string)

	for _, intf := range interfaces {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, nil
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			// ip is not available
			if ipNet.IP.IsLoopback() || ipNet.IP.IsMulticast() {
				continue
			}

			if ipNet.IP.To4() != nil {
				ipv4nics[ipNet.IP.String()] = intf.Name
			}
			if ipNet.IP.To16() != nil {
				ipv6nics[ipNet.IP.String()] = intf.Name
			}
		}
	}

	return ipv4nics, ipv6nics
}

func GetIpType(ip string) string {
	if net.ParseIP(ip).To4() != nil {
		return "ipv4"
	}
	if net.ParseIP(ip).To16() != nil {
		return "ipv6"
	}
	return ""
}
