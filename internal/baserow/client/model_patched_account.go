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

// checks if the PatchedAccount type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedAccount{}

// PatchedAccount This serializer must be kept in sync with `UserSerializer`.
type PatchedAccount struct {
	FirstName *string `json:"first_name,omitempty"`
	// An ISO 639 language code (with optional variant) selected by the user. Ex: en-GB.
	Language *string `json:"language,omitempty"`
}

// NewPatchedAccount instantiates a new PatchedAccount object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedAccount() *PatchedAccount {
	this := PatchedAccount{}
	return &this
}

// NewPatchedAccountWithDefaults instantiates a new PatchedAccount object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedAccountWithDefaults() *PatchedAccount {
	this := PatchedAccount{}
	return &this
}

// GetFirstName returns the FirstName field value if set, zero value otherwise.
func (o *PatchedAccount) GetFirstName() string {
	if o == nil || IsNil(o.FirstName) {
		var ret string
		return ret
	}
	return *o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedAccount) GetFirstNameOk() (*string, bool) {
	if o == nil || IsNil(o.FirstName) {
		return nil, false
	}
	return o.FirstName, true
}

// HasFirstName returns a boolean if a field has been set.
func (o *PatchedAccount) HasFirstName() bool {
	if o != nil && !IsNil(o.FirstName) {
		return true
	}

	return false
}

// SetFirstName gets a reference to the given string and assigns it to the FirstName field.
func (o *PatchedAccount) SetFirstName(v string) {
	o.FirstName = &v
}

// GetLanguage returns the Language field value if set, zero value otherwise.
func (o *PatchedAccount) GetLanguage() string {
	if o == nil || IsNil(o.Language) {
		var ret string
		return ret
	}
	return *o.Language
}

// GetLanguageOk returns a tuple with the Language field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedAccount) GetLanguageOk() (*string, bool) {
	if o == nil || IsNil(o.Language) {
		return nil, false
	}
	return o.Language, true
}

// HasLanguage returns a boolean if a field has been set.
func (o *PatchedAccount) HasLanguage() bool {
	if o != nil && !IsNil(o.Language) {
		return true
	}

	return false
}

// SetLanguage gets a reference to the given string and assigns it to the Language field.
func (o *PatchedAccount) SetLanguage(v string) {
	o.Language = &v
}

func (o PatchedAccount) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedAccount) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.FirstName) {
		toSerialize["first_name"] = o.FirstName
	}
	if !IsNil(o.Language) {
		toSerialize["language"] = o.Language
	}
	return toSerialize, nil
}

type NullablePatchedAccount struct {
	value *PatchedAccount
	isSet bool
}

func (v NullablePatchedAccount) Get() *PatchedAccount {
	return v.value
}

func (v *NullablePatchedAccount) Set(val *PatchedAccount) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedAccount) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedAccount) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedAccount(val *PatchedAccount) *NullablePatchedAccount {
	return &NullablePatchedAccount{value: val, isSet: true}
}

func (v NullablePatchedAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedAccount) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


