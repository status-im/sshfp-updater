package statestore

import "github.com/sirupsen/logrus"

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// CheckIfModified - returns if host is modified or not. Err can be ignored or not depends on repository
func (s *service) CheckIfModified(hostname string, index int) (bool, error) {
	logrus.Debugf("statestore: CheckIfModified %s", hostname)
	indexDb, err := s.r.GetModifyIndex(hostname)

	if err != nil {
		return true, err
	}

	if indexDb == index {
		return false, nil
	}

	return true, err
}

func (s *service) SaveState(hostname string, index int) error {
	return s.r.SetModifyIndex(hostname, index)

}
