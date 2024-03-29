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

// GridViewFieldOptionsAggregationRawType - Indicates how to compute the raw aggregation value from database. This type must be registered in the backend prior to use it.  * `empty_count` - empty_count * `not_empty_count` - not_empty_count * `unique_count` - unique_count * `min` - min * `max` - max * `sum` - sum * `average` - average * `median` - median * `decile` - decile * `variance` - variance * `std_dev` - std_dev
type GridViewFieldOptionsAggregationRawType struct {
	AggregationRawTypeEnum *AggregationRawTypeEnum
	BlankEnum *BlankEnum
}

// AggregationRawTypeEnumAsGridViewFieldOptionsAggregationRawType is a convenience function that returns AggregationRawTypeEnum wrapped in GridViewFieldOptionsAggregationRawType
func AggregationRawTypeEnumAsGridViewFieldOptionsAggregationRawType(v *AggregationRawTypeEnum) GridViewFieldOptionsAggregationRawType {
	return GridViewFieldOptionsAggregationRawType{
		AggregationRawTypeEnum: v,
	}
}

// BlankEnumAsGridViewFieldOptionsAggregationRawType is a convenience function that returns BlankEnum wrapped in GridViewFieldOptionsAggregationRawType
func BlankEnumAsGridViewFieldOptionsAggregationRawType(v *BlankEnum) GridViewFieldOptionsAggregationRawType {
	return GridViewFieldOptionsAggregationRawType{
		BlankEnum: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *GridViewFieldOptionsAggregationRawType) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into AggregationRawTypeEnum
	err = newStrictDecoder(data).Decode(&dst.AggregationRawTypeEnum)
	if err == nil {
		jsonAggregationRawTypeEnum, _ := json.Marshal(dst.AggregationRawTypeEnum)
		if string(jsonAggregationRawTypeEnum) == "{}" { // empty struct
			dst.AggregationRawTypeEnum = nil
		} else {
			match++
		}
	} else {
		dst.AggregationRawTypeEnum = nil
	}

	// try to unmarshal data into BlankEnum
	err = newStrictDecoder(data).Decode(&dst.BlankEnum)
	if err == nil {
		jsonBlankEnum, _ := json.Marshal(dst.BlankEnum)
		if string(jsonBlankEnum) == "{}" { // empty struct
			dst.BlankEnum = nil
		} else {
			match++
		}
	} else {
		dst.BlankEnum = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.AggregationRawTypeEnum = nil
		dst.BlankEnum = nil

		return fmt.Errorf("data matches more than one schema in oneOf(GridViewFieldOptionsAggregationRawType)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(GridViewFieldOptionsAggregationRawType)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src GridViewFieldOptionsAggregationRawType) MarshalJSON() ([]byte, error) {
	if src.AggregationRawTypeEnum != nil {
		return json.Marshal(&src.AggregationRawTypeEnum)
	}

	if src.BlankEnum != nil {
		return json.Marshal(&src.BlankEnum)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *GridViewFieldOptionsAggregationRawType) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.AggregationRawTypeEnum != nil {
		return obj.AggregationRawTypeEnum
	}

	if obj.BlankEnum != nil {
		return obj.BlankEnum
	}

	// all schemas are nil
	return nil
}

type NullableGridViewFieldOptionsAggregationRawType struct {
	value *GridViewFieldOptionsAggregationRawType
	isSet bool
}

func (v NullableGridViewFieldOptionsAggregationRawType) Get() *GridViewFieldOptionsAggregationRawType {
	return v.value
}

func (v *NullableGridViewFieldOptionsAggregationRawType) Set(val *GridViewFieldOptionsAggregationRawType) {
	v.value = val
	v.isSet = true
}

func (v NullableGridViewFieldOptionsAggregationRawType) IsSet() bool {
	return v.isSet
}

func (v *NullableGridViewFieldOptionsAggregationRawType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGridViewFieldOptionsAggregationRawType(val *GridViewFieldOptionsAggregationRawType) *NullableGridViewFieldOptionsAggregationRawType {
	return &NullableGridViewFieldOptionsAggregationRawType{value: val, isSet: true}
}

func (v NullableGridViewFieldOptionsAggregationRawType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGridViewFieldOptionsAggregationRawType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


