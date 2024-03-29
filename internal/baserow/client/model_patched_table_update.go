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

// checks if the PatchedTableUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedTableUpdate{}

// PatchedTableUpdate struct for PatchedTableUpdate
type PatchedTableUpdate struct {
	Name *string `json:"name,omitempty"`
}

// NewPatchedTableUpdate instantiates a new PatchedTableUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedTableUpdate() *PatchedTableUpdate {
	this := PatchedTableUpdate{}
	return &this
}

// NewPatchedTableUpdateWithDefaults instantiates a new PatchedTableUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedTableUpdateWithDefaults() *PatchedTableUpdate {
	this := PatchedTableUpdate{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedTableUpdate) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableUpdate) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedTableUpdate) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedTableUpdate) SetName(v string) {
	o.Name = &v
}

func (o PatchedTableUpdate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedTableUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	return toSerialize, nil
}

type NullablePatchedTableUpdate struct {
	value *PatchedTableUpdate
	isSet bool
}

func (v NullablePatchedTableUpdate) Get() *PatchedTableUpdate {
	return v.value
}

func (v *NullablePatchedTableUpdate) Set(val *PatchedTableUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedTableUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedTableUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedTableUpdate(val *PatchedTableUpdate) *NullablePatchedTableUpdate {
	return &NullablePatchedTableUpdate{value: val, isSet: true}
}

func (v NullablePatchedTableUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedTableUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


