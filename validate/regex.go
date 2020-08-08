package validate

import (
	"fmt"
	"regexp"
	"strings"
)

//see https://www.rfc-editor.org/errata_search.php?rfc=3696&eid=1690
const _maxEmailLength = 254
const _maxLocalLength = 64

var (
	_rfc5322       = "^(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])*$"
	_rfc5322Regexp = regexp.MustCompile(_rfc5322)

	_commonRegexp = regexp.MustCompile("^(?i)([A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,24})*$")
)

func emailIsTooLong(emailStr string) bool {
	if len(emailStr) > _maxEmailLength {
		return true
	}
	return false
}

func localPartIsTooLong(localPart string) bool {
	if len(localPart) > _maxLocalLength {
		return true
	}
	return false
}

func ParseRfc5322(emailStr string) (*Email, error) {
	if emailIsTooLong(emailStr) {
		return nil, fmt.Errorf("maximum email length exceeded")
	}

	if !_rfc5322Regexp.MatchString(emailStr) {
		return nil, fmt.Errorf("doesn't match RFC 5322")
	}

	i := strings.LastIndexByte(emailStr, '@')
	e := &Email{
		Local:  emailStr[:i],
		Domain: emailStr[i+1:],
	}

	if localPartIsTooLong(e.Local) {
		return nil, fmt.Errorf("maximum local part length exceeded")
	}

	return e, nil
}

func ParseCommon(emailStr string) (*Email, error) {
	if emailIsTooLong(emailStr) {
		return nil, fmt.Errorf("maximum email length exceeded")
	}

	if !_commonRegexp.MatchString(emailStr) {
		return nil, fmt.Errorf("doesn't match common email regexp")
	}

	i := strings.LastIndexByte(emailStr, '@')
	e := &Email{
		Local:  emailStr[:i],
		Domain: emailStr[i+1:],
	}

	if localPartIsTooLong(e.Local) {
		return nil, fmt.Errorf("maximum local part length exceeded")
	}

	return e, nil
}
