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

// checks if the Token type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Token{}

// Token A mixin that allows us to rename the `group` field to `workspace` when serializing.
type Token struct {
	Id int32 `json:"id"`
	// The human readable name of the database token for the user.
	Name string `json:"name"`
	Group string `json:"group"`
	// Only the tables of the workspace can be accessed.
	Workspace int32 `json:"workspace"`
	// The unique token key that can be used to authorize for the table row endpoints.
	Key string `json:"key"`
	Permissions PatchedTokenUpdatePermissions `json:"permissions"`
}

// NewToken instantiates a new Token object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewToken(id int32, name string, group string, workspace int32, key string, permissions PatchedTokenUpdatePermissions) *Token {
	this := Token{}
	this.Id = id
	this.Name = name
	this.Group = group
	this.Workspace = workspace
	this.Key = key
	this.Permissions = permissions
	return &this
}

// NewTokenWithDefaults instantiates a new Token object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTokenWithDefaults() *Token {
	this := Token{}
	return &this
}

// GetId returns the Id field value
func (o *Token) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Token) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Token) SetId(v int32) {
	o.Id = v
}

// GetName returns the Name field value
func (o *Token) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Token) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Token) SetName(v string) {
	o.Name = v
}

// GetGroup returns the Group field value
func (o *Token) GetGroup() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Group
}

// GetGroupOk returns a tuple with the Group field value
// and a boolean to check if the value has been set.
func (o *Token) GetGroupOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Group, true
}

// SetGroup sets field value
func (o *Token) SetGroup(v string) {
	o.Group = v
}

// GetWorkspace returns the Workspace field value
func (o *Token) GetWorkspace() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Workspace
}

// GetWorkspaceOk returns a tuple with the Workspace field value
// and a boolean to check if the value has been set.
func (o *Token) GetWorkspaceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Workspace, true
}

// SetWorkspace sets field value
func (o *Token) SetWorkspace(v int32) {
	o.Workspace = v
}

// GetKey returns the Key field value
func (o *Token) GetKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Key
}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
func (o *Token) GetKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Key, true
}

// SetKey sets field value
func (o *Token) SetKey(v string) {
	o.Key = v
}

// GetPermissions returns the Permissions field value
func (o *Token) GetPermissions() PatchedTokenUpdatePermissions {
	if o == nil {
		var ret PatchedTokenUpdatePermissions
		return ret
	}

	return o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value
// and a boolean to check if the value has been set.
func (o *Token) GetPermissionsOk() (*PatchedTokenUpdatePermissions, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Permissions, true
}

// SetPermissions sets field value
func (o *Token) SetPermissions(v PatchedTokenUpdatePermissions) {
	o.Permissions = v
}

func (o Token) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Token) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["name"] = o.Name
	// skip: group is readOnly
	toSerialize["workspace"] = o.Workspace
	toSerialize["key"] = o.Key
	toSerialize["permissions"] = o.Permissions
	return toSerialize, nil
}

type NullableToken struct {
	value *Token
	isSet bool
}

func (v NullableToken) Get() *Token {
	return v.value
}

func (v *NullableToken) Set(val *Token) {
	v.value = val
	v.isSet = true
}

func (v NullableToken) IsSet() bool {
	return v.isSet
}

func (v *NullableToken) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableToken(val *Token) *NullableToken {
	return &NullableToken{value: val, isSet: true}
}

func (v NullableToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableToken) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

