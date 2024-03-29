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

// checks if the OrderDatabaseTables400Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OrderDatabaseTables400Response{}

// OrderDatabaseTables400Response struct for OrderDatabaseTables400Response
type OrderDatabaseTables400Response struct {
	// Machine readable error indicating what went wrong.
	Error *string `json:"error,omitempty"`
	Detail *AdminListUsers400ResponseDetail `json:"detail,omitempty"`
}

// NewOrderDatabaseTables400Response instantiates a new OrderDatabaseTables400Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrderDatabaseTables400Response() *OrderDatabaseTables400Response {
	this := OrderDatabaseTables400Response{}
	return &this
}

// NewOrderDatabaseTables400ResponseWithDefaults instantiates a new OrderDatabaseTables400Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrderDatabaseTables400ResponseWithDefaults() *OrderDatabaseTables400Response {
	this := OrderDatabaseTables400Response{}
	return &this
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *OrderDatabaseTables400Response) GetError() string {
	if o == nil || IsNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrderDatabaseTables400Response) GetErrorOk() (*string, bool) {
	if o == nil || IsNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *OrderDatabaseTables400Response) HasError() bool {
	if o != nil && !IsNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *OrderDatabaseTables400Response) SetError(v string) {
	o.Error = &v
}

// GetDetail returns the Detail field value if set, zero value otherwise.
func (o *OrderDatabaseTables400Response) GetDetail() AdminListUsers400ResponseDetail {
	if o == nil || IsNil(o.Detail) {
		var ret AdminListUsers400ResponseDetail
		return ret
	}
	return *o.Detail
}

// GetDetailOk returns a tuple with the Detail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrderDatabaseTables400Response) GetDetailOk() (*AdminListUsers400ResponseDetail, bool) {
	if o == nil || IsNil(o.Detail) {
		return nil, false
	}
	return o.Detail, true
}

// HasDetail returns a boolean if a field has been set.
func (o *OrderDatabaseTables400Response) HasDetail() bool {
	if o != nil && !IsNil(o.Detail) {
		return true
	}

	return false
}

// SetDetail gets a reference to the given AdminListUsers400ResponseDetail and assigns it to the Detail field.
func (o *OrderDatabaseTables400Response) SetDetail(v AdminListUsers400ResponseDetail) {
	o.Detail = &v
}

func (o OrderDatabaseTables400Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OrderDatabaseTables400Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	if !IsNil(o.Detail) {
		toSerialize["detail"] = o.Detail
	}
	return toSerialize, nil
}

type NullableOrderDatabaseTables400Response struct {
	value *OrderDatabaseTables400Response
	isSet bool
}

func (v NullableOrderDatabaseTables400Response) Get() *OrderDatabaseTables400Response {
	return v.value
}

func (v *NullableOrderDatabaseTables400Response) Set(val *OrderDatabaseTables400Response) {
	v.value = val
	v.isSet = true
}

func (v NullableOrderDatabaseTables400Response) IsSet() bool {
	return v.isSet
}

func (v *NullableOrderDatabaseTables400Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOrderDatabaseTables400Response(val *OrderDatabaseTables400Response) *NullableOrderDatabaseTables400Response {
	return &NullableOrderDatabaseTables400Response{value: val, isSet: true}
}

func (v NullableOrderDatabaseTables400Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOrderDatabaseTables400Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


