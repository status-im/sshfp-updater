package sshfp

type Service interface {
	ParseConsulSSHRecord(key string, value string) (*SSHFPRecord, error)
	ParseConsulSSHRecords(records map[string]string) []*SSHFPRecord
	PrepareConfiguration(hostname string, current []*SSHFPRecord, new []*SSHFPRecord) ConfigPlan
	PrintConfigPlan(configPlan ConfigPlan)
}
