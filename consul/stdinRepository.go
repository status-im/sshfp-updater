package consul

import (
	"encoding/json"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type repo struct {
	rawData []byte
}

//Create new file repository satisfying Repository interface
func NewStdinRepository() Repository {
	logrus.Debug("consul: Creating stdin reader")
	return &repo{}
}

//GetData - Load data from stdin and store in the memory
func (r *repo) GetData() error {
	var err error
	logrus.Debug("consul: GetData: Opening stdin")

	rawData, err := io.ReadAll(os.Stdin)
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
