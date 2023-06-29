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

// checks if the TrashStructure type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TrashStructure{}

// TrashStructure struct for TrashStructure
type TrashStructure struct {
	// An array of group trash structure records. Deprecated, please use `workspaces`.
	Groups []TrashStructureGroup `json:"groups"`
	// An array of workspace trash structure records
	Workspaces []TrashStructureGroup `json:"workspaces"`
}

// NewTrashStructure instantiates a new TrashStructure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTrashStructure(groups []TrashStructureGroup, workspaces []TrashStructureGroup) *TrashStructure {
	this := TrashStructure{}
	this.Groups = groups
	this.Workspaces = workspaces
	return &this
}

// NewTrashStructureWithDefaults instantiates a new TrashStructure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTrashStructureWithDefaults() *TrashStructure {
	this := TrashStructure{}
	return &this
}

// GetGroups returns the Groups field value
func (o *TrashStructure) GetGroups() []TrashStructureGroup {
	if o == nil {
		var ret []TrashStructureGroup
		return ret
	}

	return o.Groups
}

// GetGroupsOk returns a tuple with the Groups field value
// and a boolean to check if the value has been set.
func (o *TrashStructure) GetGroupsOk() ([]TrashStructureGroup, bool) {
	if o == nil {
		return nil, false
	}
	return o.Groups, true
}

// SetGroups sets field value
func (o *TrashStructure) SetGroups(v []TrashStructureGroup) {
	o.Groups = v
}

// GetWorkspaces returns the Workspaces field value
func (o *TrashStructure) GetWorkspaces() []TrashStructureGroup {
	if o == nil {
		var ret []TrashStructureGroup
		return ret
	}

	return o.Workspaces
}

// GetWorkspacesOk returns a tuple with the Workspaces field value
// and a boolean to check if the value has been set.
func (o *TrashStructure) GetWorkspacesOk() ([]TrashStructureGroup, bool) {
	if o == nil {
		return nil, false
	}
	return o.Workspaces, true
}

// SetWorkspaces sets field value
func (o *TrashStructure) SetWorkspaces(v []TrashStructureGroup) {
	o.Workspaces = v
}

func (o TrashStructure) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TrashStructure) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["groups"] = o.Groups
	toSerialize["workspaces"] = o.Workspaces
	return toSerialize, nil
}

type NullableTrashStructure struct {
	value *TrashStructure
	isSet bool
}

func (v NullableTrashStructure) Get() *TrashStructure {
	return v.value
}

func (v *NullableTrashStructure) Set(val *TrashStructure) {
	v.value = val
	v.isSet = true
}

func (v NullableTrashStructure) IsSet() bool {
	return v.isSet
}

func (v *NullableTrashStructure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTrashStructure(val *TrashStructure) *NullableTrashStructure {
	return &NullableTrashStructure{value: val, isSet: true}
}

func (v NullableTrashStructure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTrashStructure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

