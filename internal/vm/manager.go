// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package vm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	fcmodels "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/firecracker-microvm/firectl/internal/logging"
	"github.com/firecracker-microvm/firectl/internal/storage"
	"github.com/firecracker-microvm/firectl/pkg/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	executableMask         = 0111
	firecrackerDefaultPath = "firecracker"
)

// RunningVM holds the state of a running VM
type RunningVM struct {
	Machine    *firecracker.Machine
	Cancel     context.CancelFunc
	SocketPath string
}

// Manager manages multiple Firecracker VMs
type Manager struct {
	store       *storage.VMStore
	runningVMs  map[string]*RunningVM
	mu          sync.RWMutex
	logger      *logging.Logger
	fcBinary    string
	logStreamer *LogStreamer
}

// NewManager creates a new VM Manager
func NewManager(store *storage.Store) *Manager {
	return &Manager{
		store:       storage.NewVMStore(store),
		runningVMs:  make(map[string]*RunningVM),
		logger:      logging.GetLogger().WithComponent("vm-manager"),
		logStreamer: NewLogStreamer(),
	}
}

// SetFirecrackerBinary sets the path to the firecracker binary
func (m *Manager) SetFirecrackerBinary(path string) {
	m.fcBinary = path
}

// Create creates a new VM configuration (does not start)
func (m *Manager) Create(config models.VMConfig) (*models.VM, error) {
	vm := &models.VM{
		ID:        uuid.New().String(),
		Name:      config.Name,
		Status:    models.VMStatusStopped,
		Config:    config,
		CreatedAt: time.Now(),
	}

	if err := m.store.Create(vm); err != nil {
		return nil, fmt.Errorf("failed to create VM: %w", err)
	}

	m.logger.Info().Str("vm_id", vm.ID).Str("name", vm.Name).Msg("VM created")
	return vm, nil
}

// Start starts a VM
func (m *Manager) Start(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already running
	if _, exists := m.runningVMs[id]; exists {
		return models.ErrVMAlreadyRunning
	}

	vm, err := m.store.Get(id)
	if err != nil {
		return err
	}

	// Update status to starting
	vm.Status = models.VMStatusStarting
	if err := m.store.Update(vm); err != nil {
		return err
	}

	// Build firecracker config
	fcConfig, err := m.buildFirecrackerConfig(vm)
	if err != nil {
		vm.Status = models.VMStatusError
		vm.Error = err.Error()
		m.store.Update(vm)
		return fmt.Errorf("failed to build config: %w", err)
	}

	// Get firecracker binary
	fcBinary, err := m.getFirecrackerBinary()
	if err != nil {
		vm.Status = models.VMStatusError
		vm.Error = err.Error()
		m.store.Update(vm)
		return err
	}

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Build command
	cmd := firecracker.VMCommandBuilder{}.
		WithBin(fcBinary).
		WithSocketPath(fcConfig.SocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	logger := log.New()
	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(logger)),
		firecracker.WithProcessRunner(cmd),
	}

	// Create machine
	machine, err := firecracker.NewMachine(ctx, fcConfig, machineOpts...)
	if err != nil {
		cancel()
		vm.Status = models.VMStatusError
		vm.Error = err.Error()
		m.store.Update(vm)
		return fmt.Errorf("failed to create machine: %w", err)
	}

	// Start machine
	if err := machine.Start(ctx); err != nil {
		cancel()
		vm.Status = models.VMStatusError
		vm.Error = err.Error()
		m.store.Update(vm)
		return fmt.Errorf("failed to start machine: %w", err)
	}

	// Update VM state
	now := time.Now()
	vm.Status = models.VMStatusRunning
	vm.StartedAt = &now
	vm.SocketPath = fcConfig.SocketPath
	vm.Error = ""
	if err := m.store.Update(vm); err != nil {
		m.logger.Error().Err(err).Msg("Failed to update VM state")
	}

	// Store running VM
	m.runningVMs[id] = &RunningVM{
		Machine:    machine,
		Cancel:     cancel,
		SocketPath: fcConfig.SocketPath,
	}

	m.logger.Info().Str("vm_id", id).Msg("VM started")

	// Start goroutine to wait for VM and handle cleanup
	go m.waitForVM(id, machine, ctx)

	return nil
}

// waitForVM waits for a VM to terminate and cleans up
func (m *Manager) waitForVM(id string, machine *firecracker.Machine, ctx context.Context) {
	err := machine.Wait(ctx)
	if err != nil {
		m.logger.Error().Err(err).Str("vm_id", id).Msg("VM wait error")
	}

	m.mu.Lock()
	delete(m.runningVMs, id)
	m.mu.Unlock()

	// Update status
	vm, err := m.store.Get(id)
	if err == nil {
		now := time.Now()
		vm.Status = models.VMStatusStopped
		vm.StoppedAt = &now
		m.store.Update(vm)
	}

	m.logger.Info().Str("vm_id", id).Msg("VM stopped")
}

// Stop force stops a VM
func (m *Manager) Stop(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	running, exists := m.runningVMs[id]
	if !exists {
		return models.ErrVMNotRunning
	}

	vm, err := m.store.Get(id)
	if err != nil {
		return err
	}

	vm.Status = models.VMStatusStopping
	m.store.Update(vm)

	if err := running.Machine.StopVMM(); err != nil {
		return fmt.Errorf("failed to stop VMM: %w", err)
	}

	running.Cancel()
	delete(m.runningVMs, id)

	now := time.Now()
	vm.Status = models.VMStatusStopped
	vm.StoppedAt = &now
	m.store.Update(vm)

	m.logger.Info().Str("vm_id", id).Msg("VM force stopped")
	return nil
}

// Shutdown gracefully shuts down a VM
func (m *Manager) Shutdown(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	running, exists := m.runningVMs[id]
	if !exists {
		return models.ErrVMNotRunning
	}

	vm, err := m.store.Get(id)
	if err != nil {
		return err
	}

	vm.Status = models.VMStatusStopping
	m.store.Update(vm)

	ctx := context.Background()
	if err := running.Machine.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown VM: %w", err)
	}

	m.logger.Info().Str("vm_id", id).Msg("VM shutdown requested")
	return nil
}

// Delete removes a VM (must be stopped first)
func (m *Manager) Delete(id string) error {
	m.mu.RLock()
	_, isRunning := m.runningVMs[id]
	m.mu.RUnlock()

	if isRunning {
		return fmt.Errorf("cannot delete running VM, stop it first")
	}

	if err := m.store.Delete(id); err != nil {
		return err
	}

	m.logger.Info().Str("vm_id", id).Msg("VM deleted")
	return nil
}

// Get retrieves a VM by ID
func (m *Manager) Get(id string) (*models.VM, error) {
	return m.store.Get(id)
}

// List returns all VMs
func (m *Manager) List() ([]*models.VM, error) {
	return m.store.List()
}

// GetMetrics retrieves metrics for a running VM
func (m *Manager) GetMetrics(id string) (*models.VMMetrics, error) {
	m.mu.RLock()
	_, exists := m.runningVMs[id]
	m.mu.RUnlock()

	if !exists {
		return nil, models.ErrVMNotRunning
	}

	// For now, return basic metrics
	// In a full implementation, we would query the Firecracker metrics FIFO
	return &models.VMMetrics{
		Timestamp: time.Now(),
	}, nil
}

// IsRunning checks if a VM is running
func (m *Manager) IsRunning(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.runningVMs[id]
	return exists
}

// GetLogStreamer returns the log streamer
func (m *Manager) GetLogStreamer() *LogStreamer {
	return m.logStreamer
}

// StopAll stops all running VMs
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, running := range m.runningVMs {
		m.logger.Info().Str("vm_id", id).Msg("Stopping VM during shutdown")
		running.Machine.StopVMM()
		running.Cancel()
	}
	m.runningVMs = make(map[string]*RunningVM)
}

// buildFirecrackerConfig converts a VM config to Firecracker config
func (m *Manager) buildFirecrackerConfig(vm *models.VM) (firecracker.Config, error) {
	cfg := vm.Config

	// Build drives
	drives := []fcmodels.Drive{
		{
			DriveID:      firecracker.String("1"),
			PathOnHost:   firecracker.String(cfg.RootDrive.Path),
			IsReadOnly:   firecracker.Bool(cfg.RootDrive.ReadOnly),
			IsRootDevice: firecracker.Bool(true),
			Partuuid:     cfg.RootDrive.PartUUID,
		},
	}

	for i, drive := range cfg.AdditionalDrives {
		drives = append(drives, fcmodels.Drive{
			DriveID:      firecracker.String(fmt.Sprintf("%d", i+2)),
			PathOnHost:   firecracker.String(drive.Path),
			IsReadOnly:   firecracker.Bool(drive.ReadOnly),
			IsRootDevice: firecracker.Bool(false),
		})
	}

	// Build network interfaces
	var nics []firecracker.NetworkInterface
	for _, nic := range cfg.NetworkInterfaces {
		nics = append(nics, firecracker.NetworkInterface{
			StaticConfiguration: &firecracker.StaticNetworkConfiguration{
				MacAddress:  nic.MacAddress,
				HostDevName: nic.Device,
			},
			AllowMMDS: nic.AllowMMDS,
		})
	}

	// Build vsock devices
	var vsocks []firecracker.VsockDevice
	for _, vs := range cfg.VsockDevices {
		vsocks = append(vsocks, firecracker.VsockDevice{
			Path: vs.Path,
			CID:  vs.CID,
		})
	}

	// Generate socket path
	socketPath := fmt.Sprintf("/tmp/firecracker-%s.sock", vm.ID)

	return firecracker.Config{
		SocketPath:        socketPath,
		KernelImagePath:   cfg.KernelPath,
		KernelArgs:        cfg.KernelOpts,
		InitrdPath:        cfg.InitrdPath,
		Drives:            drives,
		NetworkInterfaces: nics,
		VsockDevices:      vsocks,
		LogLevel:          cfg.LogLevel,
		MachineCfg: fcmodels.MachineConfiguration{
			VcpuCount:   firecracker.Int64(cfg.CPUs),
			CPUTemplate: fcmodels.CPUTemplate(cfg.CPUTemplate),
			Smt:         firecracker.Bool(!cfg.DisableSMT),
			MemSizeMib:  firecracker.Int64(cfg.MemoryMB),
		},
	}, nil
}

// getFirecrackerBinary finds the firecracker binary
func (m *Manager) getFirecrackerBinary() (string, error) {
	if m.fcBinary != "" {
		return m.validateBinary(m.fcBinary)
	}

	path, err := exec.LookPath(firecrackerDefaultPath)
	if err != nil {
		return "", fmt.Errorf("firecracker binary not found: %w", err)
	}

	return m.validateBinary(path)
}

// validateBinary validates that a binary exists and is executable
func (m *Manager) validateBinary(path string) (string, error) {
	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("binary %q does not exist", path)
	}
	if err != nil {
		return "", fmt.Errorf("failed to stat binary %q: %w", path, err)
	}
	if finfo.IsDir() {
		return "", fmt.Errorf("binary %q is a directory", path)
	}
	if finfo.Mode()&executableMask == 0 {
		return "", fmt.Errorf("binary %q is not executable", path)
	}
	return path, nil
}
