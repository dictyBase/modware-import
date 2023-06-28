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

// checks if the TableWebhookCall type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TableWebhookCall{}

// TableWebhookCall struct for TableWebhookCall
type TableWebhookCall struct {
	Id int32 `json:"id"`
	// Event ID where the call originated from.
	EventId string `json:"event_id"`
	EventType string `json:"event_type"`
	CalledTime NullableTime `json:"called_time,omitempty"`
	CalledUrl string `json:"called_url"`
	// A text copy of the request headers and body.
	Request NullableString `json:"request,omitempty"`
	// A text copy of the response headers and body.
	Response NullableString `json:"response,omitempty"`
	// The HTTP response status code.
	ResponseStatus NullableInt32 `json:"response_status,omitempty"`
	// An internal error reflecting what went wrong.
	Error NullableString `json:"error,omitempty"`
}

// NewTableWebhookCall instantiates a new TableWebhookCall object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTableWebhookCall(id int32, eventId string, eventType string, calledUrl string) *TableWebhookCall {
	this := TableWebhookCall{}
	this.Id = id
	this.EventId = eventId
	this.EventType = eventType
	this.CalledUrl = calledUrl
	return &this
}

// NewTableWebhookCallWithDefaults instantiates a new TableWebhookCall object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTableWebhookCallWithDefaults() *TableWebhookCall {
	this := TableWebhookCall{}
	return &this
}

// GetId returns the Id field value
func (o *TableWebhookCall) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *TableWebhookCall) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *TableWebhookCall) SetId(v int32) {
	o.Id = v
}

// GetEventId returns the EventId field value
func (o *TableWebhookCall) GetEventId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EventId
}

// GetEventIdOk returns a tuple with the EventId field value
// and a boolean to check if the value has been set.
func (o *TableWebhookCall) GetEventIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventId, true
}

// SetEventId sets field value
func (o *TableWebhookCall) SetEventId(v string) {
	o.EventId = v
}

// GetEventType returns the EventType field value
func (o *TableWebhookCall) GetEventType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EventType
}

// GetEventTypeOk returns a tuple with the EventType field value
// and a boolean to check if the value has been set.
func (o *TableWebhookCall) GetEventTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventType, true
}

// SetEventType sets field value
func (o *TableWebhookCall) SetEventType(v string) {
	o.EventType = v
}

// GetCalledTime returns the CalledTime field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TableWebhookCall) GetCalledTime() time.Time {
	if o == nil || IsNil(o.CalledTime.Get()) {
		var ret time.Time
		return ret
	}
	return *o.CalledTime.Get()
}

// GetCalledTimeOk returns a tuple with the CalledTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TableWebhookCall) GetCalledTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.CalledTime.Get(), o.CalledTime.IsSet()
}

// HasCalledTime returns a boolean if a field has been set.
func (o *TableWebhookCall) HasCalledTime() bool {
	if o != nil && o.CalledTime.IsSet() {
		return true
	}

	return false
}

// SetCalledTime gets a reference to the given NullableTime and assigns it to the CalledTime field.
func (o *TableWebhookCall) SetCalledTime(v time.Time) {
	o.CalledTime.Set(&v)
}
// SetCalledTimeNil sets the value for CalledTime to be an explicit nil
func (o *TableWebhookCall) SetCalledTimeNil() {
	o.CalledTime.Set(nil)
}

// UnsetCalledTime ensures that no value is present for CalledTime, not even an explicit nil
func (o *TableWebhookCall) UnsetCalledTime() {
	o.CalledTime.Unset()
}

// GetCalledUrl returns the CalledUrl field value
func (o *TableWebhookCall) GetCalledUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CalledUrl
}

// GetCalledUrlOk returns a tuple with the CalledUrl field value
// and a boolean to check if the value has been set.
func (o *TableWebhookCall) GetCalledUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CalledUrl, true
}

// SetCalledUrl sets field value
func (o *TableWebhookCall) SetCalledUrl(v string) {
	o.CalledUrl = v
}

// GetRequest returns the Request field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TableWebhookCall) GetRequest() string {
	if o == nil || IsNil(o.Request.Get()) {
		var ret string
		return ret
	}
	return *o.Request.Get()
}

// GetRequestOk returns a tuple with the Request field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TableWebhookCall) GetRequestOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Request.Get(), o.Request.IsSet()
}

// HasRequest returns a boolean if a field has been set.
func (o *TableWebhookCall) HasRequest() bool {
	if o != nil && o.Request.IsSet() {
		return true
	}

	return false
}

// SetRequest gets a reference to the given NullableString and assigns it to the Request field.
func (o *TableWebhookCall) SetRequest(v string) {
	o.Request.Set(&v)
}
// SetRequestNil sets the value for Request to be an explicit nil
func (o *TableWebhookCall) SetRequestNil() {
	o.Request.Set(nil)
}

// UnsetRequest ensures that no value is present for Request, not even an explicit nil
func (o *TableWebhookCall) UnsetRequest() {
	o.Request.Unset()
}

// GetResponse returns the Response field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TableWebhookCall) GetResponse() string {
	if o == nil || IsNil(o.Response.Get()) {
		var ret string
		return ret
	}
	return *o.Response.Get()
}

// GetResponseOk returns a tuple with the Response field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TableWebhookCall) GetResponseOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Response.Get(), o.Response.IsSet()
}

// HasResponse returns a boolean if a field has been set.
func (o *TableWebhookCall) HasResponse() bool {
	if o != nil && o.Response.IsSet() {
		return true
	}

	return false
}

// SetResponse gets a reference to the given NullableString and assigns it to the Response field.
func (o *TableWebhookCall) SetResponse(v string) {
	o.Response.Set(&v)
}
// SetResponseNil sets the value for Response to be an explicit nil
func (o *TableWebhookCall) SetResponseNil() {
	o.Response.Set(nil)
}

// UnsetResponse ensures that no value is present for Response, not even an explicit nil
func (o *TableWebhookCall) UnsetResponse() {
	o.Response.Unset()
}

// GetResponseStatus returns the ResponseStatus field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TableWebhookCall) GetResponseStatus() int32 {
	if o == nil || IsNil(o.ResponseStatus.Get()) {
		var ret int32
		return ret
	}
	return *o.ResponseStatus.Get()
}

// GetResponseStatusOk returns a tuple with the ResponseStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TableWebhookCall) GetResponseStatusOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResponseStatus.Get(), o.ResponseStatus.IsSet()
}

// HasResponseStatus returns a boolean if a field has been set.
func (o *TableWebhookCall) HasResponseStatus() bool {
	if o != nil && o.ResponseStatus.IsSet() {
		return true
	}

	return false
}

// SetResponseStatus gets a reference to the given NullableInt32 and assigns it to the ResponseStatus field.
func (o *TableWebhookCall) SetResponseStatus(v int32) {
	o.ResponseStatus.Set(&v)
}
// SetResponseStatusNil sets the value for ResponseStatus to be an explicit nil
func (o *TableWebhookCall) SetResponseStatusNil() {
	o.ResponseStatus.Set(nil)
}

// UnsetResponseStatus ensures that no value is present for ResponseStatus, not even an explicit nil
func (o *TableWebhookCall) UnsetResponseStatus() {
	o.ResponseStatus.Unset()
}

// GetError returns the Error field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TableWebhookCall) GetError() string {
	if o == nil || IsNil(o.Error.Get()) {
		var ret string
		return ret
	}
	return *o.Error.Get()
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TableWebhookCall) GetErrorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Error.Get(), o.Error.IsSet()
}

// HasError returns a boolean if a field has been set.
func (o *TableWebhookCall) HasError() bool {
	if o != nil && o.Error.IsSet() {
		return true
	}

	return false
}

// SetError gets a reference to the given NullableString and assigns it to the Error field.
func (o *TableWebhookCall) SetError(v string) {
	o.Error.Set(&v)
}
// SetErrorNil sets the value for Error to be an explicit nil
func (o *TableWebhookCall) SetErrorNil() {
	o.Error.Set(nil)
}

// UnsetError ensures that no value is present for Error, not even an explicit nil
func (o *TableWebhookCall) UnsetError() {
	o.Error.Unset()
}

func (o TableWebhookCall) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TableWebhookCall) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: event_id is readOnly
	toSerialize["event_type"] = o.EventType
	if o.CalledTime.IsSet() {
		toSerialize["called_time"] = o.CalledTime.Get()
	}
	toSerialize["called_url"] = o.CalledUrl
	if o.Request.IsSet() {
		toSerialize["request"] = o.Request.Get()
	}
	if o.Response.IsSet() {
		toSerialize["response"] = o.Response.Get()
	}
	if o.ResponseStatus.IsSet() {
		toSerialize["response_status"] = o.ResponseStatus.Get()
	}
	if o.Error.IsSet() {
		toSerialize["error"] = o.Error.Get()
	}
	return toSerialize, nil
}

type NullableTableWebhookCall struct {
	value *TableWebhookCall
	isSet bool
}

func (v NullableTableWebhookCall) Get() *TableWebhookCall {
	return v.value
}

func (v *NullableTableWebhookCall) Set(val *TableWebhookCall) {
	v.value = val
	v.isSet = true
}

func (v NullableTableWebhookCall) IsSet() bool {
	return v.isSet
}

func (v *NullableTableWebhookCall) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTableWebhookCall(val *TableWebhookCall) *NullableTableWebhookCall {
	return &NullableTableWebhookCall{value: val, isSet: true}
}

func (v NullableTableWebhookCall) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTableWebhookCall) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


