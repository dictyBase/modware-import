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

// checks if the SingleDuplicateFieldJobType type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SingleDuplicateFieldJobType{}

// SingleDuplicateFieldJobType struct for SingleDuplicateFieldJobType
type SingleDuplicateFieldJobType struct {
	Id int32 `json:"id"`
	// The type of the job.
	Type string `json:"type"`
	// A percentage indicating how far along the job is. 100 means that it's finished.
	ProgressPercentage int32 `json:"progress_percentage"`
	// Indicates the state of the import job.
	State string `json:"state"`
	// A human readable error message indicating what went wrong.
	HumanReadableError *string `json:"human_readable_error,omitempty"`
	OriginalField Field `json:"original_field"`
	DuplicatedField FieldSerializerWithRelatedFields `json:"duplicated_field"`
}

// NewSingleDuplicateFieldJobType instantiates a new SingleDuplicateFieldJobType object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSingleDuplicateFieldJobType(id int32, type_ string, progressPercentage int32, state string, originalField Field, duplicatedField FieldSerializerWithRelatedFields) *SingleDuplicateFieldJobType {
	this := SingleDuplicateFieldJobType{}
	this.Id = id
	this.Type = type_
	this.ProgressPercentage = progressPercentage
	this.State = state
	this.OriginalField = originalField
	this.DuplicatedField = duplicatedField
	return &this
}

// NewSingleDuplicateFieldJobTypeWithDefaults instantiates a new SingleDuplicateFieldJobType object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSingleDuplicateFieldJobTypeWithDefaults() *SingleDuplicateFieldJobType {
	this := SingleDuplicateFieldJobType{}
	return &this
}

// GetId returns the Id field value
func (o *SingleDuplicateFieldJobType) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *SingleDuplicateFieldJobType) SetId(v int32) {
	o.Id = v
}

// GetType returns the Type field value
func (o *SingleDuplicateFieldJobType) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *SingleDuplicateFieldJobType) SetType(v string) {
	o.Type = v
}

// GetProgressPercentage returns the ProgressPercentage field value
func (o *SingleDuplicateFieldJobType) GetProgressPercentage() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ProgressPercentage
}

// GetProgressPercentageOk returns a tuple with the ProgressPercentage field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetProgressPercentageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProgressPercentage, true
}

// SetProgressPercentage sets field value
func (o *SingleDuplicateFieldJobType) SetProgressPercentage(v int32) {
	o.ProgressPercentage = v
}

// GetState returns the State field value
func (o *SingleDuplicateFieldJobType) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *SingleDuplicateFieldJobType) SetState(v string) {
	o.State = v
}

// GetHumanReadableError returns the HumanReadableError field value if set, zero value otherwise.
func (o *SingleDuplicateFieldJobType) GetHumanReadableError() string {
	if o == nil || IsNil(o.HumanReadableError) {
		var ret string
		return ret
	}
	return *o.HumanReadableError
}

// GetHumanReadableErrorOk returns a tuple with the HumanReadableError field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetHumanReadableErrorOk() (*string, bool) {
	if o == nil || IsNil(o.HumanReadableError) {
		return nil, false
	}
	return o.HumanReadableError, true
}

// HasHumanReadableError returns a boolean if a field has been set.
func (o *SingleDuplicateFieldJobType) HasHumanReadableError() bool {
	if o != nil && !IsNil(o.HumanReadableError) {
		return true
	}

	return false
}

// SetHumanReadableError gets a reference to the given string and assigns it to the HumanReadableError field.
func (o *SingleDuplicateFieldJobType) SetHumanReadableError(v string) {
	o.HumanReadableError = &v
}

// GetOriginalField returns the OriginalField field value
func (o *SingleDuplicateFieldJobType) GetOriginalField() Field {
	if o == nil {
		var ret Field
		return ret
	}

	return o.OriginalField
}

// GetOriginalFieldOk returns a tuple with the OriginalField field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetOriginalFieldOk() (*Field, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OriginalField, true
}

// SetOriginalField sets field value
func (o *SingleDuplicateFieldJobType) SetOriginalField(v Field) {
	o.OriginalField = v
}

// GetDuplicatedField returns the DuplicatedField field value
func (o *SingleDuplicateFieldJobType) GetDuplicatedField() FieldSerializerWithRelatedFields {
	if o == nil {
		var ret FieldSerializerWithRelatedFields
		return ret
	}

	return o.DuplicatedField
}

// GetDuplicatedFieldOk returns a tuple with the DuplicatedField field value
// and a boolean to check if the value has been set.
func (o *SingleDuplicateFieldJobType) GetDuplicatedFieldOk() (*FieldSerializerWithRelatedFields, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DuplicatedField, true
}

// SetDuplicatedField sets field value
func (o *SingleDuplicateFieldJobType) SetDuplicatedField(v FieldSerializerWithRelatedFields) {
	o.DuplicatedField = v
}

func (o SingleDuplicateFieldJobType) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SingleDuplicateFieldJobType) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: type is readOnly
	toSerialize["progress_percentage"] = o.ProgressPercentage
	toSerialize["state"] = o.State
	if !IsNil(o.HumanReadableError) {
		toSerialize["human_readable_error"] = o.HumanReadableError
	}
	// skip: original_field is readOnly
	// skip: duplicated_field is readOnly
	return toSerialize, nil
}

type NullableSingleDuplicateFieldJobType struct {
	value *SingleDuplicateFieldJobType
	isSet bool
}

func (v NullableSingleDuplicateFieldJobType) Get() *SingleDuplicateFieldJobType {
	return v.value
}

func (v *NullableSingleDuplicateFieldJobType) Set(val *SingleDuplicateFieldJobType) {
	v.value = val
	v.isSet = true
}

func (v NullableSingleDuplicateFieldJobType) IsSet() bool {
	return v.isSet
}

func (v *NullableSingleDuplicateFieldJobType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSingleDuplicateFieldJobType(val *SingleDuplicateFieldJobType) *NullableSingleDuplicateFieldJobType {
	return &NullableSingleDuplicateFieldJobType{value: val, isSet: true}
}

func (v NullableSingleDuplicateFieldJobType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSingleDuplicateFieldJobType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


