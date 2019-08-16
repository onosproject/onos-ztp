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
	"github.com/onosproject/onos-config/pkg/utils"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	log "k8s.io/klog"
	"strings"
)

const (
	configAddress = "onos-config:5150"
)

// GNMIProvisioner handles provisioning of device configuration via gNMI interface.
type GNMIProvisioner struct {
	gnmi gnmi.GNMIClient
}

// Init initializes the gNMI provisioner
func (p *GNMIProvisioner) Init(opts ...grpc.DialOption) error {
	gnmiConn, err := grpc.Dial(configAddress, opts...)
	if err != nil {
		log.Error("Unable to connect to onos-config", err)
		return err
	}
	p.gnmi = gnmi.NewGNMIClient(gnmiConn)
	return nil
}

// Provision runs the gNMI provisioning task
func (p *GNMIProvisioner) Provision(d *device.Device, cfg *proto.DeviceRoleConfig) error {
	// TODO: implement this fully
	_ = makeSetRequest(cfg)
	return nil
}

func makeSetRequest(config *proto.DeviceRoleConfig) *gnmi.SetRequest {
	updatedPaths := make([]*gnmi.Update, len(config.Config.Properties))
	for _, property := range config.Config.Properties {

		path, _ := utils.ParseGNMIElements([]string{property.Path})

		updatedPaths = append(updatedPaths,
			&gnmi.Update{
				Path: &gnmi.Path{
					Elem:   path.GetElem(),
					Target: path.GetTarget(),
				},
				Val: parseVal(*property),
			})
	}

	setRequest := &gnmi.SetRequest{
		Update: updatedPaths,
	}
	return setRequest
}

func parseVal(prop proto.DeviceProperty) *gnmi.TypedValue {
	switch prop.Type {
	//TODO: support other types of values
	case "string_val":
		return &gnmi.TypedValue{
			Value: &gnmi.TypedValue_StringVal{StringVal: prop.Value},
		}
	case "bool_val":
		return &gnmi.TypedValue{
			Value: &gnmi.TypedValue_BoolVal{BoolVal: strings.ToLower(prop.Value) == "true"},
		}
	default:
		return &gnmi.TypedValue{
			Value: &gnmi.TypedValue_StringVal{StringVal: ""},
		}
	}
}
