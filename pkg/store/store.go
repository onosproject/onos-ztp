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

// Package store houses implementation of the various device role configurations
package store

import (
	"encoding/json"
	"errors"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"io/ioutil"
	log "k8s.io/klog"
	"os"
	"path/filepath"
	"strings"
)

// RoleStore provides services to persist and retrieve role configuration records
type RoleStore struct {
	Dir string
}

func (s *RoleStore) path(roleName string) string {
	// TODO: sanitize roleName to make sure it's a valid file name
	return filepath.Join(s.Dir, roleName+".json")
}

// ListRoles returns an array of existing role name
func (s *RoleStore) ListRoles() ([]string, error) {
	names := make([]string, 0)
	err := filepath.Walk(s.Dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "json") {
			names = append(names, strings.TrimSuffix(filepath.Base(path), ".json"))
		}
		return nil
	})
	return names, err
}

// WriteRole stores the specified role configuration
func (s *RoleStore) WriteRole(roleConfig *proto.DeviceRoleConfig, overwrite bool) error {
	jsonBlob, err := json.Marshal(roleConfig)
	if err != nil {
		return err
	}

	if _, err := os.Stat(s.path(roleConfig.Role)); err == nil {
		if !overwrite {
			return errors.New("Overwrite was set to false but role" + s.path(roleConfig.Role) + " already exists")
		}
	}
	log.Infof("Writing record for role %s", roleConfig.GetRole())
	return ioutil.WriteFile(s.path(roleConfig.Role), jsonBlob, 0644)
}

// ReadRole reads the named role configuration
func (s *RoleStore) ReadRole(roleName string) (*proto.DeviceRoleConfig, error) {
	jsonFile, err := os.OpenFile(s.path(roleName), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	jsonBlob, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	roleConfig := proto.DeviceRoleConfig{}
	err = json.Unmarshal(jsonBlob, &roleConfig)
	if err != nil {
		return nil, err
	}
	return &roleConfig, nil
}

// DeleteRole removes the named role configuration
func (s *RoleStore) DeleteRole(roleName string) (*proto.DeviceRoleConfig, error) {
	role, err := s.ReadRole(roleName)
	if err != nil {
		return nil, err
	}
	log.Infof("Removing record for role %s", roleName)
	err = os.Remove(s.path(roleName))
	if err != nil {
		return nil, err
	}
	return role, nil
}
