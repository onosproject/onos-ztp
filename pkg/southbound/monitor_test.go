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
	"github.com/onosproject/onos-ztp/pkg/southbound/mock"
	"gotest.tools/assert"
	"io"
	"testing"
	"time"
)

func Test_Basics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockDeviceServiceClient(ctrl)
	stream := mock.NewMockDeviceService_ListClient(ctrl)

	m.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(stream, nil)
	stream.EXPECT().Recv().
		Return(&device.ListResponse{
			Type: device.ListResponse_ADDED,
			Device: &device.Device{
				ID: "foobar",
			},
		}, nil)
	stream.EXPECT().Recv().
		Return(&device.ListResponse{
			Type: device.ListResponse_UPDATED,
			Device: &device.Device{
				ID: "barfoo",
			},
		}, nil)
	stream.EXPECT().Recv().
		Return(nil, io.ErrClosedPipe)
	stream.EXPECT().Recv().
		Return(nil, io.EOF).
		AnyTimes()

	dispatchAddDelay = 1 * time.Microsecond
	dispatchUpdateDelay = 1 * time.Microsecond
	monitor := DeviceMonitor{m, nil}
	ch := make(chan *device.Device)
	go monitor.Start(ch)

	dev1 := <-ch
	dev2 := <-ch

	// TODO: This assertion was implemented due to the non-deterministic nature of device events with the delay hack.
	// Fix this test when the delay hack is replaced!
	assert.Assert(t, (dev1.GetID() == "foobar" && dev2.GetID() != "foobar") || (dev2.GetID() == "foobar" && dev1.GetID() != "foobar"), "incorrect device")
	assert.Assert(t, (dev1.GetID() == "barfoo" && dev2.GetID() != "barfoo") || (dev2.GetID() == "barfoo" && dev1.GetID() != "barfoo"), "incorrect device")

	time.Sleep(100 * time.Millisecond)
	monitor.Stop()
}

func Test_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockDeviceServiceClient(ctrl)

	m.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(nil, io.ErrUnexpectedEOF).
		AnyTimes()

	monitor := DeviceMonitor{m, nil}
	ch := make(chan *device.Device)
	go monitor.Start(ch)
	time.Sleep(100 * time.Millisecond)
	monitor.Stop()
}
