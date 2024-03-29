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

// checks if the UserAdminResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserAdminResponse{}

// UserAdminResponse Serializes the safe user attributes to expose for a response back to the user.
type UserAdminResponse struct {
	Id int32 `json:"id"`
	Username string `json:"username"`
	Name string `json:"name"`
	Groups []UserAdminGroups `json:"groups"`
	Workspaces []UserAdminGroups `json:"workspaces"`
	LastLogin NullableTime `json:"last_login,omitempty"`
	DateJoined *time.Time `json:"date_joined,omitempty"`
	// Designates whether this user should be treated as active. Set this to false instead of deleting accounts.
	IsActive *bool `json:"is_active,omitempty"`
	// Designates whether this user is an admin and has access to all workspaces and Baserow's admin areas. 
	IsStaff *bool `json:"is_staff,omitempty"`
}

// NewUserAdminResponse instantiates a new UserAdminResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAdminResponse(id int32, username string, name string, groups []UserAdminGroups, workspaces []UserAdminGroups) *UserAdminResponse {
	this := UserAdminResponse{}
	this.Id = id
	this.Username = username
	this.Name = name
	this.Groups = groups
	this.Workspaces = workspaces
	return &this
}

// NewUserAdminResponseWithDefaults instantiates a new UserAdminResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAdminResponseWithDefaults() *UserAdminResponse {
	this := UserAdminResponse{}
	return &this
}

// GetId returns the Id field value
func (o *UserAdminResponse) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *UserAdminResponse) SetId(v int32) {
	o.Id = v
}

// GetUsername returns the Username field value
func (o *UserAdminResponse) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *UserAdminResponse) SetUsername(v string) {
	o.Username = v
}

// GetName returns the Name field value
func (o *UserAdminResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *UserAdminResponse) SetName(v string) {
	o.Name = v
}

// GetGroups returns the Groups field value
func (o *UserAdminResponse) GetGroups() []UserAdminGroups {
	if o == nil {
		var ret []UserAdminGroups
		return ret
	}

	return o.Groups
}

// GetGroupsOk returns a tuple with the Groups field value
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetGroupsOk() ([]UserAdminGroups, bool) {
	if o == nil {
		return nil, false
	}
	return o.Groups, true
}

// SetGroups sets field value
func (o *UserAdminResponse) SetGroups(v []UserAdminGroups) {
	o.Groups = v
}

// GetWorkspaces returns the Workspaces field value
func (o *UserAdminResponse) GetWorkspaces() []UserAdminGroups {
	if o == nil {
		var ret []UserAdminGroups
		return ret
	}

	return o.Workspaces
}

// GetWorkspacesOk returns a tuple with the Workspaces field value
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetWorkspacesOk() ([]UserAdminGroups, bool) {
	if o == nil {
		return nil, false
	}
	return o.Workspaces, true
}

// SetWorkspaces sets field value
func (o *UserAdminResponse) SetWorkspaces(v []UserAdminGroups) {
	o.Workspaces = v
}

// GetLastLogin returns the LastLogin field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UserAdminResponse) GetLastLogin() time.Time {
	if o == nil || IsNil(o.LastLogin.Get()) {
		var ret time.Time
		return ret
	}
	return *o.LastLogin.Get()
}

// GetLastLoginOk returns a tuple with the LastLogin field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserAdminResponse) GetLastLoginOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.LastLogin.Get(), o.LastLogin.IsSet()
}

// HasLastLogin returns a boolean if a field has been set.
func (o *UserAdminResponse) HasLastLogin() bool {
	if o != nil && o.LastLogin.IsSet() {
		return true
	}

	return false
}

// SetLastLogin gets a reference to the given NullableTime and assigns it to the LastLogin field.
func (o *UserAdminResponse) SetLastLogin(v time.Time) {
	o.LastLogin.Set(&v)
}
// SetLastLoginNil sets the value for LastLogin to be an explicit nil
func (o *UserAdminResponse) SetLastLoginNil() {
	o.LastLogin.Set(nil)
}

// UnsetLastLogin ensures that no value is present for LastLogin, not even an explicit nil
func (o *UserAdminResponse) UnsetLastLogin() {
	o.LastLogin.Unset()
}

// GetDateJoined returns the DateJoined field value if set, zero value otherwise.
func (o *UserAdminResponse) GetDateJoined() time.Time {
	if o == nil || IsNil(o.DateJoined) {
		var ret time.Time
		return ret
	}
	return *o.DateJoined
}

// GetDateJoinedOk returns a tuple with the DateJoined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetDateJoinedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DateJoined) {
		return nil, false
	}
	return o.DateJoined, true
}

// HasDateJoined returns a boolean if a field has been set.
func (o *UserAdminResponse) HasDateJoined() bool {
	if o != nil && !IsNil(o.DateJoined) {
		return true
	}

	return false
}

// SetDateJoined gets a reference to the given time.Time and assigns it to the DateJoined field.
func (o *UserAdminResponse) SetDateJoined(v time.Time) {
	o.DateJoined = &v
}

// GetIsActive returns the IsActive field value if set, zero value otherwise.
func (o *UserAdminResponse) GetIsActive() bool {
	if o == nil || IsNil(o.IsActive) {
		var ret bool
		return ret
	}
	return *o.IsActive
}

// GetIsActiveOk returns a tuple with the IsActive field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetIsActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.IsActive) {
		return nil, false
	}
	return o.IsActive, true
}

// HasIsActive returns a boolean if a field has been set.
func (o *UserAdminResponse) HasIsActive() bool {
	if o != nil && !IsNil(o.IsActive) {
		return true
	}

	return false
}

// SetIsActive gets a reference to the given bool and assigns it to the IsActive field.
func (o *UserAdminResponse) SetIsActive(v bool) {
	o.IsActive = &v
}

// GetIsStaff returns the IsStaff field value if set, zero value otherwise.
func (o *UserAdminResponse) GetIsStaff() bool {
	if o == nil || IsNil(o.IsStaff) {
		var ret bool
		return ret
	}
	return *o.IsStaff
}

// GetIsStaffOk returns a tuple with the IsStaff field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAdminResponse) GetIsStaffOk() (*bool, bool) {
	if o == nil || IsNil(o.IsStaff) {
		return nil, false
	}
	return o.IsStaff, true
}

// HasIsStaff returns a boolean if a field has been set.
func (o *UserAdminResponse) HasIsStaff() bool {
	if o != nil && !IsNil(o.IsStaff) {
		return true
	}

	return false
}

// SetIsStaff gets a reference to the given bool and assigns it to the IsStaff field.
func (o *UserAdminResponse) SetIsStaff(v bool) {
	o.IsStaff = &v
}

func (o UserAdminResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserAdminResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["username"] = o.Username
	toSerialize["name"] = o.Name
	toSerialize["groups"] = o.Groups
	toSerialize["workspaces"] = o.Workspaces
	if o.LastLogin.IsSet() {
		toSerialize["last_login"] = o.LastLogin.Get()
	}
	if !IsNil(o.DateJoined) {
		toSerialize["date_joined"] = o.DateJoined
	}
	if !IsNil(o.IsActive) {
		toSerialize["is_active"] = o.IsActive
	}
	if !IsNil(o.IsStaff) {
		toSerialize["is_staff"] = o.IsStaff
	}
	return toSerialize, nil
}

type NullableUserAdminResponse struct {
	value *UserAdminResponse
	isSet bool
}

func (v NullableUserAdminResponse) Get() *UserAdminResponse {
	return v.value
}

func (v *NullableUserAdminResponse) Set(val *UserAdminResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableUserAdminResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableUserAdminResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserAdminResponse(val *UserAdminResponse) *NullableUserAdminResponse {
	return &NullableUserAdminResponse{value: val, isSet: true}
}

func (v NullableUserAdminResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserAdminResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


