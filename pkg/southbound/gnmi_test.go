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
	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/southbound/mock"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi_ext"
	"gotest.tools/assert"
	"io"
	"testing"
)

func Test_Provision(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockGNMIClient(ctrl)
	gnmiTask := GNMIProvisioner{}
	gnmiTask.gnmi = client

	client.EXPECT().Set(gomock.Any(), gomock.Any()).
		Return(&gnmi.SetResponse{Extension: []*gnmi_ext.Extension{}}, nil)

	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	role.GetConfig().Properties = append(role.GetConfig().Properties,
		&proto.DeviceProperty{Path: "/foo/string", Type: "string_val", Value: "totally fubar"},
		&proto.DeviceProperty{Path: "/foo/bool", Type: "bool_val", Value: "true"},
	)

	d := device.Device{ID: "foo", Version: "leaf", Type: "bar"}
	err := gnmiTask.Provision(&d, &role)
	assert.NilError(t, err, "unable to provision device")
}

func Test_BadProvision(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockGNMIClient(ctrl)
	gnmiTask := GNMIProvisioner{}
	gnmiTask.gnmi = client

	client.EXPECT().Set(gomock.Any(), gomock.Any()).
		Return(nil, io.ErrClosedPipe)

	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}

	d := device.Device{ID: "foo", Version: "leaf", Type: "bar"}
	err := gnmiTask.Provision(&d, &role)
	assert.Error(t, err, "io: read/write on closed pipe")
}

func Test_Types(t *testing.T) {
	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	role.GetConfig().Properties = append(role.GetConfig().Properties,
		&proto.DeviceProperty{Path: "/foo/string", Type: "string_val", Value: "totally fubar"},
		&proto.DeviceProperty{Path: "/foo/bool", Type: "bool_val", Value: "true"},
		&proto.DeviceProperty{Path: "/foo/int", Type: "int_val", Value: "-123"},
		&proto.DeviceProperty{Path: "/foo/uint", Type: "uint_val", Value: "123"},
		&proto.DeviceProperty{Path: "/foo/float", Type: "float_val", Value: "123567890.655431"},
		&proto.DeviceProperty{Path: "/foo/huh", Type: "wut", Value: "123567890.655431"},
	)
	d := device.Device{ID: "foo", Version: "leaf", Type: "bar"}
	v := makeSetRequest(&d, &role)
	assert.Equal(t, len(v.Update), 6, "wrong number of properties")
}

func Test_BadTypes(t *testing.T) {
	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	role.GetConfig().Properties = append(role.GetConfig().Properties,
		&proto.DeviceProperty{Path: "/foo/bool", Type: "bool_val", Value: "x"},
		&proto.DeviceProperty{Path: "/foo/int", Type: "int_val", Value: "!123"},
		&proto.DeviceProperty{Path: "/foo/uint", Type: "uint_val", Value: "%123"},
		&proto.DeviceProperty{Path: "/foo/float", Type: "float_val", Value: "1|390.65x"},
	)
	d := device.Device{ID: "foo", Version: "leaf", Type: "bar"}
	v := makeSetRequest(&d, &role)
	assert.Equal(t, len(v.Update), 0, "no properties expected")
}
