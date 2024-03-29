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

// CountFieldCreateFieldNumberDecimalPlaces - The amount of digits allowed after the point.  * `0` - 1 * `1` - 1.0 * `2` - 1.00 * `3` - 1.000 * `4` - 1.0000 * `5` - 1.00000 * `6` - 1.000000 * `7` - 1.0000000 * `8` - 1.00000000 * `9` - 1.000000000 * `10` - 1.0000000000
type CountFieldCreateFieldNumberDecimalPlaces struct {
	NullEnum *NullEnum
	NumberDecimalPlacesEnum *NumberDecimalPlacesEnum
}

// NullEnumAsCountFieldCreateFieldNumberDecimalPlaces is a convenience function that returns NullEnum wrapped in CountFieldCreateFieldNumberDecimalPlaces
func NullEnumAsCountFieldCreateFieldNumberDecimalPlaces(v *NullEnum) CountFieldCreateFieldNumberDecimalPlaces {
	return CountFieldCreateFieldNumberDecimalPlaces{
		NullEnum: v,
	}
}

// NumberDecimalPlacesEnumAsCountFieldCreateFieldNumberDecimalPlaces is a convenience function that returns NumberDecimalPlacesEnum wrapped in CountFieldCreateFieldNumberDecimalPlaces
func NumberDecimalPlacesEnumAsCountFieldCreateFieldNumberDecimalPlaces(v *NumberDecimalPlacesEnum) CountFieldCreateFieldNumberDecimalPlaces {
	return CountFieldCreateFieldNumberDecimalPlaces{
		NumberDecimalPlacesEnum: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *CountFieldCreateFieldNumberDecimalPlaces) UnmarshalJSON(data []byte) error {
	var err error
	// this object is nullable so check if the payload is null or empty string
	if string(data) == "" || string(data) == "{}" {
		return nil
	}

	match := 0
	// try to unmarshal data into NullEnum
	err = newStrictDecoder(data).Decode(&dst.NullEnum)
	if err == nil {
		jsonNullEnum, _ := json.Marshal(dst.NullEnum)
		if string(jsonNullEnum) == "{}" { // empty struct
			dst.NullEnum = nil
		} else {
			match++
		}
	} else {
		dst.NullEnum = nil
	}

	// try to unmarshal data into NumberDecimalPlacesEnum
	err = newStrictDecoder(data).Decode(&dst.NumberDecimalPlacesEnum)
	if err == nil {
		jsonNumberDecimalPlacesEnum, _ := json.Marshal(dst.NumberDecimalPlacesEnum)
		if string(jsonNumberDecimalPlacesEnum) == "{}" { // empty struct
			dst.NumberDecimalPlacesEnum = nil
		} else {
			match++
		}
	} else {
		dst.NumberDecimalPlacesEnum = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.NullEnum = nil
		dst.NumberDecimalPlacesEnum = nil

		return fmt.Errorf("data matches more than one schema in oneOf(CountFieldCreateFieldNumberDecimalPlaces)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(CountFieldCreateFieldNumberDecimalPlaces)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src CountFieldCreateFieldNumberDecimalPlaces) MarshalJSON() ([]byte, error) {
	if src.NullEnum != nil {
		return json.Marshal(&src.NullEnum)
	}

	if src.NumberDecimalPlacesEnum != nil {
		return json.Marshal(&src.NumberDecimalPlacesEnum)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *CountFieldCreateFieldNumberDecimalPlaces) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.NullEnum != nil {
		return obj.NullEnum
	}

	if obj.NumberDecimalPlacesEnum != nil {
		return obj.NumberDecimalPlacesEnum
	}

	// all schemas are nil
	return nil
}

type NullableCountFieldCreateFieldNumberDecimalPlaces struct {
	value *CountFieldCreateFieldNumberDecimalPlaces
	isSet bool
}

func (v NullableCountFieldCreateFieldNumberDecimalPlaces) Get() *CountFieldCreateFieldNumberDecimalPlaces {
	return v.value
}

func (v *NullableCountFieldCreateFieldNumberDecimalPlaces) Set(val *CountFieldCreateFieldNumberDecimalPlaces) {
	v.value = val
	v.isSet = true
}

func (v NullableCountFieldCreateFieldNumberDecimalPlaces) IsSet() bool {
	return v.isSet
}

func (v *NullableCountFieldCreateFieldNumberDecimalPlaces) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCountFieldCreateFieldNumberDecimalPlaces(val *CountFieldCreateFieldNumberDecimalPlaces) *NullableCountFieldCreateFieldNumberDecimalPlaces {
	return &NullableCountFieldCreateFieldNumberDecimalPlaces{value: val, isSet: true}
}

func (v NullableCountFieldCreateFieldNumberDecimalPlaces) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCountFieldCreateFieldNumberDecimalPlaces) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


