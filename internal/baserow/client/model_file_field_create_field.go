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

// checks if the FileFieldCreateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FileFieldCreateField{}

// FileFieldCreateField struct for FileFieldCreateField
type FileFieldCreateField struct {
	Name string `json:"name"`
	Type Type712Enum `json:"type"`
}

// NewFileFieldCreateField instantiates a new FileFieldCreateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFileFieldCreateField(name string, type_ Type712Enum) *FileFieldCreateField {
	this := FileFieldCreateField{}
	this.Name = name
	this.Type = type_
	return &this
}

// NewFileFieldCreateFieldWithDefaults instantiates a new FileFieldCreateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFileFieldCreateFieldWithDefaults() *FileFieldCreateField {
	this := FileFieldCreateField{}
	return &this
}

// GetName returns the Name field value
func (o *FileFieldCreateField) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *FileFieldCreateField) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *FileFieldCreateField) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *FileFieldCreateField) GetType() Type712Enum {
	if o == nil {
		var ret Type712Enum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *FileFieldCreateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *FileFieldCreateField) SetType(v Type712Enum) {
	o.Type = v
}

func (o FileFieldCreateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FileFieldCreateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	return toSerialize, nil
}

type NullableFileFieldCreateField struct {
	value *FileFieldCreateField
	isSet bool
}

func (v NullableFileFieldCreateField) Get() *FileFieldCreateField {
	return v.value
}

func (v *NullableFileFieldCreateField) Set(val *FileFieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableFileFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableFileFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFileFieldCreateField(val *FileFieldCreateField) *NullableFileFieldCreateField {
	return &NullableFileFieldCreateField{value: val, isSet: true}
}

func (v NullableFileFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFileFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


