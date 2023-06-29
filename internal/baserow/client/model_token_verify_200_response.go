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

// checks if the TokenVerify200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TokenVerify200Response{}

// TokenVerify200Response struct for TokenVerify200Response
type TokenVerify200Response struct {
	User *AdminImpersonateUser200ResponseUser `json:"user,omitempty"`
}

// NewTokenVerify200Response instantiates a new TokenVerify200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTokenVerify200Response() *TokenVerify200Response {
	this := TokenVerify200Response{}
	return &this
}

// NewTokenVerify200ResponseWithDefaults instantiates a new TokenVerify200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTokenVerify200ResponseWithDefaults() *TokenVerify200Response {
	this := TokenVerify200Response{}
	return &this
}

// GetUser returns the User field value if set, zero value otherwise.
func (o *TokenVerify200Response) GetUser() AdminImpersonateUser200ResponseUser {
	if o == nil || IsNil(o.User) {
		var ret AdminImpersonateUser200ResponseUser
		return ret
	}
	return *o.User
}

// GetUserOk returns a tuple with the User field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TokenVerify200Response) GetUserOk() (*AdminImpersonateUser200ResponseUser, bool) {
	if o == nil || IsNil(o.User) {
		return nil, false
	}
	return o.User, true
}

// HasUser returns a boolean if a field has been set.
func (o *TokenVerify200Response) HasUser() bool {
	if o != nil && !IsNil(o.User) {
		return true
	}

	return false
}

// SetUser gets a reference to the given AdminImpersonateUser200ResponseUser and assigns it to the User field.
func (o *TokenVerify200Response) SetUser(v AdminImpersonateUser200ResponseUser) {
	o.User = &v
}

func (o TokenVerify200Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TokenVerify200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.User) {
		toSerialize["user"] = o.User
	}
	return toSerialize, nil
}

type NullableTokenVerify200Response struct {
	value *TokenVerify200Response
	isSet bool
}

func (v NullableTokenVerify200Response) Get() *TokenVerify200Response {
	return v.value
}

func (v *NullableTokenVerify200Response) Set(val *TokenVerify200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableTokenVerify200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableTokenVerify200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTokenVerify200Response(val *TokenVerify200Response) *NullableTokenVerify200Response {
	return &NullableTokenVerify200Response{value: val, isSet: true}
}

func (v NullableTokenVerify200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTokenVerify200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

