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
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Bucket names
var (
	BucketVMs      = []byte("vms")
	BucketConfigs  = []byte("configs")
	BucketUsers    = []byte("users")
	BucketSessions = []byte("sessions")
	BucketSettings = []byte("settings")
)

// Store wraps a BoltDB database
type Store struct {
	db   *bolt.DB
	path string
}

// NewStore creates a new Store instance
func NewStore(path string) (*Store, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{
		db:   db,
		path: path,
	}

	// Initialize buckets
	if err := store.initBuckets(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize buckets: %w", err)
	}

	return store, nil
}

// initBuckets creates all required buckets
func (s *Store) initBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			BucketVMs,
			BucketConfigs,
			BucketUsers,
			BucketSessions,
			BucketSettings,
		}

		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

// Close closes the database
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Path returns the database file path
func (s *Store) Path() string {
	return s.path
}

// Put stores a value in a bucket
func (s *Store) Put(bucket []byte, key string, value any) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		return b.Put([]byte(key), data)
	})
}

// Get retrieves a value from a bucket
func (s *Store) Get(bucket []byte, key string, dest any) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s not found", key)
		}

		return json.Unmarshal(data, dest)
	})
}

// Delete removes a value from a bucket
func (s *Store) Delete(bucket []byte, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Delete([]byte(key))
	})
}

// List retrieves all values from a bucket
func (s *Store) List(bucket []byte, factory func() any) ([]any, error) {
	var results []any

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		return b.ForEach(func(k, v []byte) error {
			item := factory()
			if err := json.Unmarshal(v, item); err != nil {
				return fmt.Errorf("failed to unmarshal item: %w", err)
			}
			results = append(results, item)
			return nil
		})
	})

	return results, err
}

// Count returns the number of items in a bucket
func (s *Store) Count(bucket []byte) (int, error) {
	var count int

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		count = b.Stats().KeyN
		return nil
	})

	return count, err
}

// Exists checks if a key exists in a bucket
func (s *Store) Exists(bucket []byte, key string) (bool, error) {
	var exists bool

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		exists = b.Get([]byte(key)) != nil
		return nil
	})

	return exists, err
}

// Transaction executes a function within a transaction
func (s *Store) Transaction(fn func(tx *bolt.Tx) error) error {
	return s.db.Update(fn)
}

// ViewTransaction executes a read-only function within a transaction
func (s *Store) ViewTransaction(fn func(tx *bolt.Tx) error) error {
	return s.db.View(fn)
}
