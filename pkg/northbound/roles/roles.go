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

// Package service implements the gRPC service for the zero-touch provisioning subsystem.
package roles

import (
	"context"
	"github.com/onosproject/onos-ztp/pkg/manager"
	"github.com/onosproject/onos-ztp/pkg/northbound"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"google.golang.org/grpc"
)

// Service is a Service implementation for administration.
type Service struct {
	northbound.Service
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	proto.RegisterDeviceRoleServiceServer(r, Server{})
}

// Server implements the gRPC service for zero-touch provisioning facilities.
type Server struct {
}

// Set provides means to add, update or delete device role configuration.
func (s Server) Set(ctx context.Context, r *proto.DeviceRoleChangeRequest) (*proto.DeviceRoleChangeResponse, error) {
	var err error
	var cfg = r.GetConfig()
	var changeType = proto.DeviceRoleChange_UPDATED

	switch r.GetChange() {
	case proto.DeviceRoleChangeRequest_ADD:
		changeType = proto.DeviceRoleChange_ADDED
		err = manager.GetManager().RoleStore.WriteRole(cfg, false)
	case proto.DeviceRoleChangeRequest_UPDATE:
		err = manager.GetManager().RoleStore.WriteRole(cfg, true)
	case proto.DeviceRoleChangeRequest_DELETE:
		changeType = proto.DeviceRoleChange_DELETED
		cfg, err = manager.GetManager().RoleStore.DeleteRole(cfg.Role)
	}

	if err != nil {
		return nil, err
	}

	change := proto.DeviceRoleChange{
		Change: changeType,
		Config: cfg,
	}

	// Queue up device role change for subscribers and return it to the caller
	//manager.GetManager().ChangesChannel <- change
	return &proto.DeviceRoleChangeResponse{
		Change: &change,
	}, nil
}

// Get provides means to query device role configuration.
func (s Server) Get(req *proto.DeviceRoleRequest, stream proto.DeviceRoleService_GetServer) error {
	roleName := req.GetRole()
	if len(roleName) > 0 {
		return s.sendRoleConfig(roleName, stream)
	}

	roles, err := manager.GetManager().RoleStore.ListRoles()
	for i := 0; err == nil && i < len(roles); i++ {
		err = s.sendRoleConfig(roles[i], stream)
	}
	return err
}

func (s Server) sendRoleConfig(roleName string, stream proto.DeviceRoleService_GetServer) error {
	role, err := manager.GetManager().RoleStore.ReadRole(roleName)
	if err == nil {
		err = stream.Send(role)
	}
	return err
}

// Subscribe provides means to monitor changes in the device role configuration.
func (s Server) Subscribe(req *proto.DeviceRoleRequest, stream proto.DeviceRoleService_SubscribeServer) error {
	// TODO: implement this
	return nil
}
