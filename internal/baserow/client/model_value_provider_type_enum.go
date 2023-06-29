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

// ValueProviderTypeEnum * `` -  * `single_select_color` - single_select_color * `conditional_color` - conditional_color
type ValueProviderTypeEnum string

// List of ValueProviderTypeEnum
const (
	SINGLE_SELECT_COLOR ValueProviderTypeEnum = "single_select_color"
	CONDITIONAL_COLOR ValueProviderTypeEnum = "conditional_color"
)

// All allowed values of ValueProviderTypeEnum enum
var AllowedValueProviderTypeEnumEnumValues = []ValueProviderTypeEnum{
	"single_select_color",
	"conditional_color",
}

func (v *ValueProviderTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ValueProviderTypeEnum(value)
	for _, existing := range AllowedValueProviderTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ValueProviderTypeEnum", value)
}

// NewValueProviderTypeEnumFromValue returns a pointer to a valid ValueProviderTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewValueProviderTypeEnumFromValue(v string) (*ValueProviderTypeEnum, error) {
	ev := ValueProviderTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ValueProviderTypeEnum: valid values are %v", v, AllowedValueProviderTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ValueProviderTypeEnum) IsValid() bool {
	for _, existing := range AllowedValueProviderTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ValueProviderTypeEnum value
func (v ValueProviderTypeEnum) Ptr() *ValueProviderTypeEnum {
	return &v
}

type NullableValueProviderTypeEnum struct {
	value *ValueProviderTypeEnum
	isSet bool
}

func (v NullableValueProviderTypeEnum) Get() *ValueProviderTypeEnum {
	return v.value
}

func (v *NullableValueProviderTypeEnum) Set(val *ValueProviderTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableValueProviderTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableValueProviderTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableValueProviderTypeEnum(val *ValueProviderTypeEnum) *NullableValueProviderTypeEnum {
	return &NullableValueProviderTypeEnum{value: val, isSet: true}
}

func (v NullableValueProviderTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableValueProviderTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
