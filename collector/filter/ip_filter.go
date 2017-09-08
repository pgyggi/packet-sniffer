package filter

import (
	"net"

	"../../common"
	"../../logp"
)

type IPFieldsConfig struct {
	Fields []string `config:"fields"`
}

type IPFields struct {
	Fields []string
}

func NewIPFields(fields []string) *IPFields {
	return &IPFields{fields}
}

func (f *IPFields) Filter(event common.MapStr) (common.MapStr, error) {

	logp.Debug("filter", "Before filter :%v", event)

	for _, field := range f.Fields {
		if val, exist := event[field]; exist {
			ip := net.ParseIP(val.(string))
			mask := net.IPv4Mask(255, 255, 255, 0)
			c := &net.IPNet{ip, mask}
			_, net, _ := net.ParseCIDR(c.String())
			event[field+"_net"] = net.IP
		}
	}

	logp.Debug("filter", "After filter :%v", event)

	return event, nil
}

func (f *IPFields) String() string {
	str := "ip_fields="
	for _, v := range f.Fields {
		str += v + ","
	}
	return str
}
