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

package southbound

import (
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/store"
	log "k8s.io/klog"
)

// ProvisionerTask defines a contract of an activity that provisions an aspect of device operation.
type ProvisionerTask interface {
	// Provision sets up a device for an aspect of device operation.
	Provision(d *device.Device, cfg *proto.DeviceRoleConfig) error
}

// DeviceProvisioner is responsible for provisioning devices with the role-specific configurations.
type DeviceProvisioner struct {
	Tasks []ProvisionerTask
	Store *store.RoleStore
}

// Start starts the provisioner
func (p *DeviceProvisioner) Start(devices chan *device.Device) {
	go func() {
		log.Info("Ready to provision devices")
		for {
			d := <-devices
			if d != nil {
				cfg, err := p.Store.ReadRole(string(d.GetRole()))
				if err == nil {
					p.provisionDevice(d, cfg)
				} else {
					log.Errorf("Unable to find role %s", d.GetRole())
				}
			}
		}
	}()
}

func (p *DeviceProvisioner) provisionDevice(d *device.Device, cfg *proto.DeviceRoleConfig) {
	log.Infof("Provisioning device %s with config for role %s...", d.GetID(), cfg.GetRole())
	for _, t := range p.Tasks {
		err := t.Provision(d, cfg)
		if err != nil {
			log.Errorf("Unable to provision %s due to %v", d.GetID(), err)
		}
	}
}
