package cloudflare

import (
	"context"
	"errors"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

type repository struct {
	api        *cloudflare.API
	domainId   string
	domainName string
	ctx        context.Context
}

func NewRepository(token string, domain string) (Repository, error) {
	logrus.Infof("cloudflare: Creating Repository")
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		logrus.Errorf("cloudflare: token init failure: %s", err)
		return nil, err
	}

	id, err := api.ZoneIDByName(domain)
	if err != nil {
		logrus.Errorf("cloudflare: zone init failure: %s", err)
		return nil, err
	}
	repo := repository{
		api:        api,
		domainId:   id,
		domainName: domain,
		ctx:        context.Background(),
	}
	return &repo, nil
}

func (r *repository) FindRecords(hostname string, recordType string) ([]cloudflare.DNSRecord, error) {
	logrus.Debugf("cloudflare: FindRecords: %s", hostname)

	filterSet := cloudflare.DNSRecord{Type: recordType, Name: hostname + "." + r.domainName}

	logrus.Debugf("%+v", filterSet)
	return r.api.DNSRecords(r.ctx, r.domainId, filterSet)
}

func (r *repository) CreateDNSRecord(hostname string, recordType string, payload cloudflare.DNSRecord) (int, error) {
	logrus.Infof("cloudflare: CreateDNSRecord: %s", hostname)

	payload.Name = hostname + "." + r.domainName

	logrus.Debugf("cloudflare: CreateDNSRecord - payload: %s", payload)

	resp, err := r.api.CreateDNSRecord(r.ctx, r.domainId, payload)

	logrus.Debugf("%+v", resp)

	if err != nil {
		return -1, err
	}
	if !resp.Success {
		return -1, errors.New("cannot create record")
	}
	return resp.Count, nil
}

func (r *repository) DeleteDNSRecord(recordID string) error {
	logrus.Debugf("cloudflare: DeleteDNSRecord - %s", recordID)

	return r.api.DeleteDNSRecord(r.ctx, r.domainId, recordID)
}

func (r *repository) UpdateDNSRecord(hostname string, recordID string, payload cloudflare.DNSRecord) error {
	logrus.Debugf("cloudflare: UpdateeDNSRecord - %s", recordID)
	return r.api.UpdateDNSRecord(r.ctx, r.domainId, recordID, payload)
}
