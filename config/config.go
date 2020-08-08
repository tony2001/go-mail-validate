package config

import (
	"fmt"
	"strconv"

	"go-mail-validate/log"
)

const defaultSmtpTimeoutMsec = 1000

type Config struct {
	port            int
	logLevel        log.LogLevelT
	smtpTimeoutMsec int
}

var gConfig *Config

func validatePort(portStr string) (int, error) {
	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", err)
	}

	if portNum < 1 || portNum > 65535 {
		return 0, fmt.Errorf("invalid port number: %d, port has to be > 0 and < 65536", portNum)
	}
	return portNum, nil
}

func NewConfig(portStr string, options ...func(*Config) error) error {

	if gConfig != nil {
		return fmt.Errorf("initializing global config twice")
	}

	port, err := validatePort(portStr)
	if err != nil {
		return err
	}

	gConfig = &Config{
		port:            port,
		logLevel:        log.Info,
		smtpTimeoutMsec: defaultSmtpTimeoutMsec,
	}

	for _, option := range options {
		err := option(gConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func LogLevel(logLevel log.LogLevelT) func(*Config) error {
	return func(c *Config) error {
		c.logLevel = logLevel
		return nil
	}
}

func SmtpTimeout(smtpTimeout int) func(*Config) error {
	return func(c *Config) error {
		if smtpTimeout <= 0 {
			return fmt.Errorf("SMTP timeout must be greater than 0")
		}
		gConfig.smtpTimeoutMsec = smtpTimeout
		return nil
	}
}

func GetPort() int {
	return gConfig.port
}

func GetLogLevel() log.LogLevelT {
	return gConfig.logLevel
}

func GetSmtpTimeout() int {
	return gConfig.smtpTimeoutMsec
}

func String() string {
	return fmt.Sprintf("port: %d, logLevel: %s, smtpTimeout: %dms", gConfig.port, log.GetLogLevelStr(gConfig.logLevel), gConfig.smtpTimeoutMsec)
}
