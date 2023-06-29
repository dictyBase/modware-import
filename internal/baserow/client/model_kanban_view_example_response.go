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

// checks if the KanbanViewExampleResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KanbanViewExampleResponse{}

// KanbanViewExampleResponse struct for KanbanViewExampleResponse
type KanbanViewExampleResponse struct {
	OPTION_ID KanbanViewExampleResponseOPTIONID `json:"OPTION_ID"`
	FieldOptions []KanbanViewFieldOptions `json:"field_options"`
}

// NewKanbanViewExampleResponse instantiates a new KanbanViewExampleResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKanbanViewExampleResponse(oPTIONID KanbanViewExampleResponseOPTIONID, fieldOptions []KanbanViewFieldOptions) *KanbanViewExampleResponse {
	this := KanbanViewExampleResponse{}
	this.OPTION_ID = oPTIONID
	this.FieldOptions = fieldOptions
	return &this
}

// NewKanbanViewExampleResponseWithDefaults instantiates a new KanbanViewExampleResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKanbanViewExampleResponseWithDefaults() *KanbanViewExampleResponse {
	this := KanbanViewExampleResponse{}
	return &this
}

// GetOPTION_ID returns the OPTION_ID field value
func (o *KanbanViewExampleResponse) GetOPTION_ID() KanbanViewExampleResponseOPTIONID {
	if o == nil {
		var ret KanbanViewExampleResponseOPTIONID
		return ret
	}

	return o.OPTION_ID
}

// GetOPTION_IDOk returns a tuple with the OPTION_ID field value
// and a boolean to check if the value has been set.
func (o *KanbanViewExampleResponse) GetOPTION_IDOk() (*KanbanViewExampleResponseOPTIONID, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OPTION_ID, true
}

// SetOPTION_ID sets field value
func (o *KanbanViewExampleResponse) SetOPTION_ID(v KanbanViewExampleResponseOPTIONID) {
	o.OPTION_ID = v
}

// GetFieldOptions returns the FieldOptions field value
func (o *KanbanViewExampleResponse) GetFieldOptions() []KanbanViewFieldOptions {
	if o == nil {
		var ret []KanbanViewFieldOptions
		return ret
	}

	return o.FieldOptions
}

// GetFieldOptionsOk returns a tuple with the FieldOptions field value
// and a boolean to check if the value has been set.
func (o *KanbanViewExampleResponse) GetFieldOptionsOk() ([]KanbanViewFieldOptions, bool) {
	if o == nil {
		return nil, false
	}
	return o.FieldOptions, true
}

// SetFieldOptions sets field value
func (o *KanbanViewExampleResponse) SetFieldOptions(v []KanbanViewFieldOptions) {
	o.FieldOptions = v
}

func (o KanbanViewExampleResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o KanbanViewExampleResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["OPTION_ID"] = o.OPTION_ID
	toSerialize["field_options"] = o.FieldOptions
	return toSerialize, nil
}

type NullableKanbanViewExampleResponse struct {
	value *KanbanViewExampleResponse
	isSet bool
}

func (v NullableKanbanViewExampleResponse) Get() *KanbanViewExampleResponse {
	return v.value
}

func (v *NullableKanbanViewExampleResponse) Set(val *KanbanViewExampleResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableKanbanViewExampleResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableKanbanViewExampleResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKanbanViewExampleResponse(val *KanbanViewExampleResponse) *NullableKanbanViewExampleResponse {
	return &NullableKanbanViewExampleResponse{value: val, isSet: true}
}

func (v NullableKanbanViewExampleResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKanbanViewExampleResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

