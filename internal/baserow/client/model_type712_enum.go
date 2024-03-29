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

// Type712Enum * `text` - text * `long_text` - long_text * `url` - url * `email` - email * `number` - number * `rating` - rating * `boolean` - boolean * `date` - date * `last_modified` - last_modified * `created_on` - created_on * `link_row` - link_row * `file` - file * `single_select` - single_select * `multiple_select` - multiple_select * `phone_number` - phone_number * `formula` - formula * `count` - count * `rollup` - rollup * `lookup` - lookup * `multiple_collaborators` - multiple_collaborators
type Type712Enum string

// List of Type712Enum
const (
	TEXT Type712Enum = "text"
	LONG_TEXT Type712Enum = "long_text"
	URL Type712Enum = "url"
	EMAIL Type712Enum = "email"
	NUMBER Type712Enum = "number"
	RATING Type712Enum = "rating"
	BOOLEAN Type712Enum = "boolean"
	DATE Type712Enum = "date"
	LAST_MODIFIED Type712Enum = "last_modified"
	CREATED_ON Type712Enum = "created_on"
	LINK_ROW Type712Enum = "link_row"
	FILE Type712Enum = "file"
	SINGLE_SELECT Type712Enum = "single_select"
	MULTIPLE_SELECT Type712Enum = "multiple_select"
	PHONE_NUMBER Type712Enum = "phone_number"
	FORMULA Type712Enum = "formula"
	COUNT Type712Enum = "count"
	ROLLUP Type712Enum = "rollup"
	LOOKUP Type712Enum = "lookup"
	MULTIPLE_COLLABORATORS Type712Enum = "multiple_collaborators"
)

// All allowed values of Type712Enum enum
var AllowedType712EnumEnumValues = []Type712Enum{
	"text",
	"long_text",
	"url",
	"email",
	"number",
	"rating",
	"boolean",
	"date",
	"last_modified",
	"created_on",
	"link_row",
	"file",
	"single_select",
	"multiple_select",
	"phone_number",
	"formula",
	"count",
	"rollup",
	"lookup",
	"multiple_collaborators",
}

func (v *Type712Enum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Type712Enum(value)
	for _, existing := range AllowedType712EnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Type712Enum", value)
}

// NewType712EnumFromValue returns a pointer to a valid Type712Enum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewType712EnumFromValue(v string) (*Type712Enum, error) {
	ev := Type712Enum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Type712Enum: valid values are %v", v, AllowedType712EnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Type712Enum) IsValid() bool {
	for _, existing := range AllowedType712EnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Type712Enum value
func (v Type712Enum) Ptr() *Type712Enum {
	return &v
}

type NullableType712Enum struct {
	value *Type712Enum
	isSet bool
}

func (v NullableType712Enum) Get() *Type712Enum {
	return v.value
}

func (v *NullableType712Enum) Set(val *Type712Enum) {
	v.value = val
	v.isSet = true
}

func (v NullableType712Enum) IsSet() bool {
	return v.isSet
}

func (v *NullableType712Enum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableType712Enum(val *Type712Enum) *NullableType712Enum {
	return &NullableType712Enum{value: val, isSet: true}
}

func (v NullableType712Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableType712Enum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

