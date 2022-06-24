package config

import (
	"errors"
	"os"

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
		consulToken, exists := os.LookupEnv("CONSUL_TOKEN")
		if !exists {
			return nil, errors.New("cannot find env variable CONSUL_TOKEN")
		}

		domainaName, exists := os.LookupEnv("DOMAIN_NAME")
		if exists {
			return nil, errors.New("cannot find env variable DOMAIN_NAME")
		}

		return &Config{ConsulToken: consulToken, CloudflareToken: cfToken, DomainName: domainaName}, nil

	}

	err := s.r.LoadFile(fileName)
	if err != nil {
		return nil, err
	}

	return s.r.GetConfig(), nil
}
