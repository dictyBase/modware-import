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

// checks if the RollupFieldFieldSerializerWithRelatedFields type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RollupFieldFieldSerializerWithRelatedFields{}

// RollupFieldFieldSerializerWithRelatedFields struct for RollupFieldFieldSerializerWithRelatedFields
type RollupFieldFieldSerializerWithRelatedFields struct {
	Id int32 `json:"id"`
	TableId int32 `json:"table_id"`
	Name string `json:"name"`
	// Lowest first.
	Order int32 `json:"order"`
	// The type of the related field.
	Type string `json:"type"`
	// Indicates if the field is a primary field. If `true` the field cannot be deleted and the value should represent the whole row.
	Primary *bool `json:"primary,omitempty"`
	// Indicates whether the field is a read only field. If true, it's not possible to update the cell value.
	ReadOnly bool `json:"read_only"`
	// A list of related fields which also changed.
	RelatedFields []Field `json:"related_fields"`
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
	// The id of the link row field to rollup values for.
	ThroughFieldId NullableInt32 `json:"through_field_id,omitempty"`
	// The id of the field in the table linked to by the through_field to rollup.
	TargetFieldId NullableInt32 `json:"target_field_id,omitempty"`
	// The rollup formula function that must be applied.
	RollupFunction *string `json:"rollup_function,omitempty"`
	FormulaType *FormulaTypeEnum `json:"formula_type,omitempty"`
}

// NewRollupFieldFieldSerializerWithRelatedFields instantiates a new RollupFieldFieldSerializerWithRelatedFields object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRollupFieldFieldSerializerWithRelatedFields(id int32, tableId int32, name string, order int32, type_ string, readOnly bool, relatedFields []Field, nullable bool) *RollupFieldFieldSerializerWithRelatedFields {
	this := RollupFieldFieldSerializerWithRelatedFields{}
	this.Id = id
	this.TableId = tableId
	this.Name = name
	this.Order = order
	this.Type = type_
	this.ReadOnly = readOnly
	this.RelatedFields = relatedFields
	this.Nullable = nullable
	return &this
}

// NewRollupFieldFieldSerializerWithRelatedFieldsWithDefaults instantiates a new RollupFieldFieldSerializerWithRelatedFields object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRollupFieldFieldSerializerWithRelatedFieldsWithDefaults() *RollupFieldFieldSerializerWithRelatedFields {
	this := RollupFieldFieldSerializerWithRelatedFields{}
	return &this
}

// GetId returns the Id field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetType(v string) {
	o.Type = v
}

// GetPrimary returns the Primary field value if set, zero value otherwise.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetPrimary() bool {
	if o == nil || IsNil(o.Primary) {
		var ret bool
		return ret
	}
	return *o.Primary
}

// GetPrimaryOk returns a tuple with the Primary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetPrimaryOk() (*bool, bool) {
	if o == nil || IsNil(o.Primary) {
		return nil, false
	}
	return o.Primary, true
}

// HasPrimary returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasPrimary() bool {
	if o != nil && !IsNil(o.Primary) {
		return true
	}

	return false
}

// SetPrimary gets a reference to the given bool and assigns it to the Primary field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetPrimary(v bool) {
	o.Primary = &v
}

// GetReadOnly returns the ReadOnly field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetReadOnly() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ReadOnly
}

// GetReadOnlyOk returns a tuple with the ReadOnly field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetReadOnlyOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReadOnly, true
}

// SetReadOnly sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetReadOnly(v bool) {
	o.ReadOnly = v
}

// GetRelatedFields returns the RelatedFields field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetRelatedFields() []Field {
	if o == nil {
		var ret []Field
		return ret
	}

	return o.RelatedFields
}

// GetRelatedFieldsOk returns a tuple with the RelatedFields field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetRelatedFieldsOk() ([]Field, bool) {
	if o == nil {
		return nil, false
	}
	return o.RelatedFields, true
}

// SetRelatedFields sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetRelatedFields(v []Field) {
	o.RelatedFields = v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo.Get()) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo.Get()
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateShowTzinfo.Get(), o.DateShowTzinfo.IsSet()
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasDateShowTzinfo() bool {
	if o != nil && o.DateShowTzinfo.IsSet() {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given NullableBool and assigns it to the DateShowTzinfo field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo.Set(&v)
}
// SetDateShowTzinfoNil sets the value for DateShowTzinfo to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateShowTzinfoNil() {
	o.DateShowTzinfo.Set(nil)
}

// UnsetDateShowTzinfo ensures that no value is present for DateShowTzinfo, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetDateShowTzinfo() {
	o.DateShowTzinfo.Unset()
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateFormat() CountFieldCreateFieldDateFormat {
	if o == nil || IsNil(o.DateFormat.Get()) {
		var ret CountFieldCreateFieldDateFormat
		return ret
	}
	return *o.DateFormat.Get()
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateFormatOk() (*CountFieldCreateFieldDateFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateFormat.Get(), o.DateFormat.IsSet()
}

// HasDateFormat returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasDateFormat() bool {
	if o != nil && o.DateFormat.IsSet() {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given NullableCountFieldCreateFieldDateFormat and assigns it to the DateFormat field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateFormat(v CountFieldCreateFieldDateFormat) {
	o.DateFormat.Set(&v)
}
// SetDateFormatNil sets the value for DateFormat to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateFormatNil() {
	o.DateFormat.Set(nil)
}

// UnsetDateFormat ensures that no value is present for DateFormat, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetDateFormat() {
	o.DateFormat.Unset()
}

// GetArrayFormulaType returns the ArrayFormulaType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetArrayFormulaType() CountFieldCreateFieldArrayFormulaType {
	if o == nil || IsNil(o.ArrayFormulaType.Get()) {
		var ret CountFieldCreateFieldArrayFormulaType
		return ret
	}
	return *o.ArrayFormulaType.Get()
}

// GetArrayFormulaTypeOk returns a tuple with the ArrayFormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetArrayFormulaTypeOk() (*CountFieldCreateFieldArrayFormulaType, bool) {
	if o == nil {
		return nil, false
	}
	return o.ArrayFormulaType.Get(), o.ArrayFormulaType.IsSet()
}

// HasArrayFormulaType returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasArrayFormulaType() bool {
	if o != nil && o.ArrayFormulaType.IsSet() {
		return true
	}

	return false
}

// SetArrayFormulaType gets a reference to the given NullableCountFieldCreateFieldArrayFormulaType and assigns it to the ArrayFormulaType field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetArrayFormulaType(v CountFieldCreateFieldArrayFormulaType) {
	o.ArrayFormulaType.Set(&v)
}
// SetArrayFormulaTypeNil sets the value for ArrayFormulaType to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetArrayFormulaTypeNil() {
	o.ArrayFormulaType.Set(nil)
}

// UnsetArrayFormulaType ensures that no value is present for ArrayFormulaType, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetArrayFormulaType() {
	o.ArrayFormulaType.Unset()
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime.Get()) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime.Get()
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateIncludeTime.Get(), o.DateIncludeTime.IsSet()
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasDateIncludeTime() bool {
	if o != nil && o.DateIncludeTime.IsSet() {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given NullableBool and assigns it to the DateIncludeTime field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateIncludeTime(v bool) {
	o.DateIncludeTime.Set(&v)
}
// SetDateIncludeTimeNil sets the value for DateIncludeTime to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateIncludeTimeNil() {
	o.DateIncludeTime.Set(nil)
}

// UnsetDateIncludeTime ensures that no value is present for DateIncludeTime, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetDateIncludeTime() {
	o.DateIncludeTime.Unset()
}

// GetNullable returns the Nullable field value
func (o *RollupFieldFieldSerializerWithRelatedFields) GetNullable() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Nullable
}

// GetNullableOk returns a tuple with the Nullable field value
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetNullableOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Nullable, true
}

// SetNullable sets field value
func (o *RollupFieldFieldSerializerWithRelatedFields) SetNullable(v bool) {
	o.Nullable = v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateTimeFormat() CountFieldCreateFieldDateTimeFormat {
	if o == nil || IsNil(o.DateTimeFormat.Get()) {
		var ret CountFieldCreateFieldDateTimeFormat
		return ret
	}
	return *o.DateTimeFormat.Get()
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetDateTimeFormatOk() (*CountFieldCreateFieldDateTimeFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateTimeFormat.Get(), o.DateTimeFormat.IsSet()
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasDateTimeFormat() bool {
	if o != nil && o.DateTimeFormat.IsSet() {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given NullableCountFieldCreateFieldDateTimeFormat and assigns it to the DateTimeFormat field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateTimeFormat(v CountFieldCreateFieldDateTimeFormat) {
	o.DateTimeFormat.Set(&v)
}
// SetDateTimeFormatNil sets the value for DateTimeFormat to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetDateTimeFormatNil() {
	o.DateTimeFormat.Set(nil)
}

// UnsetDateTimeFormat ensures that no value is present for DateTimeFormat, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetDateTimeFormat() {
	o.DateTimeFormat.Unset()
}

// GetNumberDecimalPlaces returns the NumberDecimalPlaces field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetNumberDecimalPlaces() CountFieldCreateFieldNumberDecimalPlaces {
	if o == nil || IsNil(o.NumberDecimalPlaces.Get()) {
		var ret CountFieldCreateFieldNumberDecimalPlaces
		return ret
	}
	return *o.NumberDecimalPlaces.Get()
}

// GetNumberDecimalPlacesOk returns a tuple with the NumberDecimalPlaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetNumberDecimalPlacesOk() (*CountFieldCreateFieldNumberDecimalPlaces, bool) {
	if o == nil {
		return nil, false
	}
	return o.NumberDecimalPlaces.Get(), o.NumberDecimalPlaces.IsSet()
}

// HasNumberDecimalPlaces returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasNumberDecimalPlaces() bool {
	if o != nil && o.NumberDecimalPlaces.IsSet() {
		return true
	}

	return false
}

// SetNumberDecimalPlaces gets a reference to the given NullableCountFieldCreateFieldNumberDecimalPlaces and assigns it to the NumberDecimalPlaces field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetNumberDecimalPlaces(v CountFieldCreateFieldNumberDecimalPlaces) {
	o.NumberDecimalPlaces.Set(&v)
}
// SetNumberDecimalPlacesNil sets the value for NumberDecimalPlaces to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetNumberDecimalPlacesNil() {
	o.NumberDecimalPlaces.Set(nil)
}

// UnsetNumberDecimalPlaces ensures that no value is present for NumberDecimalPlaces, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetNumberDecimalPlaces() {
	o.NumberDecimalPlaces.Unset()
}

// GetError returns the Error field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetError() string {
	if o == nil || IsNil(o.Error.Get()) {
		var ret string
		return ret
	}
	return *o.Error.Get()
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetErrorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Error.Get(), o.Error.IsSet()
}

// HasError returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasError() bool {
	if o != nil && o.Error.IsSet() {
		return true
	}

	return false
}

// SetError gets a reference to the given NullableString and assigns it to the Error field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetError(v string) {
	o.Error.Set(&v)
}
// SetErrorNil sets the value for Error to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetErrorNil() {
	o.Error.Set(nil)
}

// UnsetError ensures that no value is present for Error, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetError() {
	o.Error.Unset()
}

// GetThroughFieldId returns the ThroughFieldId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetThroughFieldId() int32 {
	if o == nil || IsNil(o.ThroughFieldId.Get()) {
		var ret int32
		return ret
	}
	return *o.ThroughFieldId.Get()
}

// GetThroughFieldIdOk returns a tuple with the ThroughFieldId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetThroughFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ThroughFieldId.Get(), o.ThroughFieldId.IsSet()
}

// HasThroughFieldId returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasThroughFieldId() bool {
	if o != nil && o.ThroughFieldId.IsSet() {
		return true
	}

	return false
}

// SetThroughFieldId gets a reference to the given NullableInt32 and assigns it to the ThroughFieldId field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetThroughFieldId(v int32) {
	o.ThroughFieldId.Set(&v)
}
// SetThroughFieldIdNil sets the value for ThroughFieldId to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetThroughFieldIdNil() {
	o.ThroughFieldId.Set(nil)
}

// UnsetThroughFieldId ensures that no value is present for ThroughFieldId, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetThroughFieldId() {
	o.ThroughFieldId.Unset()
}

// GetTargetFieldId returns the TargetFieldId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *RollupFieldFieldSerializerWithRelatedFields) GetTargetFieldId() int32 {
	if o == nil || IsNil(o.TargetFieldId.Get()) {
		var ret int32
		return ret
	}
	return *o.TargetFieldId.Get()
}

// GetTargetFieldIdOk returns a tuple with the TargetFieldId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RollupFieldFieldSerializerWithRelatedFields) GetTargetFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.TargetFieldId.Get(), o.TargetFieldId.IsSet()
}

// HasTargetFieldId returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasTargetFieldId() bool {
	if o != nil && o.TargetFieldId.IsSet() {
		return true
	}

	return false
}

// SetTargetFieldId gets a reference to the given NullableInt32 and assigns it to the TargetFieldId field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetTargetFieldId(v int32) {
	o.TargetFieldId.Set(&v)
}
// SetTargetFieldIdNil sets the value for TargetFieldId to be an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) SetTargetFieldIdNil() {
	o.TargetFieldId.Set(nil)
}

// UnsetTargetFieldId ensures that no value is present for TargetFieldId, not even an explicit nil
func (o *RollupFieldFieldSerializerWithRelatedFields) UnsetTargetFieldId() {
	o.TargetFieldId.Unset()
}

// GetRollupFunction returns the RollupFunction field value if set, zero value otherwise.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetRollupFunction() string {
	if o == nil || IsNil(o.RollupFunction) {
		var ret string
		return ret
	}
	return *o.RollupFunction
}

// GetRollupFunctionOk returns a tuple with the RollupFunction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetRollupFunctionOk() (*string, bool) {
	if o == nil || IsNil(o.RollupFunction) {
		return nil, false
	}
	return o.RollupFunction, true
}

// HasRollupFunction returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasRollupFunction() bool {
	if o != nil && !IsNil(o.RollupFunction) {
		return true
	}

	return false
}

// SetRollupFunction gets a reference to the given string and assigns it to the RollupFunction field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetRollupFunction(v string) {
	o.RollupFunction = &v
}

// GetFormulaType returns the FormulaType field value if set, zero value otherwise.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetFormulaType() FormulaTypeEnum {
	if o == nil || IsNil(o.FormulaType) {
		var ret FormulaTypeEnum
		return ret
	}
	return *o.FormulaType
}

// GetFormulaTypeOk returns a tuple with the FormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) GetFormulaTypeOk() (*FormulaTypeEnum, bool) {
	if o == nil || IsNil(o.FormulaType) {
		return nil, false
	}
	return o.FormulaType, true
}

// HasFormulaType returns a boolean if a field has been set.
func (o *RollupFieldFieldSerializerWithRelatedFields) HasFormulaType() bool {
	if o != nil && !IsNil(o.FormulaType) {
		return true
	}

	return false
}

// SetFormulaType gets a reference to the given FormulaTypeEnum and assigns it to the FormulaType field.
func (o *RollupFieldFieldSerializerWithRelatedFields) SetFormulaType(v FormulaTypeEnum) {
	o.FormulaType = &v
}

func (o RollupFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RollupFieldFieldSerializerWithRelatedFields) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: table_id is readOnly
	toSerialize["name"] = o.Name
	toSerialize["order"] = o.Order
	// skip: type is readOnly
	if !IsNil(o.Primary) {
		toSerialize["primary"] = o.Primary
	}
	// skip: read_only is readOnly
	// skip: related_fields is readOnly
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
	if o.TargetFieldId.IsSet() {
		toSerialize["target_field_id"] = o.TargetFieldId.Get()
	}
	if !IsNil(o.RollupFunction) {
		toSerialize["rollup_function"] = o.RollupFunction
	}
	if !IsNil(o.FormulaType) {
		toSerialize["formula_type"] = o.FormulaType
	}
	return toSerialize, nil
}

type NullableRollupFieldFieldSerializerWithRelatedFields struct {
	value *RollupFieldFieldSerializerWithRelatedFields
	isSet bool
}

func (v NullableRollupFieldFieldSerializerWithRelatedFields) Get() *RollupFieldFieldSerializerWithRelatedFields {
	return v.value
}

func (v *NullableRollupFieldFieldSerializerWithRelatedFields) Set(val *RollupFieldFieldSerializerWithRelatedFields) {
	v.value = val
	v.isSet = true
}

func (v NullableRollupFieldFieldSerializerWithRelatedFields) IsSet() bool {
	return v.isSet
}

func (v *NullableRollupFieldFieldSerializerWithRelatedFields) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRollupFieldFieldSerializerWithRelatedFields(val *RollupFieldFieldSerializerWithRelatedFields) *NullableRollupFieldFieldSerializerWithRelatedFields {
	return &NullableRollupFieldFieldSerializerWithRelatedFields{value: val, isSet: true}
}

func (v NullableRollupFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRollupFieldFieldSerializerWithRelatedFields) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


