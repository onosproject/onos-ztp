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

package command

import (
	"encoding/json"
	"fmt"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"github.com/onosproject/onos-ztp/pkg/store"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

const DefaultRoleStorePath = "stores"

func getDB() store.RoleStore {
	return store.RoleStore{
		Dir: DefaultRoleStorePath,
	}
}
func newRolesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles <subcommand>",
		Short: "Read and write roles",
	}
	cmd.AddCommand(newRolesListCommand())
	cmd.AddCommand(newRolesGetCmd())
	cmd.AddCommand(newRolesSetCmd())
	cmd.AddCommand(newRolesRemoveCmd())

	return cmd
}

func newRolesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <roleName>",
		Short: "Remove a single role",
		Run:   runRemoveRolesCmd,
	}
	return cmd
}

func runRemoveRolesCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	db := getDB()
	_, err := db.DeleteRole(args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted role")
}

func newRolesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <roleName>",
		Short: "Get a single role",
		Run:   runGetRoleCmd,
	}
	return cmd
}

func runGetRoleCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	db := getDB()
	config, err := db.ReadRole(args[0])
	if err != nil {
		log.Fatal(err)
	}
	jsonConfig, marshalErr := json.MarshalIndent(config, "", "  ")

	if marshalErr != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonConfig))
}

func newRolesSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <fileName>",
		Short: "Set a single role",
		Run:   runSetRoleCmd,
	}
	cmd.PersistentFlags().Bool("overwrite", false, "Overwrite the role if it already exists")

	return cmd
}

func runSetRoleCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Not enough arguments")
	}

	db := getDB()
	jsonFile, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, readErr := ioutil.ReadAll(jsonFile)

	if readErr != nil {
		log.Fatal(readErr)
	}
	var config proto.DeviceRoleConfig

	marshalErr := json.Unmarshal(byteValue, &config)

	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	overwrite, err := cmd.PersistentFlags().GetBool("overwrite")
	if err != nil {
		log.Fatal(err)
	}
	writeErr := db.WriteRole(&config, overwrite)
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	fmt.Printf("Role successfully written to %s\n", args[0])
}

func newRolesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List roles",
		Run:   runListRolesCmd,
	}
	return cmd
}

func runListRolesCmd(cmd *cobra.Command, args []string) {
	db := getDB()

	roles, err := db.ListRoles()
	if err != nil {
		fmt.Printf("Could not get roles: %v", err)
		return
	}
	for idx, role := range roles {
		fmt.Printf("%d - %s\n", idx+1, role)
	}
}
