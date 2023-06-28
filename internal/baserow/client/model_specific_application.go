/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the SpecificApplication type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SpecificApplication{}

// SpecificApplication struct for SpecificApplication
type SpecificApplication struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
	Order int32 `json:"order"`
	Type string `json:"type"`
	Group ApplicationGroup `json:"group"`
	Workspace ApplicationWorkspace `json:"workspace"`
}

// NewSpecificApplication instantiates a new SpecificApplication object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSpecificApplication(id int32, name string, order int32, type_ string, group ApplicationGroup, workspace ApplicationWorkspace) *SpecificApplication {
	this := SpecificApplication{}
	this.Id = id
	this.Name = name
	this.Order = order
	this.Type = type_
	this.Group = group
	this.Workspace = workspace
	return &this
}

// NewSpecificApplicationWithDefaults instantiates a new SpecificApplication object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSpecificApplicationWithDefaults() *SpecificApplication {
	this := SpecificApplication{}
	return &this
}

// GetId returns the Id field value
func (o *SpecificApplication) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *SpecificApplication) SetId(v int32) {
	o.Id = v
}

// GetName returns the Name field value
func (o *SpecificApplication) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *SpecificApplication) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *SpecificApplication) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *SpecificApplication) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *SpecificApplication) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *SpecificApplication) SetType(v string) {
	o.Type = v
}

// GetGroup returns the Group field value
func (o *SpecificApplication) GetGroup() ApplicationGroup {
	if o == nil {
		var ret ApplicationGroup
		return ret
	}

	return o.Group
}

// GetGroupOk returns a tuple with the Group field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetGroupOk() (*ApplicationGroup, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Group, true
}

// SetGroup sets field value
func (o *SpecificApplication) SetGroup(v ApplicationGroup) {
	o.Group = v
}

// GetWorkspace returns the Workspace field value
func (o *SpecificApplication) GetWorkspace() ApplicationWorkspace {
	if o == nil {
		var ret ApplicationWorkspace
		return ret
	}

	return o.Workspace
}

// GetWorkspaceOk returns a tuple with the Workspace field value
// and a boolean to check if the value has been set.
func (o *SpecificApplication) GetWorkspaceOk() (*ApplicationWorkspace, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Workspace, true
}

// SetWorkspace sets field value
func (o *SpecificApplication) SetWorkspace(v ApplicationWorkspace) {
	o.Workspace = v
}

func (o SpecificApplication) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SpecificApplication) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["name"] = o.Name
	toSerialize["order"] = o.Order
	// skip: type is readOnly
	toSerialize["group"] = o.Group
	toSerialize["workspace"] = o.Workspace
	return toSerialize, nil
}

type NullableSpecificApplication struct {
	value *SpecificApplication
	isSet bool
}

func (v NullableSpecificApplication) Get() *SpecificApplication {
	return v.value
}

func (v *NullableSpecificApplication) Set(val *SpecificApplication) {
	v.value = val
	v.isSet = true
}

func (v NullableSpecificApplication) IsSet() bool {
	return v.isSet
}

func (v *NullableSpecificApplication) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSpecificApplication(val *SpecificApplication) *NullableSpecificApplication {
	return &NullableSpecificApplication{value: val, isSet: true}
}

func (v NullableSpecificApplication) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSpecificApplication) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


