package sshfp

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type service struct {
	Algorithms AlgorithmTypes
	Hashes     HashTypes
}

func NewService() Service {
	//Based on: https://www.iana.org/assignments/dns-sshfp-rr-parameters/dns-sshfp-rr-parameters.xhtml

	hashes := make(HashTypes)
	hashes["sha1"] = "1"
	hashes["sha256"] = "2"

	algorithms := make(AlgorithmTypes)
	algorithms["rsa"] = "1"
	algorithms["dsa"] = "2"
	algorithms["ecdsa"] = "3"
	algorithms["ed25519"] = "4"
	algorithms["ed449"] = "6"

	return &service{Algorithms: algorithms, Hashes: hashes}

}

func (s *service) ParseConsulSSHRecord(key string, value string) (*SSHFPRecord, error) {
	logrus.Debugf("SSHFP: ParseConsulSSHRecord: %s %s", key, value)

	output := &SSHFPRecord{}

	key = strings.ToLower(key)
	value = strings.ToLower(value)

	output.Fingerprint = value

	splittedKey := strings.Split(key, "-")

	if len(splittedKey) < 3 {
		return nil, errors.New("incorrect format")
	}

	//Assumption: Last field is hash type
	output.Type = s.Hashes[splittedKey[len(splittedKey)-1]]

	//Check first field - for "ecdsa" we can return value, for "ssh", we have to check second field.

	switch splittedKey[0] {
	case "ecdsa":
		output.Algorithm = s.Algorithms[splittedKey[0]]

	case "ssh":
		output.Algorithm = s.Algorithms[splittedKey[1]]

	}

	if output.Type == "" || output.Algorithm == "" || output.Fingerprint == "" {
		return nil, errors.New("cannot parse record")
	}
	return output, nil
}

func (s *service) ParseConsulSSHRecords(records map[string]string) []*SSHFPRecord {
	output := make([]*SSHFPRecord, 0)

	for k, v := range records {
		sshRecord, err := s.ParseConsulSSHRecord(k, v)
		if err != nil {
			continue
		}

		output = append(output, sshRecord)
	}

	return output
}

func (s *service) PrepareConfiguration(hostname string, current []*SSHFPRecord, new []*SSHFPRecord) ConfigPlan {
	logrus.Debug("SSHFP: PrepareConfiguration")

	/* Config plan have to cover the following scenarios:
	0. Both configs are empty - do nothing.
	1. New config is empty - DELETE everything
	2. Old config is empty - CREATE everything
	3. Record in new config doesn't exist in old config - CREATE record
	4. Record in new config has different fingerprint than the old one - UPDATE and remove from oldMap
	5. Record in new config and record in old config are the same - just remove from oldMap to avoid unwanted removal

	Records left in oldMap should be qualified for removal
	*/

	configPlan := make(ConfigPlan, 0)

	// Scenario 0
	if len(current) == 0 && len(new) == 0 {
		logrus.Debug("Both configs are empty - do nothing")
		return configPlan
	}

	//Scenario 1
	if len(current) > 0 && len(new) == 0 {
		logrus.Debug("New config is empty - remove config")
		for _, v := range current {
			configPlan = append(configPlan, ConfigPlanElement{Operation: DELETE, Record: v, Hostname: hostname})
		}
		return configPlan
	}

	//Scenario 2
	if len(current) == 0 && len(new) > 0 {
		logrus.Debug("Old config is empty - create config")
		for _, v := range new {
			configPlan = append(configPlan, ConfigPlanElement{Operation: CREATE, Record: v, Hostname: hostname})
		}
		return configPlan
	}

	//To handle scenarios 3-5 we have to create temporary maps with string "<algorithm><type>" as key.
	//Assumption: Pair algorithm+hash type is unique

	oldMap := make(map[string]*SSHFPRecord)
	newMap := make(map[string]*SSHFPRecord)

	//Create temporary maps for better searching
	for _, v := range current {
		oldMap[v.Algorithm+v.Type] = v
	}

	for _, v := range new {
		newMap[v.Algorithm+v.Type] = v
	}

	for k := range newMap {
		//Scenario 3
		if _, ok := oldMap[k]; !ok {
			logrus.Debugf("Record not found in current config: %s", k)
			configPlan = append(configPlan, ConfigPlanElement{Operation: CREATE, Record: newMap[k], Hostname: hostname})
			continue
		}
		//Scenario 4
		if oldMap[k].Fingerprint != newMap[k].Fingerprint {
			logrus.Debugf("Updating record in current config: %s", k)
			newMap[k].RecordID = oldMap[k].RecordID
			configPlan = append(configPlan, ConfigPlanElement{Operation: UPDATE, Record: newMap[k], Hostname: hostname})
			delete(oldMap, k)
			continue
		}
		//Scenario 5
		if oldMap[k].Fingerprint == newMap[k].Fingerprint {
			delete(oldMap, k)
			continue
		}
	}

	//Cleanup
	for _, v := range oldMap {
		configPlan = append(configPlan, ConfigPlanElement{Operation: DELETE, Record: v, Hostname: hostname})
	}
	return configPlan
}
func (s *service) PrintConfigPlan(configPlan ConfigPlan) {
	logrus.Debug("SSHFP: PrintConfigPlan")
	logrus.Infof("Config Plan:")
	for _, v := range configPlan {
		logrus.Infof("Hostname: %s", v.Hostname)
		logrus.Infof("Operation: %v", v.Operation)
		logrus.Infof("Record: %+v", v.Record)
		logrus.Infof("---")
	}
}
