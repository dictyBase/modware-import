/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"time"
)

// checks if the FormViewCreateViewLogoImage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &FormViewCreateViewLogoImage{}

// FormViewCreateViewLogoImage The logo image that must be displayed at the top of the form.
type FormViewCreateViewLogoImage struct {
	Size int32 `json:"size"`
	MimeType *string `json:"mime_type,omitempty"`
	IsImage *bool `json:"is_image,omitempty"`
	ImageWidth NullableInt32 `json:"image_width,omitempty"`
	ImageHeight NullableInt32 `json:"image_height,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Url string `json:"url"`
	Thumbnails map[string]interface{} `json:"thumbnails"`
	Name string `json:"name"`
	OriginalName string `json:"original_name"`
}

// NewFormViewCreateViewLogoImage instantiates a new FormViewCreateViewLogoImage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFormViewCreateViewLogoImage(size int32, uploadedAt time.Time, url string, thumbnails map[string]interface{}, name string, originalName string) *FormViewCreateViewLogoImage {
	this := FormViewCreateViewLogoImage{}
	this.Size = size
	this.UploadedAt = uploadedAt
	this.Url = url
	this.Thumbnails = thumbnails
	this.Name = name
	this.OriginalName = originalName
	return &this
}

// NewFormViewCreateViewLogoImageWithDefaults instantiates a new FormViewCreateViewLogoImage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFormViewCreateViewLogoImageWithDefaults() *FormViewCreateViewLogoImage {
	this := FormViewCreateViewLogoImage{}
	return &this
}

// GetSize returns the Size field value
func (o *FormViewCreateViewLogoImage) GetSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *FormViewCreateViewLogoImage) SetSize(v int32) {
	o.Size = v
}

// GetMimeType returns the MimeType field value if set, zero value otherwise.
func (o *FormViewCreateViewLogoImage) GetMimeType() string {
	if o == nil || IsNil(o.MimeType) {
		var ret string
		return ret
	}
	return *o.MimeType
}

// GetMimeTypeOk returns a tuple with the MimeType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetMimeTypeOk() (*string, bool) {
	if o == nil || IsNil(o.MimeType) {
		return nil, false
	}
	return o.MimeType, true
}

// HasMimeType returns a boolean if a field has been set.
func (o *FormViewCreateViewLogoImage) HasMimeType() bool {
	if o != nil && !IsNil(o.MimeType) {
		return true
	}

	return false
}

// SetMimeType gets a reference to the given string and assigns it to the MimeType field.
func (o *FormViewCreateViewLogoImage) SetMimeType(v string) {
	o.MimeType = &v
}

// GetIsImage returns the IsImage field value if set, zero value otherwise.
func (o *FormViewCreateViewLogoImage) GetIsImage() bool {
	if o == nil || IsNil(o.IsImage) {
		var ret bool
		return ret
	}
	return *o.IsImage
}

// GetIsImageOk returns a tuple with the IsImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetIsImageOk() (*bool, bool) {
	if o == nil || IsNil(o.IsImage) {
		return nil, false
	}
	return o.IsImage, true
}

// HasIsImage returns a boolean if a field has been set.
func (o *FormViewCreateViewLogoImage) HasIsImage() bool {
	if o != nil && !IsNil(o.IsImage) {
		return true
	}

	return false
}

// SetIsImage gets a reference to the given bool and assigns it to the IsImage field.
func (o *FormViewCreateViewLogoImage) SetIsImage(v bool) {
	o.IsImage = &v
}

// GetImageWidth returns the ImageWidth field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *FormViewCreateViewLogoImage) GetImageWidth() int32 {
	if o == nil || IsNil(o.ImageWidth.Get()) {
		var ret int32
		return ret
	}
	return *o.ImageWidth.Get()
}

// GetImageWidthOk returns a tuple with the ImageWidth field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FormViewCreateViewLogoImage) GetImageWidthOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ImageWidth.Get(), o.ImageWidth.IsSet()
}

// HasImageWidth returns a boolean if a field has been set.
func (o *FormViewCreateViewLogoImage) HasImageWidth() bool {
	if o != nil && o.ImageWidth.IsSet() {
		return true
	}

	return false
}

// SetImageWidth gets a reference to the given NullableInt32 and assigns it to the ImageWidth field.
func (o *FormViewCreateViewLogoImage) SetImageWidth(v int32) {
	o.ImageWidth.Set(&v)
}
// SetImageWidthNil sets the value for ImageWidth to be an explicit nil
func (o *FormViewCreateViewLogoImage) SetImageWidthNil() {
	o.ImageWidth.Set(nil)
}

// UnsetImageWidth ensures that no value is present for ImageWidth, not even an explicit nil
func (o *FormViewCreateViewLogoImage) UnsetImageWidth() {
	o.ImageWidth.Unset()
}

// GetImageHeight returns the ImageHeight field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *FormViewCreateViewLogoImage) GetImageHeight() int32 {
	if o == nil || IsNil(o.ImageHeight.Get()) {
		var ret int32
		return ret
	}
	return *o.ImageHeight.Get()
}

// GetImageHeightOk returns a tuple with the ImageHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FormViewCreateViewLogoImage) GetImageHeightOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ImageHeight.Get(), o.ImageHeight.IsSet()
}

// HasImageHeight returns a boolean if a field has been set.
func (o *FormViewCreateViewLogoImage) HasImageHeight() bool {
	if o != nil && o.ImageHeight.IsSet() {
		return true
	}

	return false
}

// SetImageHeight gets a reference to the given NullableInt32 and assigns it to the ImageHeight field.
func (o *FormViewCreateViewLogoImage) SetImageHeight(v int32) {
	o.ImageHeight.Set(&v)
}
// SetImageHeightNil sets the value for ImageHeight to be an explicit nil
func (o *FormViewCreateViewLogoImage) SetImageHeightNil() {
	o.ImageHeight.Set(nil)
}

// UnsetImageHeight ensures that no value is present for ImageHeight, not even an explicit nil
func (o *FormViewCreateViewLogoImage) UnsetImageHeight() {
	o.ImageHeight.Unset()
}

// GetUploadedAt returns the UploadedAt field value
func (o *FormViewCreateViewLogoImage) GetUploadedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.UploadedAt
}

// GetUploadedAtOk returns a tuple with the UploadedAt field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetUploadedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UploadedAt, true
}

// SetUploadedAt sets field value
func (o *FormViewCreateViewLogoImage) SetUploadedAt(v time.Time) {
	o.UploadedAt = v
}

// GetUrl returns the Url field value
func (o *FormViewCreateViewLogoImage) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *FormViewCreateViewLogoImage) SetUrl(v string) {
	o.Url = v
}

// GetThumbnails returns the Thumbnails field value
func (o *FormViewCreateViewLogoImage) GetThumbnails() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Thumbnails
}

// GetThumbnailsOk returns a tuple with the Thumbnails field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetThumbnailsOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.Thumbnails, true
}

// SetThumbnails sets field value
func (o *FormViewCreateViewLogoImage) SetThumbnails(v map[string]interface{}) {
	o.Thumbnails = v
}

// GetName returns the Name field value
func (o *FormViewCreateViewLogoImage) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *FormViewCreateViewLogoImage) SetName(v string) {
	o.Name = v
}

// GetOriginalName returns the OriginalName field value
func (o *FormViewCreateViewLogoImage) GetOriginalName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OriginalName
}

// GetOriginalNameOk returns a tuple with the OriginalName field value
// and a boolean to check if the value has been set.
func (o *FormViewCreateViewLogoImage) GetOriginalNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OriginalName, true
}

// SetOriginalName sets field value
func (o *FormViewCreateViewLogoImage) SetOriginalName(v string) {
	o.OriginalName = v
}

func (o FormViewCreateViewLogoImage) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o FormViewCreateViewLogoImage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["size"] = o.Size
	if !IsNil(o.MimeType) {
		toSerialize["mime_type"] = o.MimeType
	}
	if !IsNil(o.IsImage) {
		toSerialize["is_image"] = o.IsImage
	}
	if o.ImageWidth.IsSet() {
		toSerialize["image_width"] = o.ImageWidth.Get()
	}
	if o.ImageHeight.IsSet() {
		toSerialize["image_height"] = o.ImageHeight.Get()
	}
	// skip: uploaded_at is readOnly
	// skip: url is readOnly
	// skip: thumbnails is readOnly
	// skip: name is readOnly
	toSerialize["original_name"] = o.OriginalName
	return toSerialize, nil
}

type NullableFormViewCreateViewLogoImage struct {
	value *FormViewCreateViewLogoImage
	isSet bool
}

func (v NullableFormViewCreateViewLogoImage) Get() *FormViewCreateViewLogoImage {
	return v.value
}

func (v *NullableFormViewCreateViewLogoImage) Set(val *FormViewCreateViewLogoImage) {
	v.value = val
	v.isSet = true
}

func (v NullableFormViewCreateViewLogoImage) IsSet() bool {
	return v.isSet
}

func (v *NullableFormViewCreateViewLogoImage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFormViewCreateViewLogoImage(val *FormViewCreateViewLogoImage) *NullableFormViewCreateViewLogoImage {
	return &NullableFormViewCreateViewLogoImage{value: val, isSet: true}
}

func (v NullableFormViewCreateViewLogoImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFormViewCreateViewLogoImage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

