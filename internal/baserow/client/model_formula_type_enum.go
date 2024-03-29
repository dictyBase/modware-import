/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"fmt"
)

// FormulaTypeEnum * `invalid` - invalid * `text` - text * `char` - char * `link` - link * `date_interval` - date_interval * `date` - date * `boolean` - boolean * `number` - number * `array` - array * `single_select` - single_select
type FormulaTypeEnum string

// List of FormulaTypeEnum
const (
	INVALID FormulaTypeEnum = "invalid"
	CHAR FormulaTypeEnum = "char"
	LINK FormulaTypeEnum = "link"
	DATE_INTERVAL FormulaTypeEnum = "date_interval"
	ARRAY FormulaTypeEnum = "array"
)

// All allowed values of FormulaTypeEnum enum
var AllowedFormulaTypeEnumEnumValues = []FormulaTypeEnum{
	"invalid",
	"text",
	"char",
	"link",
	"date_interval",
	"date",
	"boolean",
	"number",
	"array",
	"single_select",
}

func (v *FormulaTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := FormulaTypeEnum(value)
	for _, existing := range AllowedFormulaTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid FormulaTypeEnum", value)
}

// NewFormulaTypeEnumFromValue returns a pointer to a valid FormulaTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewFormulaTypeEnumFromValue(v string) (*FormulaTypeEnum, error) {
	ev := FormulaTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for FormulaTypeEnum: valid values are %v", v, AllowedFormulaTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v FormulaTypeEnum) IsValid() bool {
	for _, existing := range AllowedFormulaTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to FormulaTypeEnum value
func (v FormulaTypeEnum) Ptr() *FormulaTypeEnum {
	return &v
}

type NullableFormulaTypeEnum struct {
	value *FormulaTypeEnum
	isSet bool
}

func (v NullableFormulaTypeEnum) Get() *FormulaTypeEnum {
	return v.value
}

func (v *NullableFormulaTypeEnum) Set(val *FormulaTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableFormulaTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableFormulaTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFormulaTypeEnum(val *FormulaTypeEnum) *NullableFormulaTypeEnum {
	return &NullableFormulaTypeEnum{value: val, isSet: true}
}

func (v NullableFormulaTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFormulaTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

