// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/firecracker-microvm/firectl/pkg/models"
	bolt "go.etcd.io/bbolt"
)

// UserStore provides user-specific storage operations
type UserStore struct {
	store *Store
}

// NewUserStore creates a new UserStore
func NewUserStore(store *Store) *UserStore {
	return &UserStore{store: store}
}

// Create stores a new user
func (us *UserStore) Create(user *models.User) error {
	// Check if username already exists
	existing, _ := us.GetByUsername(user.Username)
	if existing != nil {
		return models.ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	return us.store.Put(BucketUsers, user.ID, user)
}

// Get retrieves a user by ID
func (us *UserStore) Get(id string) (*models.User, error) {
	var user models.User
	if err := us.store.Get(BucketUsers, id, &user); err != nil {
		return nil, models.ErrUserNotFound
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (us *UserStore) GetByUsername(username string) (*models.User, error) {
	var found *models.User

	err := us.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketUsers)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var user models.User
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			if user.Username == username {
				found = &user
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, models.ErrUserNotFound
	}
	return found, nil
}

// Update updates an existing user
func (us *UserStore) Update(user *models.User) error {
	exists, err := us.store.Exists(BucketUsers, user.ID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	return us.store.Put(BucketUsers, user.ID, user)
}

// Delete removes a user
func (us *UserStore) Delete(id string) error {
	exists, err := us.store.Exists(BucketUsers, id)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrUserNotFound
	}
	return us.store.Delete(BucketUsers, id)
}

// List returns all users
func (us *UserStore) List() ([]*models.User, error) {
	var users []*models.User

	err := us.store.ViewTransaction(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketUsers)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var user models.User
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			users = append(users, &user)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return users, nil
}

// Count returns the total number of users
func (us *UserStore) Count() (int, error) {
	return us.store.Count(BucketUsers)
}

// UpdateLastLogin updates the last login timestamp
func (us *UserStore) UpdateLastLogin(id string) error {
	user, err := us.Get(id)
	if err != nil {
		return err
	}
	now := time.Now()
	user.LastLoginAt = &now
	return us.Update(user)
}

// IsSetupRequired checks if initial setup is needed
func (us *UserStore) IsSetupRequired() (bool, error) {
	count, err := us.Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
