package gin

import (
	"net"
	"strings"

	"github.com/ucarion/urlpath"
)

func isPrivateIP(ip string) bool {
	var privateIPBlocks []*net.IPNet

	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}

	clientIp := net.ParseIP(ip)

	if clientIp.IsLoopback() || clientIp.IsLinkLocalUnicast() || clientIp.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(clientIp) {
			return true
		}
	}

	return false
}

func validatePublicApi(inputPath string, publicApi string) bool {
	var pathPattern = urlpath.New(publicApi)

	// get path without query params
	pathSplit := strings.Split(inputPath, "?")
	if len(pathSplit) == 0 {
		return false
	}
	_, ok := pathPattern.Match(pathSplit[0])
	if !ok {
		return false
	}
	return true
}
