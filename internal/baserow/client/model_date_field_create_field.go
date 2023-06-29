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

// checks if the DateFieldCreateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DateFieldCreateField{}

// DateFieldCreateField struct for DateFieldCreateField
type DateFieldCreateField struct {
	Name *string `json:"name,omitempty"`
	Type Type712Enum `json:"type"`
	DateFormat *DateFormatEnum `json:"date_format,omitempty"`
	// Indicates if the field also includes a time.
	DateIncludeTime *bool `json:"date_include_time,omitempty"`
	DateTimeFormat *DateTimeFormatEnum `json:"date_time_format,omitempty"`
	// Indicates if the timezone should be shown.
	DateShowTzinfo *bool `json:"date_show_tzinfo,omitempty"`
	// Force a timezone for the field overriding user profile settings.
	DateForceTimezone NullableString `json:"date_force_timezone,omitempty"`
}

// NewDateFieldCreateField instantiates a new DateFieldCreateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDateFieldCreateField(type_ Type712Enum) *DateFieldCreateField {
	this := DateFieldCreateField{}
	this.Type = type_
	return &this
}

// NewDateFieldCreateFieldWithDefaults instantiates a new DateFieldCreateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDateFieldCreateFieldWithDefaults() *DateFieldCreateField {
	this := DateFieldCreateField{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *DateFieldCreateField) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DateFieldCreateField) SetName(v string) {
	o.Name = &v
}

// GetType returns the Type field value
func (o *DateFieldCreateField) GetType() Type712Enum {
	if o == nil {
		var ret Type712Enum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *DateFieldCreateField) SetType(v Type712Enum) {
	o.Type = v
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise.
func (o *DateFieldCreateField) GetDateFormat() DateFormatEnum {
	if o == nil || IsNil(o.DateFormat) {
		var ret DateFormatEnum
		return ret
	}
	return *o.DateFormat
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetDateFormatOk() (*DateFormatEnum, bool) {
	if o == nil || IsNil(o.DateFormat) {
		return nil, false
	}
	return o.DateFormat, true
}

// HasDateFormat returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasDateFormat() bool {
	if o != nil && !IsNil(o.DateFormat) {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given DateFormatEnum and assigns it to the DateFormat field.
func (o *DateFieldCreateField) SetDateFormat(v DateFormatEnum) {
	o.DateFormat = &v
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise.
func (o *DateFieldCreateField) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil || IsNil(o.DateIncludeTime) {
		return nil, false
	}
	return o.DateIncludeTime, true
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasDateIncludeTime() bool {
	if o != nil && !IsNil(o.DateIncludeTime) {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given bool and assigns it to the DateIncludeTime field.
func (o *DateFieldCreateField) SetDateIncludeTime(v bool) {
	o.DateIncludeTime = &v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise.
func (o *DateFieldCreateField) GetDateTimeFormat() DateTimeFormatEnum {
	if o == nil || IsNil(o.DateTimeFormat) {
		var ret DateTimeFormatEnum
		return ret
	}
	return *o.DateTimeFormat
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetDateTimeFormatOk() (*DateTimeFormatEnum, bool) {
	if o == nil || IsNil(o.DateTimeFormat) {
		return nil, false
	}
	return o.DateTimeFormat, true
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasDateTimeFormat() bool {
	if o != nil && !IsNil(o.DateTimeFormat) {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given DateTimeFormatEnum and assigns it to the DateTimeFormat field.
func (o *DateFieldCreateField) SetDateTimeFormat(v DateTimeFormatEnum) {
	o.DateTimeFormat = &v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise.
func (o *DateFieldCreateField) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DateFieldCreateField) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil || IsNil(o.DateShowTzinfo) {
		return nil, false
	}
	return o.DateShowTzinfo, true
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasDateShowTzinfo() bool {
	if o != nil && !IsNil(o.DateShowTzinfo) {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given bool and assigns it to the DateShowTzinfo field.
func (o *DateFieldCreateField) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo = &v
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *DateFieldCreateField) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DateFieldCreateField) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *DateFieldCreateField) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *DateFieldCreateField) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *DateFieldCreateField) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *DateFieldCreateField) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

func (o DateFieldCreateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DateFieldCreateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	toSerialize["type"] = o.Type
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
	return toSerialize, nil
}

type NullableDateFieldCreateField struct {
	value *DateFieldCreateField
	isSet bool
}

func (v NullableDateFieldCreateField) Get() *DateFieldCreateField {
	return v.value
}

func (v *NullableDateFieldCreateField) Set(val *DateFieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableDateFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableDateFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDateFieldCreateField(val *DateFieldCreateField) *NullableDateFieldCreateField {
	return &NullableDateFieldCreateField{value: val, isSet: true}
}

func (v NullableDateFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDateFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

