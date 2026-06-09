package storage

import (
	"encoding/json"

	"github.com/go-oidfed/lib/jwx/keymanagement/kms"
	"github.com/go-oidfed/lighthouse/storage/model"
	"github.com/pkg/errors"
)

// DBStateStorer implements kms.KMSStateStorer using the key-value storage.
type DBStateStorer struct {
	kv  model.KeyValueStore
	key string
}

// NewDBStateStorer creates a new DBStateStorer.
func NewDBStateStorer(kv model.KeyValueStore, typeID string) *DBStateStorer {
	stateKey := "kms_state"
	if typeID != "" {
		stateKey = stateKey + "_" + typeID
	}
	return &DBStateStorer{
		kv:  kv,
		key: stateKey,
	}
}

// LoadScheduledState loads the scheduled KMS state from the key-value store.
func (d *DBStateStorer) LoadScheduledState() (kms.ScheduledState, error) {
	val, err := d.kv.Get(model.KeyValueScopeSigning, d.key)
	if err != nil {
		return kms.ScheduledState{}, errors.WithStack(err)
	}
	if len(val) == 0 {
		return kms.ScheduledState{}, nil
	}
	var st kms.ScheduledState
	err = json.Unmarshal(val, &st)
	return st, errors.WithStack(err)
}

// SaveScheduledState saves the scheduled KMS state to the key-value store.
func (d *DBStateStorer) SaveScheduledState(st kms.ScheduledState) error {
	data, err := json.Marshal(st)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(d.kv.Set(model.KeyValueScopeSigning, d.key, data))
}
