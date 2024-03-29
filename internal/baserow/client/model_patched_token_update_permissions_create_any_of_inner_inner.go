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

// PatchedTokenUpdatePermissionsCreateAnyOfInnerInner struct for PatchedTokenUpdatePermissionsCreateAnyOfInnerInner
type PatchedTokenUpdatePermissionsCreateAnyOfInnerInner struct {
	float32 *float32
	string *string
}

// Unmarshal JSON data into any of the pointers in the struct
func (dst *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner) UnmarshalJSON(data []byte) error {
	var err error
	// try to unmarshal JSON data into float32
	err = json.Unmarshal(data, &dst.float32);
	if err == nil {
		jsonfloat32, _ := json.Marshal(dst.float32)
		if string(jsonfloat32) == "{}" { // empty struct
			dst.float32 = nil
		} else {
			return nil // data stored in dst.float32, return on the first match
		}
	} else {
		dst.float32 = nil
	}

	// try to unmarshal JSON data into string
	err = json.Unmarshal(data, &dst.string);
	if err == nil {
		jsonstring, _ := json.Marshal(dst.string)
		if string(jsonstring) == "{}" { // empty struct
			dst.string = nil
		} else {
			return nil // data stored in dst.string, return on the first match
		}
	} else {
		dst.string = nil
	}

	return fmt.Errorf("data failed to match schemas in anyOf(PatchedTokenUpdatePermissionsCreateAnyOfInnerInner)")
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner) MarshalJSON() ([]byte, error) {
	if src.float32 != nil {
		return json.Marshal(&src.float32)
	}

	if src.string != nil {
		return json.Marshal(&src.string)
	}

	return nil, nil // no data in anyOf schemas
}

type NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner struct {
	value *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner
	isSet bool
}

func (v NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) Get() *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner {
	return v.value
}

func (v *NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) Set(val *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner(val *PatchedTokenUpdatePermissionsCreateAnyOfInnerInner) *NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner {
	return &NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner{value: val, isSet: true}
}

func (v NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedTokenUpdatePermissionsCreateAnyOfInnerInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


