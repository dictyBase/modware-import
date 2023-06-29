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

// checks if the WorkspaceUser type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkspaceUser{}

// WorkspaceUser struct for WorkspaceUser
type WorkspaceUser struct {
	Id int32 `json:"id"`
	// User defined name.
	Name string `json:"name"`
	// User email.
	Email string `json:"email"`
	// DEPRECATED: Please use the functionally identical `workspace` instead as this field is being removed in the future.
	Group int32 `json:"group"`
	// The workspace that the user has access to.
	Workspace int32 `json:"workspace"`
	// The permissions that the user has within the workspace.
	Permissions *string `json:"permissions,omitempty"`
	CreatedOn time.Time `json:"created_on"`
	// The user that has access to the workspace.
	UserId int32 `json:"user_id"`
	// True if user account is pending deletion.
	ToBeDeleted bool `json:"to_be_deleted"`
}

// NewWorkspaceUser instantiates a new WorkspaceUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkspaceUser(id int32, name string, email string, group int32, workspace int32, createdOn time.Time, userId int32, toBeDeleted bool) *WorkspaceUser {
	this := WorkspaceUser{}
	this.Id = id
	this.Name = name
	this.Email = email
	this.Group = group
	this.Workspace = workspace
	this.CreatedOn = createdOn
	this.UserId = userId
	this.ToBeDeleted = toBeDeleted
	return &this
}

// NewWorkspaceUserWithDefaults instantiates a new WorkspaceUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkspaceUserWithDefaults() *WorkspaceUser {
	this := WorkspaceUser{}
	return &this
}

// GetId returns the Id field value
func (o *WorkspaceUser) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkspaceUser) SetId(v int32) {
	o.Id = v
}

// GetName returns the Name field value
func (o *WorkspaceUser) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WorkspaceUser) SetName(v string) {
	o.Name = v
}

// GetEmail returns the Email field value
func (o *WorkspaceUser) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *WorkspaceUser) SetEmail(v string) {
	o.Email = v
}

// GetGroup returns the Group field value
func (o *WorkspaceUser) GetGroup() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Group
}

// GetGroupOk returns a tuple with the Group field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetGroupOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Group, true
}

// SetGroup sets field value
func (o *WorkspaceUser) SetGroup(v int32) {
	o.Group = v
}

// GetWorkspace returns the Workspace field value
func (o *WorkspaceUser) GetWorkspace() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Workspace
}

// GetWorkspaceOk returns a tuple with the Workspace field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetWorkspaceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Workspace, true
}

// SetWorkspace sets field value
func (o *WorkspaceUser) SetWorkspace(v int32) {
	o.Workspace = v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *WorkspaceUser) GetPermissions() string {
	if o == nil || IsNil(o.Permissions) {
		var ret string
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetPermissionsOk() (*string, bool) {
	if o == nil || IsNil(o.Permissions) {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *WorkspaceUser) HasPermissions() bool {
	if o != nil && !IsNil(o.Permissions) {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given string and assigns it to the Permissions field.
func (o *WorkspaceUser) SetPermissions(v string) {
	o.Permissions = &v
}

// GetCreatedOn returns the CreatedOn field value
func (o *WorkspaceUser) GetCreatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedOn
}

// GetCreatedOnOk returns a tuple with the CreatedOn field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetCreatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedOn, true
}

// SetCreatedOn sets field value
func (o *WorkspaceUser) SetCreatedOn(v time.Time) {
	o.CreatedOn = v
}

// GetUserId returns the UserId field value
func (o *WorkspaceUser) GetUserId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetUserIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UserId, true
}

// SetUserId sets field value
func (o *WorkspaceUser) SetUserId(v int32) {
	o.UserId = v
}

// GetToBeDeleted returns the ToBeDeleted field value
func (o *WorkspaceUser) GetToBeDeleted() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ToBeDeleted
}

// GetToBeDeletedOk returns a tuple with the ToBeDeleted field value
// and a boolean to check if the value has been set.
func (o *WorkspaceUser) GetToBeDeletedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ToBeDeleted, true
}

// SetToBeDeleted sets field value
func (o *WorkspaceUser) SetToBeDeleted(v bool) {
	o.ToBeDeleted = v
}

func (o WorkspaceUser) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkspaceUser) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: name is readOnly
	// skip: email is readOnly
	toSerialize["group"] = o.Group
	toSerialize["workspace"] = o.Workspace
	if !IsNil(o.Permissions) {
		toSerialize["permissions"] = o.Permissions
	}
	// skip: created_on is readOnly
	// skip: user_id is readOnly
	toSerialize["to_be_deleted"] = o.ToBeDeleted
	return toSerialize, nil
}

type NullableWorkspaceUser struct {
	value *WorkspaceUser
	isSet bool
}

func (v NullableWorkspaceUser) Get() *WorkspaceUser {
	return v.value
}

func (v *NullableWorkspaceUser) Set(val *WorkspaceUser) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkspaceUser) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkspaceUser) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkspaceUser(val *WorkspaceUser) *NullableWorkspaceUser {
	return &NullableWorkspaceUser{value: val, isSet: true}
}

func (v NullableWorkspaceUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkspaceUser) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


