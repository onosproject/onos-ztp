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
	return &proto.DeviceRoleChangeResponse{}, nil
}

// Get provides means to query device role configuration.
func (s Server) Get(req *proto.DeviceRoleRequest, stream proto.DeviceRoleService_GetServer) error {
	return nil
}

// Subscribe provides means to monitor changes in the device role configuration.
func (s Server) Subscribe(req *proto.DeviceRoleRequest, stream proto.DeviceRoleService_SubscribeServer) error {
	return nil
}
