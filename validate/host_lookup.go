package validate

import (
	"context"
	"fmt"
	"net"
)

func lookupDomain(ctx context.Context, domainStr string) (string, error) {
	mx, err := net.DefaultResolver.LookupMX(ctx, domainStr)
	if err == nil {
		return mx[0].Host, nil
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, domainStr)
	if err == nil {
		return ips[0].String(), nil
	}
	return "", fmt.Errorf("failed to find MX or A records for domain %s", domainStr)
}
