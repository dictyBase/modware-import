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

// ViewTypesEnum * `grid` - grid * `gallery` - gallery * `form` - form * `kanban` - kanban * `calendar` - calendar
type ViewTypesEnum string

// List of ViewTypesEnum
const (
	GRID ViewTypesEnum = "grid"
	GALLERY ViewTypesEnum = "gallery"
	KANBAN ViewTypesEnum = "kanban"
	CALENDAR ViewTypesEnum = "calendar"
)

// All allowed values of ViewTypesEnum enum
var AllowedViewTypesEnumEnumValues = []ViewTypesEnum{
	"grid",
	"gallery",
	"form",
	"kanban",
	"calendar",
}

func (v *ViewTypesEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ViewTypesEnum(value)
	for _, existing := range AllowedViewTypesEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ViewTypesEnum", value)
}

// NewViewTypesEnumFromValue returns a pointer to a valid ViewTypesEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewViewTypesEnumFromValue(v string) (*ViewTypesEnum, error) {
	ev := ViewTypesEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ViewTypesEnum: valid values are %v", v, AllowedViewTypesEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ViewTypesEnum) IsValid() bool {
	for _, existing := range AllowedViewTypesEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ViewTypesEnum value
func (v ViewTypesEnum) Ptr() *ViewTypesEnum {
	return &v
}

type NullableViewTypesEnum struct {
	value *ViewTypesEnum
	isSet bool
}

func (v NullableViewTypesEnum) Get() *ViewTypesEnum {
	return v.value
}

func (v *NullableViewTypesEnum) Set(val *ViewTypesEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableViewTypesEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableViewTypesEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableViewTypesEnum(val *ViewTypesEnum) *NullableViewTypesEnum {
	return &NullableViewTypesEnum{value: val, isSet: true}
}

func (v NullableViewTypesEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableViewTypesEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

