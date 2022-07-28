package config

type Config struct {
	ConsulToken     string `json:"consulToken"`
	CloudflareToken string `json:"cloudflareKey"`
	DomainName      string `json:"domain"`
	HostTimeout     int    `json:"hostTimeout"`
	LogLevel        string `json:"logLevel"`
	StorageFilePath string `json:"storageFilePath"`
}
