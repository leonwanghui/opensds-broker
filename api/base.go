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
This module implements the common data structure.

*/

package api

type BaseModel struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createAt"`
	UpdatedAt string `json:"updateAt"`
}

func (base *BaseModel) GetId() string {
	return base.Id
}

func (base *BaseModel) GetCreatedTime() string {
	return base.CreatedAt
}

func (base *BaseModel) GetUpdatedTime() string {
	return base.UpdatedAt
}

func (base *BaseModel) SetCreatedTime(createdAt string) {
	base.CreatedAt = createdAt
}

func (base *BaseModel) SetUpdatedTime(updatedAt string) {
	base.UpdatedAt = updatedAt
}
