package sshfp

type SSHFPRecord struct {
	Algorithm   string
	Type        string
	Fingerprint string
	RecordID    string
}

type AlgorithmTypes map[string]string
type HashTypes map[string]string

type ConfigPlanElement struct {
	Hostname  string
	Operation Operation
	Record    *SSHFPRecord
}

type ConfigPlan []ConfigPlanElement

type Operation int

const (
	UPDATE Operation = iota
	CREATE
	DELETE
)
