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

// checks if the GridViewCreateView type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GridViewCreateView{}

// GridViewCreateView struct for GridViewCreateView
type GridViewCreateView struct {
	Name string `json:"name"`
	Type ViewTypesEnum `json:"type"`
	OwnershipType *OwnershipTypeEnum `json:"ownership_type,omitempty"`
	FilterType *ConditionTypeEnum `json:"filter_type,omitempty"`
	// Allows users to see results unfiltered while still keeping the filters saved for the view.
	FiltersDisabled *bool `json:"filters_disabled,omitempty"`
	RowIdentifierType *RowIdentifierTypeEnum `json:"row_identifier_type,omitempty"`
	// Indicates whether the view is publicly accessible to visitors.
	Public *bool `json:"public,omitempty"`
	// The unique slug that can be used to construct a public URL.
	Slug string `json:"slug"`
}

// NewGridViewCreateView instantiates a new GridViewCreateView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGridViewCreateView(name string, type_ ViewTypesEnum, slug string) *GridViewCreateView {
	this := GridViewCreateView{}
	this.Name = name
	this.Type = type_
	var ownershipType OwnershipTypeEnum = COLLABORATIVE
	this.OwnershipType = &ownershipType
	this.Slug = slug
	return &this
}

// NewGridViewCreateViewWithDefaults instantiates a new GridViewCreateView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGridViewCreateViewWithDefaults() *GridViewCreateView {
	this := GridViewCreateView{}
	var ownershipType OwnershipTypeEnum = COLLABORATIVE
	this.OwnershipType = &ownershipType
	return &this
}

// GetName returns the Name field value
func (o *GridViewCreateView) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GridViewCreateView) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *GridViewCreateView) GetType() ViewTypesEnum {
	if o == nil {
		var ret ViewTypesEnum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetTypeOk() (*ViewTypesEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GridViewCreateView) SetType(v ViewTypesEnum) {
	o.Type = v
}

// GetOwnershipType returns the OwnershipType field value if set, zero value otherwise.
func (o *GridViewCreateView) GetOwnershipType() OwnershipTypeEnum {
	if o == nil || IsNil(o.OwnershipType) {
		var ret OwnershipTypeEnum
		return ret
	}
	return *o.OwnershipType
}

// GetOwnershipTypeOk returns a tuple with the OwnershipType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetOwnershipTypeOk() (*OwnershipTypeEnum, bool) {
	if o == nil || IsNil(o.OwnershipType) {
		return nil, false
	}
	return o.OwnershipType, true
}

// HasOwnershipType returns a boolean if a field has been set.
func (o *GridViewCreateView) HasOwnershipType() bool {
	if o != nil && !IsNil(o.OwnershipType) {
		return true
	}

	return false
}

// SetOwnershipType gets a reference to the given OwnershipTypeEnum and assigns it to the OwnershipType field.
func (o *GridViewCreateView) SetOwnershipType(v OwnershipTypeEnum) {
	o.OwnershipType = &v
}

// GetFilterType returns the FilterType field value if set, zero value otherwise.
func (o *GridViewCreateView) GetFilterType() ConditionTypeEnum {
	if o == nil || IsNil(o.FilterType) {
		var ret ConditionTypeEnum
		return ret
	}
	return *o.FilterType
}

// GetFilterTypeOk returns a tuple with the FilterType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetFilterTypeOk() (*ConditionTypeEnum, bool) {
	if o == nil || IsNil(o.FilterType) {
		return nil, false
	}
	return o.FilterType, true
}

// HasFilterType returns a boolean if a field has been set.
func (o *GridViewCreateView) HasFilterType() bool {
	if o != nil && !IsNil(o.FilterType) {
		return true
	}

	return false
}

// SetFilterType gets a reference to the given ConditionTypeEnum and assigns it to the FilterType field.
func (o *GridViewCreateView) SetFilterType(v ConditionTypeEnum) {
	o.FilterType = &v
}

// GetFiltersDisabled returns the FiltersDisabled field value if set, zero value otherwise.
func (o *GridViewCreateView) GetFiltersDisabled() bool {
	if o == nil || IsNil(o.FiltersDisabled) {
		var ret bool
		return ret
	}
	return *o.FiltersDisabled
}

// GetFiltersDisabledOk returns a tuple with the FiltersDisabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetFiltersDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.FiltersDisabled) {
		return nil, false
	}
	return o.FiltersDisabled, true
}

// HasFiltersDisabled returns a boolean if a field has been set.
func (o *GridViewCreateView) HasFiltersDisabled() bool {
	if o != nil && !IsNil(o.FiltersDisabled) {
		return true
	}

	return false
}

// SetFiltersDisabled gets a reference to the given bool and assigns it to the FiltersDisabled field.
func (o *GridViewCreateView) SetFiltersDisabled(v bool) {
	o.FiltersDisabled = &v
}

// GetRowIdentifierType returns the RowIdentifierType field value if set, zero value otherwise.
func (o *GridViewCreateView) GetRowIdentifierType() RowIdentifierTypeEnum {
	if o == nil || IsNil(o.RowIdentifierType) {
		var ret RowIdentifierTypeEnum
		return ret
	}
	return *o.RowIdentifierType
}

// GetRowIdentifierTypeOk returns a tuple with the RowIdentifierType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetRowIdentifierTypeOk() (*RowIdentifierTypeEnum, bool) {
	if o == nil || IsNil(o.RowIdentifierType) {
		return nil, false
	}
	return o.RowIdentifierType, true
}

// HasRowIdentifierType returns a boolean if a field has been set.
func (o *GridViewCreateView) HasRowIdentifierType() bool {
	if o != nil && !IsNil(o.RowIdentifierType) {
		return true
	}

	return false
}

// SetRowIdentifierType gets a reference to the given RowIdentifierTypeEnum and assigns it to the RowIdentifierType field.
func (o *GridViewCreateView) SetRowIdentifierType(v RowIdentifierTypeEnum) {
	o.RowIdentifierType = &v
}

// GetPublic returns the Public field value if set, zero value otherwise.
func (o *GridViewCreateView) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}
	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *GridViewCreateView) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *GridViewCreateView) SetPublic(v bool) {
	o.Public = &v
}

// GetSlug returns the Slug field value
func (o *GridViewCreateView) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *GridViewCreateView) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *GridViewCreateView) SetSlug(v string) {
	o.Slug = v
}

func (o GridViewCreateView) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GridViewCreateView) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	if !IsNil(o.OwnershipType) {
		toSerialize["ownership_type"] = o.OwnershipType
	}
	if !IsNil(o.FilterType) {
		toSerialize["filter_type"] = o.FilterType
	}
	if !IsNil(o.FiltersDisabled) {
		toSerialize["filters_disabled"] = o.FiltersDisabled
	}
	if !IsNil(o.RowIdentifierType) {
		toSerialize["row_identifier_type"] = o.RowIdentifierType
	}
	if !IsNil(o.Public) {
		toSerialize["public"] = o.Public
	}
	// skip: slug is readOnly
	return toSerialize, nil
}

type NullableGridViewCreateView struct {
	value *GridViewCreateView
	isSet bool
}

func (v NullableGridViewCreateView) Get() *GridViewCreateView {
	return v.value
}

func (v *NullableGridViewCreateView) Set(val *GridViewCreateView) {
	v.value = val
	v.isSet = true
}

func (v NullableGridViewCreateView) IsSet() bool {
	return v.isSet
}

func (v *NullableGridViewCreateView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGridViewCreateView(val *GridViewCreateView) *NullableGridViewCreateView {
	return &NullableGridViewCreateView{value: val, isSet: true}
}

func (v NullableGridViewCreateView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGridViewCreateView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


