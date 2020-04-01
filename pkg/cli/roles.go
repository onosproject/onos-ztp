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
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ztp/api/admin"
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
		Args:  cobra.MaximumNArgs(0),
		RunE:  runListRolesCommand,
	}
	return cmd
}

func runListRolesCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Get(ctx, &admin.DeviceRoleRequest{})
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Printf("%s\n", response.GetRole())
	}
	return nil
}

func getGetRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <roleName>",
		Short: "Get a device role",
		RunE:  runGetRoleCommand,
	}
	return cmd
}

func runGetRoleCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := admin.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stream, err := client.Get(ctx, &admin.DeviceRoleRequest{Role: args[0]})
	if err != nil {
		return err
	}

	response, err := stream.Recv()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	json, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		cli.Output("Unable to receive role")
		return err
	}

	cli.Output("%s", json)
	if err != nil {
		cli.Output("Unable to write")
		return err
	}
	return nil
}

func getAddRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <fileName>",
		Short: "Add a device role",
		RunE:  runAddRoleCommand,
	}
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite the role if it already exists")
	return cmd
}

func runAddRoleCommand(cmd *cobra.Command, args []string) error {
	return runAddOrUpdateRoleCommand(cmd, args, false)
}

func getUpdateRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <fileName>",
		Short: "Add a device role",
		RunE:  runUpdateRoleCommand,
	}
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite the role if it already exists")
	return cmd
}

func runUpdateRoleCommand(cmd *cobra.Command, args []string) error {
	return runAddOrUpdateRoleCommand(cmd, args, true)
}

func runAddOrUpdateRoleCommand(cmd *cobra.Command, args []string, overwrite bool) error {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := admin.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	roleFile := args[0]
	jsonFile, err := os.OpenFile(roleFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	jsonBlob, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	roleConfig := admin.DeviceRoleConfig{}
	err = json.Unmarshal(jsonBlob, &roleConfig)
	if err != nil {
		return err
	}

	change := admin.DeviceRoleChangeRequest_ADD
	if overwrite {
		change = admin.DeviceRoleChangeRequest_UPDATE
	}

	_, err = client.Set(ctx, &admin.DeviceRoleChangeRequest{
		Config: &roleConfig,
		Change: change,
	})
	return err
}

func getRemoveRoleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role <roleName>",
		Short: "Remove a device role",
		RunE:  runRemoveRolesCommand,
	}
	return cmd
}

func runRemoveRolesCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := admin.NewDeviceRoleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Set(ctx, &admin.DeviceRoleChangeRequest{
		Config: &admin.DeviceRoleConfig{Role: args[0]},
		Change: admin.DeviceRoleChangeRequest_DELETE,
	})
	return err
}
