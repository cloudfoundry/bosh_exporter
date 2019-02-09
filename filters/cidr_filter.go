package filters

import (
	"net"
)

type CidrFilter struct {
	cidrFilters []*net.IPNet
}

func NewCidrFilter(filters []string) (*CidrFilter, error) {
	cidrFilters := []*net.IPNet{}

	for _, filter := range filters {
		_, net, err := net.ParseCIDR(filter)
		if err != nil {
			return nil, err
		}
		cidrFilters = append(cidrFilters, net)
	}

	return &CidrFilter{cidrFilters: cidrFilters}, nil
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
