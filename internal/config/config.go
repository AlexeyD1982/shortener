package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
)

var defaultHost = "localhost:8080"
var defaultResultHost = "http://localhost:8080"

type Conf struct {
	Host       string
	ResultHost string
}

func Get(withFlags bool) *Conf {
	var haveHostEnv, haveResultHostEnv bool

	conf := Conf{
		Host:       defaultHost,
		ResultHost: defaultResultHost,
	}

	if hostEnv := os.Getenv("SERVER_ADDRESS"); hostEnv != "" {
		addr, err := validateServerAddrAndGetHost(hostEnv)
		if err == nil {
			conf.Host = addr
		}
		haveHostEnv = true
	}

	if resultHostEnv := os.Getenv("BASE_URL"); resultHostEnv != "" {
		_, err := url.Parse(resultHostEnv)
		if err == nil {
			conf.ResultHost = resultHostEnv
		}
		haveResultHostEnv = true
	}

	if withFlags {
		if !haveHostEnv {
			flag.Func("a", "HTTP server URL", func(flagValue string) error {
				addr, err := validateServerAddrAndGetHost(flagValue)
				if err != nil {
					return fmt.Errorf("%w: %s", err, flagValue)
				}
				conf.Host = addr
				return nil
			})
		}

		if !haveResultHostEnv {
			flag.Func("b", "result short URL", func(flagValue string) error {
				u, err := url.Parse(flagValue)
				if err != nil {
					return errors.New("invalid Result URL format")
				}

				if u.Hostname() == "" {
					return errors.New("invalid Result URL format")
				}
				conf.ResultHost = flagValue
				return nil
			})
		}

		flag.Parse()
	}

	return &conf
}

func validateServerAddrAndGetHost(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid URL format")
	}
	return u.Host, nil
}
