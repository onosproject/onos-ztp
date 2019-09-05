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
	"context"
	"fmt"
	ext "github.com/onosproject/onos-config/pkg/northbound/gnmi"
	"github.com/onosproject/onos-config/pkg/utils"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi_ext"
	"google.golang.org/grpc"
	log "k8s.io/klog"
	"strconv"
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
	request := makeSetRequest(d, cfg)
	log.Errorf("Applying device configuration %v to %s", request, d.GetID())

	response, err := p.gnmi.Set(context.Background(), request)
	if err != nil {
		log.Errorf("Unable to apply device configuration %v to %s", request, d.GetID())
		return err
	}
	log.Infof("Applied gNMI configuration to device %s; %v", d.GetID(), response)
	return nil
}

func makeSetRequest(d *device.Device, config *proto.DeviceRoleConfig) *gnmi.SetRequest {
	updatedPaths := make([]*gnmi.Update, 0)
	for _, prop := range config.Config.Properties {
		value, err := parseVal(*prop)
		if err != nil {
			log.Errorf("Unable to convert property %s to type %s using value [%s]", prop.Path, prop.Type, prop.Value)
		} else {
			path, _ := utils.ParseGNMIElements(utils.SplitPath(prop.Path))
			target := string(d.GetID())
			if len(d.GetTarget()) > 0 {
				target = d.GetTarget()
			}
			updatedPaths = append(updatedPaths,
				&gnmi.Update{
					Path: &gnmi.Path{Elem: path.GetElem(), Target: target},
					Val:  value,
				})
		}
	}

	changeID := fmt.Sprintf("%s-%s", d.GetID(), d.GetRole())
	ext100ChangeID := gnmi_ext.Extension_RegisteredExt{
		RegisteredExt: &gnmi_ext.RegisteredExtension{
			Id:  ext.GnmiExtensionNetwkChangeID,
			Msg: []byte(changeID),
		},
	}
	ext101Version := gnmi_ext.Extension_RegisteredExt{
		RegisteredExt: &gnmi_ext.RegisteredExtension{
			Id:  ext.GnmiExtensionVersion,
			Msg: []byte(d.GetVersion()),
		},
	}
	ext102Type := gnmi_ext.Extension_RegisteredExt{
		RegisteredExt: &gnmi_ext.RegisteredExtension{
			Id:  ext.GnmiExtensionDeviceType,
			Msg: []byte(d.GetType()),
		},
	}

	extensions := []*gnmi_ext.Extension{{Ext: &ext100ChangeID}, {Ext: &ext101Version}, {Ext: &ext102Type}}

	setRequest := &gnmi.SetRequest{
		Update:    updatedPaths,
		Extension: extensions,
	}
	return setRequest
}

func parseVal(prop proto.DeviceProperty) (*gnmi.TypedValue, error) {
	err := strconv.ErrSyntax
	switch prop.Type {
	case "string_val":
		return &gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: prop.Value}}, nil
	case "bool_val":
		b, err := strconv.ParseBool(prop.Value)
		if err == nil {
			return &gnmi.TypedValue{Value: &gnmi.TypedValue_BoolVal{BoolVal: b}}, nil
		}
	case "int_val":
		i, err := strconv.ParseInt(prop.Value, 10, 64)
		if err == nil {
			return &gnmi.TypedValue{Value: &gnmi.TypedValue_IntVal{IntVal: i}}, nil
		}
	case "uint_val":
		i, err := strconv.ParseUint(prop.Value, 10, 64)
		if err == nil {
			return &gnmi.TypedValue{Value: &gnmi.TypedValue_UintVal{UintVal: i}}, nil
		}
	case "float_val":
		f, err := strconv.ParseFloat(prop.Value, 32)
		if err == nil {
			return &gnmi.TypedValue{Value: &gnmi.TypedValue_FloatVal{FloatVal: float32(f)}}, nil
		}
	// TODO: add case "decimal_val":
	default:
		return &gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: prop.Value}}, nil
	}
	return nil, err
}
