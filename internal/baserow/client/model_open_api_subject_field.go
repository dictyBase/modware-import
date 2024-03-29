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

// checks if the OpenApiSubjectField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OpenApiSubjectField{}

// OpenApiSubjectField struct for OpenApiSubjectField
type OpenApiSubjectField struct {
	Id int32 `json:"id"`
}

// NewOpenApiSubjectField instantiates a new OpenApiSubjectField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOpenApiSubjectField(id int32) *OpenApiSubjectField {
	this := OpenApiSubjectField{}
	this.Id = id
	return &this
}

// NewOpenApiSubjectFieldWithDefaults instantiates a new OpenApiSubjectField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOpenApiSubjectFieldWithDefaults() *OpenApiSubjectField {
	this := OpenApiSubjectField{}
	return &this
}

// GetId returns the Id field value
func (o *OpenApiSubjectField) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *OpenApiSubjectField) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *OpenApiSubjectField) SetId(v int32) {
	o.Id = v
}

func (o OpenApiSubjectField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OpenApiSubjectField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

type NullableOpenApiSubjectField struct {
	value *OpenApiSubjectField
	isSet bool
}

func (v NullableOpenApiSubjectField) Get() *OpenApiSubjectField {
	return v.value
}

func (v *NullableOpenApiSubjectField) Set(val *OpenApiSubjectField) {
	v.value = val
	v.isSet = true
}

func (v NullableOpenApiSubjectField) IsSet() bool {
	return v.isSet
}

func (v *NullableOpenApiSubjectField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOpenApiSubjectField(val *OpenApiSubjectField) *NullableOpenApiSubjectField {
	return &NullableOpenApiSubjectField{value: val, isSet: true}
}

func (v NullableOpenApiSubjectField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOpenApiSubjectField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


