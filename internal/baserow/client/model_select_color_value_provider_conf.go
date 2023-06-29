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

// checks if the SelectColorValueProviderConf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SelectColorValueProviderConf{}

// SelectColorValueProviderConf struct for SelectColorValueProviderConf
type SelectColorValueProviderConf struct {
	// An id of a select field of the table. The value provider return the color of the selected option for each row.
	FieldId NullableInt32 `json:"field_id"`
}

// NewSelectColorValueProviderConf instantiates a new SelectColorValueProviderConf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSelectColorValueProviderConf(fieldId NullableInt32) *SelectColorValueProviderConf {
	this := SelectColorValueProviderConf{}
	this.FieldId = fieldId
	return &this
}

// NewSelectColorValueProviderConfWithDefaults instantiates a new SelectColorValueProviderConf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSelectColorValueProviderConfWithDefaults() *SelectColorValueProviderConf {
	this := SelectColorValueProviderConf{}
	return &this
}

// GetFieldId returns the FieldId field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *SelectColorValueProviderConf) GetFieldId() int32 {
	if o == nil || o.FieldId.Get() == nil {
		var ret int32
		return ret
	}

	return *o.FieldId.Get()
}

// GetFieldIdOk returns a tuple with the FieldId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SelectColorValueProviderConf) GetFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.FieldId.Get(), o.FieldId.IsSet()
}

// SetFieldId sets field value
func (o *SelectColorValueProviderConf) SetFieldId(v int32) {
	o.FieldId.Set(&v)
}

func (o SelectColorValueProviderConf) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SelectColorValueProviderConf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["field_id"] = o.FieldId.Get()
	return toSerialize, nil
}

type NullableSelectColorValueProviderConf struct {
	value *SelectColorValueProviderConf
	isSet bool
}

func (v NullableSelectColorValueProviderConf) Get() *SelectColorValueProviderConf {
	return v.value
}

func (v *NullableSelectColorValueProviderConf) Set(val *SelectColorValueProviderConf) {
	v.value = val
	v.isSet = true
}

func (v NullableSelectColorValueProviderConf) IsSet() bool {
	return v.isSet
}

func (v *NullableSelectColorValueProviderConf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSelectColorValueProviderConf(val *SelectColorValueProviderConf) *NullableSelectColorValueProviderConf {
	return &NullableSelectColorValueProviderConf{value: val, isSet: true}
}

func (v NullableSelectColorValueProviderConf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSelectColorValueProviderConf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

