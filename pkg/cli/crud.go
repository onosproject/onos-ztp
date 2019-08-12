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

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {roles,role} [args]",
		Short: "Get ZTP resources",
	}
	cmd.AddCommand(getGetRolesCommand())
	cmd.AddCommand(getGetRoleCommand())
	return cmd
}

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add {role} [args]",
		Short: "Add a ZTP resource",
	}
	cmd.AddCommand(getAddRoleCommand())
	return cmd
}

func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update {role} [args]",
		Short: "Update a ZTP resource",
	}
	cmd.AddCommand(getUpdateRoleCommand())
	return cmd
}

func getRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove {role} [args]",
		Short: "Remove a ZTP resource",
	}
	cmd.AddCommand(getRemoveRoleCommand())
	return cmd
}
