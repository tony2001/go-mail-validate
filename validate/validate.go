package validate

import (
	"fmt"
)

type Result struct {
	Name   string
	Valid  bool
	Reason string
}

type Email struct {
	Local  string
	Domain string
}

type ValidateFSM struct {
	Email    *Email
	MxServer string
}

const (
	StateRfc5322Regex = iota
	StateDomainNotReserved
	StateCommonRegex
	StateSimpleRegex
	StateMxRecord
	StateSmtpTest
	StateClearout
	StateInvalid
	StateValid
)

func (v *ValidateFSM) Validate(emailStr string) (bool, []Result) {

	var (
		err       error
		resultArr []Result
		valid     bool
	)

	state := StateRfc5322Regex
FSM:
	for {
		result := Result{}
		switch state {
		case StateRfc5322Regex:
			result.Name = "rfc5322regex"
			result.Valid, err = v.validateRfc5322Regex(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateCommonRegex
			} else {
				state = StateDomainNotReserved
			}
		case StateDomainNotReserved:
			result.Name = "domainNotReserved"
			result.Valid, err = v.validateReservedDomain(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateInvalid
			} else {
				state = StateMxRecord
			}
		case StateCommonRegex:
			result.Name = "commonRegex"
			result.Valid, err = v.validateCommonRegex(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateSimpleRegex
			} else {
				state = StateDomainNotReserved
			}
		case StateSimpleRegex:
			result.Name = "simpleRegex"
			result.Valid, err = v.validateSimpleRegex(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateInvalid
			} else {
				state = StateDomainNotReserved
			}
		case StateMxRecord:
			result.Name = "mxRecord"
			result.Valid, err = v.validateDomainLookup(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateInvalid
			} else {
				state = StateSmtpTest
			}
		case StateSmtpTest:
			result.Name = "smtpTest"
			result.Valid, err = v.validateSmtpServer(emailStr)
			if err != nil {
				result.Reason = err.Error()

				if result.Valid {
					//valid, but has errors - lets try Clearout then
					if ClearoutEnabled() {
						state = StateClearout
					} else {
						state = StateValid
					}
				} else {
					state = StateInvalid
				}
			} else {
				state = StateValid
			}
		case StateClearout:
			result.Name = "clearout"
			result.Valid, err = v.validateClearout(emailStr)
			if err != nil {
				result.Reason = err.Error()
				state = StateInvalid
			} else {
				state = StateValid
			}
		case StateInvalid:
			valid = false
			break FSM
		case StateValid:
			valid = true
			break FSM
		}
		resultArr = append(resultArr, result)
	}

	return valid, resultArr
}

func (v *ValidateFSM) validateRfc5322Regex(emailStr string) (bool, error) {
	email, err := ParseRfc5322(emailStr)
	if err != nil {
		return false, err
	}

	if v.Email == nil {
		v.Email = email
	}

	return true, nil
}

func (v *ValidateFSM) validateCommonRegex(emailStr string) (bool, error) {
	email, err := ParseCommon(emailStr)
	if err != nil {
		return false, err
	}

	if v.Email == nil {
		v.Email = email
	}

	return true, nil
}

func (v *ValidateFSM) validateSimpleRegex(emailStr string) (bool, error) {
	email, err := ParseSimple(emailStr)
	if err != nil {
		return false, err
	}

	if v.Email == nil {
		v.Email = email
	}

	return true, nil
}

func (v *ValidateFSM) validateReservedDomain(emailStr string) (bool, error) {
	if v.Email == nil {
		return false, fmt.Errorf("unparsed email")
	}

	if domainIsReserved(v.Email.Domain) {
		return false, fmt.Errorf("domain is reserved for testing")
	}
	return true, nil
}

func (v *ValidateFSM) validateDomainLookup(emailStr string) (bool, error) {
	if v.Email == nil {
		return false, fmt.Errorf("unparsed email")
	}

	mxServer, err := lookupDomain(v.Email.Domain)
	if err != nil {
		return false, err
	}

	v.MxServer = mxServer
	return true, nil
}

func (v *ValidateFSM) validateSmtpServer(emailStr string) (bool, error) {
	if v.MxServer == "" {
		return false, fmt.Errorf("mx server not found")
	}

	valid, err := trySmtp(v.MxServer, v.Email.Local, v.Email.Domain, false)
	return valid, err
}

func (v *ValidateFSM) validateClearout(emailStr string) (bool, error) {
	result, valid, err := ClearoutInstantCheck(emailStr)
	if !result {
		//query error
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("Clearout API request failed: unknown error")
	}

	if valid {
		return true, nil
	}
	return false, err
}
