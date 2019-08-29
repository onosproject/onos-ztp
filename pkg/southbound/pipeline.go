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
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/northbound/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	log "k8s.io/klog"
	"net/http"
	"strings"
)

// PipelineProvisioner handles provisioning of device pipeline via onos-netcfg
type PipelineProvisioner struct {
}

const (
	template    = `{"device:%s": {"basic": {"managementAddress":"grpc://%s?device_id=1","driver":"%s","pipeconf":"%s","locType":"grid","gridX":%s,"gridY":%s}}}`
	onosAddress = "10.128.100.91:8181"
	// onosAddress = "10.1.10.19:8181"
)

// Init initializes the pipeline provisioner
func (p *PipelineProvisioner) Init(opts ...grpc.DialOption) error {
	return nil
}

// Provision runs the pipeline provisioning task
func (p *PipelineProvisioner) Provision(d *device.Device, cfg *proto.DeviceRoleConfig) error {
	attrs := d.GetAttributes()
	ctl := cfg.GetPipeline()
	cfgString := fmt.Sprintf(template, d.GetID(), d.GetAddress(), ctl.GetDriver(), ctl.GetPipeconf(), attrs["x"], attrs["y"])
	url := fmt.Sprintf("http://%s/onos/v1/network/configuration/devices", onosAddress)
	log.Infof("Applying pipeline configuration %s to device %s", cfgString, d.GetID())

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, strings.NewReader(cfgString))
	if err != nil {
		log.Error("ONOS config request failed due to: ", err)
		return err
	}
	request.Header.Add("Content-type", "application/json")
	request.SetBasicAuth("onos", "rocks")

	response, err := client.Do(request)
	if err != nil {
		log.Error("ONOS config request failed due to: ", err)
		return err
	}
	data, _ := ioutil.ReadAll(response.Body)
	log.Infof("Applied pipeline configuration to device %s; %v", d.GetID(), string(data))
	return nil
}
