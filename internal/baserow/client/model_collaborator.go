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

// checks if the Collaborator type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Collaborator{}

// Collaborator struct for Collaborator
type Collaborator struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
}

// NewCollaborator instantiates a new Collaborator object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCollaborator(id int32, name string) *Collaborator {
	this := Collaborator{}
	this.Id = id
	this.Name = name
	return &this
}

// NewCollaboratorWithDefaults instantiates a new Collaborator object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCollaboratorWithDefaults() *Collaborator {
	this := Collaborator{}
	return &this
}

// GetId returns the Id field value
func (o *Collaborator) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Collaborator) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Collaborator) SetId(v int32) {
	o.Id = v
}

// GetName returns the Name field value
func (o *Collaborator) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Collaborator) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Collaborator) SetName(v string) {
	o.Name = v
}

func (o Collaborator) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Collaborator) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	// skip: name is readOnly
	return toSerialize, nil
}

type NullableCollaborator struct {
	value *Collaborator
	isSet bool
}

func (v NullableCollaborator) Get() *Collaborator {
	return v.value
}

func (v *NullableCollaborator) Set(val *Collaborator) {
	v.value = val
	v.isSet = true
}

func (v NullableCollaborator) IsSet() bool {
	return v.isSet
}

func (v *NullableCollaborator) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCollaborator(val *Collaborator) *NullableCollaborator {
	return &NullableCollaborator{value: val, isSet: true}
}

func (v NullableCollaborator) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCollaborator) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


