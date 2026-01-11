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

	"github.com/firecracker-microvm/firectl/pkg/models"
	bolt "go.etcd.io/bbolt"
)

// VMStore provides VM-specific storage operations
type VMStore struct {
	store *Store
}

// NewVMStore creates a new VMStore
func NewVMStore(store *Store) *VMStore {
	return &VMStore{store: store}
}

// Create stores a new VM
func (vs *VMStore) Create(vm *models.VM) error {
	exists, err := vs.store.Exists(BucketVMs, vm.ID)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrVMAlreadyExists
	}
	return vs.store.Put(BucketVMs, vm.ID, vm)
}

// Get retrieves a VM by ID
func (vs *VMStore) Get(id string) (*models.VM, error) {
	var vm models.VM
	if err := vs.store.Get(BucketVMs, id, &vm); err != nil {
		return nil, models.ErrVMNotFound
	}
	return &vm, nil
}

// Update updates an existing VM
func (vs *VMStore) Update(vm *models.VM) error {
	exists, err := vs.store.Exists(BucketVMs, vm.ID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrVMNotFound
	}
	return vs.store.Put(BucketVMs, vm.ID, vm)
}

// Delete removes a VM
func (vs *VMStore) Delete(id string) error {
	exists, err := vs.store.Exists(BucketVMs, id)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrVMNotFound
	}
	return vs.store.Delete(BucketVMs, id)
}

// List returns all VMs
func (vs *VMStore) List() ([]*models.VM, error) {
	var vms []*models.VM

	err := vs.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketVMs)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var vm models.VM
			if err := json.Unmarshal(v, &vm); err != nil {
				return err
			}
			vms = append(vms, &vm)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return vms, nil
}

// GetByName finds a VM by name
func (vs *VMStore) GetByName(name string) (*models.VM, error) {
	var found *models.VM

	err := vs.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketVMs)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var vm models.VM
			if err := json.Unmarshal(v, &vm); err != nil {
				return err
			}
			if vm.Name == name {
				found = &vm
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, models.ErrVMNotFound
	}
	return found, nil
}

// Count returns the total number of VMs
func (vs *VMStore) Count() (int, error) {
	return vs.store.Count(BucketVMs)
}

// UpdateStatus updates only the status of a VM
func (vs *VMStore) UpdateStatus(id string, status models.VMStatus) error {
	vm, err := vs.Get(id)
	if err != nil {
		return err
	}
	vm.Status = status
	return vs.Update(vm)
}
