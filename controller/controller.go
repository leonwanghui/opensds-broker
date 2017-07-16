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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
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
		return &brokerapi.Catalog{}, err
	}

	var plans = []brokerapi.ServicePlan{}
	for _, prf := range *prfs {
		plan := brokerapi.ServicePlan{
			Name:        prf.Name,
			ID:          prf.Id,
			Description: prf.Description,
			Metadata:    prf.StorageTags,
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
	/*
	           return &brokerapi.Catalog{
	                   Services: []*brokerapi.Service{
	                           {
	                                   Name:        "opensds-service",
	                                   ID:          "4f6e6cf6-ffdd-425f-a2c7-3c9258ad2468",
	                                   Description: "Policy based storage service",
	   				Plans: []brokerapi.ServicePlan{
	   					{
	   						Name:        "default",
	   						ID:          "4f6e6cf6-ffdd-425f-0000-3c9258ad2468",
	   						Description: "",
	   						Metadata:    map[string]string{},
	   						Free:        true,
	   					},
	   				},
	                                   Bindable:    true,
	                           },
	                   },
	           }, nil
	*/
}

func (c *openSDSController) CreateServiceInstance(
	id string,
	req *brokerapi.CreateServiceInstanceRequest,
) (*brokerapi.CreateServiceInstanceResponse, error) {
	capInterface, ok := req.Parameters["capacity"]
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	if ok {
		capacity := capInterface.(float64)
		vol, err := client.CreateVolume(id, int32(capacity))
		if err != nil {
			return &brokerapi.CreateServiceInstanceResponse{}, err
		}
		c.instanceMap[id] = &openSDSServiceInstance{
			Name: id,
			Credential: &brokerapi.Credential{
				"volumeId": vol.Id,
				"pool":     "rbd",
				"image":    "OPENSDS:" + vol.Name + ":" + vol.Id,
			},
		}
	} else {
		return &brokerapi.CreateServiceInstanceResponse{}, fmt.Errorf("Capacity not provided in request!")
	}

	glog.Infof("Created User Provided Service Instance:\n%v\n", c.instanceMap[id])
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

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

func (c *openSDSController) RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	instance, ok := c.instanceMap[id]
	if ok {
		volInterface, ok := (*instance.Credential)["volumeId"]
		if !ok {
			return &brokerapi.DeleteServiceInstanceResponse{}, fmt.Errorf("Volume id not provided in credential info!")
		}
		volID := volInterface.(string)
		resp, err := client.DeleteVolume(volID)
		if err != nil {
			return &brokerapi.DeleteServiceInstanceResponse{}, err
		} else if resp.Status != "Success" {
			return &brokerapi.DeleteServiceInstanceResponse{}, fmt.Errorf(resp.Error)
		}
		delete(c.instanceMap, id)
		return &brokerapi.DeleteServiceInstanceResponse{}, nil
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

func (c *openSDSController) UnBind(instanceID string, bindingID string) error {
	// Since we don't persist the binding, there's nothing to do here.
	return nil
}
