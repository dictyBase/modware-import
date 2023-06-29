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

// checks if the GridViewView type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GridViewView{}

// GridViewView struct for GridViewView
type GridViewView struct {
	Id int32 `json:"id"`
	TableId int32 `json:"table_id"`
	Name string `json:"name"`
	Order int32 `json:"order"`
	Type string `json:"type"`
	Table Table `json:"table"`
	FilterType *ConditionTypeEnum `json:"filter_type,omitempty"`
	Filters []ViewFilter `json:"filters,omitempty"`
	Sortings []ViewSort `json:"sortings,omitempty"`
	Decorations []ViewDecoration `json:"decorations,omitempty"`
	// Allows users to see results unfiltered while still keeping the filters saved for the view.
	FiltersDisabled *bool `json:"filters_disabled,omitempty"`
	// Indicates whether the public view is password protected or not.  :return: True if the public view is password protected, False otherwise.
	PublicViewHasPassword bool `json:"public_view_has_password"`
	ShowLogo *bool `json:"show_logo,omitempty"`
	OwnershipType string `json:"ownership_type"`
	RowIdentifierType *RowIdentifierTypeEnum `json:"row_identifier_type,omitempty"`
	// Indicates whether the view is publicly accessible to visitors.
	Public *bool `json:"public,omitempty"`
	// The unique slug that can be used to construct a public URL.
	Slug string `json:"slug"`
}

// NewGridViewView instantiates a new GridViewView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGridViewView(id int32, tableId int32, name string, order int32, type_ string, table Table, publicViewHasPassword bool, ownershipType string, slug string) *GridViewView {
	this := GridViewView{}
	this.Id = id
	this.TableId = tableId
	this.Name = name
	this.Order = order
	this.Type = type_
	this.Table = table
	this.PublicViewHasPassword = publicViewHasPassword
	this.OwnershipType = ownershipType
	this.Slug = slug
	return &this
}

// NewGridViewViewWithDefaults instantiates a new GridViewView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGridViewViewWithDefaults() *GridViewView {
	this := GridViewView{}
	return &this
}

// GetId returns the Id field value
func (o *GridViewView) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GridViewView) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *GridViewView) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *GridViewView) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *GridViewView) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GridViewView) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *GridViewView) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *GridViewView) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *GridViewView) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GridViewView) SetType(v string) {
	o.Type = v
}

// GetTable returns the Table field value
func (o *GridViewView) GetTable() Table {
	if o == nil {
		var ret Table
		return ret
	}

	return o.Table
}

// GetTableOk returns a tuple with the Table field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetTableOk() (*Table, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Table, true
}

// SetTable sets field value
func (o *GridViewView) SetTable(v Table) {
	o.Table = v
}

// GetFilterType returns the FilterType field value if set, zero value otherwise.
func (o *GridViewView) GetFilterType() ConditionTypeEnum {
	if o == nil || IsNil(o.FilterType) {
		var ret ConditionTypeEnum
		return ret
	}
	return *o.FilterType
}

// GetFilterTypeOk returns a tuple with the FilterType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetFilterTypeOk() (*ConditionTypeEnum, bool) {
	if o == nil || IsNil(o.FilterType) {
		return nil, false
	}
	return o.FilterType, true
}

// HasFilterType returns a boolean if a field has been set.
func (o *GridViewView) HasFilterType() bool {
	if o != nil && !IsNil(o.FilterType) {
		return true
	}

	return false
}

// SetFilterType gets a reference to the given ConditionTypeEnum and assigns it to the FilterType field.
func (o *GridViewView) SetFilterType(v ConditionTypeEnum) {
	o.FilterType = &v
}

// GetFilters returns the Filters field value if set, zero value otherwise.
func (o *GridViewView) GetFilters() []ViewFilter {
	if o == nil || IsNil(o.Filters) {
		var ret []ViewFilter
		return ret
	}
	return o.Filters
}

// GetFiltersOk returns a tuple with the Filters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetFiltersOk() ([]ViewFilter, bool) {
	if o == nil || IsNil(o.Filters) {
		return nil, false
	}
	return o.Filters, true
}

// HasFilters returns a boolean if a field has been set.
func (o *GridViewView) HasFilters() bool {
	if o != nil && !IsNil(o.Filters) {
		return true
	}

	return false
}

// SetFilters gets a reference to the given []ViewFilter and assigns it to the Filters field.
func (o *GridViewView) SetFilters(v []ViewFilter) {
	o.Filters = v
}

// GetSortings returns the Sortings field value if set, zero value otherwise.
func (o *GridViewView) GetSortings() []ViewSort {
	if o == nil || IsNil(o.Sortings) {
		var ret []ViewSort
		return ret
	}
	return o.Sortings
}

// GetSortingsOk returns a tuple with the Sortings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetSortingsOk() ([]ViewSort, bool) {
	if o == nil || IsNil(o.Sortings) {
		return nil, false
	}
	return o.Sortings, true
}

// HasSortings returns a boolean if a field has been set.
func (o *GridViewView) HasSortings() bool {
	if o != nil && !IsNil(o.Sortings) {
		return true
	}

	return false
}

// SetSortings gets a reference to the given []ViewSort and assigns it to the Sortings field.
func (o *GridViewView) SetSortings(v []ViewSort) {
	o.Sortings = v
}

// GetDecorations returns the Decorations field value if set, zero value otherwise.
func (o *GridViewView) GetDecorations() []ViewDecoration {
	if o == nil || IsNil(o.Decorations) {
		var ret []ViewDecoration
		return ret
	}
	return o.Decorations
}

// GetDecorationsOk returns a tuple with the Decorations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetDecorationsOk() ([]ViewDecoration, bool) {
	if o == nil || IsNil(o.Decorations) {
		return nil, false
	}
	return o.Decorations, true
}

// HasDecorations returns a boolean if a field has been set.
func (o *GridViewView) HasDecorations() bool {
	if o != nil && !IsNil(o.Decorations) {
		return true
	}

	return false
}

// SetDecorations gets a reference to the given []ViewDecoration and assigns it to the Decorations field.
func (o *GridViewView) SetDecorations(v []ViewDecoration) {
	o.Decorations = v
}

// GetFiltersDisabled returns the FiltersDisabled field value if set, zero value otherwise.
func (o *GridViewView) GetFiltersDisabled() bool {
	if o == nil || IsNil(o.FiltersDisabled) {
		var ret bool
		return ret
	}
	return *o.FiltersDisabled
}

// GetFiltersDisabledOk returns a tuple with the FiltersDisabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetFiltersDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.FiltersDisabled) {
		return nil, false
	}
	return o.FiltersDisabled, true
}

// HasFiltersDisabled returns a boolean if a field has been set.
func (o *GridViewView) HasFiltersDisabled() bool {
	if o != nil && !IsNil(o.FiltersDisabled) {
		return true
	}

	return false
}

// SetFiltersDisabled gets a reference to the given bool and assigns it to the FiltersDisabled field.
func (o *GridViewView) SetFiltersDisabled(v bool) {
	o.FiltersDisabled = &v
}

// GetPublicViewHasPassword returns the PublicViewHasPassword field value
func (o *GridViewView) GetPublicViewHasPassword() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.PublicViewHasPassword
}

// GetPublicViewHasPasswordOk returns a tuple with the PublicViewHasPassword field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetPublicViewHasPasswordOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicViewHasPassword, true
}

// SetPublicViewHasPassword sets field value
func (o *GridViewView) SetPublicViewHasPassword(v bool) {
	o.PublicViewHasPassword = v
}

// GetShowLogo returns the ShowLogo field value if set, zero value otherwise.
func (o *GridViewView) GetShowLogo() bool {
	if o == nil || IsNil(o.ShowLogo) {
		var ret bool
		return ret
	}
	return *o.ShowLogo
}

// GetShowLogoOk returns a tuple with the ShowLogo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetShowLogoOk() (*bool, bool) {
	if o == nil || IsNil(o.ShowLogo) {
		return nil, false
	}
	return o.ShowLogo, true
}

// HasShowLogo returns a boolean if a field has been set.
func (o *GridViewView) HasShowLogo() bool {
	if o != nil && !IsNil(o.ShowLogo) {
		return true
	}

	return false
}

// SetShowLogo gets a reference to the given bool and assigns it to the ShowLogo field.
func (o *GridViewView) SetShowLogo(v bool) {
	o.ShowLogo = &v
}

// GetOwnershipType returns the OwnershipType field value
func (o *GridViewView) GetOwnershipType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OwnershipType
}

// GetOwnershipTypeOk returns a tuple with the OwnershipType field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetOwnershipTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OwnershipType, true
}

// SetOwnershipType sets field value
func (o *GridViewView) SetOwnershipType(v string) {
	o.OwnershipType = v
}

// GetRowIdentifierType returns the RowIdentifierType field value if set, zero value otherwise.
func (o *GridViewView) GetRowIdentifierType() RowIdentifierTypeEnum {
	if o == nil || IsNil(o.RowIdentifierType) {
		var ret RowIdentifierTypeEnum
		return ret
	}
	return *o.RowIdentifierType
}

// GetRowIdentifierTypeOk returns a tuple with the RowIdentifierType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetRowIdentifierTypeOk() (*RowIdentifierTypeEnum, bool) {
	if o == nil || IsNil(o.RowIdentifierType) {
		return nil, false
	}
	return o.RowIdentifierType, true
}

// HasRowIdentifierType returns a boolean if a field has been set.
func (o *GridViewView) HasRowIdentifierType() bool {
	if o != nil && !IsNil(o.RowIdentifierType) {
		return true
	}

	return false
}

// SetRowIdentifierType gets a reference to the given RowIdentifierTypeEnum and assigns it to the RowIdentifierType field.
func (o *GridViewView) SetRowIdentifierType(v RowIdentifierTypeEnum) {
	o.RowIdentifierType = &v
}

// GetPublic returns the Public field value if set, zero value otherwise.
func (o *GridViewView) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GridViewView) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}
	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *GridViewView) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *GridViewView) SetPublic(v bool) {
	o.Public = &v
}

// GetSlug returns the Slug field value
func (o *GridViewView) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *GridViewView) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *GridViewView) SetSlug(v string) {
	o.Slug = v
}

func (o GridViewView) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GridViewView) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: table_id is readOnly
	toSerialize["name"] = o.Name
	toSerialize["order"] = o.Order
	// skip: type is readOnly
	toSerialize["table"] = o.Table
	if !IsNil(o.FilterType) {
		toSerialize["filter_type"] = o.FilterType
	}
	if !IsNil(o.Filters) {
		toSerialize["filters"] = o.Filters
	}
	if !IsNil(o.Sortings) {
		toSerialize["sortings"] = o.Sortings
	}
	if !IsNil(o.Decorations) {
		toSerialize["decorations"] = o.Decorations
	}
	if !IsNil(o.FiltersDisabled) {
		toSerialize["filters_disabled"] = o.FiltersDisabled
	}
	// skip: public_view_has_password is readOnly
	if !IsNil(o.ShowLogo) {
		toSerialize["show_logo"] = o.ShowLogo
	}
	toSerialize["ownership_type"] = o.OwnershipType
	if !IsNil(o.RowIdentifierType) {
		toSerialize["row_identifier_type"] = o.RowIdentifierType
	}
	if !IsNil(o.Public) {
		toSerialize["public"] = o.Public
	}
	// skip: slug is readOnly
	return toSerialize, nil
}

type NullableGridViewView struct {
	value *GridViewView
	isSet bool
}

func (v NullableGridViewView) Get() *GridViewView {
	return v.value
}

func (v *NullableGridViewView) Set(val *GridViewView) {
	v.value = val
	v.isSet = true
}

func (v NullableGridViewView) IsSet() bool {
	return v.isSet
}

func (v *NullableGridViewView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGridViewView(val *GridViewView) *NullableGridViewView {
	return &NullableGridViewView{value: val, isSet: true}
}

func (v NullableGridViewView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGridViewView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

