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

// checks if the TableCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TableCreate{}

// TableCreate struct for TableCreate
type TableCreate struct {
	Name string `json:"name"`
	// A list of rows that needs to be created as initial table data. Each row is a list of values that are going to be added in the new table in the same order as provided.  Ex:  ```json [   [\"row1_field1_value\", \"row1_field2_value\"],   [\"row2_field1_value\", \"row2_field2_value\"], ] ``` for creating a two rows table with two fields.  If not provided, some example data is going to be created.
	Data []interface{} `json:"data,omitempty"`
	// Indicates if the first provided row is the header. If true the field names are going to be the values of the first row. Otherwise they will be called \"Field N\"
	FirstRowHeader *bool `json:"first_row_header,omitempty"`
}

// NewTableCreate instantiates a new TableCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTableCreate(name string) *TableCreate {
	this := TableCreate{}
	this.Name = name
	var firstRowHeader bool = false
	this.FirstRowHeader = &firstRowHeader
	return &this
}

// NewTableCreateWithDefaults instantiates a new TableCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTableCreateWithDefaults() *TableCreate {
	this := TableCreate{}
	var firstRowHeader bool = false
	this.FirstRowHeader = &firstRowHeader
	return &this
}

// GetName returns the Name field value
func (o *TableCreate) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *TableCreate) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *TableCreate) SetName(v string) {
	o.Name = v
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *TableCreate) GetData() []interface{} {
	if o == nil || IsNil(o.Data) {
		var ret []interface{}
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TableCreate) GetDataOk() ([]interface{}, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *TableCreate) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given []interface{} and assigns it to the Data field.
func (o *TableCreate) SetData(v []interface{}) {
	o.Data = v
}

// GetFirstRowHeader returns the FirstRowHeader field value if set, zero value otherwise.
func (o *TableCreate) GetFirstRowHeader() bool {
	if o == nil || IsNil(o.FirstRowHeader) {
		var ret bool
		return ret
	}
	return *o.FirstRowHeader
}

// GetFirstRowHeaderOk returns a tuple with the FirstRowHeader field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TableCreate) GetFirstRowHeaderOk() (*bool, bool) {
	if o == nil || IsNil(o.FirstRowHeader) {
		return nil, false
	}
	return o.FirstRowHeader, true
}

// HasFirstRowHeader returns a boolean if a field has been set.
func (o *TableCreate) HasFirstRowHeader() bool {
	if o != nil && !IsNil(o.FirstRowHeader) {
		return true
	}

	return false
}

// SetFirstRowHeader gets a reference to the given bool and assigns it to the FirstRowHeader field.
func (o *TableCreate) SetFirstRowHeader(v bool) {
	o.FirstRowHeader = &v
}

func (o TableCreate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TableCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	if !IsNil(o.FirstRowHeader) {
		toSerialize["first_row_header"] = o.FirstRowHeader
	}
	return toSerialize, nil
}

type NullableTableCreate struct {
	value *TableCreate
	isSet bool
}

func (v NullableTableCreate) Get() *TableCreate {
	return v.value
}

func (v *NullableTableCreate) Set(val *TableCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableTableCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableTableCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTableCreate(val *TableCreate) *NullableTableCreate {
	return &NullableTableCreate{value: val, isSet: true}
}

func (v NullableTableCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTableCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


