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

package store

import (
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"gotest.tools/assert"
	"os"
	"testing"
)

// TestMain initializes the test suite context.
func TestMain(m *testing.M) {
	_ = os.RemoveAll("/tmp/roledb")
	os.Exit(m.Run())
}

func setupRepo(t *testing.T, path string) RoleStore {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Error("Unable to create db directory", err)
	}
	return RoleStore{Dir: path}
}

func Test_BadRepo(t *testing.T) {
	db := RoleStore{Dir: "/xney"}
	_, err := db.ReadRole("foo")
	assert.Assert(t, err != nil, "error expected")
	err = db.WriteRole(&proto.DeviceRoleConfig{}, true)
	assert.Assert(t, err != nil, "error expected")
}

func Test_EmptyRepo(t *testing.T) {
	db := setupRepo(t, "/tmp/roledb/1")
	roles, err := db.ListRoles()
	assert.NilError(t, err)
	assert.Assert(t, len(roles) == 0, "no roles expected")
}

func Test_Basics(t *testing.T) {
	db := setupRepo(t, "/tmp/roledb/2")
	role := proto.DeviceRoleConfig{
		Role: "leaf",
		Config: &proto.DeviceConfig{
			SoftwareVersion: "2019.08.02.c0ffee",
			Properties:      nil,
		},
		Pipeline: &proto.DevicePipeline{Pipeconf: "simple"},
	}
	role.GetConfig().Properties = append(role.GetConfig().Properties, &proto.DeviceProperty{
		Path:  "/foo/bar",
		Type:  "string_val",
		Value: "totally fubar",
	})
	err := db.WriteRole(&role, false)
	assert.NilError(t, err)

	roles, err := db.ListRoles()
	assert.NilError(t, err)
	assert.Assert(t, len(roles) == 1, "1 role expected")
	assert.Assert(t, roles[0] == "leaf")

	rr, err := db.ReadRole("leaf")
	assert.NilError(t, err)
	assert.Assert(t, rr.Role == "leaf", "got wrong role")
	assert.Assert(t, rr.Pipeline.Pipeconf == "simple", "got wrong pipeline")
	assert.Assert(t, len(rr.Config.Properties) == 1, "got wrong config")

	dr, err := db.DeleteRole("leaf")
	assert.NilError(t, err)
	assert.Assert(t, dr.Role == "leaf", "got wrong role")
	assert.Assert(t, dr.Pipeline.Pipeconf == "simple", "got wrong pipeline")

	noroles, err := db.ListRoles()
	assert.NilError(t, err)
	assert.Assert(t, len(noroles) == 0, "no roles expected")
}

func Test_ReadNonexistentRole(t *testing.T) {
	db := setupRepo(t, "/tmp/roledb/4")
	_, err := db.ReadRole("xney")
	assert.Assert(t, err != nil, "error expected")
}
