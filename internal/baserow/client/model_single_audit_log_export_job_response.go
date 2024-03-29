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

// checks if the SingleAuditLogExportJobResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SingleAuditLogExportJobResponse{}

// SingleAuditLogExportJobResponse When mixed in to a model serializer for an ExportJob this will add an url field with the actual usable url of the export job's file (if it has one).
type SingleAuditLogExportJobResponse struct {
	Url string `json:"url"`
	ExportCharset *ExportCharsetEnum `json:"export_charset,omitempty"`
	CsvColumnSeparator *CsvColumnSeparatorEnum `json:"csv_column_separator,omitempty"`
	// Whether or not to generate a header row at the top of the csv file.
	CsvFirstRowHeader *bool `json:"csv_first_row_header,omitempty"`
	// Optional: The user to filter the audit log by.
	FilterUserId *int32 `json:"filter_user_id,omitempty"`
	// Optional: The workspace to filter the audit log by.
	FilterWorkspaceId *int32 `json:"filter_workspace_id,omitempty"`
	FilterActionType *FilterActionTypeEnum `json:"filter_action_type,omitempty"`
	// Optional: The start date to filter the audit log by.
	FilterFromTimestamp *time.Time `json:"filter_from_timestamp,omitempty"`
	// Optional: The end date to filter the audit log by.
	FilterToTimestamp *time.Time `json:"filter_to_timestamp,omitempty"`
	// The date and time when the export job was created.
	CreatedOn time.Time `json:"created_on"`
}

// NewSingleAuditLogExportJobResponse instantiates a new SingleAuditLogExportJobResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSingleAuditLogExportJobResponse(url string, createdOn time.Time) *SingleAuditLogExportJobResponse {
	this := SingleAuditLogExportJobResponse{}
	this.Url = url
	var exportCharset ExportCharsetEnum = UTF_8
	this.ExportCharset = &exportCharset
	var csvColumnSeparator CsvColumnSeparatorEnum = COMMA
	this.CsvColumnSeparator = &csvColumnSeparator
	var csvFirstRowHeader bool = true
	this.CsvFirstRowHeader = &csvFirstRowHeader
	this.CreatedOn = createdOn
	return &this
}

// NewSingleAuditLogExportJobResponseWithDefaults instantiates a new SingleAuditLogExportJobResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSingleAuditLogExportJobResponseWithDefaults() *SingleAuditLogExportJobResponse {
	this := SingleAuditLogExportJobResponse{}
	var exportCharset ExportCharsetEnum = UTF_8
	this.ExportCharset = &exportCharset
	var csvColumnSeparator CsvColumnSeparatorEnum = COMMA
	this.CsvColumnSeparator = &csvColumnSeparator
	var csvFirstRowHeader bool = true
	this.CsvFirstRowHeader = &csvFirstRowHeader
	return &this
}

// GetUrl returns the Url field value
func (o *SingleAuditLogExportJobResponse) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *SingleAuditLogExportJobResponse) SetUrl(v string) {
	o.Url = v
}

// GetExportCharset returns the ExportCharset field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetExportCharset() ExportCharsetEnum {
	if o == nil || IsNil(o.ExportCharset) {
		var ret ExportCharsetEnum
		return ret
	}
	return *o.ExportCharset
}

// GetExportCharsetOk returns a tuple with the ExportCharset field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetExportCharsetOk() (*ExportCharsetEnum, bool) {
	if o == nil || IsNil(o.ExportCharset) {
		return nil, false
	}
	return o.ExportCharset, true
}

// HasExportCharset returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasExportCharset() bool {
	if o != nil && !IsNil(o.ExportCharset) {
		return true
	}

	return false
}

// SetExportCharset gets a reference to the given ExportCharsetEnum and assigns it to the ExportCharset field.
func (o *SingleAuditLogExportJobResponse) SetExportCharset(v ExportCharsetEnum) {
	o.ExportCharset = &v
}

// GetCsvColumnSeparator returns the CsvColumnSeparator field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetCsvColumnSeparator() CsvColumnSeparatorEnum {
	if o == nil || IsNil(o.CsvColumnSeparator) {
		var ret CsvColumnSeparatorEnum
		return ret
	}
	return *o.CsvColumnSeparator
}

// GetCsvColumnSeparatorOk returns a tuple with the CsvColumnSeparator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetCsvColumnSeparatorOk() (*CsvColumnSeparatorEnum, bool) {
	if o == nil || IsNil(o.CsvColumnSeparator) {
		return nil, false
	}
	return o.CsvColumnSeparator, true
}

// HasCsvColumnSeparator returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasCsvColumnSeparator() bool {
	if o != nil && !IsNil(o.CsvColumnSeparator) {
		return true
	}

	return false
}

// SetCsvColumnSeparator gets a reference to the given CsvColumnSeparatorEnum and assigns it to the CsvColumnSeparator field.
func (o *SingleAuditLogExportJobResponse) SetCsvColumnSeparator(v CsvColumnSeparatorEnum) {
	o.CsvColumnSeparator = &v
}

// GetCsvFirstRowHeader returns the CsvFirstRowHeader field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetCsvFirstRowHeader() bool {
	if o == nil || IsNil(o.CsvFirstRowHeader) {
		var ret bool
		return ret
	}
	return *o.CsvFirstRowHeader
}

// GetCsvFirstRowHeaderOk returns a tuple with the CsvFirstRowHeader field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetCsvFirstRowHeaderOk() (*bool, bool) {
	if o == nil || IsNil(o.CsvFirstRowHeader) {
		return nil, false
	}
	return o.CsvFirstRowHeader, true
}

// HasCsvFirstRowHeader returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasCsvFirstRowHeader() bool {
	if o != nil && !IsNil(o.CsvFirstRowHeader) {
		return true
	}

	return false
}

// SetCsvFirstRowHeader gets a reference to the given bool and assigns it to the CsvFirstRowHeader field.
func (o *SingleAuditLogExportJobResponse) SetCsvFirstRowHeader(v bool) {
	o.CsvFirstRowHeader = &v
}

// GetFilterUserId returns the FilterUserId field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetFilterUserId() int32 {
	if o == nil || IsNil(o.FilterUserId) {
		var ret int32
		return ret
	}
	return *o.FilterUserId
}

// GetFilterUserIdOk returns a tuple with the FilterUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetFilterUserIdOk() (*int32, bool) {
	if o == nil || IsNil(o.FilterUserId) {
		return nil, false
	}
	return o.FilterUserId, true
}

// HasFilterUserId returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasFilterUserId() bool {
	if o != nil && !IsNil(o.FilterUserId) {
		return true
	}

	return false
}

// SetFilterUserId gets a reference to the given int32 and assigns it to the FilterUserId field.
func (o *SingleAuditLogExportJobResponse) SetFilterUserId(v int32) {
	o.FilterUserId = &v
}

// GetFilterWorkspaceId returns the FilterWorkspaceId field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetFilterWorkspaceId() int32 {
	if o == nil || IsNil(o.FilterWorkspaceId) {
		var ret int32
		return ret
	}
	return *o.FilterWorkspaceId
}

// GetFilterWorkspaceIdOk returns a tuple with the FilterWorkspaceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetFilterWorkspaceIdOk() (*int32, bool) {
	if o == nil || IsNil(o.FilterWorkspaceId) {
		return nil, false
	}
	return o.FilterWorkspaceId, true
}

// HasFilterWorkspaceId returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasFilterWorkspaceId() bool {
	if o != nil && !IsNil(o.FilterWorkspaceId) {
		return true
	}

	return false
}

// SetFilterWorkspaceId gets a reference to the given int32 and assigns it to the FilterWorkspaceId field.
func (o *SingleAuditLogExportJobResponse) SetFilterWorkspaceId(v int32) {
	o.FilterWorkspaceId = &v
}

// GetFilterActionType returns the FilterActionType field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetFilterActionType() FilterActionTypeEnum {
	if o == nil || IsNil(o.FilterActionType) {
		var ret FilterActionTypeEnum
		return ret
	}
	return *o.FilterActionType
}

// GetFilterActionTypeOk returns a tuple with the FilterActionType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetFilterActionTypeOk() (*FilterActionTypeEnum, bool) {
	if o == nil || IsNil(o.FilterActionType) {
		return nil, false
	}
	return o.FilterActionType, true
}

// HasFilterActionType returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasFilterActionType() bool {
	if o != nil && !IsNil(o.FilterActionType) {
		return true
	}

	return false
}

// SetFilterActionType gets a reference to the given FilterActionTypeEnum and assigns it to the FilterActionType field.
func (o *SingleAuditLogExportJobResponse) SetFilterActionType(v FilterActionTypeEnum) {
	o.FilterActionType = &v
}

// GetFilterFromTimestamp returns the FilterFromTimestamp field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetFilterFromTimestamp() time.Time {
	if o == nil || IsNil(o.FilterFromTimestamp) {
		var ret time.Time
		return ret
	}
	return *o.FilterFromTimestamp
}

// GetFilterFromTimestampOk returns a tuple with the FilterFromTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetFilterFromTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FilterFromTimestamp) {
		return nil, false
	}
	return o.FilterFromTimestamp, true
}

// HasFilterFromTimestamp returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasFilterFromTimestamp() bool {
	if o != nil && !IsNil(o.FilterFromTimestamp) {
		return true
	}

	return false
}

// SetFilterFromTimestamp gets a reference to the given time.Time and assigns it to the FilterFromTimestamp field.
func (o *SingleAuditLogExportJobResponse) SetFilterFromTimestamp(v time.Time) {
	o.FilterFromTimestamp = &v
}

// GetFilterToTimestamp returns the FilterToTimestamp field value if set, zero value otherwise.
func (o *SingleAuditLogExportJobResponse) GetFilterToTimestamp() time.Time {
	if o == nil || IsNil(o.FilterToTimestamp) {
		var ret time.Time
		return ret
	}
	return *o.FilterToTimestamp
}

// GetFilterToTimestampOk returns a tuple with the FilterToTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetFilterToTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FilterToTimestamp) {
		return nil, false
	}
	return o.FilterToTimestamp, true
}

// HasFilterToTimestamp returns a boolean if a field has been set.
func (o *SingleAuditLogExportJobResponse) HasFilterToTimestamp() bool {
	if o != nil && !IsNil(o.FilterToTimestamp) {
		return true
	}

	return false
}

// SetFilterToTimestamp gets a reference to the given time.Time and assigns it to the FilterToTimestamp field.
func (o *SingleAuditLogExportJobResponse) SetFilterToTimestamp(v time.Time) {
	o.FilterToTimestamp = &v
}

// GetCreatedOn returns the CreatedOn field value
func (o *SingleAuditLogExportJobResponse) GetCreatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedOn
}

// GetCreatedOnOk returns a tuple with the CreatedOn field value
// and a boolean to check if the value has been set.
func (o *SingleAuditLogExportJobResponse) GetCreatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedOn, true
}

// SetCreatedOn sets field value
func (o *SingleAuditLogExportJobResponse) SetCreatedOn(v time.Time) {
	o.CreatedOn = v
}

func (o SingleAuditLogExportJobResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SingleAuditLogExportJobResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: url is readOnly
	if !IsNil(o.ExportCharset) {
		toSerialize["export_charset"] = o.ExportCharset
	}
	if !IsNil(o.CsvColumnSeparator) {
		toSerialize["csv_column_separator"] = o.CsvColumnSeparator
	}
	if !IsNil(o.CsvFirstRowHeader) {
		toSerialize["csv_first_row_header"] = o.CsvFirstRowHeader
	}
	if !IsNil(o.FilterUserId) {
		toSerialize["filter_user_id"] = o.FilterUserId
	}
	if !IsNil(o.FilterWorkspaceId) {
		toSerialize["filter_workspace_id"] = o.FilterWorkspaceId
	}
	if !IsNil(o.FilterActionType) {
		toSerialize["filter_action_type"] = o.FilterActionType
	}
	if !IsNil(o.FilterFromTimestamp) {
		toSerialize["filter_from_timestamp"] = o.FilterFromTimestamp
	}
	if !IsNil(o.FilterToTimestamp) {
		toSerialize["filter_to_timestamp"] = o.FilterToTimestamp
	}
	// skip: created_on is readOnly
	return toSerialize, nil
}

type NullableSingleAuditLogExportJobResponse struct {
	value *SingleAuditLogExportJobResponse
	isSet bool
}

func (v NullableSingleAuditLogExportJobResponse) Get() *SingleAuditLogExportJobResponse {
	return v.value
}

func (v *NullableSingleAuditLogExportJobResponse) Set(val *SingleAuditLogExportJobResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableSingleAuditLogExportJobResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableSingleAuditLogExportJobResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSingleAuditLogExportJobResponse(val *SingleAuditLogExportJobResponse) *NullableSingleAuditLogExportJobResponse {
	return &NullableSingleAuditLogExportJobResponse{value: val, isSet: true}
}

func (v NullableSingleAuditLogExportJobResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSingleAuditLogExportJobResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


