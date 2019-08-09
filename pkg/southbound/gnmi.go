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
	"fmt"
	"github.com/onosproject/onos-config/pkg/certs"
	"github.com/onosproject/onos-config/pkg/utils"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"strings"
)

func MakeSetRequest(config *proto.DeviceRoleConfig) *gnmi.SetRequest {
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

func getConfig(key string) string {
	return viper.GetString(key)
}

func GetClient(address string) *gnmi.GNMIClient {
	keyPath := getConfig("keyPath")
	certPath := getConfig("certPath")
	opts, err := certs.HandleCertArgs(&keyPath, &certPath)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		fmt.Println("Can't connect", err)
	}
	client := gnmi.NewGNMIClient(conn)
	return &client
}
