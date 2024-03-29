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

// checks if the TrashStructureGroup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TrashStructureGroup{}

// TrashStructureGroup struct for TrashStructureGroup
type TrashStructureGroup struct {
	Id int32 `json:"id"`
	Trashed bool `json:"trashed"`
	Name string `json:"name"`
	Applications []TrashStructureApplication `json:"applications"`
}

// NewTrashStructureGroup instantiates a new TrashStructureGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTrashStructureGroup(id int32, trashed bool, name string, applications []TrashStructureApplication) *TrashStructureGroup {
	this := TrashStructureGroup{}
	this.Id = id
	this.Trashed = trashed
	this.Name = name
	this.Applications = applications
	return &this
}

// NewTrashStructureGroupWithDefaults instantiates a new TrashStructureGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTrashStructureGroupWithDefaults() *TrashStructureGroup {
	this := TrashStructureGroup{}
	return &this
}

// GetId returns the Id field value
func (o *TrashStructureGroup) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *TrashStructureGroup) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *TrashStructureGroup) SetId(v int32) {
	o.Id = v
}

// GetTrashed returns the Trashed field value
func (o *TrashStructureGroup) GetTrashed() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Trashed
}

// GetTrashedOk returns a tuple with the Trashed field value
// and a boolean to check if the value has been set.
func (o *TrashStructureGroup) GetTrashedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Trashed, true
}

// SetTrashed sets field value
func (o *TrashStructureGroup) SetTrashed(v bool) {
	o.Trashed = v
}

// GetName returns the Name field value
func (o *TrashStructureGroup) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *TrashStructureGroup) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *TrashStructureGroup) SetName(v string) {
	o.Name = v
}

// GetApplications returns the Applications field value
func (o *TrashStructureGroup) GetApplications() []TrashStructureApplication {
	if o == nil {
		var ret []TrashStructureApplication
		return ret
	}

	return o.Applications
}

// GetApplicationsOk returns a tuple with the Applications field value
// and a boolean to check if the value has been set.
func (o *TrashStructureGroup) GetApplicationsOk() ([]TrashStructureApplication, bool) {
	if o == nil {
		return nil, false
	}
	return o.Applications, true
}

// SetApplications sets field value
func (o *TrashStructureGroup) SetApplications(v []TrashStructureApplication) {
	o.Applications = v
}

func (o TrashStructureGroup) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TrashStructureGroup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["trashed"] = o.Trashed
	toSerialize["name"] = o.Name
	toSerialize["applications"] = o.Applications
	return toSerialize, nil
}

type NullableTrashStructureGroup struct {
	value *TrashStructureGroup
	isSet bool
}

func (v NullableTrashStructureGroup) Get() *TrashStructureGroup {
	return v.value
}

func (v *NullableTrashStructureGroup) Set(val *TrashStructureGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableTrashStructureGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableTrashStructureGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTrashStructureGroup(val *TrashStructureGroup) *NullableTrashStructureGroup {
	return &NullableTrashStructureGroup{value: val, isSet: true}
}

func (v NullableTrashStructureGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTrashStructureGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


