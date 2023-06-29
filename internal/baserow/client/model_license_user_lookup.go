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

// checks if the LicenseUserLookup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LicenseUserLookup{}

// LicenseUserLookup struct for LicenseUserLookup
type LicenseUserLookup struct {
	Id int32 `json:"id"`
	// The name and the email address of the user.
	Value string `json:"value"`
}

// NewLicenseUserLookup instantiates a new LicenseUserLookup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLicenseUserLookup(id int32, value string) *LicenseUserLookup {
	this := LicenseUserLookup{}
	this.Id = id
	this.Value = value
	return &this
}

// NewLicenseUserLookupWithDefaults instantiates a new LicenseUserLookup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLicenseUserLookupWithDefaults() *LicenseUserLookup {
	this := LicenseUserLookup{}
	return &this
}

// GetId returns the Id field value
func (o *LicenseUserLookup) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *LicenseUserLookup) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *LicenseUserLookup) SetId(v int32) {
	o.Id = v
}

// GetValue returns the Value field value
func (o *LicenseUserLookup) GetValue() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *LicenseUserLookup) GetValueOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *LicenseUserLookup) SetValue(v string) {
	o.Value = v
}

func (o LicenseUserLookup) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LicenseUserLookup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: value is readOnly
	return toSerialize, nil
}

type NullableLicenseUserLookup struct {
	value *LicenseUserLookup
	isSet bool
}

func (v NullableLicenseUserLookup) Get() *LicenseUserLookup {
	return v.value
}

func (v *NullableLicenseUserLookup) Set(val *LicenseUserLookup) {
	v.value = val
	v.isSet = true
}

func (v NullableLicenseUserLookup) IsSet() bool {
	return v.isSet
}

func (v *NullableLicenseUserLookup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLicenseUserLookup(val *LicenseUserLookup) *NullableLicenseUserLookup {
	return &NullableLicenseUserLookup{value: val, isSet: true}
}

func (v NullableLicenseUserLookup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLicenseUserLookup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


