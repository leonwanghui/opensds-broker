/*
Copyright 2016 The Kubernetes Authors.
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

package controller

import (
	"fmt"
	"log"
	"sync"

	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/brokerapi"
	"github.com/leonwanghui/opensds-broker/client"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type openSDSServiceInstance struct {
	Name       string
	Credential *brokerapi.Credential
}

type openSDSController struct {
	rwMutex     sync.RWMutex
	instanceMap map[string]*openSDSServiceInstance
}

// CreateController creates an instance of an OpenSDS service broker controller.
func CreateController() controller.Controller {
	var instanceMap = make(map[string]*openSDSServiceInstance)
	return &openSDSController{
		instanceMap: instanceMap,
	}
}

func (c *openSDSController) Catalog() (*brokerapi.Catalog, error) {
	prfs, err := client.ListProfiles()
	if err != nil {
		return nil, err
	}

	var plans = []brokerapi.ServicePlan{}
	for _, prf := range prfs {
		plan := brokerapi.ServicePlan{
			Name:        prf.GetName(),
			ID:          prf.GetId(),
			Description: prf.GetDescription(),
			Metadata:    prf.Extra,
			Free:        true,
		}
		plans = append(plans, plan)
	}

	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				Name:        "opensds-service",
				ID:          "4f6e6cf6-ffdd-425f-a2c7-3c9258ad2468",
				Description: "Policy based storage service",
				Plans:       plans,
				Bindable:    true,
			},
		},
	}, nil
}

func (c *openSDSController) GetServiceInstanceLastOperation(
	instanceID, serviceID, planID, operation string,
) (*brokerapi.LastOperationResponse, error) {
	return nil, nil
}

func (c *openSDSController) CreateServiceInstance(
	instanceID string,
	req *brokerapi.CreateServiceInstanceRequest,
) (*brokerapi.CreateServiceInstanceResponse, error) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	var name, description string
	var capacity int64
	if nameInterface, ok := req.Parameters["name"]; ok {
		name = nameInterface.(string)
	}
	if despInterface, ok := req.Parameters["description"]; ok {
		description = despInterface.(string)
	}
	if capInterface, ok := req.Parameters["capacity"]; ok {
		capacity = int64(capInterface.(float64))
	}

	vol, err := client.CreateVolume(req.PlanID, name, description, capacity)
	if err != nil {
		return nil, err
	}
	c.instanceMap[instanceID] = &openSDSServiceInstance{
		Name: instanceID,
		Credential: &brokerapi.Credential{
			"volumeId": vol.GetId(),
			"pool":     "rbd",
			"image":    "OPENSDS:" + vol.GetName() + ":" + vol.GetId(),
		},
	}

	log.Printf("Created User Provided Service Instance:\n%v\n", c.instanceMap[instanceID])
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

/*
func (c *openSDSController) GetServiceInstance(id string) (string, error) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	instance, ok := c.instanceMap[id]
	if ok {
		body, err := json.Marshal(instance)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("Can not find instance id %v", id)
}
*/

func (c *openSDSController) RemoveServiceInstance(instanceID, serviceID, planID string, acceptsIncomplete bool) (*brokerapi.DeleteServiceInstanceResponse, error) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	instance, ok := c.instanceMap[instanceID]
	if ok {
		volInterface, ok := (*instance.Credential)["volumeId"]
		if !ok {
			return &brokerapi.DeleteServiceInstanceResponse{}, fmt.Errorf("Volume id not provided in credential info!")
		}
		volID := volInterface.(string)
		resp, err := client.DeleteVolume(volID)
		if err != nil {
			return nil, err
		} else if resp.Status != "Success" {
			return nil, fmt.Errorf(resp.Error)
		}
		delete(c.instanceMap, instanceID)
	}

	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

func (c *openSDSController) Bind(
	instanceID,
	bindingID string,
	req *brokerapi.BindingRequest,
) (*brokerapi.CreateServiceBindingResponse, error) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	instance, ok := c.instanceMap[instanceID]
	if !ok {
		return nil, errNoSuchInstance{instanceID: instanceID}
	}
	cred := instance.Credential
	return &brokerapi.CreateServiceBindingResponse{Credentials: *cred}, nil
}

func (c *openSDSController) UnBind(instanceID, bindingID, serviceID, planID string) error {
	// Since we don't persist the binding, there's nothing to do here.
	return nil
}
