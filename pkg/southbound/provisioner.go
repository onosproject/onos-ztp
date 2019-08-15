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
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	log "k8s.io/klog"
)

const (
	configAddress = "onos-config:5150"
	//controlAddress = "onos-control:5150"
)

// DeviceProvisioner is responsible for provisioning devices with the role-specific configurations.
type DeviceProvisioner struct {
	gnmi gnmi.GNMIClient
}

func (p *DeviceProvisioner) Start(devices chan *device.Device, dialOptions ...grpc.DialOption) error {
	gnmiConn, err := grpc.Dial(configAddress, dialOptions...)
	if err != nil {
		log.Error("Unable to connect to onos-config", err)
		return err
	}

	//ctlConn, err := grpc.Dial(controlAddress, dialOptions...)
	//if err != nil {
	//	log.Error("Unable to connect to onos-control", err)
	//	return err
	//}

	p.gnmi = gnmi.NewGNMIClient(gnmiConn)
	go func() {
		for {
			d := <-devices
			p.provisionDevice(d)
		}
	}()
	return nil
}

func (p *DeviceProvisioner) provisionDevice(d *device.Device) {
	log.Infof("Provisioning device %s with config for role %s...", d.GetID(), d.GetVersion())
	err := p.provisionConfig(d)
	if err != nil {
		log.Errorf("Unable to provision device %s configuration due to %v", d.GetID(), err)
	} else {
		err = p.provisionControl(d)
		if err != nil {
			log.Errorf("Unable to provision device %s pipeline due to %v", d.GetID(), err)
		}
	}
}

func (p *DeviceProvisioner) provisionConfig(d *device.Device) error {
	return nil
}

func (p *DeviceProvisioner) provisionControl(d *device.Device) error {
	return nil
}
