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

// checks if the CountFieldCreateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CountFieldCreateField{}

// CountFieldCreateField struct for CountFieldCreateField
type CountFieldCreateField struct {
	Name string `json:"name"`
	Type Type712Enum `json:"type"`
	// Indicates if the time zone should be shown.
	DateShowTzinfo NullableBool `json:"date_show_tzinfo,omitempty"`
	DateFormat NullableCountFieldCreateFieldDateFormat `json:"date_format,omitempty"`
	ArrayFormulaType NullableCountFieldCreateFieldArrayFormulaType `json:"array_formula_type,omitempty"`
	// Force a timezone for the field overriding user profile settings.
	DateForceTimezone NullableString `json:"date_force_timezone,omitempty"`
	// Indicates if the field also includes a time.
	DateIncludeTime NullableBool `json:"date_include_time,omitempty"`
	Nullable bool `json:"nullable"`
	DateTimeFormat NullableCountFieldCreateFieldDateTimeFormat `json:"date_time_format,omitempty"`
	NumberDecimalPlaces NullableCountFieldCreateFieldNumberDecimalPlaces `json:"number_decimal_places,omitempty"`
	Error NullableString `json:"error,omitempty"`
	// The id of the link row field to count values for.
	ThroughFieldId NullableInt32 `json:"through_field_id,omitempty"`
	FormulaType *FormulaTypeEnum `json:"formula_type,omitempty"`
}

// NewCountFieldCreateField instantiates a new CountFieldCreateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCountFieldCreateField(name string, type_ Type712Enum, nullable bool) *CountFieldCreateField {
	this := CountFieldCreateField{}
	this.Name = name
	this.Type = type_
	this.Nullable = nullable
	return &this
}

// NewCountFieldCreateFieldWithDefaults instantiates a new CountFieldCreateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCountFieldCreateFieldWithDefaults() *CountFieldCreateField {
	this := CountFieldCreateField{}
	return &this
}

// GetName returns the Name field value
func (o *CountFieldCreateField) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CountFieldCreateField) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CountFieldCreateField) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *CountFieldCreateField) GetType() Type712Enum {
	if o == nil {
		var ret Type712Enum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *CountFieldCreateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *CountFieldCreateField) SetType(v Type712Enum) {
	o.Type = v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo.Get()) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo.Get()
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateShowTzinfo.Get(), o.DateShowTzinfo.IsSet()
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasDateShowTzinfo() bool {
	if o != nil && o.DateShowTzinfo.IsSet() {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given NullableBool and assigns it to the DateShowTzinfo field.
func (o *CountFieldCreateField) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo.Set(&v)
}
// SetDateShowTzinfoNil sets the value for DateShowTzinfo to be an explicit nil
func (o *CountFieldCreateField) SetDateShowTzinfoNil() {
	o.DateShowTzinfo.Set(nil)
}

// UnsetDateShowTzinfo ensures that no value is present for DateShowTzinfo, not even an explicit nil
func (o *CountFieldCreateField) UnsetDateShowTzinfo() {
	o.DateShowTzinfo.Unset()
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetDateFormat() CountFieldCreateFieldDateFormat {
	if o == nil || IsNil(o.DateFormat.Get()) {
		var ret CountFieldCreateFieldDateFormat
		return ret
	}
	return *o.DateFormat.Get()
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetDateFormatOk() (*CountFieldCreateFieldDateFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateFormat.Get(), o.DateFormat.IsSet()
}

// HasDateFormat returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasDateFormat() bool {
	if o != nil && o.DateFormat.IsSet() {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given NullableCountFieldCreateFieldDateFormat and assigns it to the DateFormat field.
func (o *CountFieldCreateField) SetDateFormat(v CountFieldCreateFieldDateFormat) {
	o.DateFormat.Set(&v)
}
// SetDateFormatNil sets the value for DateFormat to be an explicit nil
func (o *CountFieldCreateField) SetDateFormatNil() {
	o.DateFormat.Set(nil)
}

// UnsetDateFormat ensures that no value is present for DateFormat, not even an explicit nil
func (o *CountFieldCreateField) UnsetDateFormat() {
	o.DateFormat.Unset()
}

// GetArrayFormulaType returns the ArrayFormulaType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetArrayFormulaType() CountFieldCreateFieldArrayFormulaType {
	if o == nil || IsNil(o.ArrayFormulaType.Get()) {
		var ret CountFieldCreateFieldArrayFormulaType
		return ret
	}
	return *o.ArrayFormulaType.Get()
}

// GetArrayFormulaTypeOk returns a tuple with the ArrayFormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetArrayFormulaTypeOk() (*CountFieldCreateFieldArrayFormulaType, bool) {
	if o == nil {
		return nil, false
	}
	return o.ArrayFormulaType.Get(), o.ArrayFormulaType.IsSet()
}

// HasArrayFormulaType returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasArrayFormulaType() bool {
	if o != nil && o.ArrayFormulaType.IsSet() {
		return true
	}

	return false
}

// SetArrayFormulaType gets a reference to the given NullableCountFieldCreateFieldArrayFormulaType and assigns it to the ArrayFormulaType field.
func (o *CountFieldCreateField) SetArrayFormulaType(v CountFieldCreateFieldArrayFormulaType) {
	o.ArrayFormulaType.Set(&v)
}
// SetArrayFormulaTypeNil sets the value for ArrayFormulaType to be an explicit nil
func (o *CountFieldCreateField) SetArrayFormulaTypeNil() {
	o.ArrayFormulaType.Set(nil)
}

// UnsetArrayFormulaType ensures that no value is present for ArrayFormulaType, not even an explicit nil
func (o *CountFieldCreateField) UnsetArrayFormulaType() {
	o.ArrayFormulaType.Unset()
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *CountFieldCreateField) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *CountFieldCreateField) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *CountFieldCreateField) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime.Get()) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime.Get()
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateIncludeTime.Get(), o.DateIncludeTime.IsSet()
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasDateIncludeTime() bool {
	if o != nil && o.DateIncludeTime.IsSet() {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given NullableBool and assigns it to the DateIncludeTime field.
func (o *CountFieldCreateField) SetDateIncludeTime(v bool) {
	o.DateIncludeTime.Set(&v)
}
// SetDateIncludeTimeNil sets the value for DateIncludeTime to be an explicit nil
func (o *CountFieldCreateField) SetDateIncludeTimeNil() {
	o.DateIncludeTime.Set(nil)
}

// UnsetDateIncludeTime ensures that no value is present for DateIncludeTime, not even an explicit nil
func (o *CountFieldCreateField) UnsetDateIncludeTime() {
	o.DateIncludeTime.Unset()
}

// GetNullable returns the Nullable field value
func (o *CountFieldCreateField) GetNullable() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Nullable
}

// GetNullableOk returns a tuple with the Nullable field value
// and a boolean to check if the value has been set.
func (o *CountFieldCreateField) GetNullableOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Nullable, true
}

// SetNullable sets field value
func (o *CountFieldCreateField) SetNullable(v bool) {
	o.Nullable = v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetDateTimeFormat() CountFieldCreateFieldDateTimeFormat {
	if o == nil || IsNil(o.DateTimeFormat.Get()) {
		var ret CountFieldCreateFieldDateTimeFormat
		return ret
	}
	return *o.DateTimeFormat.Get()
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetDateTimeFormatOk() (*CountFieldCreateFieldDateTimeFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateTimeFormat.Get(), o.DateTimeFormat.IsSet()
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasDateTimeFormat() bool {
	if o != nil && o.DateTimeFormat.IsSet() {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given NullableCountFieldCreateFieldDateTimeFormat and assigns it to the DateTimeFormat field.
func (o *CountFieldCreateField) SetDateTimeFormat(v CountFieldCreateFieldDateTimeFormat) {
	o.DateTimeFormat.Set(&v)
}
// SetDateTimeFormatNil sets the value for DateTimeFormat to be an explicit nil
func (o *CountFieldCreateField) SetDateTimeFormatNil() {
	o.DateTimeFormat.Set(nil)
}

// UnsetDateTimeFormat ensures that no value is present for DateTimeFormat, not even an explicit nil
func (o *CountFieldCreateField) UnsetDateTimeFormat() {
	o.DateTimeFormat.Unset()
}

// GetNumberDecimalPlaces returns the NumberDecimalPlaces field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetNumberDecimalPlaces() CountFieldCreateFieldNumberDecimalPlaces {
	if o == nil || IsNil(o.NumberDecimalPlaces.Get()) {
		var ret CountFieldCreateFieldNumberDecimalPlaces
		return ret
	}
	return *o.NumberDecimalPlaces.Get()
}

// GetNumberDecimalPlacesOk returns a tuple with the NumberDecimalPlaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetNumberDecimalPlacesOk() (*CountFieldCreateFieldNumberDecimalPlaces, bool) {
	if o == nil {
		return nil, false
	}
	return o.NumberDecimalPlaces.Get(), o.NumberDecimalPlaces.IsSet()
}

// HasNumberDecimalPlaces returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasNumberDecimalPlaces() bool {
	if o != nil && o.NumberDecimalPlaces.IsSet() {
		return true
	}

	return false
}

// SetNumberDecimalPlaces gets a reference to the given NullableCountFieldCreateFieldNumberDecimalPlaces and assigns it to the NumberDecimalPlaces field.
func (o *CountFieldCreateField) SetNumberDecimalPlaces(v CountFieldCreateFieldNumberDecimalPlaces) {
	o.NumberDecimalPlaces.Set(&v)
}
// SetNumberDecimalPlacesNil sets the value for NumberDecimalPlaces to be an explicit nil
func (o *CountFieldCreateField) SetNumberDecimalPlacesNil() {
	o.NumberDecimalPlaces.Set(nil)
}

// UnsetNumberDecimalPlaces ensures that no value is present for NumberDecimalPlaces, not even an explicit nil
func (o *CountFieldCreateField) UnsetNumberDecimalPlaces() {
	o.NumberDecimalPlaces.Unset()
}

// GetError returns the Error field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetError() string {
	if o == nil || IsNil(o.Error.Get()) {
		var ret string
		return ret
	}
	return *o.Error.Get()
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetErrorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Error.Get(), o.Error.IsSet()
}

// HasError returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasError() bool {
	if o != nil && o.Error.IsSet() {
		return true
	}

	return false
}

// SetError gets a reference to the given NullableString and assigns it to the Error field.
func (o *CountFieldCreateField) SetError(v string) {
	o.Error.Set(&v)
}
// SetErrorNil sets the value for Error to be an explicit nil
func (o *CountFieldCreateField) SetErrorNil() {
	o.Error.Set(nil)
}

// UnsetError ensures that no value is present for Error, not even an explicit nil
func (o *CountFieldCreateField) UnsetError() {
	o.Error.Unset()
}

// GetThroughFieldId returns the ThroughFieldId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CountFieldCreateField) GetThroughFieldId() int32 {
	if o == nil || IsNil(o.ThroughFieldId.Get()) {
		var ret int32
		return ret
	}
	return *o.ThroughFieldId.Get()
}

// GetThroughFieldIdOk returns a tuple with the ThroughFieldId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CountFieldCreateField) GetThroughFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ThroughFieldId.Get(), o.ThroughFieldId.IsSet()
}

// HasThroughFieldId returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasThroughFieldId() bool {
	if o != nil && o.ThroughFieldId.IsSet() {
		return true
	}

	return false
}

// SetThroughFieldId gets a reference to the given NullableInt32 and assigns it to the ThroughFieldId field.
func (o *CountFieldCreateField) SetThroughFieldId(v int32) {
	o.ThroughFieldId.Set(&v)
}
// SetThroughFieldIdNil sets the value for ThroughFieldId to be an explicit nil
func (o *CountFieldCreateField) SetThroughFieldIdNil() {
	o.ThroughFieldId.Set(nil)
}

// UnsetThroughFieldId ensures that no value is present for ThroughFieldId, not even an explicit nil
func (o *CountFieldCreateField) UnsetThroughFieldId() {
	o.ThroughFieldId.Unset()
}

// GetFormulaType returns the FormulaType field value if set, zero value otherwise.
func (o *CountFieldCreateField) GetFormulaType() FormulaTypeEnum {
	if o == nil || IsNil(o.FormulaType) {
		var ret FormulaTypeEnum
		return ret
	}
	return *o.FormulaType
}

// GetFormulaTypeOk returns a tuple with the FormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CountFieldCreateField) GetFormulaTypeOk() (*FormulaTypeEnum, bool) {
	if o == nil || IsNil(o.FormulaType) {
		return nil, false
	}
	return o.FormulaType, true
}

// HasFormulaType returns a boolean if a field has been set.
func (o *CountFieldCreateField) HasFormulaType() bool {
	if o != nil && !IsNil(o.FormulaType) {
		return true
	}

	return false
}

// SetFormulaType gets a reference to the given FormulaTypeEnum and assigns it to the FormulaType field.
func (o *CountFieldCreateField) SetFormulaType(v FormulaTypeEnum) {
	o.FormulaType = &v
}

func (o CountFieldCreateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CountFieldCreateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	if o.DateShowTzinfo.IsSet() {
		toSerialize["date_show_tzinfo"] = o.DateShowTzinfo.Get()
	}
	if o.DateFormat.IsSet() {
		toSerialize["date_format"] = o.DateFormat.Get()
	}
	if o.ArrayFormulaType.IsSet() {
		toSerialize["array_formula_type"] = o.ArrayFormulaType.Get()
	}
	if o.DateForceTimezone.IsSet() {
		toSerialize["date_force_timezone"] = o.DateForceTimezone.Get()
	}
	if o.DateIncludeTime.IsSet() {
		toSerialize["date_include_time"] = o.DateIncludeTime.Get()
	}
	// skip: nullable is readOnly
	if o.DateTimeFormat.IsSet() {
		toSerialize["date_time_format"] = o.DateTimeFormat.Get()
	}
	if o.NumberDecimalPlaces.IsSet() {
		toSerialize["number_decimal_places"] = o.NumberDecimalPlaces.Get()
	}
	if o.Error.IsSet() {
		toSerialize["error"] = o.Error.Get()
	}
	if o.ThroughFieldId.IsSet() {
		toSerialize["through_field_id"] = o.ThroughFieldId.Get()
	}
	if !IsNil(o.FormulaType) {
		toSerialize["formula_type"] = o.FormulaType
	}
	return toSerialize, nil
}

type NullableCountFieldCreateField struct {
	value *CountFieldCreateField
	isSet bool
}

func (v NullableCountFieldCreateField) Get() *CountFieldCreateField {
	return v.value
}

func (v *NullableCountFieldCreateField) Set(val *CountFieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableCountFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableCountFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCountFieldCreateField(val *CountFieldCreateField) *NullableCountFieldCreateField {
	return &NullableCountFieldCreateField{value: val, isSet: true}
}

func (v NullableCountFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCountFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


