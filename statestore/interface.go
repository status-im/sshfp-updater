package statestore

type Repository interface {
	GetModifyIndex(hostname string) (int, error)
	SetModifyIndex(hostname string, index int) error
}

type Service interface {
	CheckIfModified(hostname string, index int) (bool, error)
	SaveState(hostname string, index int) error
}
