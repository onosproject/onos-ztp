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

// Package command holds ONOS command-line command implementations.
package command

import (
	"fmt"
	"github.com/onosproject/onos-config/pkg/certs"
	"github.com/onosproject/onos-config/pkg/northbound"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
)

// GetRootCommand returns the root CLI command.
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onos",
		Short: "ONOS command line client",
	}

	viper.SetDefault("address", ":5150")
	viper.SetDefault("keyPath", certs.Client1Key)
	viper.SetDefault("certPath", certs.Client1Crt)

	cmd.PersistentFlags().StringP("address", "a", viper.GetString("address"), "the controller address")
	cmd.PersistentFlags().StringP("keyPath", "k", viper.GetString("keyPath"), "path to client private key")
	cmd.PersistentFlags().StringP("certPath", "c", viper.GetString("certPath"), "path to client certificate")
	cmd.PersistentFlags().String("config", "", "config file (default: $HOME/.onos/config.yaml)")

	_ = viper.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag("keyPath", cmd.PersistentFlags().Lookup("keyPath"))
	_ = viper.BindPFlag("certPath", cmd.PersistentFlags().Lookup("certPath"))

	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newCompletionCommand())
	cmd.AddCommand(newRolesCommand())

	// TOTO: add the commands for the following usage
	// roles list				// dumps out list of role names
	// roles set "roleName" "jsonConfigFile"
	// roles get "roleName"   	// dumps out JSON config
	// roles remove "roleName"

	// TODO: remove
	//test()
	getConnection(nil)
	return cmd
}

func test(){
	file := store.RoleStore{Dir: "stores"}
	err := file.WriteRole(&proto.DeviceRoleConfig{
		Role:                 "testRole2",
		Config:               &proto.DeviceConfig{
			SoftwareVersion:      "234234",
			Properties:           []*proto.DeviceProperty{
				{
					Path:                 "asfasfsa",
					Type:                 "asdfsaf",
					Value:                "sadfasdfsfsdf",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		Pipeline:             &proto.DevicePipeline{
			Pipeline:             "test",
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	},true)
	if err != nil {
		log.Fatal(err)
	}
	roles, err := file.ListRoles()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(roles)
}

func getConnection(cmd *cobra.Command) *grpc.ClientConn {
	keyPath := getConfig("keyPath")
	certPath := getConfig("certPath")
	address := getConfig("address")
	opts, err := certs.HandleCertArgs(&keyPath, &certPath)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return northbound.Connect(address, opts...)
}
