package validate

import (
	"testing"
)

func TestRfc5322Valid(t *testing.T) {
	var testsValid = []string{
		"prettyandsimple@google.com",
		"very.common@google.com",
		"disposable.style.email.with+symbol@google.com",
		"other.email-with-dash@google.com",
		"fully-qualified-domain@google.com",
		"x@google.com",
		"\"very.unusual.@.unusual.com\"@google.com",
		`"very.(),:;<>[]\".VERY.\"very@\ \"very\".unusual"@strange.google.com`,
		"google-indeed@strange-google.com",
		"/#!$%&'*+-/=?^_`{}|~@google.org",
		"google@s.solutions",
		"reserved@example.com",
		"reserved@example.org",
		"reserved@example.net",
	}

	for _, test := range testsValid {
		_, err := ParseRfc5322(test)
		if err != nil {
			t.Errorf("failed to parse email address: %s", test)
		}
	}
}

func TestRfc5322Invalid(t *testing.T) {
	var testsInvalid = []string{
		"admin@mailserver1",
		`()<>[]:,;@\\"!#$%&'-/=?^_` + "`" + `{}| ~.a"@google.org`,
		`" "@google.org`,
		"user@localserver",
		"user@tt",
		"user@[IPv6:2001:DB8::1]",
		"Abc.google.com",
		"A@b@c@google.com",
		"a\"b(c)d,e:f;gi[j\\k]l@google.com",
		"just\"not\"right@google.com",
		"this is\"not\\allowed@google.com",
		"this\\ still\\\"not\\allowed@google.com",
		"1234567890123456789012345678901234567890123456789012345678901234+x@google.com",
		"john..doe@google.com",
		"google@localhost",
		"john.doe@google..com",
		"\"much.more unusual\"@google.com",
		" abc@google.com",
		"abc@google.com ",
		"i_like_underscore@but_its_not_allowed_in _this_part.google.com",
		"reserved@test",
		"reserved@localhost",
		"reserved@invalid",
		"reserved@example",
	}

	for _, test := range testsInvalid {
		_, err := ParseRfc5322(test)
		if err == nil {
			t.Errorf("no error on invalid email address: %s", test)
		}
	}
}

func TestCommonValid(t *testing.T) {
	var testsValid = []string{
		"prettyandsimple@google.com",
		"very.common@google.com",
		"disposable.style.email.with+symbol@google.com",
		"other.email-with-dash@google.com",
		"fully-qualified-domain@google.com",
		"x@google.com",
		"google-indeed@strange-google.com",
		"google@s.solutions",
		"john.doe@google..com",
		"john..doe@google.com",
		"reserved@example.com",
		"reserved@example.org",
		"reserved@example.net",
	}

	for _, test := range testsValid {
		_, err := ParseCommon(test)
		if err != nil {
			t.Errorf("failed to parse email address: %s", test)
		}
	}
}

func TestCommonInvalid(t *testing.T) {
	var testsInvalid = []string{
		"admin@mailserver1",
		`()<>[]:,;@\\"!#$%&'-/=?^_` + "`" + `{}| ~.a"@google.org`,
		`" "@google.org`,
		"user@localserver",
		"user@tt",
		"user@[IPv6:2001:DB8::1]",
		"Abc.google.com",
		"A@b@c@google.com",
		"a\"b(c)d,e:f;gi[j\\k]l@google.com",
		"just\"not\"right@google.com",
		"this is\"not\\allowed@google.com",
		"this\\ still\\\"not\\allowed@google.com",
		"\"very.unusual.@.unusual.com\"@google.com",
		`"very.(),:;<>[]\".VERY.\"very@\ \"very\".unusual"@strange.google.com`,
		"1234567890123456789012345678901234567890123456789012345678901234+x@google.com",
		"google@localhost",
		"/#!$%&'*+-/=?^_`{}|~@google.org",
		"\"much.more unusual\"@google.com",
		" abc@google.com",
		"abc@google.com ",
		"i_like_underscore@but_its_not_allowed_in _this_part.google.com",
		"reserved@test",
		"reserved@localhost",
		"reserved@invalid",
		"reserved@example",
	}

	for _, test := range testsInvalid {
		_, err := ParseCommon(test)
		if err == nil {
			t.Errorf("no error on invalid email address: %s", test)
		}
	}
}

var benchTests = []string{
	"prettyandsimple@google.com",
	"very.common@google.com",
	"disposable.style.email.with+symbol@google.com",
	"other.email-with-dash@google.com",
	"fully-qualified-domain@google.com",
	"x@google.com",
	"\"very.unusual.@.unusual.com\"@google.com",
	`"very.(),:;<>[]\".VERY.\"very@\ \"very\".unusual"@strange.google.com`,
	"google-indeed@strange-google.com",
	"/#!$%&'*+-/=?^_`{}|~@google.org",
	"google@s.solutions",
	"admin@mailserver1",
	`()<>[]:,;@\\"!#$%&'-/=?^_` + "`" + `{}| ~.a"@google.org`,
	`" "@google.org`,
	"user@localserver",
	"user@tt",
	"user@[IPv6:2001:DB8::1]",
	"Abc.google.com",
	"A@b@c@google.com",
	"a\"b(c)d,e:f;gi[j\\k]l@google.com",
	"just\"not\"right@google.com",
	"this is\"not\\allowed@google.com",
	"this\\ still\\\"not\\allowed@google.com",
	"1234567890123456789012345678901234567890123456789012345678901234+x@google.com",
	"john..doe@google.com",
	"google@localhost",
	"john.doe@google..com",
	"\"much.more unusual\"@google.com",
	" abc@google.com",
	"abc@google.com ",
	"i_like_underscore@but_its_not_allowed_in _this_part.google.com",
	"reserved@test",
	"reserved@localhost",
	"reserved@invalid",
	"reserved@example",
	"reserved@example.com",
	"reserved@example.org",
	"reserved@example.net",
}

func BenchmarkRfc5322(b *testing.B) {
	testN := 0
	for i := 0; i < b.N; i++ {
		ParseRfc5322(benchTests[testN])

		testN++
		if testN == len(benchTests) {
			testN = 0
		}
	}
}

func BenchmarkCommon(b *testing.B) {
	testN := 0
	for i := 0; i < b.N; i++ {
		ParseCommon(benchTests[testN])

		testN++
		if testN == len(benchTests) {
			testN = 0
		}
	}
}
