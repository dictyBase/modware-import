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

// checks if the WorkspaceAdminUsers type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkspaceAdminUsers{}

// WorkspaceAdminUsers struct for WorkspaceAdminUsers
type WorkspaceAdminUsers struct {
	Id int32 `json:"id"`
	Email string `json:"email"`
	// The permissions that the user has within the workspace.
	Permissions *string `json:"permissions,omitempty"`
}

// NewWorkspaceAdminUsers instantiates a new WorkspaceAdminUsers object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkspaceAdminUsers(id int32, email string) *WorkspaceAdminUsers {
	this := WorkspaceAdminUsers{}
	this.Id = id
	this.Email = email
	return &this
}

// NewWorkspaceAdminUsersWithDefaults instantiates a new WorkspaceAdminUsers object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkspaceAdminUsersWithDefaults() *WorkspaceAdminUsers {
	this := WorkspaceAdminUsers{}
	return &this
}

// GetId returns the Id field value
func (o *WorkspaceAdminUsers) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkspaceAdminUsers) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkspaceAdminUsers) SetId(v int32) {
	o.Id = v
}

// GetEmail returns the Email field value
func (o *WorkspaceAdminUsers) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *WorkspaceAdminUsers) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *WorkspaceAdminUsers) SetEmail(v string) {
	o.Email = v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *WorkspaceAdminUsers) GetPermissions() string {
	if o == nil || IsNil(o.Permissions) {
		var ret string
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkspaceAdminUsers) GetPermissionsOk() (*string, bool) {
	if o == nil || IsNil(o.Permissions) {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *WorkspaceAdminUsers) HasPermissions() bool {
	if o != nil && !IsNil(o.Permissions) {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given string and assigns it to the Permissions field.
func (o *WorkspaceAdminUsers) SetPermissions(v string) {
	o.Permissions = &v
}

func (o WorkspaceAdminUsers) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkspaceAdminUsers) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["email"] = o.Email
	if !IsNil(o.Permissions) {
		toSerialize["permissions"] = o.Permissions
	}
	return toSerialize, nil
}

type NullableWorkspaceAdminUsers struct {
	value *WorkspaceAdminUsers
	isSet bool
}

func (v NullableWorkspaceAdminUsers) Get() *WorkspaceAdminUsers {
	return v.value
}

func (v *NullableWorkspaceAdminUsers) Set(val *WorkspaceAdminUsers) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkspaceAdminUsers) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkspaceAdminUsers) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkspaceAdminUsers(val *WorkspaceAdminUsers) *NullableWorkspaceAdminUsers {
	return &NullableWorkspaceAdminUsers{value: val, isSet: true}
}

func (v NullableWorkspaceAdminUsers) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkspaceAdminUsers) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

