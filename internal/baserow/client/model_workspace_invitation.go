/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"time"
)

// checks if the WorkspaceInvitation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkspaceInvitation{}

// WorkspaceInvitation struct for WorkspaceInvitation
type WorkspaceInvitation struct {
	Id int32 `json:"id"`
	// DEPRECATED: Please use the functionally identical `workspace` instead as this field is being removed in the future.
	Group int32 `json:"group"`
	// The workspace that the user will get access to once the invitation is accepted.
	Workspace int32 `json:"workspace"`
	// The email address of the user that the invitation is meant for. Only a user with that email address can accept it.
	Email string `json:"email"`
	// The permissions that the user is going to get within the workspace after accepting the invitation.
	Permissions *string `json:"permissions,omitempty"`
	// An optional message that the invitor can provide. This will be visible to the receiver of the invitation.
	Message *string `json:"message,omitempty"`
	CreatedOn time.Time `json:"created_on"`
}

// NewWorkspaceInvitation instantiates a new WorkspaceInvitation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkspaceInvitation(id int32, group int32, workspace int32, email string, createdOn time.Time) *WorkspaceInvitation {
	this := WorkspaceInvitation{}
	this.Id = id
	this.Group = group
	this.Workspace = workspace
	this.Email = email
	this.CreatedOn = createdOn
	return &this
}

// NewWorkspaceInvitationWithDefaults instantiates a new WorkspaceInvitation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkspaceInvitationWithDefaults() *WorkspaceInvitation {
	this := WorkspaceInvitation{}
	return &this
}

// GetId returns the Id field value
func (o *WorkspaceInvitation) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkspaceInvitation) SetId(v int32) {
	o.Id = v
}

// GetGroup returns the Group field value
func (o *WorkspaceInvitation) GetGroup() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Group
}

// GetGroupOk returns a tuple with the Group field value
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetGroupOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Group, true
}

// SetGroup sets field value
func (o *WorkspaceInvitation) SetGroup(v int32) {
	o.Group = v
}

// GetWorkspace returns the Workspace field value
func (o *WorkspaceInvitation) GetWorkspace() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Workspace
}

// GetWorkspaceOk returns a tuple with the Workspace field value
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetWorkspaceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Workspace, true
}

// SetWorkspace sets field value
func (o *WorkspaceInvitation) SetWorkspace(v int32) {
	o.Workspace = v
}

// GetEmail returns the Email field value
func (o *WorkspaceInvitation) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *WorkspaceInvitation) SetEmail(v string) {
	o.Email = v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *WorkspaceInvitation) GetPermissions() string {
	if o == nil || IsNil(o.Permissions) {
		var ret string
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetPermissionsOk() (*string, bool) {
	if o == nil || IsNil(o.Permissions) {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *WorkspaceInvitation) HasPermissions() bool {
	if o != nil && !IsNil(o.Permissions) {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given string and assigns it to the Permissions field.
func (o *WorkspaceInvitation) SetPermissions(v string) {
	o.Permissions = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *WorkspaceInvitation) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *WorkspaceInvitation) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *WorkspaceInvitation) SetMessage(v string) {
	o.Message = &v
}

// GetCreatedOn returns the CreatedOn field value
func (o *WorkspaceInvitation) GetCreatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedOn
}

// GetCreatedOnOk returns a tuple with the CreatedOn field value
// and a boolean to check if the value has been set.
func (o *WorkspaceInvitation) GetCreatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedOn, true
}

// SetCreatedOn sets field value
func (o *WorkspaceInvitation) SetCreatedOn(v time.Time) {
	o.CreatedOn = v
}

func (o WorkspaceInvitation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkspaceInvitation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["group"] = o.Group
	toSerialize["workspace"] = o.Workspace
	toSerialize["email"] = o.Email
	if !IsNil(o.Permissions) {
		toSerialize["permissions"] = o.Permissions
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	// skip: created_on is readOnly
	return toSerialize, nil
}

type NullableWorkspaceInvitation struct {
	value *WorkspaceInvitation
	isSet bool
}

func (v NullableWorkspaceInvitation) Get() *WorkspaceInvitation {
	return v.value
}

func (v *NullableWorkspaceInvitation) Set(val *WorkspaceInvitation) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkspaceInvitation) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkspaceInvitation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkspaceInvitation(val *WorkspaceInvitation) *NullableWorkspaceInvitation {
	return &NullableWorkspaceInvitation{value: val, isSet: true}
}

func (v NullableWorkspaceInvitation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkspaceInvitation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


