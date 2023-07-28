package netutil

import (
	"net"

	"github.com/wangweihong/eazycloud/pkg/sets"
)

func GetIPAddrNotError(wantIpv6 bool) string {
	ips, err := GetIPAddrs(wantIpv6)
	if err != nil {
		return ""
	}
	if len(ips) == 0 {
		return ""
	}
	return ips[0]
}

func GetIPAddrs(wantIpv6 bool) ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, iface := range ifaces {
		if sets.NewString(iface.Name).HasAnyPrefix("e", "br") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok {
					if p4 := ipnet.IP.To4(); len(p4) == net.IPv4len {
						ips = append(ips, ipnet.IP.String())
					} else if len(ipnet.IP) == net.IPv6len && wantIpv6 {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}

	return ips, nil
}
