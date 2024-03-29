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

// checks if the AirtableImportJobJob type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AirtableImportJobJob{}

// AirtableImportJobJob struct for AirtableImportJobJob
type AirtableImportJobJob struct {
	Id int32 `json:"id"`
	// The type of the job.
	Type string `json:"type"`
	// A percentage indicating how far along the job is. 100 means that it's finished.
	ProgressPercentage int32 `json:"progress_percentage"`
	// Indicates the state of the import job.
	State string `json:"state"`
	// A human readable error message indicating what went wrong.
	HumanReadableError *string `json:"human_readable_error,omitempty"`
	// The group ID where the Airtable base must be imported into.
	GroupId int32 `json:"group_id"`
	// The workspace ID where the Airtable base must be imported into.
	WorkspaceId int32 `json:"workspace_id"`
	Database Application `json:"database"`
	// Public ID of the shared Airtable base that must be imported.
	AirtableShareId string `json:"airtable_share_id"`
}

// NewAirtableImportJobJob instantiates a new AirtableImportJobJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAirtableImportJobJob(id int32, type_ string, progressPercentage int32, state string, groupId int32, workspaceId int32, database Application, airtableShareId string) *AirtableImportJobJob {
	this := AirtableImportJobJob{}
	this.Id = id
	this.Type = type_
	this.ProgressPercentage = progressPercentage
	this.State = state
	this.GroupId = groupId
	this.WorkspaceId = workspaceId
	this.Database = database
	this.AirtableShareId = airtableShareId
	return &this
}

// NewAirtableImportJobJobWithDefaults instantiates a new AirtableImportJobJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAirtableImportJobJobWithDefaults() *AirtableImportJobJob {
	this := AirtableImportJobJob{}
	return &this
}

// GetId returns the Id field value
func (o *AirtableImportJobJob) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AirtableImportJobJob) SetId(v int32) {
	o.Id = v
}

// GetType returns the Type field value
func (o *AirtableImportJobJob) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *AirtableImportJobJob) SetType(v string) {
	o.Type = v
}

// GetProgressPercentage returns the ProgressPercentage field value
func (o *AirtableImportJobJob) GetProgressPercentage() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ProgressPercentage
}

// GetProgressPercentageOk returns a tuple with the ProgressPercentage field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetProgressPercentageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProgressPercentage, true
}

// SetProgressPercentage sets field value
func (o *AirtableImportJobJob) SetProgressPercentage(v int32) {
	o.ProgressPercentage = v
}

// GetState returns the State field value
func (o *AirtableImportJobJob) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *AirtableImportJobJob) SetState(v string) {
	o.State = v
}

// GetHumanReadableError returns the HumanReadableError field value if set, zero value otherwise.
func (o *AirtableImportJobJob) GetHumanReadableError() string {
	if o == nil || IsNil(o.HumanReadableError) {
		var ret string
		return ret
	}
	return *o.HumanReadableError
}

// GetHumanReadableErrorOk returns a tuple with the HumanReadableError field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetHumanReadableErrorOk() (*string, bool) {
	if o == nil || IsNil(o.HumanReadableError) {
		return nil, false
	}
	return o.HumanReadableError, true
}

// HasHumanReadableError returns a boolean if a field has been set.
func (o *AirtableImportJobJob) HasHumanReadableError() bool {
	if o != nil && !IsNil(o.HumanReadableError) {
		return true
	}

	return false
}

// SetHumanReadableError gets a reference to the given string and assigns it to the HumanReadableError field.
func (o *AirtableImportJobJob) SetHumanReadableError(v string) {
	o.HumanReadableError = &v
}

// GetGroupId returns the GroupId field value
func (o *AirtableImportJobJob) GetGroupId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetGroupIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupId, true
}

// SetGroupId sets field value
func (o *AirtableImportJobJob) SetGroupId(v int32) {
	o.GroupId = v
}

// GetWorkspaceId returns the WorkspaceId field value
func (o *AirtableImportJobJob) GetWorkspaceId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.WorkspaceId
}

// GetWorkspaceIdOk returns a tuple with the WorkspaceId field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetWorkspaceIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkspaceId, true
}

// SetWorkspaceId sets field value
func (o *AirtableImportJobJob) SetWorkspaceId(v int32) {
	o.WorkspaceId = v
}

// GetDatabase returns the Database field value
func (o *AirtableImportJobJob) GetDatabase() Application {
	if o == nil {
		var ret Application
		return ret
	}

	return o.Database
}

// GetDatabaseOk returns a tuple with the Database field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetDatabaseOk() (*Application, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Database, true
}

// SetDatabase sets field value
func (o *AirtableImportJobJob) SetDatabase(v Application) {
	o.Database = v
}

// GetAirtableShareId returns the AirtableShareId field value
func (o *AirtableImportJobJob) GetAirtableShareId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AirtableShareId
}

// GetAirtableShareIdOk returns a tuple with the AirtableShareId field value
// and a boolean to check if the value has been set.
func (o *AirtableImportJobJob) GetAirtableShareIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AirtableShareId, true
}

// SetAirtableShareId sets field value
func (o *AirtableImportJobJob) SetAirtableShareId(v string) {
	o.AirtableShareId = v
}

func (o AirtableImportJobJob) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AirtableImportJobJob) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: type is readOnly
	toSerialize["progress_percentage"] = o.ProgressPercentage
	toSerialize["state"] = o.State
	if !IsNil(o.HumanReadableError) {
		toSerialize["human_readable_error"] = o.HumanReadableError
	}
	toSerialize["group_id"] = o.GroupId
	toSerialize["workspace_id"] = o.WorkspaceId
	toSerialize["database"] = o.Database
	toSerialize["airtable_share_id"] = o.AirtableShareId
	return toSerialize, nil
}

type NullableAirtableImportJobJob struct {
	value *AirtableImportJobJob
	isSet bool
}

func (v NullableAirtableImportJobJob) Get() *AirtableImportJobJob {
	return v.value
}

func (v *NullableAirtableImportJobJob) Set(val *AirtableImportJobJob) {
	v.value = val
	v.isSet = true
}

func (v NullableAirtableImportJobJob) IsSet() bool {
	return v.isSet
}

func (v *NullableAirtableImportJobJob) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAirtableImportJobJob(val *AirtableImportJobJob) *NullableAirtableImportJobJob {
	return &NullableAirtableImportJobJob{value: val, isSet: true}
}

func (v NullableAirtableImportJobJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAirtableImportJobJob) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


