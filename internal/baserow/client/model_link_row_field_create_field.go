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

// checks if the LinkRowFieldCreateField type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LinkRowFieldCreateField{}

// LinkRowFieldCreateField struct for LinkRowFieldCreateField
type LinkRowFieldCreateField struct {
	Name string `json:"name"`
	Type Type712Enum `json:"type"`
	// The id of the linked table.
	LinkRowTableId NullableInt32 `json:"link_row_table_id,omitempty"`
	// The id of the related field.
	LinkRowRelatedFieldId NullableInt32 `json:"link_row_related_field_id"`
	// (Deprecated) The id of the linked table.
	LinkRowTable NullableInt32 `json:"link_row_table,omitempty"`
	// (Deprecated) The id of the related field.
	LinkRowRelatedField int32 `json:"link_row_related_field"`
}

// NewLinkRowFieldCreateField instantiates a new LinkRowFieldCreateField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinkRowFieldCreateField(name string, type_ Type712Enum, linkRowRelatedFieldId NullableInt32, linkRowRelatedField int32) *LinkRowFieldCreateField {
	this := LinkRowFieldCreateField{}
	this.Name = name
	this.Type = type_
	this.LinkRowRelatedFieldId = linkRowRelatedFieldId
	this.LinkRowRelatedField = linkRowRelatedField
	return &this
}

// NewLinkRowFieldCreateFieldWithDefaults instantiates a new LinkRowFieldCreateField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinkRowFieldCreateFieldWithDefaults() *LinkRowFieldCreateField {
	this := LinkRowFieldCreateField{}
	return &this
}

// GetName returns the Name field value
func (o *LinkRowFieldCreateField) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *LinkRowFieldCreateField) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *LinkRowFieldCreateField) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *LinkRowFieldCreateField) GetType() Type712Enum {
	if o == nil {
		var ret Type712Enum
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *LinkRowFieldCreateField) GetTypeOk() (*Type712Enum, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *LinkRowFieldCreateField) SetType(v Type712Enum) {
	o.Type = v
}

// GetLinkRowTableId returns the LinkRowTableId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LinkRowFieldCreateField) GetLinkRowTableId() int32 {
	if o == nil || IsNil(o.LinkRowTableId.Get()) {
		var ret int32
		return ret
	}
	return *o.LinkRowTableId.Get()
}

// GetLinkRowTableIdOk returns a tuple with the LinkRowTableId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LinkRowFieldCreateField) GetLinkRowTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.LinkRowTableId.Get(), o.LinkRowTableId.IsSet()
}

// HasLinkRowTableId returns a boolean if a field has been set.
func (o *LinkRowFieldCreateField) HasLinkRowTableId() bool {
	if o != nil && o.LinkRowTableId.IsSet() {
		return true
	}

	return false
}

// SetLinkRowTableId gets a reference to the given NullableInt32 and assigns it to the LinkRowTableId field.
func (o *LinkRowFieldCreateField) SetLinkRowTableId(v int32) {
	o.LinkRowTableId.Set(&v)
}
// SetLinkRowTableIdNil sets the value for LinkRowTableId to be an explicit nil
func (o *LinkRowFieldCreateField) SetLinkRowTableIdNil() {
	o.LinkRowTableId.Set(nil)
}

// UnsetLinkRowTableId ensures that no value is present for LinkRowTableId, not even an explicit nil
func (o *LinkRowFieldCreateField) UnsetLinkRowTableId() {
	o.LinkRowTableId.Unset()
}

// GetLinkRowRelatedFieldId returns the LinkRowRelatedFieldId field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *LinkRowFieldCreateField) GetLinkRowRelatedFieldId() int32 {
	if o == nil || o.LinkRowRelatedFieldId.Get() == nil {
		var ret int32
		return ret
	}

	return *o.LinkRowRelatedFieldId.Get()
}

// GetLinkRowRelatedFieldIdOk returns a tuple with the LinkRowRelatedFieldId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LinkRowFieldCreateField) GetLinkRowRelatedFieldIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.LinkRowRelatedFieldId.Get(), o.LinkRowRelatedFieldId.IsSet()
}

// SetLinkRowRelatedFieldId sets field value
func (o *LinkRowFieldCreateField) SetLinkRowRelatedFieldId(v int32) {
	o.LinkRowRelatedFieldId.Set(&v)
}

// GetLinkRowTable returns the LinkRowTable field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *LinkRowFieldCreateField) GetLinkRowTable() int32 {
	if o == nil || IsNil(o.LinkRowTable.Get()) {
		var ret int32
		return ret
	}
	return *o.LinkRowTable.Get()
}

// GetLinkRowTableOk returns a tuple with the LinkRowTable field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LinkRowFieldCreateField) GetLinkRowTableOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.LinkRowTable.Get(), o.LinkRowTable.IsSet()
}

// HasLinkRowTable returns a boolean if a field has been set.
func (o *LinkRowFieldCreateField) HasLinkRowTable() bool {
	if o != nil && o.LinkRowTable.IsSet() {
		return true
	}

	return false
}

// SetLinkRowTable gets a reference to the given NullableInt32 and assigns it to the LinkRowTable field.
func (o *LinkRowFieldCreateField) SetLinkRowTable(v int32) {
	o.LinkRowTable.Set(&v)
}
// SetLinkRowTableNil sets the value for LinkRowTable to be an explicit nil
func (o *LinkRowFieldCreateField) SetLinkRowTableNil() {
	o.LinkRowTable.Set(nil)
}

// UnsetLinkRowTable ensures that no value is present for LinkRowTable, not even an explicit nil
func (o *LinkRowFieldCreateField) UnsetLinkRowTable() {
	o.LinkRowTable.Unset()
}

// GetLinkRowRelatedField returns the LinkRowRelatedField field value
func (o *LinkRowFieldCreateField) GetLinkRowRelatedField() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.LinkRowRelatedField
}

// GetLinkRowRelatedFieldOk returns a tuple with the LinkRowRelatedField field value
// and a boolean to check if the value has been set.
func (o *LinkRowFieldCreateField) GetLinkRowRelatedFieldOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LinkRowRelatedField, true
}

// SetLinkRowRelatedField sets field value
func (o *LinkRowFieldCreateField) SetLinkRowRelatedField(v int32) {
	o.LinkRowRelatedField = v
}

func (o LinkRowFieldCreateField) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LinkRowFieldCreateField) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	if o.LinkRowTableId.IsSet() {
		toSerialize["link_row_table_id"] = o.LinkRowTableId.Get()
	}
	toSerialize["link_row_related_field_id"] = o.LinkRowRelatedFieldId.Get()
	if o.LinkRowTable.IsSet() {
		toSerialize["link_row_table"] = o.LinkRowTable.Get()
	}
	// skip: link_row_related_field is readOnly
	return toSerialize, nil
}

type NullableLinkRowFieldCreateField struct {
	value *LinkRowFieldCreateField
	isSet bool
}

func (v NullableLinkRowFieldCreateField) Get() *LinkRowFieldCreateField {
	return v.value
}

func (v *NullableLinkRowFieldCreateField) Set(val *LinkRowFieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableLinkRowFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableLinkRowFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinkRowFieldCreateField(val *LinkRowFieldCreateField) *NullableLinkRowFieldCreateField {
	return &NullableLinkRowFieldCreateField{value: val, isSet: true}
}

func (v NullableLinkRowFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinkRowFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


