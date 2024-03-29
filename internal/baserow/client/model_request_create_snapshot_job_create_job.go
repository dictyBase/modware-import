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

// checks if the RequestCreateSnapshotJobCreateJob type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RequestCreateSnapshotJobCreateJob{}

// RequestCreateSnapshotJobCreateJob struct for RequestCreateSnapshotJobCreateJob
type RequestCreateSnapshotJobCreateJob struct {
	Type Type4afEnum `json:"type"`
}

// NewRequestCreateSnapshotJobCreateJob instantiates a new RequestCreateSnapshotJobCreateJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequestCreateSnapshotJobCreateJob(type_ Type4afEnum) *RequestCreateSnapshotJobCreateJob {
	this := RequestCreateSnapshotJobCreateJob{}
	this.Type = type_
	return &this
}

// NewRequestCreateSnapshotJobCreateJobWithDefaults instantiates a new RequestCreateSnapshotJobCreateJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequestCreateSnapshotJobCreateJobWithDefaults() *RequestCreateSnapshotJobCreateJob {
	this := RequestCreateSnapshotJobCreateJob{}
	return &this
}

// GetType returns the Type field value
func (o *RequestCreateSnapshotJobCreateJob) GetType() Type4afEnum {
	if o == nil {
		var ret Type4afEnum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *RequestCreateSnapshotJobCreateJob) GetTypeOk() (*Type4afEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *RequestCreateSnapshotJobCreateJob) SetType(v Type4afEnum) {
	o.Type = v
}

func (o RequestCreateSnapshotJobCreateJob) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RequestCreateSnapshotJobCreateJob) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["type"] = o.Type
	return toSerialize, nil
}

type NullableRequestCreateSnapshotJobCreateJob struct {
	value *RequestCreateSnapshotJobCreateJob
	isSet bool
}

func (v NullableRequestCreateSnapshotJobCreateJob) Get() *RequestCreateSnapshotJobCreateJob {
	return v.value
}

func (v *NullableRequestCreateSnapshotJobCreateJob) Set(val *RequestCreateSnapshotJobCreateJob) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestCreateSnapshotJobCreateJob) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestCreateSnapshotJobCreateJob) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestCreateSnapshotJobCreateJob(val *RequestCreateSnapshotJobCreateJob) *NullableRequestCreateSnapshotJobCreateJob {
	return &NullableRequestCreateSnapshotJobCreateJob{value: val, isSet: true}
}

func (v NullableRequestCreateSnapshotJobCreateJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestCreateSnapshotJobCreateJob) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


