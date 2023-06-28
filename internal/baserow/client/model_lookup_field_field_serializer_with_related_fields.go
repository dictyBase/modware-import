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

// checks if the LookupFieldFieldSerializerWithRelatedFields type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LookupFieldFieldSerializerWithRelatedFields{}

// LookupFieldFieldSerializerWithRelatedFields struct for LookupFieldFieldSerializerWithRelatedFields
type LookupFieldFieldSerializerWithRelatedFields struct {
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
	// The id of the link row field to lookup values for. Will override the `through_field_name` parameter if both are provided, however only one is required.
	ThroughFieldId NullableInt32 `json:"through_field_id,omitempty"`
	// The name of the link row field to lookup values for.
	ThroughFieldName NullableString `json:"through_field_name,omitempty"`
	// The id of the field in the table linked to by the through_field to lookup. Will override the `target_field_id` parameter if both are provided, however only one is required.
	TargetFieldId NullableInt32 `json:"target_field_id,omitempty"`
	// The name of the field in the table linked to by the through_field to lookup.
	TargetFieldName NullableString `json:"target_field_name,omitempty"`
	FormulaType *FormulaTypeEnum `json:"formula_type,omitempty"`
}

// NewLookupFieldFieldSerializerWithRelatedFields instantiates a new LookupFieldFieldSerializerWithRelatedFields object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLookupFieldFieldSerializerWithRelatedFields(id int32, tableId int32, name string, order int32, type_ string, readOnly bool, relatedFields []Field, nullable bool) *LookupFieldFieldSerializerWithRelatedFields {
	this := LookupFieldFieldSerializerWithRelatedFields{}
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

// NewLookupFieldFieldSerializerWithRelatedFieldsWithDefaults instantiates a new LookupFieldFieldSerializerWithRelatedFields object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLookupFieldFieldSerializerWithRelatedFieldsWithDefaults() *LookupFieldFieldSerializerWithRelatedFields {
	this := LookupFieldFieldSerializerWithRelatedFields{}
	return &this
}

// GetId returns the Id field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetType(v string) {
	o.Type = v
}

// GetPrimary returns the Primary field value if set, zero value otherwise.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetPrimary() bool {
	if o == nil || IsNil(o.Primary) {
		var ret bool
		return ret
	}
	return *o.Primary
}

// GetPrimaryOk returns a tuple with the Primary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetPrimaryOk() (*bool, bool) {
	if o == nil || IsNil(o.Primary) {
		return nil, false
	}
	return o.Primary, true
}

// HasPrimary returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasPrimary() bool {
	if o != nil && !IsNil(o.Primary) {
		return true
	}

	return false
}

// SetPrimary gets a reference to the given bool and assigns it to the Primary field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetPrimary(v bool) {
	o.Primary = &v
}

// GetReadOnly returns the ReadOnly field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetReadOnly() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ReadOnly
}

// GetReadOnlyOk returns a tuple with the ReadOnly field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetReadOnlyOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReadOnly, true
}

// SetReadOnly sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetReadOnly(v bool) {
	o.ReadOnly = v
}

// GetRelatedFields returns the RelatedFields field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetRelatedFields() []Field {
	if o == nil {
		var ret []Field
		return ret
	}

	return o.RelatedFields
}

// GetRelatedFieldsOk returns a tuple with the RelatedFields field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetRelatedFieldsOk() ([]Field, bool) {
	if o == nil {
		return nil, false
	}
	return o.RelatedFields, true
}

// SetRelatedFields sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetRelatedFields(v []Field) {
	o.RelatedFields = v
}

// GetDateShowTzinfo returns the DateShowTzinfo field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateShowTzinfo() bool {
	if o == nil || IsNil(o.DateShowTzinfo.Get()) {
		var ret bool
		return ret
	}
	return *o.DateShowTzinfo.Get()
}

// GetDateShowTzinfoOk returns a tuple with the DateShowTzinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateShowTzinfoOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateShowTzinfo.Get(), o.DateShowTzinfo.IsSet()
}

// HasDateShowTzinfo returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasDateShowTzinfo() bool {
	if o != nil && o.DateShowTzinfo.IsSet() {
		return true
	}

	return false
}

// SetDateShowTzinfo gets a reference to the given NullableBool and assigns it to the DateShowTzinfo field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateShowTzinfo(v bool) {
	o.DateShowTzinfo.Set(&v)
}
// SetDateShowTzinfoNil sets the value for DateShowTzinfo to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateShowTzinfoNil() {
	o.DateShowTzinfo.Set(nil)
}

// UnsetDateShowTzinfo ensures that no value is present for DateShowTzinfo, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetDateShowTzinfo() {
	o.DateShowTzinfo.Unset()
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateFormat() CountFieldCreateFieldDateFormat {
	if o == nil || IsNil(o.DateFormat.Get()) {
		var ret CountFieldCreateFieldDateFormat
		return ret
	}
	return *o.DateFormat.Get()
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateFormatOk() (*CountFieldCreateFieldDateFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateFormat.Get(), o.DateFormat.IsSet()
}

// HasDateFormat returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasDateFormat() bool {
	if o != nil && o.DateFormat.IsSet() {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given NullableCountFieldCreateFieldDateFormat and assigns it to the DateFormat field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateFormat(v CountFieldCreateFieldDateFormat) {
	o.DateFormat.Set(&v)
}
// SetDateFormatNil sets the value for DateFormat to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateFormatNil() {
	o.DateFormat.Set(nil)
}

// UnsetDateFormat ensures that no value is present for DateFormat, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetDateFormat() {
	o.DateFormat.Unset()
}

// GetArrayFormulaType returns the ArrayFormulaType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetArrayFormulaType() CountFieldCreateFieldArrayFormulaType {
	if o == nil || IsNil(o.ArrayFormulaType.Get()) {
		var ret CountFieldCreateFieldArrayFormulaType
		return ret
	}
	return *o.ArrayFormulaType.Get()
}

// GetArrayFormulaTypeOk returns a tuple with the ArrayFormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetArrayFormulaTypeOk() (*CountFieldCreateFieldArrayFormulaType, bool) {
	if o == nil {
		return nil, false
	}
	return o.ArrayFormulaType.Get(), o.ArrayFormulaType.IsSet()
}

// HasArrayFormulaType returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasArrayFormulaType() bool {
	if o != nil && o.ArrayFormulaType.IsSet() {
		return true
	}

	return false
}

// SetArrayFormulaType gets a reference to the given NullableCountFieldCreateFieldArrayFormulaType and assigns it to the ArrayFormulaType field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetArrayFormulaType(v CountFieldCreateFieldArrayFormulaType) {
	o.ArrayFormulaType.Set(&v)
}
// SetArrayFormulaTypeNil sets the value for ArrayFormulaType to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetArrayFormulaTypeNil() {
	o.ArrayFormulaType.Set(nil)
}

// UnsetArrayFormulaType ensures that no value is present for ArrayFormulaType, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetArrayFormulaType() {
	o.ArrayFormulaType.Unset()
}

// GetDateForceTimezone returns the DateForceTimezone field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateForceTimezone() string {
	if o == nil || IsNil(o.DateForceTimezone.Get()) {
		var ret string
		return ret
	}
	return *o.DateForceTimezone.Get()
}

// GetDateForceTimezoneOk returns a tuple with the DateForceTimezone field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateForceTimezoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateForceTimezone.Get(), o.DateForceTimezone.IsSet()
}

// HasDateForceTimezone returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasDateForceTimezone() bool {
	if o != nil && o.DateForceTimezone.IsSet() {
		return true
	}

	return false
}

// SetDateForceTimezone gets a reference to the given NullableString and assigns it to the DateForceTimezone field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateForceTimezone(v string) {
	o.DateForceTimezone.Set(&v)
}
// SetDateForceTimezoneNil sets the value for DateForceTimezone to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateForceTimezoneNil() {
	o.DateForceTimezone.Set(nil)
}

// UnsetDateForceTimezone ensures that no value is present for DateForceTimezone, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetDateForceTimezone() {
	o.DateForceTimezone.Unset()
}

// GetDateIncludeTime returns the DateIncludeTime field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateIncludeTime() bool {
	if o == nil || IsNil(o.DateIncludeTime.Get()) {
		var ret bool
		return ret
	}
	return *o.DateIncludeTime.Get()
}

// GetDateIncludeTimeOk returns a tuple with the DateIncludeTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateIncludeTimeOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateIncludeTime.Get(), o.DateIncludeTime.IsSet()
}

// HasDateIncludeTime returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasDateIncludeTime() bool {
	if o != nil && o.DateIncludeTime.IsSet() {
		return true
	}

	return false
}

// SetDateIncludeTime gets a reference to the given NullableBool and assigns it to the DateIncludeTime field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateIncludeTime(v bool) {
	o.DateIncludeTime.Set(&v)
}
// SetDateIncludeTimeNil sets the value for DateIncludeTime to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateIncludeTimeNil() {
	o.DateIncludeTime.Set(nil)
}

// UnsetDateIncludeTime ensures that no value is present for DateIncludeTime, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetDateIncludeTime() {
	o.DateIncludeTime.Unset()
}

// GetNullable returns the Nullable field value
func (o *LookupFieldFieldSerializerWithRelatedFields) GetNullable() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Nullable
}

// GetNullableOk returns a tuple with the Nullable field value
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetNullableOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Nullable, true
}

// SetNullable sets field value
func (o *LookupFieldFieldSerializerWithRelatedFields) SetNullable(v bool) {
	o.Nullable = v
}

// GetDateTimeFormat returns the DateTimeFormat field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateTimeFormat() CountFieldCreateFieldDateTimeFormat {
	if o == nil || IsNil(o.DateTimeFormat.Get()) {
		var ret CountFieldCreateFieldDateTimeFormat
		return ret
	}
	return *o.DateTimeFormat.Get()
}

// GetDateTimeFormatOk returns a tuple with the DateTimeFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetDateTimeFormatOk() (*CountFieldCreateFieldDateTimeFormat, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateTimeFormat.Get(), o.DateTimeFormat.IsSet()
}

// HasDateTimeFormat returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasDateTimeFormat() bool {
	if o != nil && o.DateTimeFormat.IsSet() {
		return true
	}

	return false
}

// SetDateTimeFormat gets a reference to the given NullableCountFieldCreateFieldDateTimeFormat and assigns it to the DateTimeFormat field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateTimeFormat(v CountFieldCreateFieldDateTimeFormat) {
	o.DateTimeFormat.Set(&v)
}
// SetDateTimeFormatNil sets the value for DateTimeFormat to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetDateTimeFormatNil() {
	o.DateTimeFormat.Set(nil)
}

// UnsetDateTimeFormat ensures that no value is present for DateTimeFormat, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetDateTimeFormat() {
	o.DateTimeFormat.Unset()
}

// GetNumberDecimalPlaces returns the NumberDecimalPlaces field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetNumberDecimalPlaces() CountFieldCreateFieldNumberDecimalPlaces {
	if o == nil || IsNil(o.NumberDecimalPlaces.Get()) {
		var ret CountFieldCreateFieldNumberDecimalPlaces
		return ret
	}
	return *o.NumberDecimalPlaces.Get()
}

// GetNumberDecimalPlacesOk returns a tuple with the NumberDecimalPlaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetNumberDecimalPlacesOk() (*CountFieldCreateFieldNumberDecimalPlaces, bool) {
	if o == nil {
		return nil, false
	}
	return o.NumberDecimalPlaces.Get(), o.NumberDecimalPlaces.IsSet()
}

// HasNumberDecimalPlaces returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasNumberDecimalPlaces() bool {
	if o != nil && o.NumberDecimalPlaces.IsSet() {
		return true
	}

	return false
}

// SetNumberDecimalPlaces gets a reference to the given NullableCountFieldCreateFieldNumberDecimalPlaces and assigns it to the NumberDecimalPlaces field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetNumberDecimalPlaces(v CountFieldCreateFieldNumberDecimalPlaces) {
	o.NumberDecimalPlaces.Set(&v)
}
// SetNumberDecimalPlacesNil sets the value for NumberDecimalPlaces to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetNumberDecimalPlacesNil() {
	o.NumberDecimalPlaces.Set(nil)
}

// UnsetNumberDecimalPlaces ensures that no value is present for NumberDecimalPlaces, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetNumberDecimalPlaces() {
	o.NumberDecimalPlaces.Unset()
}

// GetError returns the Error field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetError() string {
	if o == nil || IsNil(o.Error.Get()) {
		var ret string
		return ret
	}
	return *o.Error.Get()
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetErrorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Error.Get(), o.Error.IsSet()
}

// HasError returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasError() bool {
	if o != nil && o.Error.IsSet() {
		return true
	}

	return false
}

// SetError gets a reference to the given NullableString and assigns it to the Error field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetError(v string) {
	o.Error.Set(&v)
}
// SetErrorNil sets the value for Error to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetErrorNil() {
	o.Error.Set(nil)
}

// UnsetError ensures that no value is present for Error, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetError() {
	o.Error.Unset()
}

// GetThroughFieldId returns the ThroughFieldId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetThroughFieldId() int32 {
	if o == nil || IsNil(o.ThroughFieldId.Get()) {
		var ret int32
		return ret
	}
	return *o.ThroughFieldId.Get()
}

// GetThroughFieldIdOk returns a tuple with the ThroughFieldId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetThroughFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ThroughFieldId.Get(), o.ThroughFieldId.IsSet()
}

// HasThroughFieldId returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasThroughFieldId() bool {
	if o != nil && o.ThroughFieldId.IsSet() {
		return true
	}

	return false
}

// SetThroughFieldId gets a reference to the given NullableInt32 and assigns it to the ThroughFieldId field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetThroughFieldId(v int32) {
	o.ThroughFieldId.Set(&v)
}
// SetThroughFieldIdNil sets the value for ThroughFieldId to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetThroughFieldIdNil() {
	o.ThroughFieldId.Set(nil)
}

// UnsetThroughFieldId ensures that no value is present for ThroughFieldId, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetThroughFieldId() {
	o.ThroughFieldId.Unset()
}

// GetThroughFieldName returns the ThroughFieldName field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetThroughFieldName() string {
	if o == nil || IsNil(o.ThroughFieldName.Get()) {
		var ret string
		return ret
	}
	return *o.ThroughFieldName.Get()
}

// GetThroughFieldNameOk returns a tuple with the ThroughFieldName field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetThroughFieldNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ThroughFieldName.Get(), o.ThroughFieldName.IsSet()
}

// HasThroughFieldName returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasThroughFieldName() bool {
	if o != nil && o.ThroughFieldName.IsSet() {
		return true
	}

	return false
}

// SetThroughFieldName gets a reference to the given NullableString and assigns it to the ThroughFieldName field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetThroughFieldName(v string) {
	o.ThroughFieldName.Set(&v)
}
// SetThroughFieldNameNil sets the value for ThroughFieldName to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetThroughFieldNameNil() {
	o.ThroughFieldName.Set(nil)
}

// UnsetThroughFieldName ensures that no value is present for ThroughFieldName, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetThroughFieldName() {
	o.ThroughFieldName.Unset()
}

// GetTargetFieldId returns the TargetFieldId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTargetFieldId() int32 {
	if o == nil || IsNil(o.TargetFieldId.Get()) {
		var ret int32
		return ret
	}
	return *o.TargetFieldId.Get()
}

// GetTargetFieldIdOk returns a tuple with the TargetFieldId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTargetFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.TargetFieldId.Get(), o.TargetFieldId.IsSet()
}

// HasTargetFieldId returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasTargetFieldId() bool {
	if o != nil && o.TargetFieldId.IsSet() {
		return true
	}

	return false
}

// SetTargetFieldId gets a reference to the given NullableInt32 and assigns it to the TargetFieldId field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetTargetFieldId(v int32) {
	o.TargetFieldId.Set(&v)
}
// SetTargetFieldIdNil sets the value for TargetFieldId to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetTargetFieldIdNil() {
	o.TargetFieldId.Set(nil)
}

// UnsetTargetFieldId ensures that no value is present for TargetFieldId, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetTargetFieldId() {
	o.TargetFieldId.Unset()
}

// GetTargetFieldName returns the TargetFieldName field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTargetFieldName() string {
	if o == nil || IsNil(o.TargetFieldName.Get()) {
		var ret string
		return ret
	}
	return *o.TargetFieldName.Get()
}

// GetTargetFieldNameOk returns a tuple with the TargetFieldName field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LookupFieldFieldSerializerWithRelatedFields) GetTargetFieldNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.TargetFieldName.Get(), o.TargetFieldName.IsSet()
}

// HasTargetFieldName returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasTargetFieldName() bool {
	if o != nil && o.TargetFieldName.IsSet() {
		return true
	}

	return false
}

// SetTargetFieldName gets a reference to the given NullableString and assigns it to the TargetFieldName field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetTargetFieldName(v string) {
	o.TargetFieldName.Set(&v)
}
// SetTargetFieldNameNil sets the value for TargetFieldName to be an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) SetTargetFieldNameNil() {
	o.TargetFieldName.Set(nil)
}

// UnsetTargetFieldName ensures that no value is present for TargetFieldName, not even an explicit nil
func (o *LookupFieldFieldSerializerWithRelatedFields) UnsetTargetFieldName() {
	o.TargetFieldName.Unset()
}

// GetFormulaType returns the FormulaType field value if set, zero value otherwise.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetFormulaType() FormulaTypeEnum {
	if o == nil || IsNil(o.FormulaType) {
		var ret FormulaTypeEnum
		return ret
	}
	return *o.FormulaType
}

// GetFormulaTypeOk returns a tuple with the FormulaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) GetFormulaTypeOk() (*FormulaTypeEnum, bool) {
	if o == nil || IsNil(o.FormulaType) {
		return nil, false
	}
	return o.FormulaType, true
}

// HasFormulaType returns a boolean if a field has been set.
func (o *LookupFieldFieldSerializerWithRelatedFields) HasFormulaType() bool {
	if o != nil && !IsNil(o.FormulaType) {
		return true
	}

	return false
}

// SetFormulaType gets a reference to the given FormulaTypeEnum and assigns it to the FormulaType field.
func (o *LookupFieldFieldSerializerWithRelatedFields) SetFormulaType(v FormulaTypeEnum) {
	o.FormulaType = &v
}

func (o LookupFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LookupFieldFieldSerializerWithRelatedFields) ToMap() (map[string]interface{}, error) {
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
	if o.ThroughFieldName.IsSet() {
		toSerialize["through_field_name"] = o.ThroughFieldName.Get()
	}
	if o.TargetFieldId.IsSet() {
		toSerialize["target_field_id"] = o.TargetFieldId.Get()
	}
	if o.TargetFieldName.IsSet() {
		toSerialize["target_field_name"] = o.TargetFieldName.Get()
	}
	if !IsNil(o.FormulaType) {
		toSerialize["formula_type"] = o.FormulaType
	}
	return toSerialize, nil
}

type NullableLookupFieldFieldSerializerWithRelatedFields struct {
	value *LookupFieldFieldSerializerWithRelatedFields
	isSet bool
}

func (v NullableLookupFieldFieldSerializerWithRelatedFields) Get() *LookupFieldFieldSerializerWithRelatedFields {
	return v.value
}

func (v *NullableLookupFieldFieldSerializerWithRelatedFields) Set(val *LookupFieldFieldSerializerWithRelatedFields) {
	v.value = val
	v.isSet = true
}

func (v NullableLookupFieldFieldSerializerWithRelatedFields) IsSet() bool {
	return v.isSet
}

func (v *NullableLookupFieldFieldSerializerWithRelatedFields) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLookupFieldFieldSerializerWithRelatedFields(val *LookupFieldFieldSerializerWithRelatedFields) *NullableLookupFieldFieldSerializerWithRelatedFields {
	return &NullableLookupFieldFieldSerializerWithRelatedFields{value: val, isSet: true}
}

func (v NullableLookupFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLookupFieldFieldSerializerWithRelatedFields) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


