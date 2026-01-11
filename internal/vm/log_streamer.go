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

package vm

import (
	"sync"
	"time"

	"github.com/anubhavg-icpl/agni/pkg/models"
)

// LogSubscriber represents a subscriber to VM logs
type LogSubscriber struct {
	ID      string
	VMID    string
	Channel chan *models.LogEntry
	Filter  LogFilter
	done    chan struct{}
}

// LogFilter defines log filtering options
type LogFilter struct {
	Level    string   // Filter by log level
	Contains string   // Filter by message content
	Sources  []string // Filter by source
}

// LogStreamer manages log streaming to subscribers
type LogStreamer struct {
	subscribers map[string]*LogSubscriber
	mu          sync.RWMutex
	buffers     map[string][]*models.LogEntry // Recent logs per VM
	bufferSize  int
}

// NewLogStreamer creates a new LogStreamer
func NewLogStreamer() *LogStreamer {
	return &LogStreamer{
		subscribers: make(map[string]*LogSubscriber),
		buffers:     make(map[string][]*models.LogEntry),
		bufferSize:  1000, // Keep last 1000 logs per VM
	}
}

// Subscribe creates a new subscription for VM logs
func (ls *LogStreamer) Subscribe(vmID string, filter LogFilter) *LogSubscriber {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	sub := &LogSubscriber{
		ID:      generateSubscriberID(),
		VMID:    vmID,
		Channel: make(chan *models.LogEntry, 100),
		Filter:  filter,
		done:    make(chan struct{}),
	}

	ls.subscribers[sub.ID] = sub

	// Send recent logs from buffer
	if logs, ok := ls.buffers[vmID]; ok {
		go func() {
			for _, log := range logs {
				if ls.matchesFilter(log, filter) {
					select {
					case sub.Channel <- log:
					case <-sub.done:
						return
					default:
						// Channel full, skip
					}
				}
			}
		}()
	}

	return sub
}

// Unsubscribe removes a subscription
func (ls *LogStreamer) Unsubscribe(subID string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if sub, ok := ls.subscribers[subID]; ok {
		close(sub.done)
		close(sub.Channel)
		delete(ls.subscribers, subID)
	}
}

// Publish sends a log entry to all relevant subscribers
func (ls *LogStreamer) Publish(vmID string, entry *models.LogEntry) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	// Add to buffer
	ls.addToBuffer(vmID, entry)

	// Send to subscribers
	for _, sub := range ls.subscribers {
		if sub.VMID == vmID && ls.matchesFilter(entry, sub.Filter) {
			select {
			case sub.Channel <- entry:
			case <-sub.done:
			default:
				// Channel full, skip
			}
		}
	}
}

// GetRecentLogs returns recent logs for a VM
func (ls *LogStreamer) GetRecentLogs(vmID string, limit int) []*models.LogEntry {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	logs, ok := ls.buffers[vmID]
	if !ok {
		return nil
	}

	if limit <= 0 || limit > len(logs) {
		limit = len(logs)
	}

	start := len(logs) - limit
	result := make([]*models.LogEntry, limit)
	copy(result, logs[start:])
	return result
}

// ClearBuffer clears the log buffer for a VM
func (ls *LogStreamer) ClearBuffer(vmID string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	delete(ls.buffers, vmID)
}

// addToBuffer adds a log entry to the buffer
func (ls *LogStreamer) addToBuffer(vmID string, entry *models.LogEntry) {
	if _, ok := ls.buffers[vmID]; !ok {
		ls.buffers[vmID] = make([]*models.LogEntry, 0, ls.bufferSize)
	}

	ls.buffers[vmID] = append(ls.buffers[vmID], entry)

	// Trim if over buffer size
	if len(ls.buffers[vmID]) > ls.bufferSize {
		ls.buffers[vmID] = ls.buffers[vmID][1:]
	}
}

// matchesFilter checks if a log entry matches the filter
func (ls *LogStreamer) matchesFilter(entry *models.LogEntry, filter LogFilter) bool {
	// Check level filter
	if filter.Level != "" && entry.Level != filter.Level {
		levelPriority := map[string]int{
			"debug": 0,
			"info":  1,
			"warn":  2,
			"error": 3,
		}
		entryPriority, ok1 := levelPriority[entry.Level]
		filterPriority, ok2 := levelPriority[filter.Level]
		if ok1 && ok2 && entryPriority < filterPriority {
			return false
		}
	}

	// Check contains filter
	if filter.Contains != "" {
		// Simple substring check
		// In production, could use regex
		if len(entry.Message) > 0 && !containsString(entry.Message, filter.Contains) {
			return false
		}
	}

	// Check source filter
	if len(filter.Sources) > 0 {
		found := false
		for _, src := range filter.Sources {
			if entry.Source == src {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// containsString is a simple substring check
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// generateSubscriberID generates a unique subscriber ID
func generateSubscriberID() string {
	return time.Now().Format("20060102150405.000000")
}
