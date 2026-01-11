// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anubhavg-icpl/agni/pkg/models"
	bolt "go.etcd.io/bbolt"
)

// ConfigStore provides configuration template storage operations
type ConfigStore struct {
	store *Store
}

// NewConfigStore creates a new ConfigStore
func NewConfigStore(store *Store) *ConfigStore {
	return &ConfigStore{store: store}
}

// Create stores a new configuration template
func (cs *ConfigStore) Create(config *models.ConfigTemplate) error {
	exists, err := cs.store.Exists(BucketConfigs, config.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("configuration already exists")
	}
	config.CreatedAt = time.Now()
	config.UpdatedAt = config.CreatedAt
	return cs.store.Put(BucketConfigs, config.ID, config)
}

// Get retrieves a configuration template by ID
func (cs *ConfigStore) Get(id string) (*models.ConfigTemplate, error) {
	var config models.ConfigTemplate
	if err := cs.store.Get(BucketConfigs, id, &config); err != nil {
		return nil, models.ErrConfigNotFound
	}
	return &config, nil
}

// Update updates an existing configuration template
func (cs *ConfigStore) Update(config *models.ConfigTemplate) error {
	exists, err := cs.store.Exists(BucketConfigs, config.ID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrConfigNotFound
	}
	config.UpdatedAt = time.Now()
	return cs.store.Put(BucketConfigs, config.ID, config)
}

// Delete removes a configuration template
func (cs *ConfigStore) Delete(id string) error {
	exists, err := cs.store.Exists(BucketConfigs, id)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrConfigNotFound
	}
	return cs.store.Delete(BucketConfigs, id)
}

// List returns all configuration templates
func (cs *ConfigStore) List() ([]*models.ConfigTemplate, error) {
	configs := make([]*models.ConfigTemplate, 0)

	err := cs.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketConfigs)
		if b == nil {
			// Bucket doesn't exist yet, return empty list
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var config models.ConfigTemplate
			if err := json.Unmarshal(v, &config); err != nil {
				return err
			}
			configs = append(configs, &config)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return configs, nil
}

// GetByName finds a configuration template by name
func (cs *ConfigStore) GetByName(name string) (*models.ConfigTemplate, error) {
	var found *models.ConfigTemplate

	err := cs.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketConfigs)
		if b == nil {
			// Bucket doesn't exist yet
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var config models.ConfigTemplate
			if err := json.Unmarshal(v, &config); err != nil {
				return err
			}
			if config.Name == name {
				found = &config
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, models.ErrConfigNotFound
	}
	return found, nil
}

// Count returns the total number of configuration templates
func (cs *ConfigStore) Count() (int, error) {
	return cs.store.Count(BucketConfigs)
}
