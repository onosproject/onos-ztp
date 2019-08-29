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
	"time"
)

// DeviceMonitor is responsible for monitoring topology for new device events.
type DeviceMonitor struct {
	client device.DeviceServiceClient
	events *chan *device.Device
}

const (
	topoAddress = "onos-topo:5150"
)

var (
	dispatchAddDelay    = 5 * time.Second
	dispatchUpdateDelay = 1 * time.Second
)

// Init initializes the connection to the topo server
func (m *DeviceMonitor) Init(dialOptions ...grpc.DialOption) error {
	conn, err := grpc.Dial(topoAddress, dialOptions...)
	if err != nil {
		log.Error("Unable to connect to topology server: ", err)
		return err
	}
	m.client = device.NewDeviceServiceClient(conn)
	return nil
}

// Start kicks off the device monitor listening for the topology device add events.
func (m *DeviceMonitor) Start(deviceEvents chan *device.Device) error {
	m.events = &deviceEvents
	topoEvents, err := m.client.List(context.Background(), &device.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		return err
	}

	go func() {
		log.Info("Listening for device events")
		for {
			event, err := topoEvents.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error("Unable to receive device event: ", err)
			} else if event.Type == device.ListResponse_ADDED || event.Type == device.ListResponse_UPDATED {
				log.Infof("Detected addition or update of device %s", event.Device.GetID())
				queueDevice(deviceEvents, event.Device, event.Type == device.ListResponse_UPDATED)
			}
		}
	}()
	return nil
}

// Stop stops the device monitor and associated resources
func (m *DeviceMonitor) Stop() {
	defer close(*m.events)
}

func queueDevice(devices chan *device.Device, d *device.Device, updated bool) {
	// HACK:  Induce delay before delivering the event onto the channel
	var t *time.Timer
	if updated {
		t = time.NewTimer(dispatchUpdateDelay)
	} else {
		t = time.NewTimer(dispatchAddDelay)
	}
	go func() {
		<-t.C
		log.Infof("Queueing event new device %s", d.GetID())
		devices <- d
	}()
}
