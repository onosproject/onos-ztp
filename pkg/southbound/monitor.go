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

// Package southbound is has facilities for monitoring topology and provisioning new devices.
package southbound

import (
	"context"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"google.golang.org/grpc"
	"io"
	log "k8s.io/klog"
)

// DeviceMonitor is responsible for monitoring topology for new device events.
type DeviceMonitor struct {
	client device.DeviceServiceClient
}

const (
	topoAddress = "onos-topo:5150"
)

func (m *DeviceMonitor) Init(dialOptions ...grpc.DialOption) error {
	conn, err := grpc.Dial(topoAddress, dialOptions...)
	if err != nil {
		log.Error("Unable to connect to topology server", err)
		return err
	}
	m.client = device.NewDeviceServiceClient(conn)
	return nil
}

// Start kicks off the device monitor listening for the topology device add events.
func (m *DeviceMonitor) Start(deviceEvents chan *device.Device) error {
	topoEvents, err := m.client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		return err
	}

	go func() {
		defer close(deviceEvents)
		for {
			event, err := topoEvents.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error("Unable to receive device event", err)
			} else if event.Type == device.ListResponse_ADDED {
				deviceEvents <- event.Device
			}
		}
	}()
	return nil
}
