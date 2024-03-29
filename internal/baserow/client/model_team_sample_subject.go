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

// checks if the TeamSampleSubject type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TeamSampleSubject{}

// TeamSampleSubject struct for TeamSampleSubject
type TeamSampleSubject struct {
	// The subject's unique identifier.
	SubjectId int32 `json:"subject_id"`
	SubjectType SubjectType3ffEnum `json:"subject_type"`
	// The subject's label. Defaults to a user's first name.
	SubjectLabel string `json:"subject_label"`
	// The team subject unique identifier.
	TeamSubjectId int32 `json:"team_subject_id"`
}

// NewTeamSampleSubject instantiates a new TeamSampleSubject object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTeamSampleSubject(subjectId int32, subjectType SubjectType3ffEnum, subjectLabel string, teamSubjectId int32) *TeamSampleSubject {
	this := TeamSampleSubject{}
	this.SubjectId = subjectId
	this.SubjectType = subjectType
	this.SubjectLabel = subjectLabel
	this.TeamSubjectId = teamSubjectId
	return &this
}

// NewTeamSampleSubjectWithDefaults instantiates a new TeamSampleSubject object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTeamSampleSubjectWithDefaults() *TeamSampleSubject {
	this := TeamSampleSubject{}
	return &this
}

// GetSubjectId returns the SubjectId field value
func (o *TeamSampleSubject) GetSubjectId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.SubjectId
}

// GetSubjectIdOk returns a tuple with the SubjectId field value
// and a boolean to check if the value has been set.
func (o *TeamSampleSubject) GetSubjectIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubjectId, true
}

// SetSubjectId sets field value
func (o *TeamSampleSubject) SetSubjectId(v int32) {
	o.SubjectId = v
}

// GetSubjectType returns the SubjectType field value
func (o *TeamSampleSubject) GetSubjectType() SubjectType3ffEnum {
	if o == nil {
		var ret SubjectType3ffEnum
		return ret
	}

	return o.SubjectType
}

// GetSubjectTypeOk returns a tuple with the SubjectType field value
// and a boolean to check if the value has been set.
func (o *TeamSampleSubject) GetSubjectTypeOk() (*SubjectType3ffEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubjectType, true
}

// SetSubjectType sets field value
func (o *TeamSampleSubject) SetSubjectType(v SubjectType3ffEnum) {
	o.SubjectType = v
}

// GetSubjectLabel returns the SubjectLabel field value
func (o *TeamSampleSubject) GetSubjectLabel() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SubjectLabel
}

// GetSubjectLabelOk returns a tuple with the SubjectLabel field value
// and a boolean to check if the value has been set.
func (o *TeamSampleSubject) GetSubjectLabelOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubjectLabel, true
}

// SetSubjectLabel sets field value
func (o *TeamSampleSubject) SetSubjectLabel(v string) {
	o.SubjectLabel = v
}

// GetTeamSubjectId returns the TeamSubjectId field value
func (o *TeamSampleSubject) GetTeamSubjectId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TeamSubjectId
}

// GetTeamSubjectIdOk returns a tuple with the TeamSubjectId field value
// and a boolean to check if the value has been set.
func (o *TeamSampleSubject) GetTeamSubjectIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TeamSubjectId, true
}

// SetTeamSubjectId sets field value
func (o *TeamSampleSubject) SetTeamSubjectId(v int32) {
	o.TeamSubjectId = v
}

func (o TeamSampleSubject) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TeamSampleSubject) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["subject_id"] = o.SubjectId
	toSerialize["subject_type"] = o.SubjectType
	toSerialize["subject_label"] = o.SubjectLabel
	toSerialize["team_subject_id"] = o.TeamSubjectId
	return toSerialize, nil
}

type NullableTeamSampleSubject struct {
	value *TeamSampleSubject
	isSet bool
}

func (v NullableTeamSampleSubject) Get() *TeamSampleSubject {
	return v.value
}

func (v *NullableTeamSampleSubject) Set(val *TeamSampleSubject) {
	v.value = val
	v.isSet = true
}

func (v NullableTeamSampleSubject) IsSet() bool {
	return v.isSet
}

func (v *NullableTeamSampleSubject) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTeamSampleSubject(val *TeamSampleSubject) *NullableTeamSampleSubject {
	return &NullableTeamSampleSubject{value: val, isSet: true}
}

func (v NullableTeamSampleSubject) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTeamSampleSubject) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


