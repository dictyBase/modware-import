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

// checks if the PatchedUpdateWorkspaceInvitation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedUpdateWorkspaceInvitation{}

// PatchedUpdateWorkspaceInvitation struct for PatchedUpdateWorkspaceInvitation
type PatchedUpdateWorkspaceInvitation struct {
	// The permissions that the user is going to get within the workspace after accepting the invitation.
	Permissions *string `json:"permissions,omitempty"`
}

// NewPatchedUpdateWorkspaceInvitation instantiates a new PatchedUpdateWorkspaceInvitation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedUpdateWorkspaceInvitation() *PatchedUpdateWorkspaceInvitation {
	this := PatchedUpdateWorkspaceInvitation{}
	return &this
}

// NewPatchedUpdateWorkspaceInvitationWithDefaults instantiates a new PatchedUpdateWorkspaceInvitation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedUpdateWorkspaceInvitationWithDefaults() *PatchedUpdateWorkspaceInvitation {
	this := PatchedUpdateWorkspaceInvitation{}
	return &this
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *PatchedUpdateWorkspaceInvitation) GetPermissions() string {
	if o == nil || IsNil(o.Permissions) {
		var ret string
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUpdateWorkspaceInvitation) GetPermissionsOk() (*string, bool) {
	if o == nil || IsNil(o.Permissions) {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *PatchedUpdateWorkspaceInvitation) HasPermissions() bool {
	if o != nil && !IsNil(o.Permissions) {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given string and assigns it to the Permissions field.
func (o *PatchedUpdateWorkspaceInvitation) SetPermissions(v string) {
	o.Permissions = &v
}

func (o PatchedUpdateWorkspaceInvitation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedUpdateWorkspaceInvitation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Permissions) {
		toSerialize["permissions"] = o.Permissions
	}
	return toSerialize, nil
}

type NullablePatchedUpdateWorkspaceInvitation struct {
	value *PatchedUpdateWorkspaceInvitation
	isSet bool
}

func (v NullablePatchedUpdateWorkspaceInvitation) Get() *PatchedUpdateWorkspaceInvitation {
	return v.value
}

func (v *NullablePatchedUpdateWorkspaceInvitation) Set(val *PatchedUpdateWorkspaceInvitation) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedUpdateWorkspaceInvitation) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedUpdateWorkspaceInvitation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedUpdateWorkspaceInvitation(val *PatchedUpdateWorkspaceInvitation) *NullablePatchedUpdateWorkspaceInvitation {
	return &NullablePatchedUpdateWorkspaceInvitation{value: val, isSet: true}
}

func (v NullablePatchedUpdateWorkspaceInvitation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedUpdateWorkspaceInvitation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

