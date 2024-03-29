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

// checks if the CsvExporterOptions type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CsvExporterOptions{}

// CsvExporterOptions struct for CsvExporterOptions
type CsvExporterOptions struct {
	// Optional: The view for this table to export using its filters, sorts and other view specific settings.
	ViewId NullableInt32 `json:"view_id,omitempty"`
	ExporterType ExporterTypeEnum `json:"exporter_type"`
	ExportCharset *ExportCharsetEnum `json:"export_charset,omitempty"`
	CsvColumnSeparator *CsvColumnSeparatorEnum `json:"csv_column_separator,omitempty"`
	// Whether or not to generate a header row at the top of the csv file.
	CsvIncludeHeader *bool `json:"csv_include_header,omitempty"`
}

// NewCsvExporterOptions instantiates a new CsvExporterOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCsvExporterOptions(exporterType ExporterTypeEnum) *CsvExporterOptions {
	this := CsvExporterOptions{}
	this.ExporterType = exporterType
	var exportCharset ExportCharsetEnum = UTF_8
	this.ExportCharset = &exportCharset
	var csvColumnSeparator CsvColumnSeparatorEnum = COMMA
	this.CsvColumnSeparator = &csvColumnSeparator
	var csvIncludeHeader bool = true
	this.CsvIncludeHeader = &csvIncludeHeader
	return &this
}

// NewCsvExporterOptionsWithDefaults instantiates a new CsvExporterOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCsvExporterOptionsWithDefaults() *CsvExporterOptions {
	this := CsvExporterOptions{}
	var exportCharset ExportCharsetEnum = UTF_8
	this.ExportCharset = &exportCharset
	var csvColumnSeparator CsvColumnSeparatorEnum = COMMA
	this.CsvColumnSeparator = &csvColumnSeparator
	var csvIncludeHeader bool = true
	this.CsvIncludeHeader = &csvIncludeHeader
	return &this
}

// GetViewId returns the ViewId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CsvExporterOptions) GetViewId() int32 {
	if o == nil || IsNil(o.ViewId.Get()) {
		var ret int32
		return ret
	}
	return *o.ViewId.Get()
}

// GetViewIdOk returns a tuple with the ViewId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CsvExporterOptions) GetViewIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ViewId.Get(), o.ViewId.IsSet()
}

// HasViewId returns a boolean if a field has been set.
func (o *CsvExporterOptions) HasViewId() bool {
	if o != nil && o.ViewId.IsSet() {
		return true
	}

	return false
}

// SetViewId gets a reference to the given NullableInt32 and assigns it to the ViewId field.
func (o *CsvExporterOptions) SetViewId(v int32) {
	o.ViewId.Set(&v)
}
// SetViewIdNil sets the value for ViewId to be an explicit nil
func (o *CsvExporterOptions) SetViewIdNil() {
	o.ViewId.Set(nil)
}

// UnsetViewId ensures that no value is present for ViewId, not even an explicit nil
func (o *CsvExporterOptions) UnsetViewId() {
	o.ViewId.Unset()
}

// GetExporterType returns the ExporterType field value
func (o *CsvExporterOptions) GetExporterType() ExporterTypeEnum {
	if o == nil {
		var ret ExporterTypeEnum
		return ret
	}

	return o.ExporterType
}

// GetExporterTypeOk returns a tuple with the ExporterType field value
// and a boolean to check if the value has been set.
func (o *CsvExporterOptions) GetExporterTypeOk() (*ExporterTypeEnum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExporterType, true
}

// SetExporterType sets field value
func (o *CsvExporterOptions) SetExporterType(v ExporterTypeEnum) {
	o.ExporterType = v
}

// GetExportCharset returns the ExportCharset field value if set, zero value otherwise.
func (o *CsvExporterOptions) GetExportCharset() ExportCharsetEnum {
	if o == nil || IsNil(o.ExportCharset) {
		var ret ExportCharsetEnum
		return ret
	}
	return *o.ExportCharset
}

// GetExportCharsetOk returns a tuple with the ExportCharset field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CsvExporterOptions) GetExportCharsetOk() (*ExportCharsetEnum, bool) {
	if o == nil || IsNil(o.ExportCharset) {
		return nil, false
	}
	return o.ExportCharset, true
}

// HasExportCharset returns a boolean if a field has been set.
func (o *CsvExporterOptions) HasExportCharset() bool {
	if o != nil && !IsNil(o.ExportCharset) {
		return true
	}

	return false
}

// SetExportCharset gets a reference to the given ExportCharsetEnum and assigns it to the ExportCharset field.
func (o *CsvExporterOptions) SetExportCharset(v ExportCharsetEnum) {
	o.ExportCharset = &v
}

// GetCsvColumnSeparator returns the CsvColumnSeparator field value if set, zero value otherwise.
func (o *CsvExporterOptions) GetCsvColumnSeparator() CsvColumnSeparatorEnum {
	if o == nil || IsNil(o.CsvColumnSeparator) {
		var ret CsvColumnSeparatorEnum
		return ret
	}
	return *o.CsvColumnSeparator
}

// GetCsvColumnSeparatorOk returns a tuple with the CsvColumnSeparator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CsvExporterOptions) GetCsvColumnSeparatorOk() (*CsvColumnSeparatorEnum, bool) {
	if o == nil || IsNil(o.CsvColumnSeparator) {
		return nil, false
	}
	return o.CsvColumnSeparator, true
}

// HasCsvColumnSeparator returns a boolean if a field has been set.
func (o *CsvExporterOptions) HasCsvColumnSeparator() bool {
	if o != nil && !IsNil(o.CsvColumnSeparator) {
		return true
	}

	return false
}

// SetCsvColumnSeparator gets a reference to the given CsvColumnSeparatorEnum and assigns it to the CsvColumnSeparator field.
func (o *CsvExporterOptions) SetCsvColumnSeparator(v CsvColumnSeparatorEnum) {
	o.CsvColumnSeparator = &v
}

// GetCsvIncludeHeader returns the CsvIncludeHeader field value if set, zero value otherwise.
func (o *CsvExporterOptions) GetCsvIncludeHeader() bool {
	if o == nil || IsNil(o.CsvIncludeHeader) {
		var ret bool
		return ret
	}
	return *o.CsvIncludeHeader
}

// GetCsvIncludeHeaderOk returns a tuple with the CsvIncludeHeader field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CsvExporterOptions) GetCsvIncludeHeaderOk() (*bool, bool) {
	if o == nil || IsNil(o.CsvIncludeHeader) {
		return nil, false
	}
	return o.CsvIncludeHeader, true
}

// HasCsvIncludeHeader returns a boolean if a field has been set.
func (o *CsvExporterOptions) HasCsvIncludeHeader() bool {
	if o != nil && !IsNil(o.CsvIncludeHeader) {
		return true
	}

	return false
}

// SetCsvIncludeHeader gets a reference to the given bool and assigns it to the CsvIncludeHeader field.
func (o *CsvExporterOptions) SetCsvIncludeHeader(v bool) {
	o.CsvIncludeHeader = &v
}

func (o CsvExporterOptions) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CsvExporterOptions) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.ViewId.IsSet() {
		toSerialize["view_id"] = o.ViewId.Get()
	}
	toSerialize["exporter_type"] = o.ExporterType
	if !IsNil(o.ExportCharset) {
		toSerialize["export_charset"] = o.ExportCharset
	}
	if !IsNil(o.CsvColumnSeparator) {
		toSerialize["csv_column_separator"] = o.CsvColumnSeparator
	}
	if !IsNil(o.CsvIncludeHeader) {
		toSerialize["csv_include_header"] = o.CsvIncludeHeader
	}
	return toSerialize, nil
}

type NullableCsvExporterOptions struct {
	value *CsvExporterOptions
	isSet bool
}

func (v NullableCsvExporterOptions) Get() *CsvExporterOptions {
	return v.value
}

func (v *NullableCsvExporterOptions) Set(val *CsvExporterOptions) {
	v.value = val
	v.isSet = true
}

func (v NullableCsvExporterOptions) IsSet() bool {
	return v.isSet
}

func (v *NullableCsvExporterOptions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCsvExporterOptions(val *CsvExporterOptions) *NullableCsvExporterOptions {
	return &NullableCsvExporterOptions{value: val, isSet: true}
}

func (v NullableCsvExporterOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCsvExporterOptions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


