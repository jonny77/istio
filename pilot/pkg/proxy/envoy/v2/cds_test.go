// Copyright 2018 Istio Authors
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
package v2_test

import (
	"context"
	"testing"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"google.golang.org/grpc"

	"istio.io/istio/tests/util"
)

func connectCDS(url string, t *testing.T) xdsapi.ClusterDiscoveryService_StreamClustersClient {
	conn, err := grpc.Dial(util.MockPilotGrpcAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatal("Connection failed", err)
	}

	xds := xdsapi.NewClusterDiscoveryServiceClient(conn)
	cdsstr, err := xds.StreamClusters(context.Background())
	if err != nil {
		t.Fatal("Rpc failed", err)
	}
	err = cdsstr.Send(&xdsapi.DiscoveryRequest{
		Node: &envoy_api_v2_core1.Node{
			Id: "sidecar~10.1.10.1~b~c",
		},
	})
	if err != nil {
		t.Fatal("Send failed", err)
	}
	return cdsstr
}

// Regression for envoy restart and overlapping connections
func TestCDS(t *testing.T) {
	initMocks()

	cdsr := connectCDS(util.MockPilotGrpcAddr, t)

	res, err := cdsr.Recv()
	if err != nil {
		t.Fatal("Failed to receive CDS", err)
		return
	}

	t.Log("CDS response", res)
	// TODO: uncomment once Shriram PR is in.
	//if len(res.Resources) == 0 {
	// 	t.Fatal("No response")
	//}

	// TODO: dump the response resources, compare with some golden once it's stable
	// check that each mocked service and destination rule has a corresponding resource

	// TODO: dynamic checks ( see EDS )
}
