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

// checks if the ResetPasswordBodyValidation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ResetPasswordBodyValidation{}

// ResetPasswordBodyValidation struct for ResetPasswordBodyValidation
type ResetPasswordBodyValidation struct {
	Token string `json:"token"`
	Password string `json:"password"`
}

// NewResetPasswordBodyValidation instantiates a new ResetPasswordBodyValidation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResetPasswordBodyValidation(token string, password string) *ResetPasswordBodyValidation {
	this := ResetPasswordBodyValidation{}
	this.Token = token
	this.Password = password
	return &this
}

// NewResetPasswordBodyValidationWithDefaults instantiates a new ResetPasswordBodyValidation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResetPasswordBodyValidationWithDefaults() *ResetPasswordBodyValidation {
	this := ResetPasswordBodyValidation{}
	return &this
}

// GetToken returns the Token field value
func (o *ResetPasswordBodyValidation) GetToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Token
}

// GetTokenOk returns a tuple with the Token field value
// and a boolean to check if the value has been set.
func (o *ResetPasswordBodyValidation) GetTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Token, true
}

// SetToken sets field value
func (o *ResetPasswordBodyValidation) SetToken(v string) {
	o.Token = v
}

// GetPassword returns the Password field value
func (o *ResetPasswordBodyValidation) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *ResetPasswordBodyValidation) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *ResetPasswordBodyValidation) SetPassword(v string) {
	o.Password = v
}

func (o ResetPasswordBodyValidation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ResetPasswordBodyValidation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["token"] = o.Token
	toSerialize["password"] = o.Password
	return toSerialize, nil
}

type NullableResetPasswordBodyValidation struct {
	value *ResetPasswordBodyValidation
	isSet bool
}

func (v NullableResetPasswordBodyValidation) Get() *ResetPasswordBodyValidation {
	return v.value
}

func (v *NullableResetPasswordBodyValidation) Set(val *ResetPasswordBodyValidation) {
	v.value = val
	v.isSet = true
}

func (v NullableResetPasswordBodyValidation) IsSet() bool {
	return v.isSet
}

func (v *NullableResetPasswordBodyValidation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResetPasswordBodyValidation(val *ResetPasswordBodyValidation) *NullableResetPasswordBodyValidation {
	return &NullableResetPasswordBodyValidation{value: val, isSet: true}
}

func (v NullableResetPasswordBodyValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResetPasswordBodyValidation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

