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

// checks if the PatchedTableWebhookUpdateRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedTableWebhookUpdateRequest{}

// PatchedTableWebhookUpdateRequest struct for PatchedTableWebhookUpdateRequest
type PatchedTableWebhookUpdateRequest struct {
	// The URL that must be called when the webhook is triggered.
	Url *string `json:"url,omitempty"`
	// Indicates whether this webhook should listen to all events.
	IncludeAllEvents *bool `json:"include_all_events,omitempty"`
	// A list containing the events that will trigger this webhook.
	Events []EventTypesEnum `json:"events,omitempty"`
	RequestMethod *RequestMethodEnum `json:"request_method,omitempty"`
	// The additional headers as an object where the key is the name and the value the value.
	Headers map[string]interface{} `json:"headers,omitempty"`
	// An internal name of the webhook.
	Name *string `json:"name,omitempty"`
	// Indicates whether the web hook is active. When a webhook has failed multiple times, it will automatically be deactivated.
	Active *bool `json:"active,omitempty"`
	// Indicates whether the field names must be used as payload key instead of the id.
	UseUserFieldNames *bool `json:"use_user_field_names,omitempty"`
}

// NewPatchedTableWebhookUpdateRequest instantiates a new PatchedTableWebhookUpdateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedTableWebhookUpdateRequest() *PatchedTableWebhookUpdateRequest {
	this := PatchedTableWebhookUpdateRequest{}
	return &this
}

// NewPatchedTableWebhookUpdateRequestWithDefaults instantiates a new PatchedTableWebhookUpdateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedTableWebhookUpdateRequestWithDefaults() *PatchedTableWebhookUpdateRequest {
	this := PatchedTableWebhookUpdateRequest{}
	return &this
}

// GetUrl returns the Url field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetUrl() string {
	if o == nil || IsNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetUrlOk() (*string, bool) {
	if o == nil || IsNil(o.Url) {
		return nil, false
	}
	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasUrl() bool {
	if o != nil && !IsNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *PatchedTableWebhookUpdateRequest) SetUrl(v string) {
	o.Url = &v
}

// GetIncludeAllEvents returns the IncludeAllEvents field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetIncludeAllEvents() bool {
	if o == nil || IsNil(o.IncludeAllEvents) {
		var ret bool
		return ret
	}
	return *o.IncludeAllEvents
}

// GetIncludeAllEventsOk returns a tuple with the IncludeAllEvents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetIncludeAllEventsOk() (*bool, bool) {
	if o == nil || IsNil(o.IncludeAllEvents) {
		return nil, false
	}
	return o.IncludeAllEvents, true
}

// HasIncludeAllEvents returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasIncludeAllEvents() bool {
	if o != nil && !IsNil(o.IncludeAllEvents) {
		return true
	}

	return false
}

// SetIncludeAllEvents gets a reference to the given bool and assigns it to the IncludeAllEvents field.
func (o *PatchedTableWebhookUpdateRequest) SetIncludeAllEvents(v bool) {
	o.IncludeAllEvents = &v
}

// GetEvents returns the Events field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetEvents() []EventTypesEnum {
	if o == nil || IsNil(o.Events) {
		var ret []EventTypesEnum
		return ret
	}
	return o.Events
}

// GetEventsOk returns a tuple with the Events field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetEventsOk() ([]EventTypesEnum, bool) {
	if o == nil || IsNil(o.Events) {
		return nil, false
	}
	return o.Events, true
}

// HasEvents returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasEvents() bool {
	if o != nil && !IsNil(o.Events) {
		return true
	}

	return false
}

// SetEvents gets a reference to the given []EventTypesEnum and assigns it to the Events field.
func (o *PatchedTableWebhookUpdateRequest) SetEvents(v []EventTypesEnum) {
	o.Events = v
}

// GetRequestMethod returns the RequestMethod field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetRequestMethod() RequestMethodEnum {
	if o == nil || IsNil(o.RequestMethod) {
		var ret RequestMethodEnum
		return ret
	}
	return *o.RequestMethod
}

// GetRequestMethodOk returns a tuple with the RequestMethod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetRequestMethodOk() (*RequestMethodEnum, bool) {
	if o == nil || IsNil(o.RequestMethod) {
		return nil, false
	}
	return o.RequestMethod, true
}

// HasRequestMethod returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasRequestMethod() bool {
	if o != nil && !IsNil(o.RequestMethod) {
		return true
	}

	return false
}

// SetRequestMethod gets a reference to the given RequestMethodEnum and assigns it to the RequestMethod field.
func (o *PatchedTableWebhookUpdateRequest) SetRequestMethod(v RequestMethodEnum) {
	o.RequestMethod = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetHeaders() map[string]interface{} {
	if o == nil || IsNil(o.Headers) {
		var ret map[string]interface{}
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetHeadersOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Headers) {
		return map[string]interface{}{}, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasHeaders() bool {
	if o != nil && !IsNil(o.Headers) {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given map[string]interface{} and assigns it to the Headers field.
func (o *PatchedTableWebhookUpdateRequest) SetHeaders(v map[string]interface{}) {
	o.Headers = v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedTableWebhookUpdateRequest) SetName(v string) {
	o.Name = &v
}

// GetActive returns the Active field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetActive() bool {
	if o == nil || IsNil(o.Active) {
		var ret bool
		return ret
	}
	return *o.Active
}

// GetActiveOk returns a tuple with the Active field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.Active) {
		return nil, false
	}
	return o.Active, true
}

// HasActive returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasActive() bool {
	if o != nil && !IsNil(o.Active) {
		return true
	}

	return false
}

// SetActive gets a reference to the given bool and assigns it to the Active field.
func (o *PatchedTableWebhookUpdateRequest) SetActive(v bool) {
	o.Active = &v
}

// GetUseUserFieldNames returns the UseUserFieldNames field value if set, zero value otherwise.
func (o *PatchedTableWebhookUpdateRequest) GetUseUserFieldNames() bool {
	if o == nil || IsNil(o.UseUserFieldNames) {
		var ret bool
		return ret
	}
	return *o.UseUserFieldNames
}

// GetUseUserFieldNamesOk returns a tuple with the UseUserFieldNames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedTableWebhookUpdateRequest) GetUseUserFieldNamesOk() (*bool, bool) {
	if o == nil || IsNil(o.UseUserFieldNames) {
		return nil, false
	}
	return o.UseUserFieldNames, true
}

// HasUseUserFieldNames returns a boolean if a field has been set.
func (o *PatchedTableWebhookUpdateRequest) HasUseUserFieldNames() bool {
	if o != nil && !IsNil(o.UseUserFieldNames) {
		return true
	}

	return false
}

// SetUseUserFieldNames gets a reference to the given bool and assigns it to the UseUserFieldNames field.
func (o *PatchedTableWebhookUpdateRequest) SetUseUserFieldNames(v bool) {
	o.UseUserFieldNames = &v
}

func (o PatchedTableWebhookUpdateRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedTableWebhookUpdateRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Url) {
		toSerialize["url"] = o.Url
	}
	if !IsNil(o.IncludeAllEvents) {
		toSerialize["include_all_events"] = o.IncludeAllEvents
	}
	if !IsNil(o.Events) {
		toSerialize["events"] = o.Events
	}
	if !IsNil(o.RequestMethod) {
		toSerialize["request_method"] = o.RequestMethod
	}
	if !IsNil(o.Headers) {
		toSerialize["headers"] = o.Headers
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Active) {
		toSerialize["active"] = o.Active
	}
	if !IsNil(o.UseUserFieldNames) {
		toSerialize["use_user_field_names"] = o.UseUserFieldNames
	}
	return toSerialize, nil
}

type NullablePatchedTableWebhookUpdateRequest struct {
	value *PatchedTableWebhookUpdateRequest
	isSet bool
}

func (v NullablePatchedTableWebhookUpdateRequest) Get() *PatchedTableWebhookUpdateRequest {
	return v.value
}

func (v *NullablePatchedTableWebhookUpdateRequest) Set(val *PatchedTableWebhookUpdateRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedTableWebhookUpdateRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedTableWebhookUpdateRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedTableWebhookUpdateRequest(val *PatchedTableWebhookUpdateRequest) *NullablePatchedTableWebhookUpdateRequest {
	return &NullablePatchedTableWebhookUpdateRequest{value: val, isSet: true}
}

func (v NullablePatchedTableWebhookUpdateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedTableWebhookUpdateRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

