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

// checks if the RequestDateFieldUpdateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RequestDateFieldUpdateField{}

// RequestDateFieldUpdateField struct for RequestDateFieldUpdateField
type RequestDateFieldUpdateField struct {
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

// NewRequestDateFieldUpdateField instantiates a new RequestDateFieldUpdateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequestDateFieldUpdateField() *RequestDateFieldUpdateField {
	this := RequestDateFieldUpdateField{}
	return &this
}

// NewRequestDateFieldUpdateFieldWithDefaults instantiates a new RequestDateFieldUpdateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequestDateFieldUpdateFieldWithDefaults() *RequestDateFieldUpdateField {
	this := RequestDateFieldUpdateField{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *RequestDateFieldUpdateField) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetType() Type712Enum {
	if o == nil || IsNil(o.Type) {
		var ret Type712Enum
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given Type712Enum and assigns it to the Type field.
func (o *RequestDateFieldUpdateField) SetType(v Type712Enum) {
	o.Type = &v
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetDateFormat() DateFormatEnum {
	if o == nil || IsNil(o.DateFormat) {
		var ret DateFormatEnum
		return ret
	}
	return *o.DateFormat
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetDateFormatOk() (*DateFormatEnum, bool) {
	if o == nil || IsNil(o.DateFormat) {
		return nil, false
	}
	return o.DateFormat, true
}

// HasDateFormat returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateFormat() bool {
	if o != nil && !IsNil(o.DateFormat) {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given DateFormatEnum and assigns it to the DateFormat field.
func (o *RequestDateFieldUpdateField) SetDateFormat(v DateFormatEnum) {
	o.DateFormat = &v
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil || IsNil(o.DateIncludeTime) {
		return nil, false
	}
	return o.DateIncludeTime, true
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateIncludeTime() bool {
	if o != nil && !IsNil(o.DateIncludeTime) {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given bool and assigns it to the DateIncludeTime field.
func (o *RequestDateFieldUpdateField) SetDateIncludeTime(v bool) {
	o.DateIncludeTime = &v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetDateTimeFormat() DateTimeFormatEnum {
	if o == nil || IsNil(o.DateTimeFormat) {
		var ret DateTimeFormatEnum
		return ret
	}
	return *o.DateTimeFormat
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetDateTimeFormatOk() (*DateTimeFormatEnum, bool) {
	if o == nil || IsNil(o.DateTimeFormat) {
		return nil, false
	}
	return o.DateTimeFormat, true
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateTimeFormat() bool {
	if o != nil && !IsNil(o.DateTimeFormat) {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given DateTimeFormatEnum and assigns it to the DateTimeFormat field.
func (o *RequestDateFieldUpdateField) SetDateTimeFormat(v DateTimeFormatEnum) {
	o.DateTimeFormat = &v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise.
func (o *RequestDateFieldUpdateField) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RequestDateFieldUpdateField) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil || IsNil(o.DateShowTzinfo) {
		return nil, false
	}
	return o.DateShowTzinfo, true
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateShowTzinfo() bool {
	if o != nil && !IsNil(o.DateShowTzinfo) {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given bool and assigns it to the DateShowTzinfo field.
func (o *RequestDateFieldUpdateField) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo = &v
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RequestDateFieldUpdateField) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestDateFieldUpdateField) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *RequestDateFieldUpdateField) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *RequestDateFieldUpdateField) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *RequestDateFieldUpdateField) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

// GetDateForceTimezoneOffset returns the DateForceTimezoneOffset field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RequestDateFieldUpdateField) GetDateForceTimezoneOffset() int32 {
	if o == nil || IsNil(o.DateForceTimezoneOffset.Get()) {
		var ret int32
		return ret
	}
	return *o.DateForceTimezoneOffset.Get()
}

// GetDateForceTimezoneOffsetOk returns a tuple with the DateForceTimezoneOffset field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestDateFieldUpdateField) GetDateForceTimezoneOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezoneOffset.Get(), o.DateForceTimezoneOffset.IsSet()
}

// HasDateForceTimezoneOffset returns a boolean if a field has been set.
func (o *RequestDateFieldUpdateField) HasDateForceTimezoneOffset() bool {
	if o != nil && o.DateForceTimezoneOffset.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezoneOffset gets a reference to the given NullableInt32 and assigns it to the DateForceTimezoneOffset field.
func (o *RequestDateFieldUpdateField) SetDateForceTimezoneOffset(v int32) {
	o.DateForceTimezoneOffset.Set(&v)
}
// SetDateForceTimezoneOffsetNil sets the value for DateForceTimezoneOffset to be an explicit nil
func (o *RequestDateFieldUpdateField) SetDateForceTimezoneOffsetNil() {
	o.DateForceTimezoneOffset.Set(nil)
}

// UnsetDateForceTimezoneOffset ensures that no value is present for DateForceTimezoneOffset, not even an explicit nil
func (o *RequestDateFieldUpdateField) UnsetDateForceTimezoneOffset() {
	o.DateForceTimezoneOffset.Unset()
}

func (o RequestDateFieldUpdateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RequestDateFieldUpdateField) ToMap() (map[string]interface{}, error) {
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

type NullableRequestDateFieldUpdateField struct {
	value *RequestDateFieldUpdateField
	isSet bool
}

func (v NullableRequestDateFieldUpdateField) Get() *RequestDateFieldUpdateField {
	return v.value
}

func (v *NullableRequestDateFieldUpdateField) Set(val *RequestDateFieldUpdateField) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestDateFieldUpdateField) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestDateFieldUpdateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestDateFieldUpdateField(val *RequestDateFieldUpdateField) *NullableRequestDateFieldUpdateField {
	return &NullableRequestDateFieldUpdateField{value: val, isSet: true}
}

func (v NullableRequestDateFieldUpdateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestDateFieldUpdateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


