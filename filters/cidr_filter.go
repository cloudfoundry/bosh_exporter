package filters

import (
	"net"
)

type CidrFilter struct {
	cidrFilters []*net.IPNet
}

func NewCidrFilter(cidrs []string) (*CidrFilter, error) {
	nets := []*net.IPNet{}
	for _, c := range cidrs {
		_, net, err := net.ParseCIDR(c)
		if err != nil {
			return nil, err
		}
		nets = append(nets, net)
	}
	return &CidrFilter{
		cidrFilters: nets,
	}, nil
}

func (f *CidrFilter) Select(ips []string) (string, bool) {
	for _, c := range f.cidrFilters {
		for _, val := range ips {
			ip := net.ParseIP(val)
			if ip == nil {
				continue
			}
			if c.Contains(ip) {
				return val, true
			}
		}
	}
	return "", false
}
