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

// Events3eaEnum * `rows.created` - rows.created * `row.created` - row.created * `rows.updated` - rows.updated * `row.updated` - row.updated * `rows.deleted` - rows.deleted * `row.deleted` - row.deleted
type Events3eaEnum string

// List of Events3eaEnum
const (
	ROWS_CREATED Events3eaEnum = "rows.created"
	ROW_CREATED Events3eaEnum = "row.created"
	ROWS_UPDATED Events3eaEnum = "rows.updated"
	ROW_UPDATED Events3eaEnum = "row.updated"
	ROWS_DELETED Events3eaEnum = "rows.deleted"
	ROW_DELETED Events3eaEnum = "row.deleted"
)

// All allowed values of Events3eaEnum enum
var AllowedEvents3eaEnumEnumValues = []Events3eaEnum{
	"rows.created",
	"row.created",
	"rows.updated",
	"row.updated",
	"rows.deleted",
	"row.deleted",
}

func (v *Events3eaEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Events3eaEnum(value)
	for _, existing := range AllowedEvents3eaEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Events3eaEnum", value)
}

// NewEvents3eaEnumFromValue returns a pointer to a valid Events3eaEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewEvents3eaEnumFromValue(v string) (*Events3eaEnum, error) {
	ev := Events3eaEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Events3eaEnum: valid values are %v", v, AllowedEvents3eaEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Events3eaEnum) IsValid() bool {
	for _, existing := range AllowedEvents3eaEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Events3eaEnum value
func (v Events3eaEnum) Ptr() *Events3eaEnum {
	return &v
}

type NullableEvents3eaEnum struct {
	value *Events3eaEnum
	isSet bool
}

func (v NullableEvents3eaEnum) Get() *Events3eaEnum {
	return v.value
}

func (v *NullableEvents3eaEnum) Set(val *Events3eaEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableEvents3eaEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableEvents3eaEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEvents3eaEnum(val *Events3eaEnum) *NullableEvents3eaEnum {
	return &NullableEvents3eaEnum{value: val, isSet: true}
}

func (v NullableEvents3eaEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEvents3eaEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

