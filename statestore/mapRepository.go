package statestore

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
)

type mapEntry struct {
	ModifyIndex int
	LastSeen    time.Time
}

type mapRepository struct {
	db       map[string]mapEntry
	filename string
}

func NewMapRepository(filename string) Repository {

	if filename == "" {
		filename = "stateStore.json"
	}

	repo := new(mapRepository)
	repo.filename = filename

	err := repo.openDatabase()
	if err == nil {
		return repo
	}

	db := make(map[string]mapEntry)
	repo.db = db

	err = repo.saveDatabase()
	if err != nil {
		logrus.Error("mapRepository: cannot save database %s", filename)
	}

	return repo
}

func (r *mapRepository) GetModifyIndex(hostname string) (int, error) {
	if value, ok := r.db[hostname]; ok {
		r.SetModifyIndex(hostname, value.ModifyIndex)
		return value.ModifyIndex, nil
	}
	return -1, nil
}

func (r *mapRepository) GetOutdatedHosts(duration time.Duration) ([]string, error) {
	var output []string
	for k, v := range r.db {

		if v.LastSeen.Add(duration).Before(time.Now()) {
			output = append(output, k)
		}

	}
	return output, nil
}

func (r *mapRepository) SetModifyIndex(hostname string, index int) error {
	if value, ok := r.db[hostname]; ok {
		value.LastSeen = time.Now()
		value.ModifyIndex = index
		r.db[hostname] = value
	} else {
		value := new(mapEntry)
		value.LastSeen = time.Now()
		value.ModifyIndex = index
		r.db[hostname] = *value
	}
	return r.saveDatabase()
}

func (r *mapRepository) DeleteHost(hostname string) error {

	logrus.Debugf("mapRepository: DeleteHost %s", hostname)

	delete(r.db, hostname)
	return r.saveDatabase()
}

func (r *mapRepository) DeleteHosts(hostnames []string) error {
	for _, host := range hostnames {
		r.DeleteHost(host)
	}
	return nil
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
