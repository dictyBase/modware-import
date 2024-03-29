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

// RequestMethodEnum * `POST` - Post * `GET` - Get * `PUT` - Put * `PATCH` - Patch * `DELETE` - Delete
type RequestMethodEnum string

// List of RequestMethodEnum
const (
	POST RequestMethodEnum = "POST"
	GET RequestMethodEnum = "GET"
	PUT RequestMethodEnum = "PUT"
	PATCH RequestMethodEnum = "PATCH"
	DELETE RequestMethodEnum = "DELETE"
)

// All allowed values of RequestMethodEnum enum
var AllowedRequestMethodEnumEnumValues = []RequestMethodEnum{
	"POST",
	"GET",
	"PUT",
	"PATCH",
	"DELETE",
}

func (v *RequestMethodEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := RequestMethodEnum(value)
	for _, existing := range AllowedRequestMethodEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid RequestMethodEnum", value)
}

// NewRequestMethodEnumFromValue returns a pointer to a valid RequestMethodEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewRequestMethodEnumFromValue(v string) (*RequestMethodEnum, error) {
	ev := RequestMethodEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for RequestMethodEnum: valid values are %v", v, AllowedRequestMethodEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v RequestMethodEnum) IsValid() bool {
	for _, existing := range AllowedRequestMethodEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to RequestMethodEnum value
func (v RequestMethodEnum) Ptr() *RequestMethodEnum {
	return &v
}

type NullableRequestMethodEnum struct {
	value *RequestMethodEnum
	isSet bool
}

func (v NullableRequestMethodEnum) Get() *RequestMethodEnum {
	return v.value
}

func (v *NullableRequestMethodEnum) Set(val *RequestMethodEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestMethodEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestMethodEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestMethodEnum(val *RequestMethodEnum) *NullableRequestMethodEnum {
	return &NullableRequestMethodEnum{value: val, isSet: true}
}

func (v NullableRequestMethodEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestMethodEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

