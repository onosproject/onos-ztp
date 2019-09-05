// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package manager is is the main coordinator for the ONOS control subsystem.
package manager

import (
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/southbound"
	"github.com/onosproject/onos-ztp/pkg/store"
	"google.golang.org/grpc"
	log "k8s.io/klog"
	"os"
)

var mgr Manager

// Manager single point of entry for the zero touch provisioning system.
type Manager struct {
	RoleStore      store.RoleStore
	ChangesChannel chan proto.DeviceRoleChange

	deviceChanel chan *device.Device
	connOptions  []grpc.DialOption
	monitor      southbound.DeviceMonitor
	provisioner  southbound.DeviceProvisioner
}

// NewManager initializes the provisioning manager subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")
	mgr = Manager{
		RoleStore:      store.RoleStore{Dir: "roledb"},
		ChangesChannel: make(chan proto.DeviceRoleChange, 10),
		monitor:        southbound.DeviceMonitor{},
		provisioner:    southbound.DeviceProvisioner{},
	}
	return &mgr, nil
}

// LoadManager creates a provisioning subsystem manager primed with stores loaded from the specified files.
func LoadManager(roleStorePath string, opts ...grpc.DialOption) (*Manager, error) {
	err := os.MkdirAll(roleStorePath, 0755)
	if err != nil {
		log.Errorf("Unable to create role store directory %s due to %v", roleStorePath, err)
		return nil, err
	}

	mgr, err := NewManager()
	if err != nil {
		return nil, err
	}

	mgr.deviceChanel = make(chan *device.Device)
	mgr.RoleStore.Dir = roleStorePath
	mgr.connOptions = opts
	mgr.provisioner.Store = &mgr.RoleStore

	err = mgr.monitor.Init(opts...)
	if err != nil {
		log.Error("Unable to setup topology monitor", err)
		return nil, err
	}

	gnmiTask := southbound.GNMIProvisioner{}
	err = gnmiTask.Init(opts...)
	if err != nil {
		log.Error("Unable to setup GNMI provisioner", err)
		return nil, err
	}

	// TODO: replace with p4Task
	pipelineTask := southbound.PipelineProvisioner{}
	err = pipelineTask.Init(opts...)
	if err != nil {
		log.Error("Unable to setup pipeline provisioner", err)
		return nil, err
	}

	mgr.provisioner.Tasks = []southbound.ProvisionerTask{&gnmiTask, &pipelineTask}
	return mgr, err
}

// Run starts any background tasks associated with the manager.
func (m *Manager) Run() {
	log.Info("Starting Manager")

	// Start the device monitor and provisioner components.
	go m.monitor.Start(m.deviceChanel)
	m.provisioner.Start(m.deviceChanel)
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}
