package validate

import (
	"fmt"
	"net"
)

func lookupDomain(domainStr string) (string, error) {
	mx, err := net.LookupMX(domainStr)
	if err == nil {
		return mx[0].Host, nil
	}

	ips, err := net.LookupIP(domainStr)
	if err == nil {
		return ips[0].String(), nil
	}
	return "", fmt.Errorf("failed to find MX or A records for domain %s", domainStr)
}
