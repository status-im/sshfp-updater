package statestore

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type mapRepository struct {
	db       map[string]int
	filename string
}

func NewMapRepository(filename string) Repository {
	repo := new(mapRepository)
	repo.filename = filename

	err := repo.openDatabase()
	if err == nil {
		return repo
	}

	db := make(map[string]int)
	repo.db = db

	err = repo.saveDatabase()
	if err != nil {
		logrus.Error("mapRepository: cannot save database %s", filename)
	}

	return repo
}

func (r *mapRepository) GetModifyIndex(hostname string) (int, error) {
	if value, ok := r.db[hostname]; ok {
		return value, nil
	}
	return -1, nil
}

func (r *mapRepository) SetModifyIndex(hostname string, index int) error {
	r.db[hostname] = index
	return r.saveDatabase()
}

func (r *mapRepository) openDatabase() error {
	logrus.Infof("mapRepository: openDatabase %s", r.filename)
	content, err := ioutil.ReadFile(r.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &r.db)
	if err != nil {
		return err
	}
	return nil

}

func (r *mapRepository) saveDatabase() error {
	logrus.Infof("mapRepository: saveDatabase %s", r.filename)
	content, err := json.MarshalIndent(&r.db, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(r.filename, content, 0644)
	if err != nil {
		return err
	}
	return nil

}
