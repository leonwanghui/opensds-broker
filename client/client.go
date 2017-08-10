// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements a standard SouthBound interface of volume resource to
storage plugins.
*/

package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego/httplib"

	"github.com/leonwanghui/opensds-broker/api"
)

const (
	URL_PREFIX string = "http://100.64.128.40:50040"
)

func ListProfiles() (*[]api.ProfileSpec, error) {
	url := URL_PREFIX + "/api/v1alpha1/block/profiles"

	// fmt.Println("Start GET request to list profiles, url =", url)
	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prfs = &[]api.ProfileSpec{}
	if err = json.Unmarshal(rbody, prfs); err != nil {
		return nil, err
	}
	return prfs, nil
}

func CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	url := URL_PREFIX + "/api/v1alpha1/block/volumes"
	vr := api.NewVolumeRequest()
	vr.Spec.Name, vr.Spec.Size = name, size

	// fmt.Println("Start POST request to create volume, url =", url)
	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vresp = &api.VolumeSpec{}
	if err = json.Unmarshal(rbody, vresp); err != nil {
		return nil, err
	}
	return vresp, nil
}

func ListVolumes() (*[]api.VolumeSpec, error) {
	url := URL_PREFIX + "/api/v1alpha1/block/volumes"

	// fmt.Println("Start GET request to list volumes, url =", url)
	req := httplib.Get(url).SetTimeout(100*time.Second, 50*time.Second)

	resp, err := req.Response()
	if err != nil {
		return &[]api.VolumeSpec{}, err
	}
	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return &[]api.VolumeSpec{}, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &[]api.VolumeSpec{}, err
	}

	var vols = &[]api.VolumeSpec{}
	if err = json.Unmarshal(rbody, vols); err != nil {
		return &[]api.VolumeSpec{}, err
	}
	return vols, nil
}

func DeleteVolume(volID string) (*api.Response, error) {
	url := URL_PREFIX + "/api/v1alpha1/block/volumes/" + volID
	vr := api.NewVolumeRequest()

	// fmt.Println("Start DELETE request to delete volume, url =", url)
	req := httplib.Delete(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		return nil, err
	}
	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vresp = &api.Response{}
	err = json.Unmarshal(rbody, vresp)
	if err != nil {
		return nil, err
	}
	return vresp, nil
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204. Partial Content: 206.
// Otherwise return error message.
func CheckHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}
