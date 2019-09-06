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

package roles

import (
	"context"
	"github.com/onosproject/onos-ztp/pkg/northbound"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"google.golang.org/grpc"
	"gotest.tools/assert"
	"io"
	"os"
	"sync"
	"testing"
)

var (
	spineRole = proto.DeviceRoleConfig{
		Role: "spine",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.06.1ea",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	leafRole = proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "complex"},
	}
)

// TestMain initializes the test suite context.
func TestMain(m *testing.M) {
	var waitGroup sync.WaitGroup

	spineRole.GetConfig().Properties = append(leafRole.GetConfig().Properties, &proto.DeviceProperty{
		Path:  "/spine/weak",
		Type:  "bool_val",
		Value: "false",
	})
	leafRole.GetConfig().Properties = append(leafRole.GetConfig().Properties, &proto.DeviceProperty{
		Path:  "/foo/bar",
		Type:  "string_val",
		Value: "totally fubar",
	})

	waitGroup.Add(1)
	northbound.SetUpServer(10123, Service{}, &waitGroup)
	waitGroup.Wait()
	os.Exit(m.Run())
}

func getClient() (*grpc.ClientConn, proto.DeviceRoleServiceClient) {
	conn := northbound.Connect(northbound.Address, northbound.Opts...)
	return conn, proto.NewDeviceRoleServiceClient(conn)
}

func Test_Set(t *testing.T) {
	conn, client := getClient()
	defer conn.Close()
	resp, err := client.Set(context.Background(), &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_ADD,
		Config: &leafRole,
	})
	assert.NilError(t, err, "unable to issue set add request")
	assert.Equal(t, resp.GetChange().Change, proto.DeviceRoleChange_ADDED, "incorrect change type")

	resp, err = client.Set(context.Background(), &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_UPDATE,
		Config: &leafRole,
	})
	assert.NilError(t, err, "unable to issue set update request")
	assert.Equal(t, resp.GetChange().Change, proto.DeviceRoleChange_UPDATED, "incorrect change type")

	resp, err = client.Set(context.Background(), &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_DELETE,
		Config: &leafRole,
	})
	assert.NilError(t, err, "unable to issue set delete request")
	assert.Equal(t, resp.GetChange().Change, proto.DeviceRoleChange_DELETED, "incorrect change type")
}

func Test_Get(t *testing.T) {
	conn, client := getClient()
	defer conn.Close()
	testSetRole(t, client, &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_ADD,
		Config: &spineRole,
	})
	testSetRole(t, client, &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_ADD,
		Config: &leafRole,
	})
	testGetRole(t, client, &proto.DeviceRoleRequest{Role: "leaf"}, 1)
	testGetRole(t, client, &proto.DeviceRoleRequest{}, 2)
}

func testSetRole(t *testing.T, client proto.DeviceRoleServiceClient, request *proto.DeviceRoleChangeRequest) {
	resp, err := client.Set(context.Background(), request)
	assert.NilError(t, err, "unable to issue set request")
	assert.Equal(t, resp.GetChange().Change, proto.DeviceRoleChange_ADDED, "incorrect change type")
}

func testGetRole(t *testing.T, client proto.DeviceRoleServiceClient, request *proto.DeviceRoleRequest, count int) {
	stream, err := client.Get(context.Background(), request)
	assert.NilError(t, err, "unable to issue get request")
	i := 0
	for {
		in, err := stream.Recv()
		if err == io.EOF || in == nil {
			break
		}
		assert.NilError(t, err, "unable to receive message")
		if len(request.GetRole()) > 0 {
			assert.Equal(t, in.GetRole(), request.GetRole(), "incorrect role name")
		}
		i++
	}
	err = stream.CloseSend()
	assert.NilError(t, err, "unable to close stream")
	assert.Equal(t, count, i, "wrong role count")
}

func Test_BadGet(t *testing.T) {
	conn, client := getClient()
	defer conn.Close()
	testGetRole(t, client, &proto.DeviceRoleRequest{Role: "none"}, 0)
}

func Test_BadSet(t *testing.T) {
	conn, client := getClient()
	defer conn.Close()

	_, err := client.Set(context.Background(), &proto.DeviceRoleChangeRequest{
		Change: proto.DeviceRoleChangeRequest_DELETE,
		Config: &proto.DeviceRoleConfig{Role: "none"},
	})
	assert.Assert(t, err != nil, "unable to issue set delete request")
}

//func Test_Subscribe(t *testing.T) {
//	conn, client := getClient()
//	defer conn.Close()
//	stream, err := client.Subscribe(context.Background(), &proto.DeviceRoleRequest{})
//	assert.NilError(t, err, "unable to issue subscribe request")
//
//	w := sync.WaitGroup{}
//	w.Add(1)
//	i := 0
//	go func() {
//		in, err := stream.Recv()
//		assert.Assert(t, err != io.EOF && in != nil, "expecting message; not EOF")
//		assert.NilError(t, err, "unable to receive message")
//		assert.Assert(t,
//			in.GetChange() == proto.DeviceRoleChange_ADDED || in.GetChange() == proto.DeviceRoleChange_UPDATED,
//			"incorrect change type")
//		i++
//		w.Done()
//	}()
//
//	go func() {
//		time.Sleep(100 * time.Millisecond)
//		conn1, client1 := getClient()
//		defer conn1.Close()
//		testSetRole(t, client1, &proto.DeviceRoleChangeRequest{
//			Change: proto.DeviceRoleChangeRequest_ADD,
//			Config: &spineRole,
//		})
//	}()
//
//	w.Wait()
//
//	err = stream.CloseSend()
//	assert.NilError(t, err, "unable to close stream")
//	assert.Equal(t, i, 1, "wrong role count")
//
//	// Allow time for cleanup
//	time.Sleep(100 * time.Millisecond)
//}
