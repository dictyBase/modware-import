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

// checks if the KanbanViewView type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KanbanViewView{}

// KanbanViewView struct for KanbanViewView
type KanbanViewView struct {
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
	SingleSelectField NullableInt32 `json:"single_select_field,omitempty"`
	// References a file field of which the first image must be shown as card cover image.
	CardCoverImageField NullableInt32 `json:"card_cover_image_field,omitempty"`
	// Indicates whether the view is publicly accessible to visitors.
	Public *bool `json:"public,omitempty"`
	// The unique slug that can be used to construct a public URL.
	Slug string `json:"slug"`
}

// NewKanbanViewView instantiates a new KanbanViewView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKanbanViewView(id int32, tableId int32, name string, order int32, type_ string, table Table, publicViewHasPassword bool, ownershipType string, slug string) *KanbanViewView {
	this := KanbanViewView{}
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

// NewKanbanViewViewWithDefaults instantiates a new KanbanViewView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKanbanViewViewWithDefaults() *KanbanViewView {
	this := KanbanViewView{}
	return &this
}

// GetId returns the Id field value
func (o *KanbanViewView) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *KanbanViewView) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *KanbanViewView) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *KanbanViewView) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *KanbanViewView) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *KanbanViewView) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *KanbanViewView) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *KanbanViewView) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *KanbanViewView) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *KanbanViewView) SetType(v string) {
	o.Type = v
}

// GetTable returns the Table field value
func (o *KanbanViewView) GetTable() Table {
	if o == nil {
		var ret Table
		return ret
	}

	return o.Table
}

// GetTableOk returns a tuple with the Table field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetTableOk() (*Table, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Table, true
}

// SetTable sets field value
func (o *KanbanViewView) SetTable(v Table) {
	o.Table = v
}

// GetFilterType returns the FilterType field value if set, zero value otherwise.
func (o *KanbanViewView) GetFilterType() ConditionTypeEnum {
	if o == nil || IsNil(o.FilterType) {
		var ret ConditionTypeEnum
		return ret
	}
	return *o.FilterType
}

// GetFilterTypeOk returns a tuple with the FilterType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetFilterTypeOk() (*ConditionTypeEnum, bool) {
	if o == nil || IsNil(o.FilterType) {
		return nil, false
	}
	return o.FilterType, true
}

// HasFilterType returns a boolean if a field has been set.
func (o *KanbanViewView) HasFilterType() bool {
	if o != nil && !IsNil(o.FilterType) {
		return true
	}

	return false
}

// SetFilterType gets a reference to the given ConditionTypeEnum and assigns it to the FilterType field.
func (o *KanbanViewView) SetFilterType(v ConditionTypeEnum) {
	o.FilterType = &v
}

// GetFilters returns the Filters field value if set, zero value otherwise.
func (o *KanbanViewView) GetFilters() []ViewFilter {
	if o == nil || IsNil(o.Filters) {
		var ret []ViewFilter
		return ret
	}
	return o.Filters
}

// GetFiltersOk returns a tuple with the Filters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetFiltersOk() ([]ViewFilter, bool) {
	if o == nil || IsNil(o.Filters) {
		return nil, false
	}
	return o.Filters, true
}

// HasFilters returns a boolean if a field has been set.
func (o *KanbanViewView) HasFilters() bool {
	if o != nil && !IsNil(o.Filters) {
		return true
	}

	return false
}

// SetFilters gets a reference to the given []ViewFilter and assigns it to the Filters field.
func (o *KanbanViewView) SetFilters(v []ViewFilter) {
	o.Filters = v
}

// GetSortings returns the Sortings field value if set, zero value otherwise.
func (o *KanbanViewView) GetSortings() []ViewSort {
	if o == nil || IsNil(o.Sortings) {
		var ret []ViewSort
		return ret
	}
	return o.Sortings
}

// GetSortingsOk returns a tuple with the Sortings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetSortingsOk() ([]ViewSort, bool) {
	if o == nil || IsNil(o.Sortings) {
		return nil, false
	}
	return o.Sortings, true
}

// HasSortings returns a boolean if a field has been set.
func (o *KanbanViewView) HasSortings() bool {
	if o != nil && !IsNil(o.Sortings) {
		return true
	}

	return false
}

// SetSortings gets a reference to the given []ViewSort and assigns it to the Sortings field.
func (o *KanbanViewView) SetSortings(v []ViewSort) {
	o.Sortings = v
}

// GetDecorations returns the Decorations field value if set, zero value otherwise.
func (o *KanbanViewView) GetDecorations() []ViewDecoration {
	if o == nil || IsNil(o.Decorations) {
		var ret []ViewDecoration
		return ret
	}
	return o.Decorations
}

// GetDecorationsOk returns a tuple with the Decorations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetDecorationsOk() ([]ViewDecoration, bool) {
	if o == nil || IsNil(o.Decorations) {
		return nil, false
	}
	return o.Decorations, true
}

// HasDecorations returns a boolean if a field has been set.
func (o *KanbanViewView) HasDecorations() bool {
	if o != nil && !IsNil(o.Decorations) {
		return true
	}

	return false
}

// SetDecorations gets a reference to the given []ViewDecoration and assigns it to the Decorations field.
func (o *KanbanViewView) SetDecorations(v []ViewDecoration) {
	o.Decorations = v
}

// GetFiltersDisabled returns the FiltersDisabled field value if set, zero value otherwise.
func (o *KanbanViewView) GetFiltersDisabled() bool {
	if o == nil || IsNil(o.FiltersDisabled) {
		var ret bool
		return ret
	}
	return *o.FiltersDisabled
}

// GetFiltersDisabledOk returns a tuple with the FiltersDisabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetFiltersDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.FiltersDisabled) {
		return nil, false
	}
	return o.FiltersDisabled, true
}

// HasFiltersDisabled returns a boolean if a field has been set.
func (o *KanbanViewView) HasFiltersDisabled() bool {
	if o != nil && !IsNil(o.FiltersDisabled) {
		return true
	}

	return false
}

// SetFiltersDisabled gets a reference to the given bool and assigns it to the FiltersDisabled field.
func (o *KanbanViewView) SetFiltersDisabled(v bool) {
	o.FiltersDisabled = &v
}

// GetPublicViewHasPassword returns the PublicViewHasPassword field value
func (o *KanbanViewView) GetPublicViewHasPassword() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.PublicViewHasPassword
}

// GetPublicViewHasPasswordOk returns a tuple with the PublicViewHasPassword field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetPublicViewHasPasswordOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicViewHasPassword, true
}

// SetPublicViewHasPassword sets field value
func (o *KanbanViewView) SetPublicViewHasPassword(v bool) {
	o.PublicViewHasPassword = v
}

// GetShowLogo returns the ShowLogo field value if set, zero value otherwise.
func (o *KanbanViewView) GetShowLogo() bool {
	if o == nil || IsNil(o.ShowLogo) {
		var ret bool
		return ret
	}
	return *o.ShowLogo
}

// GetShowLogoOk returns a tuple with the ShowLogo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetShowLogoOk() (*bool, bool) {
	if o == nil || IsNil(o.ShowLogo) {
		return nil, false
	}
	return o.ShowLogo, true
}

// HasShowLogo returns a boolean if a field has been set.
func (o *KanbanViewView) HasShowLogo() bool {
	if o != nil && !IsNil(o.ShowLogo) {
		return true
	}

	return false
}

// SetShowLogo gets a reference to the given bool and assigns it to the ShowLogo field.
func (o *KanbanViewView) SetShowLogo(v bool) {
	o.ShowLogo = &v
}

// GetOwnershipType returns the OwnershipType field value
func (o *KanbanViewView) GetOwnershipType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OwnershipType
}

// GetOwnershipTypeOk returns a tuple with the OwnershipType field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetOwnershipTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OwnershipType, true
}

// SetOwnershipType sets field value
func (o *KanbanViewView) SetOwnershipType(v string) {
	o.OwnershipType = v
}

// GetSingleSelectField returns the SingleSelectField field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *KanbanViewView) GetSingleSelectField() int32 {
	if o == nil || IsNil(o.SingleSelectField.Get()) {
		var ret int32
		return ret
	}
	return *o.SingleSelectField.Get()
}

// GetSingleSelectFieldOk returns a tuple with the SingleSelectField field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KanbanViewView) GetSingleSelectFieldOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.SingleSelectField.Get(), o.SingleSelectField.IsSet()
}

// HasSingleSelectField returns a boolean if a field has been set.
func (o *KanbanViewView) HasSingleSelectField() bool {
	if o != nil && o.SingleSelectField.IsSet() {
		return true
	}

	return false
}

// SetSingleSelectField gets a reference to the given NullableInt32 and assigns it to the SingleSelectField field.
func (o *KanbanViewView) SetSingleSelectField(v int32) {
	o.SingleSelectField.Set(&v)
}
// SetSingleSelectFieldNil sets the value for SingleSelectField to be an explicit nil
func (o *KanbanViewView) SetSingleSelectFieldNil() {
	o.SingleSelectField.Set(nil)
}

// UnsetSingleSelectField ensures that no value is present for SingleSelectField, not even an explicit nil
func (o *KanbanViewView) UnsetSingleSelectField() {
	o.SingleSelectField.Unset()
}

// GetCardCoverImageField returns the CardCoverImageField field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *KanbanViewView) GetCardCoverImageField() int32 {
	if o == nil || IsNil(o.CardCoverImageField.Get()) {
		var ret int32
		return ret
	}
	return *o.CardCoverImageField.Get()
}

// GetCardCoverImageFieldOk returns a tuple with the CardCoverImageField field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KanbanViewView) GetCardCoverImageFieldOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.CardCoverImageField.Get(), o.CardCoverImageField.IsSet()
}

// HasCardCoverImageField returns a boolean if a field has been set.
func (o *KanbanViewView) HasCardCoverImageField() bool {
	if o != nil && o.CardCoverImageField.IsSet() {
		return true
	}

	return false
}

// SetCardCoverImageField gets a reference to the given NullableInt32 and assigns it to the CardCoverImageField field.
func (o *KanbanViewView) SetCardCoverImageField(v int32) {
	o.CardCoverImageField.Set(&v)
}
// SetCardCoverImageFieldNil sets the value for CardCoverImageField to be an explicit nil
func (o *KanbanViewView) SetCardCoverImageFieldNil() {
	o.CardCoverImageField.Set(nil)
}

// UnsetCardCoverImageField ensures that no value is present for CardCoverImageField, not even an explicit nil
func (o *KanbanViewView) UnsetCardCoverImageField() {
	o.CardCoverImageField.Unset()
}

// GetPublic returns the Public field value if set, zero value otherwise.
func (o *KanbanViewView) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}
	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *KanbanViewView) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *KanbanViewView) SetPublic(v bool) {
	o.Public = &v
}

// GetSlug returns the Slug field value
func (o *KanbanViewView) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *KanbanViewView) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *KanbanViewView) SetSlug(v string) {
	o.Slug = v
}

func (o KanbanViewView) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o KanbanViewView) ToMap() (map[string]interface{}, error) {
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
	if o.SingleSelectField.IsSet() {
		toSerialize["single_select_field"] = o.SingleSelectField.Get()
	}
	if o.CardCoverImageField.IsSet() {
		toSerialize["card_cover_image_field"] = o.CardCoverImageField.Get()
	}
	if !IsNil(o.Public) {
		toSerialize["public"] = o.Public
	}
	// skip: slug is readOnly
	return toSerialize, nil
}

type NullableKanbanViewView struct {
	value *KanbanViewView
	isSet bool
}

func (v NullableKanbanViewView) Get() *KanbanViewView {
	return v.value
}

func (v *NullableKanbanViewView) Set(val *KanbanViewView) {
	v.value = val
	v.isSet = true
}

func (v NullableKanbanViewView) IsSet() bool {
	return v.isSet
}

func (v *NullableKanbanViewView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKanbanViewView(val *KanbanViewView) *NullableKanbanViewView {
	return &NullableKanbanViewView{value: val, isSet: true}
}

func (v NullableKanbanViewView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKanbanViewView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


