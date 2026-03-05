package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var defaultHost = "localhost:8080"
var defaultResultHost = "http://localhost:8080"

type Conf struct {
	Host       string
	ResultHost string
}

func Get() *Conf {
	conf := Conf{
		Host:       defaultHost,
		ResultHost: defaultResultHost,
	}

	flag.Func("a", "HTTP server URL", func(flagValue string) error {
		urlParts := strings.Split(flagValue, ":")
		if len(urlParts) != 2 {
			return errors.New("invalid URL format")
		}
		if _, err := strconv.Atoi(urlParts[1]); err != nil {
			return fmt.Errorf("invalid port: %s", flagValue)
		}

		conf.Host = flagValue
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
		conf.ResultHost = flagValue
		return nil
	})

	flag.Parse()

	return &conf
}
