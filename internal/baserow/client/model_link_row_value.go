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

// checks if the LinkRowValue type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LinkRowValue{}

// LinkRowValue struct for LinkRowValue
type LinkRowValue struct {
	// The unique identifier of the row in the related table.
	Id int32 `json:"id"`
	// The primary field's value as a string of the row in the related table.
	Value *string `json:"value,omitempty"`
}

// NewLinkRowValue instantiates a new LinkRowValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinkRowValue(id int32) *LinkRowValue {
	this := LinkRowValue{}
	this.Id = id
	return &this
}

// NewLinkRowValueWithDefaults instantiates a new LinkRowValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinkRowValueWithDefaults() *LinkRowValue {
	this := LinkRowValue{}
	return &this
}

// GetId returns the Id field value
func (o *LinkRowValue) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *LinkRowValue) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *LinkRowValue) SetId(v int32) {
	o.Id = v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *LinkRowValue) GetValue() string {
	if o == nil || IsNil(o.Value) {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkRowValue) GetValueOk() (*string, bool) {
	if o == nil || IsNil(o.Value) {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *LinkRowValue) HasValue() bool {
	if o != nil && !IsNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *LinkRowValue) SetValue(v string) {
	o.Value = &v
}

func (o LinkRowValue) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LinkRowValue) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	if !IsNil(o.Value) {
		toSerialize["value"] = o.Value
	}
	return toSerialize, nil
}

type NullableLinkRowValue struct {
	value *LinkRowValue
	isSet bool
}

func (v NullableLinkRowValue) Get() *LinkRowValue {
	return v.value
}

func (v *NullableLinkRowValue) Set(val *LinkRowValue) {
	v.value = val
	v.isSet = true
}

func (v NullableLinkRowValue) IsSet() bool {
	return v.isSet
}

func (v *NullableLinkRowValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinkRowValue(val *LinkRowValue) *NullableLinkRowValue {
	return &NullableLinkRowValue{value: val, isSet: true}
}

func (v NullableLinkRowValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinkRowValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


