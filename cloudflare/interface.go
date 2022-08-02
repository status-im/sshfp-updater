package cloudflare

import (
	"sshfp-updater/sshfp"

	"github.com/cloudflare/cloudflare-go"
)

type Repository interface {
	FindRecords(hostname string, recordType string) ([]cloudflare.DNSRecord, error)
	CreateDNSRecord(hostname string, recordType string, payload cloudflare.DNSRecord) (int, error)
	DeleteDNSRecord(recordID string) error
	UpdateDNSRecord(hostname string, recordID string, payload cloudflare.DNSRecord) error
}

type Service interface {
	FindHostByName(hostname string) (bool, error)
	GetSSHFPRecordsForHost(hostname string) ([]*sshfp.SSHFPRecord, error)
	DeleteSSHFPRecordsForHost(hostname string) error
	CreateSSHFPRecord(hostname string, record sshfp.SSHFPRecord) (int, error)
	DeleteSSHFPRecord(hostname string, record sshfp.SSHFPRecord) error
	UpdateSSHFPRecord(hostname string, record sshfp.SSHFPRecord) error
	ApplyConfigPlan(configPlan sshfp.ConfigPlan) (int, error)
}
