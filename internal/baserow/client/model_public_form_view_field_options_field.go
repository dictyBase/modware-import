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

// checks if the PublicFormViewFieldOptionsField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PublicFormViewFieldOptionsField{}

// PublicFormViewFieldOptionsField The properties of the related field. These can be used to construct the correct input. Additional properties could be added depending on the field type.
type PublicFormViewFieldOptionsField struct {
	Id int32 `json:"id"`
	// The type of the related field.
	Type string `json:"type"`
}

// NewPublicFormViewFieldOptionsField instantiates a new PublicFormViewFieldOptionsField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPublicFormViewFieldOptionsField(id int32, type_ string) *PublicFormViewFieldOptionsField {
	this := PublicFormViewFieldOptionsField{}
	this.Id = id
	this.Type = type_
	return &this
}

// NewPublicFormViewFieldOptionsFieldWithDefaults instantiates a new PublicFormViewFieldOptionsField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPublicFormViewFieldOptionsFieldWithDefaults() *PublicFormViewFieldOptionsField {
	this := PublicFormViewFieldOptionsField{}
	return &this
}

// GetId returns the Id field value
func (o *PublicFormViewFieldOptionsField) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *PublicFormViewFieldOptionsField) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *PublicFormViewFieldOptionsField) SetId(v int32) {
	o.Id = v
}

// GetType returns the Type field value
func (o *PublicFormViewFieldOptionsField) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *PublicFormViewFieldOptionsField) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *PublicFormViewFieldOptionsField) SetType(v string) {
	o.Type = v
}

func (o PublicFormViewFieldOptionsField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PublicFormViewFieldOptionsField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: type is readOnly
	return toSerialize, nil
}

type NullablePublicFormViewFieldOptionsField struct {
	value *PublicFormViewFieldOptionsField
	isSet bool
}

func (v NullablePublicFormViewFieldOptionsField) Get() *PublicFormViewFieldOptionsField {
	return v.value
}

func (v *NullablePublicFormViewFieldOptionsField) Set(val *PublicFormViewFieldOptionsField) {
	v.value = val
	v.isSet = true
}

func (v NullablePublicFormViewFieldOptionsField) IsSet() bool {
	return v.isSet
}

func (v *NullablePublicFormViewFieldOptionsField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePublicFormViewFieldOptionsField(val *PublicFormViewFieldOptionsField) *NullablePublicFormViewFieldOptionsField {
	return &NullablePublicFormViewFieldOptionsField{value: val, isSet: true}
}

func (v NullablePublicFormViewFieldOptionsField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePublicFormViewFieldOptionsField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


