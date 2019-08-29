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
	"errors"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/store"
	"gotest.tools/assert"
	"os"
	"sync"
	"testing"
)

// TestMain initializes the test suite context.
func TestMain(m *testing.M) {
	_ = os.RemoveAll("/tmp/roledb")
	os.Exit(m.Run())
}

func setupRepo(t *testing.T, path string) store.RoleStore {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Error("Unable to create db directory", err)
	}
	store := store.RoleStore{Dir: path}
	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	err = store.WriteRole(&role, true)
	assert.NilError(t, err, "Unable to create test role")
	return store
}

type TestProvisioner struct {
	d    *device.Device
	wg   *sync.WaitGroup
	fail bool
}

func (t *TestProvisioner) Provision(d *device.Device, cfg *proto.DeviceRoleConfig) error {
	defer t.wg.Done()
	if t.fail {
		return errors.New("boom")
	}
	t.d = d
	return nil
}

func Test_NormalEvent(t *testing.T) {
	task := TestProvisioner{wg: &sync.WaitGroup{}, fail: false}
	task.wg.Add(1)
	repo := setupRepo(t, "/tmp/roledb/p1")
	p := DeviceProvisioner{Tasks: []ProvisionerTask{&task}, Store: &repo}
	devices := make(chan *device.Device)
	p.Start(devices)

	d := device.Device{ID: "foo", Role: "leaf"}
	devices <- &d
	task.wg.Wait()

	assert.Assert(t, task.d.ID == d.ID, "incorrect device")
	close(devices)
}

func Test_BadEvent(t *testing.T) {
	task := TestProvisioner{wg: &sync.WaitGroup{}, fail: true}
	task.wg.Add(1)
	repo := setupRepo(t, "/tmp/roledb/p2")
	p := DeviceProvisioner{Tasks: []ProvisionerTask{&task}, Store: &repo}
	devices := make(chan *device.Device)
	p.Start(devices)

	d := device.Device{ID: "foo", Role: "leaf"}
	devices <- &d
	task.wg.Wait()

	assert.Assert(t, task.d == nil, "device should be nil")
	close(devices)
}
