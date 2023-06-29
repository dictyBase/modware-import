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

// checks if the CreateUser200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateUser200Response{}

// CreateUser200Response struct for CreateUser200Response
type CreateUser200Response struct {
	User *AdminImpersonateUser200ResponseUser `json:"user,omitempty"`
	// Deprecated. Use the `access_token` instead.
	// Deprecated
	Token *string `json:"token,omitempty"`
	// 'access_token' can be used to authorize for other endpoints that require authorization. This token will be valid for 10 minutes.
	AccessToken *string `json:"access_token,omitempty"`
	// 'refresh_token' can be used to get a new valid 'access_token'. This token will be valid for 168 hours.
	RefreshToken *string `json:"refresh_token,omitempty"`
}

// NewCreateUser200Response instantiates a new CreateUser200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateUser200Response() *CreateUser200Response {
	this := CreateUser200Response{}
	return &this
}

// NewCreateUser200ResponseWithDefaults instantiates a new CreateUser200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateUser200ResponseWithDefaults() *CreateUser200Response {
	this := CreateUser200Response{}
	return &this
}

// GetUser returns the User field value if set, zero value otherwise.
func (o *CreateUser200Response) GetUser() AdminImpersonateUser200ResponseUser {
	if o == nil || IsNil(o.User) {
		var ret AdminImpersonateUser200ResponseUser
		return ret
	}
	return *o.User
}

// GetUserOk returns a tuple with the User field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateUser200Response) GetUserOk() (*AdminImpersonateUser200ResponseUser, bool) {
	if o == nil || IsNil(o.User) {
		return nil, false
	}
	return o.User, true
}

// HasUser returns a boolean if a field has been set.
func (o *CreateUser200Response) HasUser() bool {
	if o != nil && !IsNil(o.User) {
		return true
	}

	return false
}

// SetUser gets a reference to the given AdminImpersonateUser200ResponseUser and assigns it to the User field.
func (o *CreateUser200Response) SetUser(v AdminImpersonateUser200ResponseUser) {
	o.User = &v
}

// GetToken returns the Token field value if set, zero value otherwise.
// Deprecated
func (o *CreateUser200Response) GetToken() string {
	if o == nil || IsNil(o.Token) {
		var ret string
		return ret
	}
	return *o.Token
}

// GetTokenOk returns a tuple with the Token field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *CreateUser200Response) GetTokenOk() (*string, bool) {
	if o == nil || IsNil(o.Token) {
		return nil, false
	}
	return o.Token, true
}

// HasToken returns a boolean if a field has been set.
func (o *CreateUser200Response) HasToken() bool {
	if o != nil && !IsNil(o.Token) {
		return true
	}

	return false
}

// SetToken gets a reference to the given string and assigns it to the Token field.
// Deprecated
func (o *CreateUser200Response) SetToken(v string) {
	o.Token = &v
}

// GetAccessToken returns the AccessToken field value if set, zero value otherwise.
func (o *CreateUser200Response) GetAccessToken() string {
	if o == nil || IsNil(o.AccessToken) {
		var ret string
		return ret
	}
	return *o.AccessToken
}

// GetAccessTokenOk returns a tuple with the AccessToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateUser200Response) GetAccessTokenOk() (*string, bool) {
	if o == nil || IsNil(o.AccessToken) {
		return nil, false
	}
	return o.AccessToken, true
}

// HasAccessToken returns a boolean if a field has been set.
func (o *CreateUser200Response) HasAccessToken() bool {
	if o != nil && !IsNil(o.AccessToken) {
		return true
	}

	return false
}

// SetAccessToken gets a reference to the given string and assigns it to the AccessToken field.
func (o *CreateUser200Response) SetAccessToken(v string) {
	o.AccessToken = &v
}

// GetRefreshToken returns the RefreshToken field value if set, zero value otherwise.
func (o *CreateUser200Response) GetRefreshToken() string {
	if o == nil || IsNil(o.RefreshToken) {
		var ret string
		return ret
	}
	return *o.RefreshToken
}

// GetRefreshTokenOk returns a tuple with the RefreshToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateUser200Response) GetRefreshTokenOk() (*string, bool) {
	if o == nil || IsNil(o.RefreshToken) {
		return nil, false
	}
	return o.RefreshToken, true
}

// HasRefreshToken returns a boolean if a field has been set.
func (o *CreateUser200Response) HasRefreshToken() bool {
	if o != nil && !IsNil(o.RefreshToken) {
		return true
	}

	return false
}

// SetRefreshToken gets a reference to the given string and assigns it to the RefreshToken field.
func (o *CreateUser200Response) SetRefreshToken(v string) {
	o.RefreshToken = &v
}

func (o CreateUser200Response) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateUser200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.User) {
		toSerialize["user"] = o.User
	}
	if !IsNil(o.Token) {
		toSerialize["token"] = o.Token
	}
	if !IsNil(o.AccessToken) {
		toSerialize["access_token"] = o.AccessToken
	}
	if !IsNil(o.RefreshToken) {
		toSerialize["refresh_token"] = o.RefreshToken
	}
	return toSerialize, nil
}

type NullableCreateUser200Response struct {
	value *CreateUser200Response
	isSet bool
}

func (v NullableCreateUser200Response) Get() *CreateUser200Response {
	return v.value
}

func (v *NullableCreateUser200Response) Set(val *CreateUser200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateUser200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateUser200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateUser200Response(val *CreateUser200Response) *NullableCreateUser200Response {
	return &NullableCreateUser200Response{value: val, isSet: true}
}

func (v NullableCreateUser200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateUser200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

