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

// checks if the AuditLogActionType type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuditLogActionType{}

// AuditLogActionType struct for AuditLogActionType
type AuditLogActionType struct {
	Id IdEnum `json:"id"`
	// Given the *incoming* primitive data, return the value for this field that should be validated and transformed to a native value.
	Value string `json:"value"`
}

// NewAuditLogActionType instantiates a new AuditLogActionType object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuditLogActionType(id IdEnum, value string) *AuditLogActionType {
	this := AuditLogActionType{}
	this.Id = id
	this.Value = value
	return &this
}

// NewAuditLogActionTypeWithDefaults instantiates a new AuditLogActionType object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuditLogActionTypeWithDefaults() *AuditLogActionType {
	this := AuditLogActionType{}
	return &this
}

// GetId returns the Id field value
func (o *AuditLogActionType) GetId() IdEnum {
	if o == nil {
		var ret IdEnum
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AuditLogActionType) GetIdOk() (*IdEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AuditLogActionType) SetId(v IdEnum) {
	o.Id = v
}

// GetValue returns the Value field value
func (o *AuditLogActionType) GetValue() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *AuditLogActionType) GetValueOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *AuditLogActionType) SetValue(v string) {
	o.Value = v
}

func (o AuditLogActionType) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuditLogActionType) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	// skip: value is readOnly
	return toSerialize, nil
}

type NullableAuditLogActionType struct {
	value *AuditLogActionType
	isSet bool
}

func (v NullableAuditLogActionType) Get() *AuditLogActionType {
	return v.value
}

func (v *NullableAuditLogActionType) Set(val *AuditLogActionType) {
	v.value = val
	v.isSet = true
}

func (v NullableAuditLogActionType) IsSet() bool {
	return v.isSet
}

func (v *NullableAuditLogActionType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuditLogActionType(val *AuditLogActionType) *NullableAuditLogActionType {
	return &NullableAuditLogActionType{value: val, isSet: true}
}

func (v NullableAuditLogActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuditLogActionType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


