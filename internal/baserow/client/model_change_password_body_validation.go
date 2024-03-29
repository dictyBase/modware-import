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

// checks if the ChangePasswordBodyValidation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ChangePasswordBodyValidation{}

// ChangePasswordBodyValidation struct for ChangePasswordBodyValidation
type ChangePasswordBodyValidation struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// NewChangePasswordBodyValidation instantiates a new ChangePasswordBodyValidation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewChangePasswordBodyValidation(oldPassword string, newPassword string) *ChangePasswordBodyValidation {
	this := ChangePasswordBodyValidation{}
	this.OldPassword = oldPassword
	this.NewPassword = newPassword
	return &this
}

// NewChangePasswordBodyValidationWithDefaults instantiates a new ChangePasswordBodyValidation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewChangePasswordBodyValidationWithDefaults() *ChangePasswordBodyValidation {
	this := ChangePasswordBodyValidation{}
	return &this
}

// GetOldPassword returns the OldPassword field value
func (o *ChangePasswordBodyValidation) GetOldPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OldPassword
}

// GetOldPasswordOk returns a tuple with the OldPassword field value
// and a boolean to check if the value has been set.
func (o *ChangePasswordBodyValidation) GetOldPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OldPassword, true
}

// SetOldPassword sets field value
func (o *ChangePasswordBodyValidation) SetOldPassword(v string) {
	o.OldPassword = v
}

// GetNewPassword returns the NewPassword field value
func (o *ChangePasswordBodyValidation) GetNewPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NewPassword
}

// GetNewPasswordOk returns a tuple with the NewPassword field value
// and a boolean to check if the value has been set.
func (o *ChangePasswordBodyValidation) GetNewPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NewPassword, true
}

// SetNewPassword sets field value
func (o *ChangePasswordBodyValidation) SetNewPassword(v string) {
	o.NewPassword = v
}

func (o ChangePasswordBodyValidation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ChangePasswordBodyValidation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["old_password"] = o.OldPassword
	toSerialize["new_password"] = o.NewPassword
	return toSerialize, nil
}

type NullableChangePasswordBodyValidation struct {
	value *ChangePasswordBodyValidation
	isSet bool
}

func (v NullableChangePasswordBodyValidation) Get() *ChangePasswordBodyValidation {
	return v.value
}

func (v *NullableChangePasswordBodyValidation) Set(val *ChangePasswordBodyValidation) {
	v.value = val
	v.isSet = true
}

func (v NullableChangePasswordBodyValidation) IsSet() bool {
	return v.isSet
}

func (v *NullableChangePasswordBodyValidation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableChangePasswordBodyValidation(val *ChangePasswordBodyValidation) *NullableChangePasswordBodyValidation {
	return &NullableChangePasswordBodyValidation{value: val, isSet: true}
}

func (v NullableChangePasswordBodyValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableChangePasswordBodyValidation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


