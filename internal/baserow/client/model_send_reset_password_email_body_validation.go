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

// checks if the SendResetPasswordEmailBodyValidation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SendResetPasswordEmailBodyValidation{}

// SendResetPasswordEmailBodyValidation struct for SendResetPasswordEmailBodyValidation
type SendResetPasswordEmailBodyValidation struct {
	// The email address of the user that has requested a password reset.
	Email string `json:"email"`
	// The base URL where the user can reset his password. The reset token is going to be appended to the base_url (base_url '/token').
	BaseUrl string `json:"base_url"`
}

// NewSendResetPasswordEmailBodyValidation instantiates a new SendResetPasswordEmailBodyValidation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSendResetPasswordEmailBodyValidation(email string, baseUrl string) *SendResetPasswordEmailBodyValidation {
	this := SendResetPasswordEmailBodyValidation{}
	this.Email = email
	this.BaseUrl = baseUrl
	return &this
}

// NewSendResetPasswordEmailBodyValidationWithDefaults instantiates a new SendResetPasswordEmailBodyValidation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSendResetPasswordEmailBodyValidationWithDefaults() *SendResetPasswordEmailBodyValidation {
	this := SendResetPasswordEmailBodyValidation{}
	return &this
}

// GetEmail returns the Email field value
func (o *SendResetPasswordEmailBodyValidation) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *SendResetPasswordEmailBodyValidation) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *SendResetPasswordEmailBodyValidation) SetEmail(v string) {
	o.Email = v
}

// GetBaseUrl returns the BaseUrl field value
func (o *SendResetPasswordEmailBodyValidation) GetBaseUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BaseUrl
}

// GetBaseUrlOk returns a tuple with the BaseUrl field value
// and a boolean to check if the value has been set.
func (o *SendResetPasswordEmailBodyValidation) GetBaseUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BaseUrl, true
}

// SetBaseUrl sets field value
func (o *SendResetPasswordEmailBodyValidation) SetBaseUrl(v string) {
	o.BaseUrl = v
}

func (o SendResetPasswordEmailBodyValidation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SendResetPasswordEmailBodyValidation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["email"] = o.Email
	toSerialize["base_url"] = o.BaseUrl
	return toSerialize, nil
}

type NullableSendResetPasswordEmailBodyValidation struct {
	value *SendResetPasswordEmailBodyValidation
	isSet bool
}

func (v NullableSendResetPasswordEmailBodyValidation) Get() *SendResetPasswordEmailBodyValidation {
	return v.value
}

func (v *NullableSendResetPasswordEmailBodyValidation) Set(val *SendResetPasswordEmailBodyValidation) {
	v.value = val
	v.isSet = true
}

func (v NullableSendResetPasswordEmailBodyValidation) IsSet() bool {
	return v.isSet
}

func (v *NullableSendResetPasswordEmailBodyValidation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSendResetPasswordEmailBodyValidation(val *SendResetPasswordEmailBodyValidation) *NullableSendResetPasswordEmailBodyValidation {
	return &NullableSendResetPasswordEmailBodyValidation{value: val, isSet: true}
}

func (v NullableSendResetPasswordEmailBodyValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSendResetPasswordEmailBodyValidation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


