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

// checks if the AuditLogUser type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuditLogUser{}

// AuditLogUser struct for AuditLogUser
type AuditLogUser struct {
	Id int32 `json:"id"`
	Value string `json:"value"`
}

// NewAuditLogUser instantiates a new AuditLogUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuditLogUser(id int32, value string) *AuditLogUser {
	this := AuditLogUser{}
	this.Id = id
	this.Value = value
	return &this
}

// NewAuditLogUserWithDefaults instantiates a new AuditLogUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuditLogUserWithDefaults() *AuditLogUser {
	this := AuditLogUser{}
	return &this
}

// GetId returns the Id field value
func (o *AuditLogUser) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AuditLogUser) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AuditLogUser) SetId(v int32) {
	o.Id = v
}

// GetValue returns the Value field value
func (o *AuditLogUser) GetValue() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *AuditLogUser) GetValueOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *AuditLogUser) SetValue(v string) {
	o.Value = v
}

func (o AuditLogUser) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuditLogUser) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["value"] = o.Value
	return toSerialize, nil
}

type NullableAuditLogUser struct {
	value *AuditLogUser
	isSet bool
}

func (v NullableAuditLogUser) Get() *AuditLogUser {
	return v.value
}

func (v *NullableAuditLogUser) Set(val *AuditLogUser) {
	v.value = val
	v.isSet = true
}

func (v NullableAuditLogUser) IsSet() bool {
	return v.isSet
}

func (v *NullableAuditLogUser) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuditLogUser(val *AuditLogUser) *NullableAuditLogUser {
	return &NullableAuditLogUser{value: val, isSet: true}
}

func (v NullableAuditLogUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuditLogUser) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


