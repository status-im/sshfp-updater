package consul

import (
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"
)

type repo struct {
	rawData []byte
	reader  io.Reader
}

//Create new file repository satisfying Repository interface
func NewStdinRepository(reader io.Reader) Repository {
	logrus.Debug("consul: Creating stdin reader")
	return &repo{reader: reader}
}

//GetData - Load data from stdin and store in the memory
func (r *repo) GetData() error {
	var err error
	logrus.Debug("consul: GetData: Opening stdin")

	rawData, err := io.ReadAll(r.reader)
	if err != nil {
		logrus.Fatal(err)
	}

	r.rawData = rawData

	return err
}

//ParseData - Parse loaded data and return as json.
func (r *repo) ParseData() (hostsMap, error) {
	logrus.Debugf("consul: ParseData: Parsing Data")
	var hosts rawHosts
	if err := json.Unmarshal(r.rawData, &hosts); err != nil {
		return nil, err
	}

	hostsMap := make(hostsMap)
	for _, v := range hosts {
		hostsMap[v.Node.Node] = v
	}

	return hostsMap, nil
}
