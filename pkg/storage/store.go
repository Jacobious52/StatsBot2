package storage

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// DataStore defines the storage structure that owns a model
type DataStore struct {
	sync.RWMutex

	Data  Model
	Dirty bool
	path  string
}

// NewDataStore creates a new datastore with an empty model
func NewDataStore(path string) *DataStore {
	return &DataStore{
		Data: make(Model),
		path: path,
	}
}

// Write writes the model to disk
func (s *DataStore) Write(w io.Writer) error {
	err := json.NewEncoder(w).Encode(&s.Data)
	if err != nil {
		return err
	}
	return nil
}

// Load reads the model to disk
func (s *DataStore) Load() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Debugln("locking store")
	s.Lock()
	defer s.Unlock()
	defer log.Debugln("unlocking store")

	err = json.NewDecoder(file).Decode(&s.Data)
	if err != nil {
		return err
	}
	return nil
}

// StartSync starts a routine to periodically save the db if there are changes
func (s *DataStore) StartSync(stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				log.Debugln("locking store")
				s.Lock()
				if s.Dirty {
					file, err := os.Create(s.path)
					if err != nil {
						log.Errorf("data not saved. %v", err)
						log.Debugln("unlocking store")
						s.Unlock()
						break
					}
					err = s.Write(file)

					if err != nil {
						log.Errorf("data not saved. %v", err)
					} else {
						log.Infoln("saved data")
						s.Dirty = false
					}
					file.Close()
				}
				log.Debugln("unlocking store")
				s.Unlock()
			case <-stop:
				log.Infoln("stopping data sync")
				ticker.Stop()
				break
			}
		}
	}()
}
