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

// checks if the RowCommentCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RowCommentCreate{}

// RowCommentCreate struct for RowCommentCreate
type RowCommentCreate struct {
	Comment string `json:"comment"`
}

// NewRowCommentCreate instantiates a new RowCommentCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRowCommentCreate(comment string) *RowCommentCreate {
	this := RowCommentCreate{}
	this.Comment = comment
	return &this
}

// NewRowCommentCreateWithDefaults instantiates a new RowCommentCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRowCommentCreateWithDefaults() *RowCommentCreate {
	this := RowCommentCreate{}
	return &this
}

// GetComment returns the Comment field value
func (o *RowCommentCreate) GetComment() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Comment
}

// GetCommentOk returns a tuple with the Comment field value
// and a boolean to check if the value has been set.
func (o *RowCommentCreate) GetCommentOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Comment, true
}

// SetComment sets field value
func (o *RowCommentCreate) SetComment(v string) {
	o.Comment = v
}

func (o RowCommentCreate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RowCommentCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["comment"] = o.Comment
	return toSerialize, nil
}

type NullableRowCommentCreate struct {
	value *RowCommentCreate
	isSet bool
}

func (v NullableRowCommentCreate) Get() *RowCommentCreate {
	return v.value
}

func (v *NullableRowCommentCreate) Set(val *RowCommentCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableRowCommentCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableRowCommentCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRowCommentCreate(val *RowCommentCreate) *NullableRowCommentCreate {
	return &NullableRowCommentCreate{value: val, isSet: true}
}

func (v NullableRowCommentCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRowCommentCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


