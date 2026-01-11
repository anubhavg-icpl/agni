// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package models

import (
	"time"
)

// VMStatus represents the current state of a VM
type VMStatus string

const (
	VMStatusStopped  VMStatus = "stopped"
	VMStatusStarting VMStatus = "starting"
	VMStatusRunning  VMStatus = "running"
	VMStatusStopping VMStatus = "stopping"
	VMStatusError    VMStatus = "error"
)

// VM represents a Firecracker microVM instance
type VM struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Status     VMStatus   `json:"status"`
	Config     VMConfig   `json:"config"`
	Metrics    *VMMetrics `json:"metrics,omitempty"`
	Error      string     `json:"error,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	StoppedAt  *time.Time `json:"stopped_at,omitempty"`
	PID        int        `json:"pid,omitempty"`
	SocketPath string     `json:"socket_path,omitempty"`
}

// VMConfig holds the configuration for a VM
type VMConfig struct {
	Name              string        `json:"name"`
	KernelPath        string        `json:"kernel_path"`
	KernelOpts        string        `json:"kernel_opts"`
	InitrdPath        string        `json:"initrd_path,omitempty"`
	RootDrive         Drive         `json:"root_drive"`
	AdditionalDrives  []Drive       `json:"additional_drives,omitempty"`
	CPUs              int64         `json:"cpus"`
	MemoryMB          int64         `json:"memory_mb"`
	CPUTemplate       string        `json:"cpu_template,omitempty"`
	DisableSMT        bool          `json:"disable_smt"`
	NetworkInterfaces []NIC         `json:"network_interfaces,omitempty"`
	VsockDevices      []Vsock       `json:"vsock_devices,omitempty"`
	Metadata          string        `json:"metadata,omitempty"`
	Jailer            *JailerConfig `json:"jailer,omitempty"`
	LogLevel          string        `json:"log_level"`
}

// Drive represents a block device
type Drive struct {
	ID       string `json:"id,omitempty"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"read_only"`
	IsRoot   bool   `json:"is_root,omitempty"`
	PartUUID string `json:"part_uuid,omitempty"`
}

// NIC represents a network interface configuration
type NIC struct {
	Device     string `json:"device"`
	MacAddress string `json:"mac_address"`
	AllowMMDS  bool   `json:"allow_mmds,omitempty"`
}

// Vsock represents a vsock device
type Vsock struct {
	Path string `json:"path"`
	CID  uint32 `json:"cid"`
}

// JailerConfig holds jailer configuration for privilege isolation
type JailerConfig struct {
	Binary        string `json:"binary"`
	ExecFile      string `json:"exec_file"`
	ID            string `json:"id"`
	UID           int    `json:"uid"`
	GID           int    `json:"gid"`
	NumaNode      int    `json:"numa_node"`
	ChrootBaseDir string `json:"chroot_base_dir"`
	Daemonize     bool   `json:"daemonize"`
}

// VMMetrics holds runtime metrics for a VM
type VMMetrics struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsed  int64     `json:"memory_used"`
	MemoryTotal int64     `json:"memory_total"`
	DiskRead    int64     `json:"disk_read_bytes"`
	DiskWrite   int64     `json:"disk_write_bytes"`
	NetRx       int64     `json:"net_rx_bytes"`
	NetTx       int64     `json:"net_tx_bytes"`
	Timestamp   time.Time `json:"timestamp"`
}

// ConfigTemplate represents a saved VM configuration template
type ConfigTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Config      VMConfig  `json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
