package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type repository struct {
	data Config
}

func NewFileRepository() Repository {
	return &repository{data: Config{}}
}

func (r *repository) LoadFile(fileName string) error {
	logrus.Infof("config: LoadFile %s", fileName)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &r.data)
	if err != nil {
		return err
	}
	return nil

}

func (r *repository) SaveFile(fileName string) error {
	logrus.Infof("config: SaveFile %s", fileName)
	content, err := json.MarshalIndent(r.data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetConfig() *Config {
	return &r.data
}
