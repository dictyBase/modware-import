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

// checks if the PaginationSerializerRowComment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PaginationSerializerRowComment{}

// PaginationSerializerRowComment struct for PaginationSerializerRowComment
type PaginationSerializerRowComment struct {
	// The total amount of results.
	Count int32 `json:"count"`
	// URL to the next page.
	Next NullableString `json:"next"`
	// URL to the previous page.
	Previous NullableString `json:"previous"`
	Results []RowComment `json:"results"`
}

// NewPaginationSerializerRowComment instantiates a new PaginationSerializerRowComment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginationSerializerRowComment(count int32, next NullableString, previous NullableString, results []RowComment) *PaginationSerializerRowComment {
	this := PaginationSerializerRowComment{}
	this.Count = count
	this.Next = next
	this.Previous = previous
	this.Results = results
	return &this
}

// NewPaginationSerializerRowCommentWithDefaults instantiates a new PaginationSerializerRowComment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginationSerializerRowCommentWithDefaults() *PaginationSerializerRowComment {
	this := PaginationSerializerRowComment{}
	return &this
}

// GetCount returns the Count field value
func (o *PaginationSerializerRowComment) GetCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Count
}

// GetCountOk returns a tuple with the Count field value
// and a boolean to check if the value has been set.
func (o *PaginationSerializerRowComment) GetCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Count, true
}

// SetCount sets field value
func (o *PaginationSerializerRowComment) SetCount(v int32) {
	o.Count = v
}

// GetNext returns the Next field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PaginationSerializerRowComment) GetNext() string {
	if o == nil || o.Next.Get() == nil {
		var ret string
		return ret
	}

	return *o.Next.Get()
}

// GetNextOk returns a tuple with the Next field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PaginationSerializerRowComment) GetNextOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Next.Get(), o.Next.IsSet()
}

// SetNext sets field value
func (o *PaginationSerializerRowComment) SetNext(v string) {
	o.Next.Set(&v)
}

// GetPrevious returns the Previous field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PaginationSerializerRowComment) GetPrevious() string {
	if o == nil || o.Previous.Get() == nil {
		var ret string
		return ret
	}

	return *o.Previous.Get()
}

// GetPreviousOk returns a tuple with the Previous field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PaginationSerializerRowComment) GetPreviousOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Previous.Get(), o.Previous.IsSet()
}

// SetPrevious sets field value
func (o *PaginationSerializerRowComment) SetPrevious(v string) {
	o.Previous.Set(&v)
}

// GetResults returns the Results field value
func (o *PaginationSerializerRowComment) GetResults() []RowComment {
	if o == nil {
		var ret []RowComment
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginationSerializerRowComment) GetResultsOk() ([]RowComment, bool) {
	if o == nil {
		return nil, false
	}
	return o.Results, true
}

// SetResults sets field value
func (o *PaginationSerializerRowComment) SetResults(v []RowComment) {
	o.Results = v
}

func (o PaginationSerializerRowComment) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PaginationSerializerRowComment) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["count"] = o.Count
	toSerialize["next"] = o.Next.Get()
	toSerialize["previous"] = o.Previous.Get()
	toSerialize["results"] = o.Results
	return toSerialize, nil
}

type NullablePaginationSerializerRowComment struct {
	value *PaginationSerializerRowComment
	isSet bool
}

func (v NullablePaginationSerializerRowComment) Get() *PaginationSerializerRowComment {
	return v.value
}

func (v *NullablePaginationSerializerRowComment) Set(val *PaginationSerializerRowComment) {
	v.value = val
	v.isSet = true
}

func (v NullablePaginationSerializerRowComment) IsSet() bool {
	return v.isSet
}

func (v *NullablePaginationSerializerRowComment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaginationSerializerRowComment(val *PaginationSerializerRowComment) *NullablePaginationSerializerRowComment {
	return &NullablePaginationSerializerRowComment{value: val, isSet: true}
}

func (v NullablePaginationSerializerRowComment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaginationSerializerRowComment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

