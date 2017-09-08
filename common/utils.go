package common

import (
	"net"
	"sort"
	"strconv"
)

func CalInterval(startSec int64, endSec int64, size int64) string {
	val := (endSec - startSec) / size
	if val == 0 {
		return "1s"
	} else {
		return ""
	}
}

// Convert uint to net.IP http://www.sharejs.com
func Inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func SortInts(ints []int) bool {
	a := sort.IntSlice(ints)
	sort.Sort(a)
	return true
}

func SortStrings(strings []string) bool {
	a := sort.StringSlice(strings)
	return sort.StringsAreSorted(a)
}

func ParseFloatShift(f float64, prec int) float64 {
	val, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', prec, 64), 64)
	return val
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func CIDRRange(cidr string) (string, string, error) {

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", err
	}
	start_ip := ""
	end_ip := ""
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if start_ip != "" {
			end_ip = ip.String()
		} else {
			start_ip = ip.String()
		}
	}
	return start_ip, end_ip, nil

}
