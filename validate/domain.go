package validate

import (
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
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

func domainIsManaged(domainStr string) error {
	eTLD, icann := publicsuffix.PublicSuffix(domainStr)

	if icann {
		//ICANN managed domain
		return nil
	} else if strings.IndexByte(eTLD, '.') >= 0 {
		//privately managed
		return nil
	}
	//no such domain
	return fmt.Errorf("such domain probably doesn't exist: %s", domainStr)
}
