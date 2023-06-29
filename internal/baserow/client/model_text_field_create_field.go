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

// checks if the TextFieldCreateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TextFieldCreateField{}

// TextFieldCreateField struct for TextFieldCreateField
type TextFieldCreateField struct {
	Name string `json:"name"`
	Type Type712Enum `json:"type"`
	// If set, this value is going to be added every time a new row created.
	TextDefault *string `json:"text_default,omitempty"`
}

// NewTextFieldCreateField instantiates a new TextFieldCreateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTextFieldCreateField(name string, type_ Type712Enum) *TextFieldCreateField {
	this := TextFieldCreateField{}
	this.Name = name
	this.Type = type_
	return &this
}

// NewTextFieldCreateFieldWithDefaults instantiates a new TextFieldCreateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTextFieldCreateFieldWithDefaults() *TextFieldCreateField {
	this := TextFieldCreateField{}
	return &this
}

// GetName returns the Name field value
func (o *TextFieldCreateField) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *TextFieldCreateField) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *TextFieldCreateField) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *TextFieldCreateField) GetType() Type712Enum {
	if o == nil {
		var ret Type712Enum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *TextFieldCreateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *TextFieldCreateField) SetType(v Type712Enum) {
	o.Type = v
}

// GetTextDefault returns the TextDefault field value if set, zero value otherwise.
func (o *TextFieldCreateField) GetTextDefault() string {
	if o == nil || IsNil(o.TextDefault) {
		var ret string
		return ret
	}
	return *o.TextDefault
}

// GetTextDefaultOk returns a tuple with the TextDefault field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TextFieldCreateField) GetTextDefaultOk() (*string, bool) {
	if o == nil || IsNil(o.TextDefault) {
		return nil, false
	}
	return o.TextDefault, true
}

// HasTextDefault returns a boolean if a field has been set.
func (o *TextFieldCreateField) HasTextDefault() bool {
	if o != nil && !IsNil(o.TextDefault) {
		return true
	}

	return false
}

// SetTextDefault gets a reference to the given string and assigns it to the TextDefault field.
func (o *TextFieldCreateField) SetTextDefault(v string) {
	o.TextDefault = &v
}

func (o TextFieldCreateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TextFieldCreateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	if !IsNil(o.TextDefault) {
		toSerialize["text_default"] = o.TextDefault
	}
	return toSerialize, nil
}

type NullableTextFieldCreateField struct {
	value *TextFieldCreateField
	isSet bool
}

func (v NullableTextFieldCreateField) Get() *TextFieldCreateField {
	return v.value
}

func (v *NullableTextFieldCreateField) Set(val *TextFieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableTextFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableTextFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTextFieldCreateField(val *TextFieldCreateField) *NullableTextFieldCreateField {
	return &NullableTextFieldCreateField{value: val, isSet: true}
}

func (v NullableTextFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTextFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

