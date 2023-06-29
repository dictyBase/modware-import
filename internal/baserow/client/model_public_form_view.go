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

// checks if the PublicFormView type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PublicFormView{}

// PublicFormView struct for PublicFormView
type PublicFormView struct {
	// The title that is displayed at the beginning of the form.
	Title *string `json:"title,omitempty"`
	// The description that is displayed at the beginning of the form.
	Description *string `json:"description,omitempty"`
	Mode *ModeEnum `json:"mode,omitempty"`
	CoverImage NullablePublicFormViewCoverImage `json:"cover_image,omitempty"`
	LogoImage NullablePublicFormViewLogoImage `json:"logo_image,omitempty"`
	// The text displayed on the submit button.
	SubmitText *string `json:"submit_text,omitempty"`
	Fields []PublicFormViewFieldOptions `json:"fields"`
	ShowLogo *bool `json:"show_logo,omitempty"`
}

// NewPublicFormView instantiates a new PublicFormView object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPublicFormView(fields []PublicFormViewFieldOptions) *PublicFormView {
	this := PublicFormView{}
	this.Fields = fields
	return &this
}

// NewPublicFormViewWithDefaults instantiates a new PublicFormView object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPublicFormViewWithDefaults() *PublicFormView {
	this := PublicFormView{}
	return &this
}

// GetTitle returns the Title field value if set, zero value otherwise.
func (o *PublicFormView) GetTitle() string {
	if o == nil || IsNil(o.Title) {
		var ret string
		return ret
	}
	return *o.Title
}

// GetTitleOk returns a tuple with the Title field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetTitleOk() (*string, bool) {
	if o == nil || IsNil(o.Title) {
		return nil, false
	}
	return o.Title, true
}

// HasTitle returns a boolean if a field has been set.
func (o *PublicFormView) HasTitle() bool {
	if o != nil && !IsNil(o.Title) {
		return true
	}

	return false
}

// SetTitle gets a reference to the given string and assigns it to the Title field.
func (o *PublicFormView) SetTitle(v string) {
	o.Title = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *PublicFormView) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *PublicFormView) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *PublicFormView) SetDescription(v string) {
	o.Description = &v
}

// GetMode returns the Mode field value if set, zero value otherwise.
func (o *PublicFormView) GetMode() ModeEnum {
	if o == nil || IsNil(o.Mode) {
		var ret ModeEnum
		return ret
	}
	return *o.Mode
}

// GetModeOk returns a tuple with the Mode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetModeOk() (*ModeEnum, bool) {
	if o == nil || IsNil(o.Mode) {
		return nil, false
	}
	return o.Mode, true
}

// HasMode returns a boolean if a field has been set.
func (o *PublicFormView) HasMode() bool {
	if o != nil && !IsNil(o.Mode) {
		return true
	}

	return false
}

// SetMode gets a reference to the given ModeEnum and assigns it to the Mode field.
func (o *PublicFormView) SetMode(v ModeEnum) {
	o.Mode = &v
}

// GetCoverImage returns the CoverImage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublicFormView) GetCoverImage() PublicFormViewCoverImage {
	if o == nil || IsNil(o.CoverImage.Get()) {
		var ret PublicFormViewCoverImage
		return ret
	}
	return *o.CoverImage.Get()
}

// GetCoverImageOk returns a tuple with the CoverImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublicFormView) GetCoverImageOk() (*PublicFormViewCoverImage, bool) {
	if o == nil {
		return nil, false
	}
	return o.CoverImage.Get(), o.CoverImage.IsSet()
}

// HasCoverImage returns a boolean if a field has been set.
func (o *PublicFormView) HasCoverImage() bool {
	if o != nil && o.CoverImage.IsSet() {
		return true
	}

	return false
}

// SetCoverImage gets a reference to the given NullablePublicFormViewCoverImage and assigns it to the CoverImage field.
func (o *PublicFormView) SetCoverImage(v PublicFormViewCoverImage) {
	o.CoverImage.Set(&v)
}
// SetCoverImageNil sets the value for CoverImage to be an explicit nil
func (o *PublicFormView) SetCoverImageNil() {
	o.CoverImage.Set(nil)
}

// UnsetCoverImage ensures that no value is present for CoverImage, not even an explicit nil
func (o *PublicFormView) UnsetCoverImage() {
	o.CoverImage.Unset()
}

// GetLogoImage returns the LogoImage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublicFormView) GetLogoImage() PublicFormViewLogoImage {
	if o == nil || IsNil(o.LogoImage.Get()) {
		var ret PublicFormViewLogoImage
		return ret
	}
	return *o.LogoImage.Get()
}

// GetLogoImageOk returns a tuple with the LogoImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublicFormView) GetLogoImageOk() (*PublicFormViewLogoImage, bool) {
	if o == nil {
		return nil, false
	}
	return o.LogoImage.Get(), o.LogoImage.IsSet()
}

// HasLogoImage returns a boolean if a field has been set.
func (o *PublicFormView) HasLogoImage() bool {
	if o != nil && o.LogoImage.IsSet() {
		return true
	}

	return false
}

// SetLogoImage gets a reference to the given NullablePublicFormViewLogoImage and assigns it to the LogoImage field.
func (o *PublicFormView) SetLogoImage(v PublicFormViewLogoImage) {
	o.LogoImage.Set(&v)
}
// SetLogoImageNil sets the value for LogoImage to be an explicit nil
func (o *PublicFormView) SetLogoImageNil() {
	o.LogoImage.Set(nil)
}

// UnsetLogoImage ensures that no value is present for LogoImage, not even an explicit nil
func (o *PublicFormView) UnsetLogoImage() {
	o.LogoImage.Unset()
}

// GetSubmitText returns the SubmitText field value if set, zero value otherwise.
func (o *PublicFormView) GetSubmitText() string {
	if o == nil || IsNil(o.SubmitText) {
		var ret string
		return ret
	}
	return *o.SubmitText
}

// GetSubmitTextOk returns a tuple with the SubmitText field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetSubmitTextOk() (*string, bool) {
	if o == nil || IsNil(o.SubmitText) {
		return nil, false
	}
	return o.SubmitText, true
}

// HasSubmitText returns a boolean if a field has been set.
func (o *PublicFormView) HasSubmitText() bool {
	if o != nil && !IsNil(o.SubmitText) {
		return true
	}

	return false
}

// SetSubmitText gets a reference to the given string and assigns it to the SubmitText field.
func (o *PublicFormView) SetSubmitText(v string) {
	o.SubmitText = &v
}

// GetFields returns the Fields field value
func (o *PublicFormView) GetFields() []PublicFormViewFieldOptions {
	if o == nil {
		var ret []PublicFormViewFieldOptions
		return ret
	}

	return o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetFieldsOk() ([]PublicFormViewFieldOptions, bool) {
	if o == nil {
		return nil, false
	}
	return o.Fields, true
}

// SetFields sets field value
func (o *PublicFormView) SetFields(v []PublicFormViewFieldOptions) {
	o.Fields = v
}

// GetShowLogo returns the ShowLogo field value if set, zero value otherwise.
func (o *PublicFormView) GetShowLogo() bool {
	if o == nil || IsNil(o.ShowLogo) {
		var ret bool
		return ret
	}
	return *o.ShowLogo
}

// GetShowLogoOk returns a tuple with the ShowLogo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PublicFormView) GetShowLogoOk() (*bool, bool) {
	if o == nil || IsNil(o.ShowLogo) {
		return nil, false
	}
	return o.ShowLogo, true
}

// HasShowLogo returns a boolean if a field has been set.
func (o *PublicFormView) HasShowLogo() bool {
	if o != nil && !IsNil(o.ShowLogo) {
		return true
	}

	return false
}

// SetShowLogo gets a reference to the given bool and assigns it to the ShowLogo field.
func (o *PublicFormView) SetShowLogo(v bool) {
	o.ShowLogo = &v
}

func (o PublicFormView) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PublicFormView) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
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
	toSerialize["fields"] = o.Fields
	if !IsNil(o.ShowLogo) {
		toSerialize["show_logo"] = o.ShowLogo
	}
	return toSerialize, nil
}

type NullablePublicFormView struct {
	value *PublicFormView
	isSet bool
}

func (v NullablePublicFormView) Get() *PublicFormView {
	return v.value
}

func (v *NullablePublicFormView) Set(val *PublicFormView) {
	v.value = val
	v.isSet = true
}

func (v NullablePublicFormView) IsSet() bool {
	return v.isSet
}

func (v *NullablePublicFormView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePublicFormView(val *PublicFormView) *NullablePublicFormView {
	return &NullablePublicFormView{value: val, isSet: true}
}

func (v NullablePublicFormView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePublicFormView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


