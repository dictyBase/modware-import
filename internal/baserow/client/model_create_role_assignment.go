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

// checks if the CreateRoleAssignment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateRoleAssignment{}

// CreateRoleAssignment The create role assignment serializer.
type CreateRoleAssignment struct {
	// The subject ID. A subject is an actor that can do operations.
	SubjectId int32 `json:"subject_id"`
	SubjectType SubjectType6dcEnum `json:"subject_type"`
	// The uid of the role you want to assign to the user or team in the given workspace. You can omit this property if you want to remove the role.
	Role NullableString `json:"role"`
	// The ID of the scope object. The scope object limit the role assignment to this scope and all its descendants.
	ScopeId int32 `json:"scope_id"`
	ScopeType ScopeTypeEnum `json:"scope_type"`
}

// NewCreateRoleAssignment instantiates a new CreateRoleAssignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateRoleAssignment(subjectId int32, subjectType SubjectType6dcEnum, role NullableString, scopeId int32, scopeType ScopeTypeEnum) *CreateRoleAssignment {
	this := CreateRoleAssignment{}
	this.SubjectId = subjectId
	this.SubjectType = subjectType
	this.Role = role
	this.ScopeId = scopeId
	this.ScopeType = scopeType
	return &this
}

// NewCreateRoleAssignmentWithDefaults instantiates a new CreateRoleAssignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateRoleAssignmentWithDefaults() *CreateRoleAssignment {
	this := CreateRoleAssignment{}
	return &this
}

// GetSubjectId returns the SubjectId field value
func (o *CreateRoleAssignment) GetSubjectId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.SubjectId
}

// GetSubjectIdOk returns a tuple with the SubjectId field value
// and a boolean to check if the value has been set.
func (o *CreateRoleAssignment) GetSubjectIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubjectId, true
}

// SetSubjectId sets field value
func (o *CreateRoleAssignment) SetSubjectId(v int32) {
	o.SubjectId = v
}

// GetSubjectType returns the SubjectType field value
func (o *CreateRoleAssignment) GetSubjectType() SubjectType6dcEnum {
	if o == nil {
		var ret SubjectType6dcEnum
		return ret
	}

	return o.SubjectType
}

// GetSubjectTypeOk returns a tuple with the SubjectType field value
// and a boolean to check if the value has been set.
func (o *CreateRoleAssignment) GetSubjectTypeOk() (*SubjectType6dcEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubjectType, true
}

// SetSubjectType sets field value
func (o *CreateRoleAssignment) SetSubjectType(v SubjectType6dcEnum) {
	o.SubjectType = v
}

// GetRole returns the Role field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateRoleAssignment) GetRole() string {
	if o == nil || o.Role.Get() == nil {
		var ret string
		return ret
	}

	return *o.Role.Get()
}

// GetRoleOk returns a tuple with the Role field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateRoleAssignment) GetRoleOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Role.Get(), o.Role.IsSet()
}

// SetRole sets field value
func (o *CreateRoleAssignment) SetRole(v string) {
	o.Role.Set(&v)
}

// GetScopeId returns the ScopeId field value
func (o *CreateRoleAssignment) GetScopeId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ScopeId
}

// GetScopeIdOk returns a tuple with the ScopeId field value
// and a boolean to check if the value has been set.
func (o *CreateRoleAssignment) GetScopeIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScopeId, true
}

// SetScopeId sets field value
func (o *CreateRoleAssignment) SetScopeId(v int32) {
	o.ScopeId = v
}

// GetScopeType returns the ScopeType field value
func (o *CreateRoleAssignment) GetScopeType() ScopeTypeEnum {
	if o == nil {
		var ret ScopeTypeEnum
		return ret
	}

	return o.ScopeType
}

// GetScopeTypeOk returns a tuple with the ScopeType field value
// and a boolean to check if the value has been set.
func (o *CreateRoleAssignment) GetScopeTypeOk() (*ScopeTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScopeType, true
}

// SetScopeType sets field value
func (o *CreateRoleAssignment) SetScopeType(v ScopeTypeEnum) {
	o.ScopeType = v
}

func (o CreateRoleAssignment) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateRoleAssignment) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["subject_id"] = o.SubjectId
	toSerialize["subject_type"] = o.SubjectType
	toSerialize["role"] = o.Role.Get()
	toSerialize["scope_id"] = o.ScopeId
	toSerialize["scope_type"] = o.ScopeType
	return toSerialize, nil
}

type NullableCreateRoleAssignment struct {
	value *CreateRoleAssignment
	isSet bool
}

func (v NullableCreateRoleAssignment) Get() *CreateRoleAssignment {
	return v.value
}

func (v *NullableCreateRoleAssignment) Set(val *CreateRoleAssignment) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateRoleAssignment) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateRoleAssignment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateRoleAssignment(val *CreateRoleAssignment) *NullableCreateRoleAssignment {
	return &NullableCreateRoleAssignment{value: val, isSet: true}
}

func (v NullableCreateRoleAssignment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateRoleAssignment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


