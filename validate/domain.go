package validate

import (
	"strings"
)

var testDomains = map[string]bool{
	"test":        true,
	"example":     true,
	"invalid":     true,
	"localhost":   true,
	"example.com": true,
	"example.net": true,
	"example.org": true,
}

func stripSubdomains(domainStr string, leaveNum int) string {
	dots := strings.Count(domainStr, ".")
	for dots > leaveNum {
		firstDot := strings.IndexByte(domainStr, '.')
		domainStr = domainStr[firstDot+1:]
		dots--
	}
	return domainStr
}

func domainIsReserved(domainStr string) bool {
	//leave only one subdomain because all tests domains are at most second-level ones
	domainStr = stripSubdomains(domainStr, 1)
	_, ok := testDomains[domainStr]
	return ok
}
