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

// checks if the PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse{}

// PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse struct for PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse
type PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse struct {
	// An object containing the field id as key and the properties related to view as value.
	FieldOptions *map[string]GalleryViewFieldOptions `json:"field_options,omitempty"`
	// The total amount of results.
	Count int32 `json:"count"`
	// URL to the next page.
	Next NullableString `json:"next"`
	// URL to the previous page.
	Previous NullableString `json:"previous"`
	Results []ExampleRowResponse `json:"results"`
}

// NewPaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse instantiates a new PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse(count int32, next NullableString, previous NullableString, results []ExampleRowResponse) *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse {
	this := PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse{}
	this.Count = count
	this.Next = next
	this.Previous = previous
	this.Results = results
	return &this
}

// NewPaginationSerializerWithGalleryViewFieldOptionsExampleRowResponseWithDefaults instantiates a new PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginationSerializerWithGalleryViewFieldOptionsExampleRowResponseWithDefaults() *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse {
	this := PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse{}
	return &this
}

// GetFieldOptions returns the FieldOptions field value if set, zero value otherwise.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetFieldOptions() map[string]GalleryViewFieldOptions {
	if o == nil || IsNil(o.FieldOptions) {
		var ret map[string]GalleryViewFieldOptions
		return ret
	}
	return *o.FieldOptions
}

// GetFieldOptionsOk returns a tuple with the FieldOptions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetFieldOptionsOk() (*map[string]GalleryViewFieldOptions, bool) {
	if o == nil || IsNil(o.FieldOptions) {
		return nil, false
	}
	return o.FieldOptions, true
}

// HasFieldOptions returns a boolean if a field has been set.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) HasFieldOptions() bool {
	if o != nil && !IsNil(o.FieldOptions) {
		return true
	}

	return false
}

// SetFieldOptions gets a reference to the given map[string]GalleryViewFieldOptions and assigns it to the FieldOptions field.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) SetFieldOptions(v map[string]GalleryViewFieldOptions) {
	o.FieldOptions = &v
}

// GetCount returns the Count field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Count
}

// GetCountOk returns a tuple with the Count field value
// and a boolean to check if the value has been set.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Count, true
}

// SetCount sets field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) SetCount(v int32) {
	o.Count = v
}

// GetNext returns the Next field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetNext() string {
	if o == nil || o.Next.Get() == nil {
		var ret string
		return ret
	}

	return *o.Next.Get()
}

// GetNextOk returns a tuple with the Next field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetNextOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Next.Get(), o.Next.IsSet()
}

// SetNext sets field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) SetNext(v string) {
	o.Next.Set(&v)
}

// GetPrevious returns the Previous field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetPrevious() string {
	if o == nil || o.Previous.Get() == nil {
		var ret string
		return ret
	}

	return *o.Previous.Get()
}

// GetPreviousOk returns a tuple with the Previous field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetPreviousOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Previous.Get(), o.Previous.IsSet()
}

// SetPrevious sets field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) SetPrevious(v string) {
	o.Previous.Set(&v)
}

// GetResults returns the Results field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetResults() []ExampleRowResponse {
	if o == nil {
		var ret []ExampleRowResponse
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) GetResultsOk() ([]ExampleRowResponse, bool) {
	if o == nil {
		return nil, false
	}
	return o.Results, true
}

// SetResults sets field value
func (o *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) SetResults(v []ExampleRowResponse) {
	o.Results = v
}

func (o PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.FieldOptions) {
		toSerialize["field_options"] = o.FieldOptions
	}
	toSerialize["count"] = o.Count
	toSerialize["next"] = o.Next.Get()
	toSerialize["previous"] = o.Previous.Get()
	toSerialize["results"] = o.Results
	return toSerialize, nil
}

type NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse struct {
	value *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse
	isSet bool
}

func (v NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) Get() *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse {
	return v.value
}

func (v *NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) Set(val *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) {
	v.value = val
	v.isSet = true
}

func (v NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) IsSet() bool {
	return v.isSet
}

func (v *NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse(val *PaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) *NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse {
	return &NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse{value: val, isSet: true}
}

func (v NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaginationSerializerWithGalleryViewFieldOptionsExampleRowResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

