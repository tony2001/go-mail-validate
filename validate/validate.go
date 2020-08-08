package validate

import (
	"fmt"
)

type Result struct {
	Name   string
	Valid  bool
	Reason string
	Weight int
}

type EmailData struct {
	Local    string
	Domain   string
	MxServer string
}

type action struct {
	name       string
	weight     int
	process    func(string) (bool, error)
	subActions []action
}

var supportedValidators = []action{
	{"rfc5322Regex", 30, validateRfc5322Regex, []action{
		{"domain", 40, validateDomain, []action{
			{"smtpServer", 70, validateSmtpServer, nil},
		},
		},
	},
	},
}

func runAction(resultArr []Result, emailStr string, validator action) []Result {
	result := Result{
		Name: validator.name,
	}

	valid, err := validator.process(emailStr)
	result.Valid = valid
	result.Weight = validator.weight
	if err != nil {
		result.Weight = -validator.weight
		result.Reason = err.Error()
	}
	resultArr = append(resultArr, result)

	if !valid {
		return resultArr
	}

	for _, subAction := range validator.subActions {
		resultArr = runAction(resultArr, emailStr, subAction)
	}
	return resultArr
}

func Validate(emailStr string) []Result {
	resultArr := make([]Result, 0, len(supportedValidators))

	for _, action := range supportedValidators {
		resultArr = runAction(resultArr, emailStr, action)
	}
	return resultArr
}

func validateRfc5322Regex(emailStr string) (bool, error) {
	_, err := ParseRfc5322(emailStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateCommonRegex(emailStr string) (bool, error) {
	_, err := ParseCommon(emailStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateDomain(emailStr string) (bool, error) {
	email, err := ParseCommon(emailStr)
	if err != nil {
		return false, err
	}

	if domainIsReserved(email.Domain) {
		return false, fmt.Errorf("domain is reserved for testing")
	}

	err = domainIsManaged(email.Domain)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateSmtpServer(emailStr string) (bool, error) {
	email, err := ParseCommon(emailStr)
	if err != nil {
		return false, err
	}

	mxServer, err := lookupDomain(email.Domain)
	if err != nil {
		return false, err
	}

	err = trySmtp(mxServer, email.Local, email.Domain)
	if err != nil {
		return false, err
	}
	return true, nil
}
