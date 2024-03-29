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

// checks if the CreateWorkspaceInvitation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateWorkspaceInvitation{}

// CreateWorkspaceInvitation struct for CreateWorkspaceInvitation
type CreateWorkspaceInvitation struct {
	// The email address of the user that the invitation is meant for. Only a user with that email address can accept it.
	Email string `json:"email"`
	// The permissions that the user is going to get within the workspace after accepting the invitation.
	Permissions *string `json:"permissions,omitempty"`
	// An optional message that the invitor can provide. This will be visible to the receiver of the invitation.
	Message *string `json:"message,omitempty"`
	// The base URL where the user can publicly accept his invitation.The accept token is going to be appended to the base_url (base_url '/token').
	BaseUrl string `json:"base_url"`
}

// NewCreateWorkspaceInvitation instantiates a new CreateWorkspaceInvitation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateWorkspaceInvitation(email string, baseUrl string) *CreateWorkspaceInvitation {
	this := CreateWorkspaceInvitation{}
	this.Email = email
	this.BaseUrl = baseUrl
	return &this
}

// NewCreateWorkspaceInvitationWithDefaults instantiates a new CreateWorkspaceInvitation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateWorkspaceInvitationWithDefaults() *CreateWorkspaceInvitation {
	this := CreateWorkspaceInvitation{}
	return &this
}

// GetEmail returns the Email field value
func (o *CreateWorkspaceInvitation) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *CreateWorkspaceInvitation) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *CreateWorkspaceInvitation) SetEmail(v string) {
	o.Email = v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *CreateWorkspaceInvitation) GetPermissions() string {
	if o == nil || IsNil(o.Permissions) {
		var ret string
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateWorkspaceInvitation) GetPermissionsOk() (*string, bool) {
	if o == nil || IsNil(o.Permissions) {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *CreateWorkspaceInvitation) HasPermissions() bool {
	if o != nil && !IsNil(o.Permissions) {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given string and assigns it to the Permissions field.
func (o *CreateWorkspaceInvitation) SetPermissions(v string) {
	o.Permissions = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *CreateWorkspaceInvitation) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateWorkspaceInvitation) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *CreateWorkspaceInvitation) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *CreateWorkspaceInvitation) SetMessage(v string) {
	o.Message = &v
}

// GetBaseUrl returns the BaseUrl field value
func (o *CreateWorkspaceInvitation) GetBaseUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BaseUrl
}

// GetBaseUrlOk returns a tuple with the BaseUrl field value
// and a boolean to check if the value has been set.
func (o *CreateWorkspaceInvitation) GetBaseUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BaseUrl, true
}

// SetBaseUrl sets field value
func (o *CreateWorkspaceInvitation) SetBaseUrl(v string) {
	o.BaseUrl = v
}

func (o CreateWorkspaceInvitation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateWorkspaceInvitation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["email"] = o.Email
	if !IsNil(o.Permissions) {
		toSerialize["permissions"] = o.Permissions
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	toSerialize["base_url"] = o.BaseUrl
	return toSerialize, nil
}

type NullableCreateWorkspaceInvitation struct {
	value *CreateWorkspaceInvitation
	isSet bool
}

func (v NullableCreateWorkspaceInvitation) Get() *CreateWorkspaceInvitation {
	return v.value
}

func (v *NullableCreateWorkspaceInvitation) Set(val *CreateWorkspaceInvitation) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateWorkspaceInvitation) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateWorkspaceInvitation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateWorkspaceInvitation(val *CreateWorkspaceInvitation) *NullableCreateWorkspaceInvitation {
	return &NullableCreateWorkspaceInvitation{value: val, isSet: true}
}

func (v NullableCreateWorkspaceInvitation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateWorkspaceInvitation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


