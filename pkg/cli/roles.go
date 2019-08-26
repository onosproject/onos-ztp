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

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// DefaultRoleStorePath :
const DefaultRoleStorePath = "stores"

func getGetRolesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Get device roles",
		Run:   runListRolesCommand,
	}
	return cmd
}

func runListRolesCommand(cmd *cobra.Command, args []string) {
	conn := getConnection()
	defer conn.Close()
	client := proto.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Get(ctx, &proto.DeviceRoleRequest{})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			ExitWithError(ExitError, err)
		}
		fmt.Printf("%s\n", response.GetRole())
	}
}

func getGetRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <roleName>",
		Short: "Get a device role",
		Run:   runGetRoleCommand,
	}
	return cmd
}

func runGetRoleCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn := getConnection()
	defer conn.Close()
	client := proto.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Get(ctx, &proto.DeviceRoleRequest{Role: args[0]})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}

	response, err := stream.Recv()
	if err == io.EOF {
		return
	} else if err != nil {
		ExitWithError(ExitError, err)
	}

	json, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		ExitWithOutput("Unable to receive role: %v", err)
	}

	_, err = os.Stdout.Write(json)
	if err != nil {
		ExitWithOutput("Unable to write: %v", err)
	}
}

func getAddRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <fileName>",
		Short: "Add a device role",
		Run:   runAddRoleCommand,
	}
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite the role if it already exists")
	return cmd
}

func runAddRoleCommand(cmd *cobra.Command, args []string) {
	runAddOrUpdateRoleCommand(cmd, args, false)
}

func getUpdateRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <fileName>",
		Short: "Add a device role",
		Run:   runUpdateRoleCommand,
	}
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite the role if it already exists")
	return cmd
}

func runUpdateRoleCommand(cmd *cobra.Command, args []string) {
	runAddOrUpdateRoleCommand(cmd, args, true)
}

func runAddOrUpdateRoleCommand(cmd *cobra.Command, args []string, overwrite bool) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn := getConnection()
	defer conn.Close()
	client := proto.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	roleFile := args[0]
	jsonFile, err := os.OpenFile(roleFile, os.O_RDONLY, 0644)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	jsonBlob, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	roleConfig := proto.DeviceRoleConfig{}
	err = json.Unmarshal(jsonBlob, &roleConfig)
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	change := proto.DeviceRoleChangeRequest_ADD
	if overwrite {
		change = proto.DeviceRoleChangeRequest_UPDATE
	}

	_, err = client.Set(ctx, &proto.DeviceRoleChangeRequest{
		Config: &roleConfig,
		Change: change,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}
}

func getRemoveRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <roleName>",
		Short: "Remove a device role",
		Run:   runRemoveRolesCommand,
	}
	return cmd
}

func runRemoveRolesCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn := getConnection()
	defer conn.Close()
	client := proto.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.Set(ctx, &proto.DeviceRoleChangeRequest{
		Config: &proto.DeviceRoleConfig{Role: args[0]},
		Change: proto.DeviceRoleChangeRequest_DELETE,
	})
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}
}
