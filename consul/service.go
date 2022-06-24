package consul

import (
	"github.com/sirupsen/logrus"
)

type service struct {
	r        Repository
	hostsMap hostsMap
}

//Generates new service instance based on the repo
func NewService(repository Repository) Service {
	logrus.Debug("consul: NewService")
	return &service{r: repository}
}

//Function calling underlying repository, retreive data and store as interface map
func (s *service) LoadData() error {
	logrus.Debug("consul: LoadData")
	err := s.r.GetData()
	if err != nil {
		return err
	}

	s.hostsMap, err = s.r.ParseData()
	logrus.Debugf("%+v", s.hostsMap)
	return err
}

func (s *service) GetHostnames() []string {
	logrus.Debug("consul: GetHostnames")
	hostnames := make([]string, 0)

	for k := range s.hostsMap {
		hostnames = append(hostnames, k)
	}
	return hostnames

}

func (s *service) GetModifiedIndex(hostname string) int {
	return s.hostsMap[hostname].Service.ModifyIndex
}

func (s *service) GetCreateIndex(hostname string) int {
	return s.hostsMap[hostname].Service.CreateIndex
}

func (s *service) GetMetaData(hostname string) map[string]string {
	return s.hostsMap[hostname].Service.Meta
}
