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

// ApplicationCreateTypeEnum * `database` - database
type ApplicationCreateTypeEnum string

// List of ApplicationCreateTypeEnum
const (
	DATABASE ApplicationCreateTypeEnum = "database"
)

// All allowed values of ApplicationCreateTypeEnum enum
var AllowedApplicationCreateTypeEnumEnumValues = []ApplicationCreateTypeEnum{
	"database",
}

func (v *ApplicationCreateTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ApplicationCreateTypeEnum(value)
	for _, existing := range AllowedApplicationCreateTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ApplicationCreateTypeEnum", value)
}

// NewApplicationCreateTypeEnumFromValue returns a pointer to a valid ApplicationCreateTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewApplicationCreateTypeEnumFromValue(v string) (*ApplicationCreateTypeEnum, error) {
	ev := ApplicationCreateTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ApplicationCreateTypeEnum: valid values are %v", v, AllowedApplicationCreateTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ApplicationCreateTypeEnum) IsValid() bool {
	for _, existing := range AllowedApplicationCreateTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ApplicationCreateTypeEnum value
func (v ApplicationCreateTypeEnum) Ptr() *ApplicationCreateTypeEnum {
	return &v
}

type NullableApplicationCreateTypeEnum struct {
	value *ApplicationCreateTypeEnum
	isSet bool
}

func (v NullableApplicationCreateTypeEnum) Get() *ApplicationCreateTypeEnum {
	return v.value
}

func (v *NullableApplicationCreateTypeEnum) Set(val *ApplicationCreateTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationCreateTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationCreateTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationCreateTypeEnum(val *ApplicationCreateTypeEnum) *NullableApplicationCreateTypeEnum {
	return &NullableApplicationCreateTypeEnum{value: val, isSet: true}
}

func (v NullableApplicationCreateTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationCreateTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
