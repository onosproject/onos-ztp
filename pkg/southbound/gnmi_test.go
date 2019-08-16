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
	"gotest.tools/assert"
	"testing"
)

func Test_MakeRequest(t *testing.T) {
	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeline: "simple"},
	}
	role.GetConfig().Properties = append(role.GetConfig().Properties,
		&proto.DeviceProperty{
			Path:  "/foo/bar",
			Type:  "string_val",
			Value: "totally fubar",
		},
		&proto.DeviceProperty{
			Path:  "/foo/enabled",
			Type:  "bool_val",
			Value: "true",
		},
		&proto.DeviceProperty{
			Path:  "/foo/something",
			Type:  "int_val",
			Value: "123",
		},
	)

	d := device.Device{ID: "foo", Version: "leaf"}

	gnmi := GNMIProvisioner{}
	err := gnmi.Provision(&d, &role)
	assert.NilError(t, err, "unable to provision device")
}
