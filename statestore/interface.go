package statestore

import "time"

type Repository interface {
	GetModifyIndex(hostname string) (int, error)
	SetModifyIndex(hostname string, index int) error
	DeleteHost(hostname string) error
	DeleteHosts(hostnames []string) error
	GetOutdatedHosts(time time.Duration) ([]string, error)
}

type Service interface {
	CheckIfModified(hostname string, index int) (bool, error)
	SaveState(hostname string, index int) error
	PurgeStalledHosts(timeTreshold int) error
	GetStalledHosts(timeTreshold int) ([]string, error)
}
