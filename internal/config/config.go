package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
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
		err := validateServerAddr(hostEnv)
		if err == nil {
			conf.Host = hostEnv
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
		flag.Func("a", "HTTP server URL", func(flagValue string) error {
			err := validateServerAddr(flagValue)
			if err != nil {
				return fmt.Errorf("%w: %s", err, flagValue)
			}

			if !haveHostEnv {
				conf.Host = flagValue
			}
			return nil
		})

		flag.Func("b", "result short URL", func(flagValue string) error {
			u, err := url.Parse(flagValue)
			if err != nil {
				return errors.New("invalid Result URL format")
			}

			if u.Hostname() == "" {
				return errors.New("invalid Result URL format")
			}

			if !haveResultHostEnv {
				conf.ResultHost = flagValue
			}
			return nil
		})

		flag.Parse()
	}

	return &conf
}

func validateServerAddr(rawURL string) error {
	urlParts := strings.Split(rawURL, ":")
	if len(urlParts) != 2 {
		return errors.New("invalid URL format")
	}
	if _, err := strconv.Atoi(urlParts[1]); err != nil {
		return fmt.Errorf("invalid port - %s", rawURL)
	}
	return nil
}
