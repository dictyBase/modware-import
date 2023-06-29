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

// checks if the BooleanFieldFieldSerializerWithRelatedFields type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BooleanFieldFieldSerializerWithRelatedFields{}

// BooleanFieldFieldSerializerWithRelatedFields struct for BooleanFieldFieldSerializerWithRelatedFields
type BooleanFieldFieldSerializerWithRelatedFields struct {
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

// NewBooleanFieldFieldSerializerWithRelatedFields instantiates a new BooleanFieldFieldSerializerWithRelatedFields object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBooleanFieldFieldSerializerWithRelatedFields(id int32, tableId int32, name string, order int32, type_ string, readOnly bool, relatedFields []Field) *BooleanFieldFieldSerializerWithRelatedFields {
	this := BooleanFieldFieldSerializerWithRelatedFields{}
	this.Id = id
	this.TableId = tableId
	this.Name = name
	this.Order = order
	this.Type = type_
	this.ReadOnly = readOnly
	this.RelatedFields = relatedFields
	return &this
}

// NewBooleanFieldFieldSerializerWithRelatedFieldsWithDefaults instantiates a new BooleanFieldFieldSerializerWithRelatedFields object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBooleanFieldFieldSerializerWithRelatedFieldsWithDefaults() *BooleanFieldFieldSerializerWithRelatedFields {
	this := BooleanFieldFieldSerializerWithRelatedFields{}
	return &this
}

// GetId returns the Id field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetTableId(v int32) {
	o.TableId = v
}

// GetName returns the Name field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetName(v string) {
	o.Name = v
}

// GetOrder returns the Order field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetOrder() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetOrderOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetOrder(v int32) {
	o.Order = v
}

// GetType returns the Type field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetType(v string) {
	o.Type = v
}

// GetPrimary returns the Primary field value if set, zero value otherwise.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetPrimary() bool {
	if o == nil || IsNil(o.Primary) {
		var ret bool
		return ret
	}
	return *o.Primary
}

// GetPrimaryOk returns a tuple with the Primary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetPrimaryOk() (*bool, bool) {
	if o == nil || IsNil(o.Primary) {
		return nil, false
	}
	return o.Primary, true
}

// HasPrimary returns a boolean if a field has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) HasPrimary() bool {
	if o != nil && !IsNil(o.Primary) {
		return true
	}

	return false
}

// SetPrimary gets a reference to the given bool and assigns it to the Primary field.
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetPrimary(v bool) {
	o.Primary = &v
}

// GetReadOnly returns the ReadOnly field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetReadOnly() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ReadOnly
}

// GetReadOnlyOk returns a tuple with the ReadOnly field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetReadOnlyOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReadOnly, true
}

// SetReadOnly sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetReadOnly(v bool) {
	o.ReadOnly = v
}

// GetRelatedFields returns the RelatedFields field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetRelatedFields() []Field {
	if o == nil {
		var ret []Field
		return ret
	}

	return o.RelatedFields
}

// GetRelatedFieldsOk returns a tuple with the RelatedFields field value
// and a boolean to check if the value has been set.
func (o *BooleanFieldFieldSerializerWithRelatedFields) GetRelatedFieldsOk() ([]Field, bool) {
	if o == nil {
		return nil, false
	}
	return o.RelatedFields, true
}

// SetRelatedFields sets field value
func (o *BooleanFieldFieldSerializerWithRelatedFields) SetRelatedFields(v []Field) {
	o.RelatedFields = v
}

func (o BooleanFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BooleanFieldFieldSerializerWithRelatedFields) ToMap() (map[string]interface{}, error) {
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

type NullableBooleanFieldFieldSerializerWithRelatedFields struct {
	value *BooleanFieldFieldSerializerWithRelatedFields
	isSet bool
}

func (v NullableBooleanFieldFieldSerializerWithRelatedFields) Get() *BooleanFieldFieldSerializerWithRelatedFields {
	return v.value
}

func (v *NullableBooleanFieldFieldSerializerWithRelatedFields) Set(val *BooleanFieldFieldSerializerWithRelatedFields) {
	v.value = val
	v.isSet = true
}

func (v NullableBooleanFieldFieldSerializerWithRelatedFields) IsSet() bool {
	return v.isSet
}

func (v *NullableBooleanFieldFieldSerializerWithRelatedFields) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBooleanFieldFieldSerializerWithRelatedFields(val *BooleanFieldFieldSerializerWithRelatedFields) *NullableBooleanFieldFieldSerializerWithRelatedFields {
	return &NullableBooleanFieldFieldSerializerWithRelatedFields{value: val, isSet: true}
}

func (v NullableBooleanFieldFieldSerializerWithRelatedFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBooleanFieldFieldSerializerWithRelatedFields) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

