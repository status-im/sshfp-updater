package consul

type Repository interface {
	GetData() error
	ParseData() (hostsMap, error)
}

type Service interface {
	LoadData() error
	GetHostnames() []string
	GetModifiedIndex(hostname string) int
	GetCreateIndex(hostname string) int
	GetMetaData(hostname string) map[string]string
}
