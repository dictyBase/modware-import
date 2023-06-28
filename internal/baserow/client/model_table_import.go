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

// checks if the TableImport type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TableImport{}

// TableImport struct for TableImport
type TableImport struct {
	// A list of rows you want to add to the specified table. Each row is a list of values, one for each **writable** field. The field values must be ordered according to the field order in the table. All values must be compatible with the corresponding field type.  Ex:  ```json [   [\"row1_field1_value\", \"row1_field2_value\"],   [\"row2_field1_value\", \"row2_field2_value\"], ] ``` for adding two rows to a table with two writable fields.
	Data []interface{} `json:"data"`
}

// NewTableImport instantiates a new TableImport object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTableImport(data []interface{}) *TableImport {
	this := TableImport{}
	this.Data = data
	return &this
}

// NewTableImportWithDefaults instantiates a new TableImport object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTableImportWithDefaults() *TableImport {
	this := TableImport{}
	return &this
}

// GetData returns the Data field value
func (o *TableImport) GetData() []interface{} {
	if o == nil {
		var ret []interface{}
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *TableImport) GetDataOk() ([]interface{}, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *TableImport) SetData(v []interface{}) {
	o.Data = v
}

func (o TableImport) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TableImport) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableTableImport struct {
	value *TableImport
	isSet bool
}

func (v NullableTableImport) Get() *TableImport {
	return v.value
}

func (v *NullableTableImport) Set(val *TableImport) {
	v.value = val
	v.isSet = true
}

func (v NullableTableImport) IsSet() bool {
	return v.isSet
}

func (v *NullableTableImport) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTableImport(val *TableImport) *NullableTableImport {
	return &NullableTableImport{value: val, isSet: true}
}

func (v NullableTableImport) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTableImport) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


