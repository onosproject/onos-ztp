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

package northbound

import (
	"fmt"
	"os"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-ztp/pkg/manager"
	"google.golang.org/grpc"
)

var (
	// Address is a test server address as "127.0.0.1:port" string
	Address string

	// Opts is a set of gRPC connection options
	Opts []grpc.DialOption

	// DbPath is the file system path to the database
	DbPath = "/tmp/rolesdb"
)

var log = logging.GetLogger("northbound")

// SetUpServer sets up a test manager and a gRPC end-point
// to which it registers the given service.
func SetUpServer(port int16, service northbound.Service, waitGroup *sync.WaitGroup) {
	_ = os.RemoveAll(DbPath)
	err := os.MkdirAll(DbPath, 0755)
	if err != nil {
		log.Error("Unable to create DB directory")
	}

	_, err = manager.LoadManager(DbPath)
	if err != nil {
		log.Error("Unable to load manager")
	}

	config := northbound.NewServerConfig("", "", "", port, false)
	s := northbound.NewServer(config)
	s.AddService(service)

	empty := ""
	Address = fmt.Sprintf(":%d", port)
	Opts, err = certs.HandleCertPaths(empty, empty, empty, true)
	if err != nil {
		log.Error("Error loading cert ", err)
	}
	go func() {
		err := s.Serve(func(started string) {
			waitGroup.Done()
			fmt.Printf("Started %v", started)
		})
		if err != nil {
			log.Error("Unable to serve ", err)
		}
	}()
}
