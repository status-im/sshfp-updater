package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) LoadConfig(fileName string) (*Config, error) {
	if s.r == nil || fileName == "" {
		logrus.Infoln("config: LoadConfig")
		cfToken, exists := os.LookupEnv("CF_TOKEN")
		if !exists {
			return nil, errors.New("cannot find env variable CF_TOKEN")
		}
		domainaName, exists := os.LookupEnv("DOMAIN_NAME")
		if !exists {
			return nil, errors.New("cannot find env variable DOMAIN_NAME")
		}

		hostTimeout, exists := os.LookupEnv("HOST_LIVENESS_TIMEOUT")
		if !exists {
			return nil, errors.New("cannot find env variable HOST_LIVENESS_TIMEOUT")
		}
		hostTimeoutInt, err := strconv.ParseInt(hostTimeout, 10, 32)
		if err != nil {
			return nil, errors.New("incorrect HOST_LIVENESS_TIMEOUT value")

		}

		logLevel, exists := os.LookupEnv(("LOG_LEVEL"))
		if !exists {
			return nil, errors.New("cannot find env variable LOG_LEVEL")
		}

		return &Config{CloudflareToken: cfToken, DomainName: domainaName, HostTimeout: int(hostTimeoutInt), LogLevel: logLevel}, nil

	}

	err := s.r.LoadFile(fileName)
	if err != nil {
		return nil, err
	}

	return s.r.GetConfig(), nil
}
