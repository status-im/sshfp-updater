package cloudflare

import (
	"errors"
	"fmt"
	"infra-sshfp-cf/sshfp"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	logrus.Debug("cloudflare: Creating Service")
	return &service{r: r}
}

func (s *service) FindHostByName(hostname string) (bool, error) {
	logrus.Debugf("cloudflare: FindHostByName %s", hostname)
	recs, err := s.r.FindRecords(hostname, "A")
	if err != nil {
		logrus.Errorf("CF: %s", err)
		return false, err
	}

	if len(recs) > 0 {
		logrus.Infof("FindHostByName: Host found: %s", hostname)
		logrus.Debugf("%+v", recs)
		return true, nil
	}
	logrus.Infof("FindHostByName: Host not found: %s", hostname)
	return false, nil
}

func (s *service) GetSSHFPRecordsForHost(hostname string) ([]*sshfp.SSHFPRecord, error) {
	logrus.Debugf("cloudflare: GetSSHFPRecordsForHost %s", hostname)
	recs, err := s.r.FindRecords(hostname, "SSHFP")
	if err != nil {
		logrus.Errorf("CF: %s", err)
		return nil, err
	}
	logrus.Debugf("%+v", recs)

	output := make([]*sshfp.SSHFPRecord, 0)

	for _, v := range recs {
		data, ok := v.Data.(map[string]interface{})
		/*
			data["type"] - float64
			data["fingerprint"] - string
			data["algorithm"] - float64
		*/

		if !ok {
			return nil, errors.New("cannot parse data")
		}
		output = append(output, &sshfp.SSHFPRecord{RecordID: v.ID, Algorithm: fmt.Sprintf("%.f", data["algorithm"]), Type: fmt.Sprintf("%.f", data["type"]), Fingerprint: data["fingerprint"].(string)})
	}

	return output, nil
}

func (s *service) DeleteSSHFPRecordsForHost(hostname string) error {
	logrus.Debugf("cloudflare: DeleteSSHFPRecordsForHost: %s", hostname)
	records, err := s.GetSSHFPRecordsForHost(hostname)
	if err != nil {
		return err
	}

	for _, record := range records {
		err := s.DeleteSSHFPRecord(hostname, *record)
		if err != nil {
			return nil
		}
	}
	return nil

}

func (s *service) CreateSSHFPRecord(hostname string, record sshfp.SSHFPRecord) (int, error) {
	logrus.Infof("cloudflare: CreateSSHFPRecord: %+v", record)

	payload := s.preparePayloadFromSSHRecord(record)

	return s.r.CreateDNSRecord(hostname, "SSHFP", payload)
}

func (s *service) DeleteSSHFPRecord(hostname string, record sshfp.SSHFPRecord) error {
	logrus.Infof("cloudflare: DeleteDNSRecord: %+v", record)

	//In case we have record ID we can just delete the record.
	if record.RecordID != "" {
		return s.r.DeleteDNSRecord(record.RecordID)
	}

	recs, err := s.GetSSHFPRecordsForHost(hostname)
	if err != nil {
		return err
	}
	//Comparing fingerprints should be more than sufficient
	for _, rec := range recs {
		if rec.Fingerprint == record.Fingerprint {
			return s.r.DeleteDNSRecord(rec.RecordID)
		}
	}
	return errors.New("record not found")

}

func (s *service) ApplyConfigPlan(configPlan sshfp.ConfigPlan) (int, error) {
	var item int = 0
	for _, v := range configPlan {

		switch v.Operation {
		case sshfp.CREATE:
			_, err := s.CreateSSHFPRecord(v.Hostname, *v.Record)
			if err != nil {
				return item, err
			}
		case sshfp.DELETE:
			err := s.DeleteSSHFPRecord(v.Hostname, *v.Record)
			if err != nil {
				return item, err
			}
		case sshfp.UPDATE:
			err := s.UpdateSSHFPRecord(v.Hostname, *v.Record)
			if err != nil {
				return item, err
			}
		}
		item++
	}

	return item, nil
}

func (s *service) UpdateSSHFPRecord(hostname string, record sshfp.SSHFPRecord) error {
	logrus.Infof("UpdateSSHFPRecord: %+v", record)
	payload := s.preparePayloadFromSSHRecord(record)
	return s.r.UpdateDNSRecord(hostname, record.RecordID, payload)
}

func (s *service) preparePayloadFromSSHRecord(record sshfp.SSHFPRecord) cloudflare.DNSRecord {
	data := make(map[string]string)
	data["algorithm"] = record.Algorithm
	data["type"] = record.Type
	data["fingerprint"] = record.Fingerprint

	return cloudflare.DNSRecord{Type: "SSHFP", Data: data, Content: fmt.Sprintf("%s %s %s", record.Algorithm, record.Type, record.Fingerprint)}
}
