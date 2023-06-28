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

// checks if the RequestLastModifiedFieldUpdateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RequestLastModifiedFieldUpdateField{}

// RequestLastModifiedFieldUpdateField struct for RequestLastModifiedFieldUpdateField
type RequestLastModifiedFieldUpdateField struct {
	Name *string `json:"name,omitempty"`
	Type *Type712Enum `json:"type,omitempty"`
	DateFormat *DateFormatEnum `json:"date_format,omitempty"`
	// Indicates if the field also includes a time.
	DateIncludeTime *bool `json:"date_include_time,omitempty"`
	DateTimeFormat *DateTimeFormatEnum `json:"date_time_format,omitempty"`
	// Indicates if the timezone should be shown.
	DateShowTzinfo *bool `json:"date_show_tzinfo,omitempty"`
	// Force a timezone for the field overriding user profile settings.
	DateForceTimezone NullableString `json:"date_force_timezone,omitempty"`
	// ('A UTC offset in minutes to add to all the field datetimes values.',)
	DateForceTimezoneOffset NullableInt32 `json:"date_force_timezone_offset,omitempty"`
}

// NewRequestLastModifiedFieldUpdateField instantiates a new RequestLastModifiedFieldUpdateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequestLastModifiedFieldUpdateField() *RequestLastModifiedFieldUpdateField {
	this := RequestLastModifiedFieldUpdateField{}
	return &this
}

// NewRequestLastModifiedFieldUpdateFieldWithDefaults instantiates a new RequestLastModifiedFieldUpdateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequestLastModifiedFieldUpdateFieldWithDefaults() *RequestLastModifiedFieldUpdateField {
	this := RequestLastModifiedFieldUpdateField{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *RequestLastModifiedFieldUpdateField) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetType() Type712Enum {
	if o == nil || IsNil(o.Type) {
		var ret Type712Enum
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given Type712Enum and assigns it to the Type field.
func (o *RequestLastModifiedFieldUpdateField) SetType(v Type712Enum) {
	o.Type = &v
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetDateFormat() DateFormatEnum {
	if o == nil || IsNil(o.DateFormat) {
		var ret DateFormatEnum
		return ret
	}
	return *o.DateFormat
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetDateFormatOk() (*DateFormatEnum, bool) {
	if o == nil || IsNil(o.DateFormat) {
		return nil, false
	}
	return o.DateFormat, true
}

// HasDateFormat returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateFormat() bool {
	if o != nil && !IsNil(o.DateFormat) {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given DateFormatEnum and assigns it to the DateFormat field.
func (o *RequestLastModifiedFieldUpdateField) SetDateFormat(v DateFormatEnum) {
	o.DateFormat = &v
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil || IsNil(o.DateIncludeTime) {
		return nil, false
	}
	return o.DateIncludeTime, true
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateIncludeTime() bool {
	if o != nil && !IsNil(o.DateIncludeTime) {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given bool and assigns it to the DateIncludeTime field.
func (o *RequestLastModifiedFieldUpdateField) SetDateIncludeTime(v bool) {
	o.DateIncludeTime = &v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetDateTimeFormat() DateTimeFormatEnum {
	if o == nil || IsNil(o.DateTimeFormat) {
		var ret DateTimeFormatEnum
		return ret
	}
	return *o.DateTimeFormat
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetDateTimeFormatOk() (*DateTimeFormatEnum, bool) {
	if o == nil || IsNil(o.DateTimeFormat) {
		return nil, false
	}
	return o.DateTimeFormat, true
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateTimeFormat() bool {
	if o != nil && !IsNil(o.DateTimeFormat) {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given DateTimeFormatEnum and assigns it to the DateTimeFormat field.
func (o *RequestLastModifiedFieldUpdateField) SetDateTimeFormat(v DateTimeFormatEnum) {
	o.DateTimeFormat = &v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise.
func (o *RequestLastModifiedFieldUpdateField) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestLastModifiedFieldUpdateField) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil || IsNil(o.DateShowTzinfo) {
		return nil, false
	}
	return o.DateShowTzinfo, true
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateShowTzinfo() bool {
	if o != nil && !IsNil(o.DateShowTzinfo) {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given bool and assigns it to the DateShowTzinfo field.
func (o *RequestLastModifiedFieldUpdateField) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo = &v
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RequestLastModifiedFieldUpdateField) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestLastModifiedFieldUpdateField) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *RequestLastModifiedFieldUpdateField) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *RequestLastModifiedFieldUpdateField) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *RequestLastModifiedFieldUpdateField) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

// GetDateForceTimezoneOffset returns the DateForceTimezoneOffset field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RequestLastModifiedFieldUpdateField) GetDateForceTimezoneOffset() int32 {
	if o == nil || IsNil(o.DateForceTimezoneOffset.Get()) {
		var ret int32
		return ret
	}
	return *o.DateForceTimezoneOffset.Get()
}

// GetDateForceTimezoneOffsetOk returns a tuple with the DateForceTimezoneOffset field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestLastModifiedFieldUpdateField) GetDateForceTimezoneOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezoneOffset.Get(), o.DateForceTimezoneOffset.IsSet()
}

// HasDateForceTimezoneOffset returns a boolean if a field has been set.
func (o *RequestLastModifiedFieldUpdateField) HasDateForceTimezoneOffset() bool {
	if o != nil && o.DateForceTimezoneOffset.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezoneOffset gets a reference to the given NullableInt32 and assigns it to the DateForceTimezoneOffset field.
func (o *RequestLastModifiedFieldUpdateField) SetDateForceTimezoneOffset(v int32) {
	o.DateForceTimezoneOffset.Set(&v)
}
// SetDateForceTimezoneOffsetNil sets the value for DateForceTimezoneOffset to be an explicit nil
func (o *RequestLastModifiedFieldUpdateField) SetDateForceTimezoneOffsetNil() {
	o.DateForceTimezoneOffset.Set(nil)
}

// UnsetDateForceTimezoneOffset ensures that no value is present for DateForceTimezoneOffset, not even an explicit nil
func (o *RequestLastModifiedFieldUpdateField) UnsetDateForceTimezoneOffset() {
	o.DateForceTimezoneOffset.Unset()
}

func (o RequestLastModifiedFieldUpdateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RequestLastModifiedFieldUpdateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.DateFormat) {
		toSerialize["date_format"] = o.DateFormat
	}
	if !IsNil(o.DateIncludeTime) {
		toSerialize["date_include_time"] = o.DateIncludeTime
	}
	if !IsNil(o.DateTimeFormat) {
		toSerialize["date_time_format"] = o.DateTimeFormat
	}
	if !IsNil(o.DateShowTzinfo) {
		toSerialize["date_show_tzinfo"] = o.DateShowTzinfo
	}
	if o.DateForceTimezone.IsSet() {
		toSerialize["date_force_timezone"] = o.DateForceTimezone.Get()
	}
	if o.DateForceTimezoneOffset.IsSet() {
		toSerialize["date_force_timezone_offset"] = o.DateForceTimezoneOffset.Get()
	}
	return toSerialize, nil
}

type NullableRequestLastModifiedFieldUpdateField struct {
	value *RequestLastModifiedFieldUpdateField
	isSet bool
}

func (v NullableRequestLastModifiedFieldUpdateField) Get() *RequestLastModifiedFieldUpdateField {
	return v.value
}

func (v *NullableRequestLastModifiedFieldUpdateField) Set(val *RequestLastModifiedFieldUpdateField) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestLastModifiedFieldUpdateField) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestLastModifiedFieldUpdateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestLastModifiedFieldUpdateField(val *RequestLastModifiedFieldUpdateField) *NullableRequestLastModifiedFieldUpdateField {
	return &NullableRequestLastModifiedFieldUpdateField{value: val, isSet: true}
}

func (v NullableRequestLastModifiedFieldUpdateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestLastModifiedFieldUpdateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


