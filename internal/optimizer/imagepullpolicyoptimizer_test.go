/*
Copyright IBM Corporation 2020

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package optimize

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konveyor/move2kube/internal/types"
	plantypes "github.com/konveyor/move2kube/types/plan"
	log "github.com/sirupsen/logrus"
	core "k8s.io/kubernetes/pkg/apis/core"
)

func TestImagePullPolicyOptimizer(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("IR with no services", func(t *testing.T) {
		// Setup
		ir := getIRWithoutServices()
		imagePullPolicyOptimizer := imagePullPolicyOptimizer{}
		want := getIRWithoutServices()

		// Test
		actual, err := imagePullPolicyOptimizer.optimize(ir)
		if err != nil {
			t.Fatal("Failed to get the expected. Error:", err)
		}
		if !cmp.Equal(actual, want) {
			t.Fatalf("Failed to get the intermediate representation properly. Differences:\n%s", cmp.Diff(want, actual))
		}
	})

	t.Run("IR containing services that have no containers", func(t *testing.T) {
		// Setup
		ir := getIRWithServicesAndWithoutContainers()
		imagePullPolicyOptimizer := imagePullPolicyOptimizer{}
		want := getIRWithServicesAndWithoutContainers()

		// Test
		actual, err := imagePullPolicyOptimizer.optimize(ir)
		if err != nil {
			t.Fatal("Failed to get the expected. Error:", err)
		}
		if !cmp.Equal(actual, want) {
			t.Fatalf("Failed to get the intermediate representation properly. Differences:\n%s", cmp.Diff(want, actual))
		}
	})

	t.Run("IR containing services and containers without image pull policy", func(t *testing.T) {
		// Setup
		c1 := core.Container{
			Name: "container-1",
		}
		c2 := core.Container{
			Name: "container-2",
		}
		svcname1 := "svcname1"
		svcname2 := "svcname2"
		svc1 := types.Service{Name: svcname1, Replicas: 2}
		svc2 := types.Service{Name: svcname2, Replicas: 4}
		svc1.Containers = append(svc1.Containers, c1)
		svc2.Containers = append(svc2.Containers, c2)

		p := plantypes.NewPlan()
		ir := types.NewIR(p)
		ir.Services[svcname1] = svc1
		ir.Services[svcname2] = svc2

		imagePullPolicyOptimizer := imagePullPolicyOptimizer{}
		want := getIRWithImagePullPolicySetAsAlways()

		// Test
		actual, err := imagePullPolicyOptimizer.optimize(ir)
		if err != nil {
			t.Fatal("Failed to get the expected. Error:", err)
		}
		if !cmp.Equal(actual, want) {
			t.Fatalf("Failed to get the intermediate representation properly. Differences:\n%s", cmp.Diff(want, actual))
		}
	})

	t.Run("IR containing services and containers with image pull policy already set as Always", func(t *testing.T) {
		// Setup
		ir := getIRWithImagePullPolicySetAsAlways()
		imagePullPolicyOptimizer := imagePullPolicyOptimizer{}
		want := getIRWithImagePullPolicySetAsAlways()

		// Test
		actual, err := imagePullPolicyOptimizer.optimize(ir)
		if err != nil {
			t.Fatal("Failed to get the expected. Error:", err)
		}
		if !cmp.Equal(actual, want) {
			t.Fatalf("Failed to get the intermediate representation properly. Differences:\n%s", cmp.Diff(want, actual))
		}
	})
}

func getIRWithServicesAndWithoutContainers() types.IR {
	svcname1 := "svcname1"
	svcname2 := "svcname2"
	svc1 := types.Service{Name: svcname1, Replicas: 2}
	svc2 := types.Service{Name: svcname2, Replicas: 2}
	p := plantypes.NewPlan()
	ir := types.NewIR(p)
	ir.Services[svcname1] = svc1
	ir.Services[svcname2] = svc2
	return ir
}

func getIRWithoutServices() types.IR {
	p := plantypes.NewPlan()
	ir := types.NewIR(p)
	return ir
}

func getIRWithImagePullPolicySetAsAlways() types.IR {
	c1 := core.Container{
		Name: "container-1",
	}
	c2 := core.Container{
		Name: "container-2",
	}
	c1.ImagePullPolicy = core.PullAlways
	c2.ImagePullPolicy = core.PullAlways
	svcname1 := "svcname1"
	svcname2 := "svcname2"
	svc1 := types.Service{Name: svcname1, Replicas: 2}
	svc2 := types.Service{Name: svcname2, Replicas: 4}
	svc1.Containers = append(svc1.Containers, c1)
	svc2.Containers = append(svc2.Containers, c2)
	p := plantypes.NewPlan()
	ir := types.NewIR(p)
	ir.Services[svcname1] = svc1
	ir.Services[svcname2] = svc2
	return ir
}
