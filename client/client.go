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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/leonwanghui/opensds-broker/model"
)

var Edp string
var httpClient = &http.Client{}

func ListProfiles() ([]*model.ProfileSpec, error) {
	url := Edp + "/api/v1alpha/profiles"

	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prfs = []*model.ProfileSpec{}
	if err = json.Unmarshal(rbody, &prfs); err != nil {
		return nil, err
	}
	return prfs, nil
}

func CreateVolume(planID, name, description string, size int64) (*model.VolumeSpec, error) {
	url := Edp + "/api/v1alpha/block/volumes"
	vr := &model.VolumeSpec{
		BaseModel:   &model.BaseModel{},
		Name:        name,
		Description: description,
		Size:        size,
		ProfileId:   planID,
	}

	vrJSON, err := json.Marshal(vr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(vrJSON))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vresp = &model.VolumeSpec{}
	if err = json.Unmarshal(rbody, vresp); err != nil {
		return nil, err
	}
	return vresp, nil
}

func ListVolumes() ([]*model.VolumeSpec, error) {
	url := Edp + "/api/v1alpha/block/volumes"

	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vols = []*model.VolumeSpec{}
	if err = json.Unmarshal(rbody, &vols); err != nil {
		return nil, err
	}
	return vols, nil
}

func DeleteVolume(volID string) (*model.Response, error) {
	url := Edp + "/api/v1alpha/block/volumes/" + volID
	vr := &model.VolumeSpec{}

	vrJSON, err := json.Marshal(vr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(vrJSON))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vresp = &model.Response{}
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
