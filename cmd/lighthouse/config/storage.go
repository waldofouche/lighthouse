package config

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/lighthouse/storage"
)

type storageConf struct {
	DataLocation    string `yaml:"data_location"`
	ReadableStorage bool   `yaml:"human_readable_storage"`
}

func (c *storageConf) validate() error {
	if c.DataLocation == "" {
		return errors.New("error in storage conf: data_location must be specified")
	}
	return nil
}

// LoadStorageBackends loads and returns the storage backends for the passed Config
func LoadStorageBackends(c storageConf) (
	subordinateStorage storage.SubordinateStorageBackend,
	trustMarkedEntitiesStorage storage.TrustMarkedEntitiesStorageBackend, err error,
) {
	if c.ReadableStorage {
		warehouse := storage.NewFileStorage(c.DataLocation)
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
	} else {
		warehouse, err := storage.NewBadgerStorage(c.DataLocation)
		if err != nil {
			return nil, nil, err
		}
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
	}
	log.Info("Loaded storage backend")
	return
}
