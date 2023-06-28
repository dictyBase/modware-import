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

// checks if the OrderViews type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OrderViews{}

// OrderViews struct for OrderViews
type OrderViews struct {
	// View ids in the desired order.
	ViewIds []int32 `json:"view_ids"`
}

// NewOrderViews instantiates a new OrderViews object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrderViews(viewIds []int32) *OrderViews {
	this := OrderViews{}
	this.ViewIds = viewIds
	return &this
}

// NewOrderViewsWithDefaults instantiates a new OrderViews object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrderViewsWithDefaults() *OrderViews {
	this := OrderViews{}
	return &this
}

// GetViewIds returns the ViewIds field value
func (o *OrderViews) GetViewIds() []int32 {
	if o == nil {
		var ret []int32
		return ret
	}

	return o.ViewIds
}

// GetViewIdsOk returns a tuple with the ViewIds field value
// and a boolean to check if the value has been set.
func (o *OrderViews) GetViewIdsOk() ([]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ViewIds, true
}

// SetViewIds sets field value
func (o *OrderViews) SetViewIds(v []int32) {
	o.ViewIds = v
}

func (o OrderViews) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OrderViews) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["view_ids"] = o.ViewIds
	return toSerialize, nil
}

type NullableOrderViews struct {
	value *OrderViews
	isSet bool
}

func (v NullableOrderViews) Get() *OrderViews {
	return v.value
}

func (v *NullableOrderViews) Set(val *OrderViews) {
	v.value = val
	v.isSet = true
}

func (v NullableOrderViews) IsSet() bool {
	return v.isSet
}

func (v *NullableOrderViews) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOrderViews(val *OrderViews) *NullableOrderViews {
	return &NullableOrderViews{value: val, isSet: true}
}

func (v NullableOrderViews) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOrderViews) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


