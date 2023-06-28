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

// checks if the PatchedUserAdminUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedUserAdminUpdate{}

// PatchedUserAdminUpdate Serializes a request body for updating a given user. Do not use for returning user data as the password will be returned also.
type PatchedUserAdminUpdate struct {
	Username *string `json:"username,omitempty"`
	Name *string `json:"name,omitempty"`
	// Designates whether this user should be treated as active. Set this to false instead of deleting accounts.
	IsActive *bool `json:"is_active,omitempty"`
	// Designates whether this user is an admin and has access to all workspaces and Baserow's admin areas. 
	IsStaff *bool `json:"is_staff,omitempty"`
	Password *string `json:"password,omitempty"`
}

// NewPatchedUserAdminUpdate instantiates a new PatchedUserAdminUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedUserAdminUpdate() *PatchedUserAdminUpdate {
	this := PatchedUserAdminUpdate{}
	return &this
}

// NewPatchedUserAdminUpdateWithDefaults instantiates a new PatchedUserAdminUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedUserAdminUpdateWithDefaults() *PatchedUserAdminUpdate {
	this := PatchedUserAdminUpdate{}
	return &this
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *PatchedUserAdminUpdate) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUserAdminUpdate) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *PatchedUserAdminUpdate) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *PatchedUserAdminUpdate) SetUsername(v string) {
	o.Username = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedUserAdminUpdate) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUserAdminUpdate) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedUserAdminUpdate) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedUserAdminUpdate) SetName(v string) {
	o.Name = &v
}

// GetIsActive returns the IsActive field value if set, zero value otherwise.
func (o *PatchedUserAdminUpdate) GetIsActive() bool {
	if o == nil || IsNil(o.IsActive) {
		var ret bool
		return ret
	}
	return *o.IsActive
}

// GetIsActiveOk returns a tuple with the IsActive field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUserAdminUpdate) GetIsActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.IsActive) {
		return nil, false
	}
	return o.IsActive, true
}

// HasIsActive returns a boolean if a field has been set.
func (o *PatchedUserAdminUpdate) HasIsActive() bool {
	if o != nil && !IsNil(o.IsActive) {
		return true
	}

	return false
}

// SetIsActive gets a reference to the given bool and assigns it to the IsActive field.
func (o *PatchedUserAdminUpdate) SetIsActive(v bool) {
	o.IsActive = &v
}

// GetIsStaff returns the IsStaff field value if set, zero value otherwise.
func (o *PatchedUserAdminUpdate) GetIsStaff() bool {
	if o == nil || IsNil(o.IsStaff) {
		var ret bool
		return ret
	}
	return *o.IsStaff
}

// GetIsStaffOk returns a tuple with the IsStaff field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUserAdminUpdate) GetIsStaffOk() (*bool, bool) {
	if o == nil || IsNil(o.IsStaff) {
		return nil, false
	}
	return o.IsStaff, true
}

// HasIsStaff returns a boolean if a field has been set.
func (o *PatchedUserAdminUpdate) HasIsStaff() bool {
	if o != nil && !IsNil(o.IsStaff) {
		return true
	}

	return false
}

// SetIsStaff gets a reference to the given bool and assigns it to the IsStaff field.
func (o *PatchedUserAdminUpdate) SetIsStaff(v bool) {
	o.IsStaff = &v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *PatchedUserAdminUpdate) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedUserAdminUpdate) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *PatchedUserAdminUpdate) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *PatchedUserAdminUpdate) SetPassword(v string) {
	o.Password = &v
}

func (o PatchedUserAdminUpdate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedUserAdminUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Username) {
		toSerialize["username"] = o.Username
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.IsActive) {
		toSerialize["is_active"] = o.IsActive
	}
	if !IsNil(o.IsStaff) {
		toSerialize["is_staff"] = o.IsStaff
	}
	if !IsNil(o.Password) {
		toSerialize["password"] = o.Password
	}
	return toSerialize, nil
}

type NullablePatchedUserAdminUpdate struct {
	value *PatchedUserAdminUpdate
	isSet bool
}

func (v NullablePatchedUserAdminUpdate) Get() *PatchedUserAdminUpdate {
	return v.value
}

func (v *NullablePatchedUserAdminUpdate) Set(val *PatchedUserAdminUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedUserAdminUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedUserAdminUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedUserAdminUpdate(val *PatchedUserAdminUpdate) *NullablePatchedUserAdminUpdate {
	return &NullablePatchedUserAdminUpdate{value: val, isSet: true}
}

func (v NullablePatchedUserAdminUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedUserAdminUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


