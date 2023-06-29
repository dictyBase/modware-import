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

// checks if the URLFieldFieldSerializerWithRelatedFields type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &URLFieldFieldSerializerWithRelatedFields{}

// URLFieldFieldSerializerWithRelatedFields struct for URLFieldFieldSerializerWithRelatedFields
type URLFieldFieldSerializerWithRelatedFields struct {
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
}

// NewURLFieldFieldSerializerWithRelatedFields instantiates a new URLFieldFieldSerializerWithRelatedFields object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewURLFieldFieldSerializerWithRelatedFields(id int32, tableId int32, name string, order int32, type_ string, readOnly bool, relatedFields []Field) *URLFieldFieldSerializerWithRelatedFields {
	this := URLFieldFieldSerializerWithRelatedFields{}
	this.Id = id
	this.TableId = tableId
	this.Name = name
	this.Order = order
	this.Type = type_
	this.ReadOnly = readOnly
	this.RelatedFields = relatedFields
	return &this
}

// NewURLFieldFieldSerializerWithRelatedFieldsWithDefaults instantiates a new URLFieldFieldSerializerWithRelatedFields object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewURLFieldFieldSerializerWithRelatedFieldsWithDefaults() *URLFieldFieldSerializerWithRelatedFields {
	this := URLFieldFieldSerializerWithRelatedFields{}
	return &this
}

// GetId returns the Id field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetType(v string) {
	o.Type = v
}

// GetPrimary returns the Primary field value if set, zero value otherwise.
func (o *URLFieldFieldSerializerWithRelatedFields) GetPrimary() bool {
	if o == nil || IsNil(o.Primary) {
		var ret bool
		return ret
	}
	return *o.Primary
}

// GetPrimaryOk returns a tuple with the Primary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetPrimaryOk() (*bool, bool) {
	if o == nil || IsNil(o.Primary) {
		return nil, false
	}
	return o.Primary, true
}

// HasPrimary returns a boolean if a field has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) HasPrimary() bool {
	if o != nil && !IsNil(o.Primary) {
		return true
	}

	return false
}

// SetPrimary gets a reference to the given bool and assigns it to the Primary field.
func (o *URLFieldFieldSerializerWithRelatedFields) SetPrimary(v bool) {
	o.Primary = &v
}

// GetReadOnly returns the ReadOnly field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetReadOnly() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ReadOnly
}

// GetReadOnlyOk returns a tuple with the ReadOnly field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetReadOnlyOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReadOnly, true
}

// SetReadOnly sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetReadOnly(v bool) {
	o.ReadOnly = v
}

// GetRelatedFields returns the RelatedFields field value
func (o *URLFieldFieldSerializerWithRelatedFields) GetRelatedFields() []Field {
	if o == nil {
		var ret []Field
		return ret
	}

	return o.RelatedFields
}

// GetRelatedFieldsOk returns a tuple with the RelatedFields field value
// and a boolean to check if the value has been set.
func (o *URLFieldFieldSerializerWithRelatedFields) GetRelatedFieldsOk() ([]Field, bool) {
	if o == nil {
		return nil, false
	}
	return o.RelatedFields, true
}

// SetRelatedFields sets field value
func (o *URLFieldFieldSerializerWithRelatedFields) SetRelatedFields(v []Field) {
	o.RelatedFields = v
}

func (o URLFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o URLFieldFieldSerializerWithRelatedFields) ToMap() (map[string]interface{}, error) {
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
	return toSerialize, nil
}

type NullableURLFieldFieldSerializerWithRelatedFields struct {
	value *URLFieldFieldSerializerWithRelatedFields
	isSet bool
}

func (v NullableURLFieldFieldSerializerWithRelatedFields) Get() *URLFieldFieldSerializerWithRelatedFields {
	return v.value
}

func (v *NullableURLFieldFieldSerializerWithRelatedFields) Set(val *URLFieldFieldSerializerWithRelatedFields) {
	v.value = val
	v.isSet = true
}

func (v NullableURLFieldFieldSerializerWithRelatedFields) IsSet() bool {
	return v.isSet
}

func (v *NullableURLFieldFieldSerializerWithRelatedFields) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableURLFieldFieldSerializerWithRelatedFields(val *URLFieldFieldSerializerWithRelatedFields) *NullableURLFieldFieldSerializerWithRelatedFields {
	return &NullableURLFieldFieldSerializerWithRelatedFields{value: val, isSet: true}
}

func (v NullableURLFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableURLFieldFieldSerializerWithRelatedFields) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

