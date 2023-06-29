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

// checks if the PublicViewInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PublicViewInfo{}

// PublicViewInfo struct for PublicViewInfo
type PublicViewInfo struct {
	Fields []PublicField `json:"fields"`
	View PublicView `json:"view"`
}

// NewPublicViewInfo instantiates a new PublicViewInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPublicViewInfo(fields []PublicField, view PublicView) *PublicViewInfo {
	this := PublicViewInfo{}
	this.Fields = fields
	this.View = view
	return &this
}

// NewPublicViewInfoWithDefaults instantiates a new PublicViewInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPublicViewInfoWithDefaults() *PublicViewInfo {
	this := PublicViewInfo{}
	return &this
}

// GetFields returns the Fields field value
func (o *PublicViewInfo) GetFields() []PublicField {
	if o == nil {
		var ret []PublicField
		return ret
	}

	return o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value
// and a boolean to check if the value has been set.
func (o *PublicViewInfo) GetFieldsOk() ([]PublicField, bool) {
	if o == nil {
		return nil, false
	}
	return o.Fields, true
}

// SetFields sets field value
func (o *PublicViewInfo) SetFields(v []PublicField) {
	o.Fields = v
}

// GetView returns the View field value
func (o *PublicViewInfo) GetView() PublicView {
	if o == nil {
		var ret PublicView
		return ret
	}

	return o.View
}

// GetViewOk returns a tuple with the View field value
// and a boolean to check if the value has been set.
func (o *PublicViewInfo) GetViewOk() (*PublicView, bool) {
	if o == nil {
		return nil, false
	}
	return &o.View, true
}

// SetView sets field value
func (o *PublicViewInfo) SetView(v PublicView) {
	o.View = v
}

func (o PublicViewInfo) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PublicViewInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: fields is readOnly
	// skip: view is readOnly
	return toSerialize, nil
}

type NullablePublicViewInfo struct {
	value *PublicViewInfo
	isSet bool
}

func (v NullablePublicViewInfo) Get() *PublicViewInfo {
	return v.value
}

func (v *NullablePublicViewInfo) Set(val *PublicViewInfo) {
	v.value = val
	v.isSet = true
}

func (v NullablePublicViewInfo) IsSet() bool {
	return v.isSet
}

func (v *NullablePublicViewInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePublicViewInfo(val *PublicViewInfo) *NullablePublicViewInfo {
	return &NullablePublicViewInfo{value: val, isSet: true}
}

func (v NullablePublicViewInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePublicViewInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

