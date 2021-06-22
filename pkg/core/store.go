package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type StatusStore interface {
	// TODO: handles Status instead of time.duration
	Save(s map[string]time.Duration) error
	GetAll() map[string]time.Duration
}

var _ StatusStore = &FileStore{}

type FileStore struct {
	Status   map[string]time.Duration `json:"status,omitempty"`
	filename string
}

func NewFileStore(filename string) (StatusStore, error) {
	store := &FileStore{filename: filename}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return store, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return store, err
	}

	err = json.Unmarshal(data, store)
	return store, err
}

func (f *FileStore) Save(s map[string]time.Duration) error {
	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	f.Status = s

	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func (f *FileStore) GetAll() map[string]time.Duration {
	return f.Status
}
