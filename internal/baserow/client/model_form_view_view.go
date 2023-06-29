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

// checks if the FormViewView type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FormViewView{}

// FormViewView struct for FormViewView
type FormViewView struct {
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
	// The title that is displayed at the beginning of the form.
	Title *string `json:"title,omitempty"`
	// The description that is displayed at the beginning of the form.
	Description *string `json:"description,omitempty"`
	Mode *ModeEnum `json:"mode,omitempty"`
	CoverImage NullableFormViewCreateViewCoverImage `json:"cover_image,omitempty"`
	LogoImage NullableFormViewCreateViewLogoImage `json:"logo_image,omitempty"`
	// The text displayed on the submit button.
	SubmitText *string `json:"submit_text,omitempty"`
	SubmitAction *SubmitActionEnum `json:"submit_action,omitempty"`
	// If the `submit_action` is MESSAGE, then this message will be shown to the visitor after submitting the form.
	SubmitActionMessage *string `json:"submit_action_message,omitempty"`
	// If the `submit_action` is REDIRECT,then the visitors will be redirected to the this URL after submitting the form.
	SubmitActionRedirectUrl *string `json:"submit_action_redirect_url,omitempty"`
	// Indicates whether the view is publicly accessible to visitors.
	Public *bool `json:"public,omitempty"`
	// The unique slug that can be used to construct a public URL.
	Slug string `json:"slug"`
}

// NewFormViewView instantiates a new FormViewView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFormViewView(id int32, tableId int32, name string, order int32, type_ string, table Table, publicViewHasPassword bool, ownershipType string, slug string) *FormViewView {
	this := FormViewView{}
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

// NewFormViewViewWithDefaults instantiates a new FormViewView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFormViewViewWithDefaults() *FormViewView {
	this := FormViewView{}
	return &this
}

// GetId returns the Id field value
func (o *FormViewView) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *FormViewView) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *FormViewView) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *FormViewView) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *FormViewView) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *FormViewView) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *FormViewView) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *FormViewView) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *FormViewView) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *FormViewView) SetType(v string) {
	o.Type = v
}

// GetTable returns the Table field value
func (o *FormViewView) GetTable() Table {
	if o == nil {
		var ret Table
		return ret
	}

	return o.Table
}

// GetTableOk returns a tuple with the Table field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetTableOk() (*Table, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Table, true
}

// SetTable sets field value
func (o *FormViewView) SetTable(v Table) {
	o.Table = v
}

// GetFilterType returns the FilterType field value if set, zero value otherwise.
func (o *FormViewView) GetFilterType() ConditionTypeEnum {
	if o == nil || IsNil(o.FilterType) {
		var ret ConditionTypeEnum
		return ret
	}
	return *o.FilterType
}

// GetFilterTypeOk returns a tuple with the FilterType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetFilterTypeOk() (*ConditionTypeEnum, bool) {
	if o == nil || IsNil(o.FilterType) {
		return nil, false
	}
	return o.FilterType, true
}

// HasFilterType returns a boolean if a field has been set.
func (o *FormViewView) HasFilterType() bool {
	if o != nil && !IsNil(o.FilterType) {
		return true
	}

	return false
}

// SetFilterType gets a reference to the given ConditionTypeEnum and assigns it to the FilterType field.
func (o *FormViewView) SetFilterType(v ConditionTypeEnum) {
	o.FilterType = &v
}

// GetFilters returns the Filters field value if set, zero value otherwise.
func (o *FormViewView) GetFilters() []ViewFilter {
	if o == nil || IsNil(o.Filters) {
		var ret []ViewFilter
		return ret
	}
	return o.Filters
}

// GetFiltersOk returns a tuple with the Filters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetFiltersOk() ([]ViewFilter, bool) {
	if o == nil || IsNil(o.Filters) {
		return nil, false
	}
	return o.Filters, true
}

// HasFilters returns a boolean if a field has been set.
func (o *FormViewView) HasFilters() bool {
	if o != nil && !IsNil(o.Filters) {
		return true
	}

	return false
}

// SetFilters gets a reference to the given []ViewFilter and assigns it to the Filters field.
func (o *FormViewView) SetFilters(v []ViewFilter) {
	o.Filters = v
}

// GetSortings returns the Sortings field value if set, zero value otherwise.
func (o *FormViewView) GetSortings() []ViewSort {
	if o == nil || IsNil(o.Sortings) {
		var ret []ViewSort
		return ret
	}
	return o.Sortings
}

// GetSortingsOk returns a tuple with the Sortings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSortingsOk() ([]ViewSort, bool) {
	if o == nil || IsNil(o.Sortings) {
		return nil, false
	}
	return o.Sortings, true
}

// HasSortings returns a boolean if a field has been set.
func (o *FormViewView) HasSortings() bool {
	if o != nil && !IsNil(o.Sortings) {
		return true
	}

	return false
}

// SetSortings gets a reference to the given []ViewSort and assigns it to the Sortings field.
func (o *FormViewView) SetSortings(v []ViewSort) {
	o.Sortings = v
}

// GetDecorations returns the Decorations field value if set, zero value otherwise.
func (o *FormViewView) GetDecorations() []ViewDecoration {
	if o == nil || IsNil(o.Decorations) {
		var ret []ViewDecoration
		return ret
	}
	return o.Decorations
}

// GetDecorationsOk returns a tuple with the Decorations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetDecorationsOk() ([]ViewDecoration, bool) {
	if o == nil || IsNil(o.Decorations) {
		return nil, false
	}
	return o.Decorations, true
}

// HasDecorations returns a boolean if a field has been set.
func (o *FormViewView) HasDecorations() bool {
	if o != nil && !IsNil(o.Decorations) {
		return true
	}

	return false
}

// SetDecorations gets a reference to the given []ViewDecoration and assigns it to the Decorations field.
func (o *FormViewView) SetDecorations(v []ViewDecoration) {
	o.Decorations = v
}

// GetFiltersDisabled returns the FiltersDisabled field value if set, zero value otherwise.
func (o *FormViewView) GetFiltersDisabled() bool {
	if o == nil || IsNil(o.FiltersDisabled) {
		var ret bool
		return ret
	}
	return *o.FiltersDisabled
}

// GetFiltersDisabledOk returns a tuple with the FiltersDisabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetFiltersDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.FiltersDisabled) {
		return nil, false
	}
	return o.FiltersDisabled, true
}

// HasFiltersDisabled returns a boolean if a field has been set.
func (o *FormViewView) HasFiltersDisabled() bool {
	if o != nil && !IsNil(o.FiltersDisabled) {
		return true
	}

	return false
}

// SetFiltersDisabled gets a reference to the given bool and assigns it to the FiltersDisabled field.
func (o *FormViewView) SetFiltersDisabled(v bool) {
	o.FiltersDisabled = &v
}

// GetPublicViewHasPassword returns the PublicViewHasPassword field value
func (o *FormViewView) GetPublicViewHasPassword() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.PublicViewHasPassword
}

// GetPublicViewHasPasswordOk returns a tuple with the PublicViewHasPassword field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetPublicViewHasPasswordOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicViewHasPassword, true
}

// SetPublicViewHasPassword sets field value
func (o *FormViewView) SetPublicViewHasPassword(v bool) {
	o.PublicViewHasPassword = v
}

// GetShowLogo returns the ShowLogo field value if set, zero value otherwise.
func (o *FormViewView) GetShowLogo() bool {
	if o == nil || IsNil(o.ShowLogo) {
		var ret bool
		return ret
	}
	return *o.ShowLogo
}

// GetShowLogoOk returns a tuple with the ShowLogo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetShowLogoOk() (*bool, bool) {
	if o == nil || IsNil(o.ShowLogo) {
		return nil, false
	}
	return o.ShowLogo, true
}

// HasShowLogo returns a boolean if a field has been set.
func (o *FormViewView) HasShowLogo() bool {
	if o != nil && !IsNil(o.ShowLogo) {
		return true
	}

	return false
}

// SetShowLogo gets a reference to the given bool and assigns it to the ShowLogo field.
func (o *FormViewView) SetShowLogo(v bool) {
	o.ShowLogo = &v
}

// GetOwnershipType returns the OwnershipType field value
func (o *FormViewView) GetOwnershipType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OwnershipType
}

// GetOwnershipTypeOk returns a tuple with the OwnershipType field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetOwnershipTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OwnershipType, true
}

// SetOwnershipType sets field value
func (o *FormViewView) SetOwnershipType(v string) {
	o.OwnershipType = v
}

// GetTitle returns the Title field value if set, zero value otherwise.
func (o *FormViewView) GetTitle() string {
	if o == nil || IsNil(o.Title) {
		var ret string
		return ret
	}
	return *o.Title
}

// GetTitleOk returns a tuple with the Title field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetTitleOk() (*string, bool) {
	if o == nil || IsNil(o.Title) {
		return nil, false
	}
	return o.Title, true
}

// HasTitle returns a boolean if a field has been set.
func (o *FormViewView) HasTitle() bool {
	if o != nil && !IsNil(o.Title) {
		return true
	}

	return false
}

// SetTitle gets a reference to the given string and assigns it to the Title field.
func (o *FormViewView) SetTitle(v string) {
	o.Title = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *FormViewView) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *FormViewView) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *FormViewView) SetDescription(v string) {
	o.Description = &v
}

// GetMode returns the Mode field value if set, zero value otherwise.
func (o *FormViewView) GetMode() ModeEnum {
	if o == nil || IsNil(o.Mode) {
		var ret ModeEnum
		return ret
	}
	return *o.Mode
}

// GetModeOk returns a tuple with the Mode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetModeOk() (*ModeEnum, bool) {
	if o == nil || IsNil(o.Mode) {
		return nil, false
	}
	return o.Mode, true
}

// HasMode returns a boolean if a field has been set.
func (o *FormViewView) HasMode() bool {
	if o != nil && !IsNil(o.Mode) {
		return true
	}

	return false
}

// SetMode gets a reference to the given ModeEnum and assigns it to the Mode field.
func (o *FormViewView) SetMode(v ModeEnum) {
	o.Mode = &v
}

// GetCoverImage returns the CoverImage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *FormViewView) GetCoverImage() FormViewCreateViewCoverImage {
	if o == nil || IsNil(o.CoverImage.Get()) {
		var ret FormViewCreateViewCoverImage
		return ret
	}
	return *o.CoverImage.Get()
}

// GetCoverImageOk returns a tuple with the CoverImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FormViewView) GetCoverImageOk() (*FormViewCreateViewCoverImage, bool) {
	if o == nil {
		return nil, false
	}
	return o.CoverImage.Get(), o.CoverImage.IsSet()
}

// HasCoverImage returns a boolean if a field has been set.
func (o *FormViewView) HasCoverImage() bool {
	if o != nil && o.CoverImage.IsSet() {
		return true
	}

	return false
}

// SetCoverImage gets a reference to the given NullableFormViewCreateViewCoverImage and assigns it to the CoverImage field.
func (o *FormViewView) SetCoverImage(v FormViewCreateViewCoverImage) {
	o.CoverImage.Set(&v)
}
// SetCoverImageNil sets the value for CoverImage to be an explicit nil
func (o *FormViewView) SetCoverImageNil() {
	o.CoverImage.Set(nil)
}

// UnsetCoverImage ensures that no value is present for CoverImage, not even an explicit nil
func (o *FormViewView) UnsetCoverImage() {
	o.CoverImage.Unset()
}

// GetLogoImage returns the LogoImage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *FormViewView) GetLogoImage() FormViewCreateViewLogoImage {
	if o == nil || IsNil(o.LogoImage.Get()) {
		var ret FormViewCreateViewLogoImage
		return ret
	}
	return *o.LogoImage.Get()
}

// GetLogoImageOk returns a tuple with the LogoImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FormViewView) GetLogoImageOk() (*FormViewCreateViewLogoImage, bool) {
	if o == nil {
		return nil, false
	}
	return o.LogoImage.Get(), o.LogoImage.IsSet()
}

// HasLogoImage returns a boolean if a field has been set.
func (o *FormViewView) HasLogoImage() bool {
	if o != nil && o.LogoImage.IsSet() {
		return true
	}

	return false
}

// SetLogoImage gets a reference to the given NullableFormViewCreateViewLogoImage and assigns it to the LogoImage field.
func (o *FormViewView) SetLogoImage(v FormViewCreateViewLogoImage) {
	o.LogoImage.Set(&v)
}
// SetLogoImageNil sets the value for LogoImage to be an explicit nil
func (o *FormViewView) SetLogoImageNil() {
	o.LogoImage.Set(nil)
}

// UnsetLogoImage ensures that no value is present for LogoImage, not even an explicit nil
func (o *FormViewView) UnsetLogoImage() {
	o.LogoImage.Unset()
}

// GetSubmitText returns the SubmitText field value if set, zero value otherwise.
func (o *FormViewView) GetSubmitText() string {
	if o == nil || IsNil(o.SubmitText) {
		var ret string
		return ret
	}
	return *o.SubmitText
}

// GetSubmitTextOk returns a tuple with the SubmitText field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSubmitTextOk() (*string, bool) {
	if o == nil || IsNil(o.SubmitText) {
		return nil, false
	}
	return o.SubmitText, true
}

// HasSubmitText returns a boolean if a field has been set.
func (o *FormViewView) HasSubmitText() bool {
	if o != nil && !IsNil(o.SubmitText) {
		return true
	}

	return false
}

// SetSubmitText gets a reference to the given string and assigns it to the SubmitText field.
func (o *FormViewView) SetSubmitText(v string) {
	o.SubmitText = &v
}

// GetSubmitAction returns the SubmitAction field value if set, zero value otherwise.
func (o *FormViewView) GetSubmitAction() SubmitActionEnum {
	if o == nil || IsNil(o.SubmitAction) {
		var ret SubmitActionEnum
		return ret
	}
	return *o.SubmitAction
}

// GetSubmitActionOk returns a tuple with the SubmitAction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSubmitActionOk() (*SubmitActionEnum, bool) {
	if o == nil || IsNil(o.SubmitAction) {
		return nil, false
	}
	return o.SubmitAction, true
}

// HasSubmitAction returns a boolean if a field has been set.
func (o *FormViewView) HasSubmitAction() bool {
	if o != nil && !IsNil(o.SubmitAction) {
		return true
	}

	return false
}

// SetSubmitAction gets a reference to the given SubmitActionEnum and assigns it to the SubmitAction field.
func (o *FormViewView) SetSubmitAction(v SubmitActionEnum) {
	o.SubmitAction = &v
}

// GetSubmitActionMessage returns the SubmitActionMessage field value if set, zero value otherwise.
func (o *FormViewView) GetSubmitActionMessage() string {
	if o == nil || IsNil(o.SubmitActionMessage) {
		var ret string
		return ret
	}
	return *o.SubmitActionMessage
}

// GetSubmitActionMessageOk returns a tuple with the SubmitActionMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSubmitActionMessageOk() (*string, bool) {
	if o == nil || IsNil(o.SubmitActionMessage) {
		return nil, false
	}
	return o.SubmitActionMessage, true
}

// HasSubmitActionMessage returns a boolean if a field has been set.
func (o *FormViewView) HasSubmitActionMessage() bool {
	if o != nil && !IsNil(o.SubmitActionMessage) {
		return true
	}

	return false
}

// SetSubmitActionMessage gets a reference to the given string and assigns it to the SubmitActionMessage field.
func (o *FormViewView) SetSubmitActionMessage(v string) {
	o.SubmitActionMessage = &v
}

// GetSubmitActionRedirectUrl returns the SubmitActionRedirectUrl field value if set, zero value otherwise.
func (o *FormViewView) GetSubmitActionRedirectUrl() string {
	if o == nil || IsNil(o.SubmitActionRedirectUrl) {
		var ret string
		return ret
	}
	return *o.SubmitActionRedirectUrl
}

// GetSubmitActionRedirectUrlOk returns a tuple with the SubmitActionRedirectUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSubmitActionRedirectUrlOk() (*string, bool) {
	if o == nil || IsNil(o.SubmitActionRedirectUrl) {
		return nil, false
	}
	return o.SubmitActionRedirectUrl, true
}

// HasSubmitActionRedirectUrl returns a boolean if a field has been set.
func (o *FormViewView) HasSubmitActionRedirectUrl() bool {
	if o != nil && !IsNil(o.SubmitActionRedirectUrl) {
		return true
	}

	return false
}

// SetSubmitActionRedirectUrl gets a reference to the given string and assigns it to the SubmitActionRedirectUrl field.
func (o *FormViewView) SetSubmitActionRedirectUrl(v string) {
	o.SubmitActionRedirectUrl = &v
}

// GetPublic returns the Public field value if set, zero value otherwise.
func (o *FormViewView) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewView) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}
	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *FormViewView) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *FormViewView) SetPublic(v bool) {
	o.Public = &v
}

// GetSlug returns the Slug field value
func (o *FormViewView) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *FormViewView) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *FormViewView) SetSlug(v string) {
	o.Slug = v
}

func (o FormViewView) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FormViewView) ToMap() (map[string]interface{}, error) {
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
	if !IsNil(o.Title) {
		toSerialize["title"] = o.Title
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Mode) {
		toSerialize["mode"] = o.Mode
	}
	if o.CoverImage.IsSet() {
		toSerialize["cover_image"] = o.CoverImage.Get()
	}
	if o.LogoImage.IsSet() {
		toSerialize["logo_image"] = o.LogoImage.Get()
	}
	if !IsNil(o.SubmitText) {
		toSerialize["submit_text"] = o.SubmitText
	}
	if !IsNil(o.SubmitAction) {
		toSerialize["submit_action"] = o.SubmitAction
	}
	if !IsNil(o.SubmitActionMessage) {
		toSerialize["submit_action_message"] = o.SubmitActionMessage
	}
	if !IsNil(o.SubmitActionRedirectUrl) {
		toSerialize["submit_action_redirect_url"] = o.SubmitActionRedirectUrl
	}
	if !IsNil(o.Public) {
		toSerialize["public"] = o.Public
	}
	// skip: slug is readOnly
	return toSerialize, nil
}

type NullableFormViewView struct {
	value *FormViewView
	isSet bool
}

func (v NullableFormViewView) Get() *FormViewView {
	return v.value
}

func (v *NullableFormViewView) Set(val *FormViewView) {
	v.value = val
	v.isSet = true
}

func (v NullableFormViewView) IsSet() bool {
	return v.isSet
}

func (v *NullableFormViewView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFormViewView(val *FormViewView) *NullableFormViewView {
	return &NullableFormViewView{value: val, isSet: true}
}

func (v NullableFormViewView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFormViewView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

