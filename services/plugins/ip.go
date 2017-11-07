package plugins

import (
	"net"
)

// IPLimit ip限制
func IPLimit(ip string, whitelist, blacklist *[]string) bool {
	if whitelist == nil && blacklist == nil {
		return false
	}

	if whitelist != nil && len(*whitelist) > 0 {
		for _, cidr := range *whitelist {
			_, network, _ := net.ParseCIDR(cidr)
			if network.Contains(net.ParseIP(ip)) == false {
				return true
			}
		}

		return false
	}

	if blacklist != nil && len(*blacklist) > 0 {
		for _, cidr := range *blacklist {
			_, network, _ := net.ParseCIDR(cidr)
			if network.Contains(net.ParseIP(ip)) == true {
				return true
			}
		}

	}

	return false
}
