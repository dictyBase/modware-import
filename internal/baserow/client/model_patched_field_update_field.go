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

// PatchedFieldUpdateField - struct for PatchedFieldUpdateField
type PatchedFieldUpdateField struct {
	RequestBooleanFieldUpdateField *RequestBooleanFieldUpdateField
	RequestCountFieldUpdateField *RequestCountFieldUpdateField
	RequestCreatedOnFieldUpdateField *RequestCreatedOnFieldUpdateField
	RequestDateFieldUpdateField *RequestDateFieldUpdateField
	RequestEmailFieldUpdateField *RequestEmailFieldUpdateField
	RequestFileFieldUpdateField *RequestFileFieldUpdateField
	RequestFormulaFieldUpdateField *RequestFormulaFieldUpdateField
	RequestLastModifiedFieldUpdateField *RequestLastModifiedFieldUpdateField
	RequestLinkRowFieldUpdateField *RequestLinkRowFieldUpdateField
	RequestLongTextFieldUpdateField *RequestLongTextFieldUpdateField
	RequestLookupFieldUpdateField *RequestLookupFieldUpdateField
	RequestMultipleCollaboratorsFieldUpdateField *RequestMultipleCollaboratorsFieldUpdateField
	RequestMultipleSelectFieldUpdateField *RequestMultipleSelectFieldUpdateField
	RequestNumberFieldUpdateField *RequestNumberFieldUpdateField
	RequestPhoneNumberFieldUpdateField *RequestPhoneNumberFieldUpdateField
	RequestRatingFieldUpdateField *RequestRatingFieldUpdateField
	RequestRollupFieldUpdateField *RequestRollupFieldUpdateField
	RequestSingleSelectFieldUpdateField *RequestSingleSelectFieldUpdateField
	RequestTextFieldUpdateField *RequestTextFieldUpdateField
	RequestURLFieldUpdateField *RequestURLFieldUpdateField
}

// RequestBooleanFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestBooleanFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestBooleanFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestBooleanFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestBooleanFieldUpdateField: v,
	}
}

// RequestCountFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestCountFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestCountFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestCountFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestCountFieldUpdateField: v,
	}
}

// RequestCreatedOnFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestCreatedOnFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestCreatedOnFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestCreatedOnFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestCreatedOnFieldUpdateField: v,
	}
}

// RequestDateFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestDateFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestDateFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestDateFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestDateFieldUpdateField: v,
	}
}

// RequestEmailFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestEmailFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestEmailFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestEmailFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestEmailFieldUpdateField: v,
	}
}

// RequestFileFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestFileFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestFileFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestFileFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestFileFieldUpdateField: v,
	}
}

// RequestFormulaFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestFormulaFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestFormulaFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestFormulaFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestFormulaFieldUpdateField: v,
	}
}

// RequestLastModifiedFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestLastModifiedFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestLastModifiedFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestLastModifiedFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestLastModifiedFieldUpdateField: v,
	}
}

// RequestLinkRowFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestLinkRowFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestLinkRowFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestLinkRowFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestLinkRowFieldUpdateField: v,
	}
}

// RequestLongTextFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestLongTextFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestLongTextFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestLongTextFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestLongTextFieldUpdateField: v,
	}
}

// RequestLookupFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestLookupFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestLookupFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestLookupFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestLookupFieldUpdateField: v,
	}
}

// RequestMultipleCollaboratorsFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestMultipleCollaboratorsFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestMultipleCollaboratorsFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestMultipleCollaboratorsFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestMultipleCollaboratorsFieldUpdateField: v,
	}
}

// RequestMultipleSelectFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestMultipleSelectFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestMultipleSelectFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestMultipleSelectFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestMultipleSelectFieldUpdateField: v,
	}
}

// RequestNumberFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestNumberFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestNumberFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestNumberFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestNumberFieldUpdateField: v,
	}
}

// RequestPhoneNumberFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestPhoneNumberFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestPhoneNumberFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestPhoneNumberFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestPhoneNumberFieldUpdateField: v,
	}
}

// RequestRatingFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestRatingFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestRatingFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestRatingFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestRatingFieldUpdateField: v,
	}
}

// RequestRollupFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestRollupFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestRollupFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestRollupFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestRollupFieldUpdateField: v,
	}
}

// RequestSingleSelectFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestSingleSelectFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestSingleSelectFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestSingleSelectFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestSingleSelectFieldUpdateField: v,
	}
}

// RequestTextFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestTextFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestTextFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestTextFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestTextFieldUpdateField: v,
	}
}

// RequestURLFieldUpdateFieldAsPatchedFieldUpdateField is a convenience function that returns RequestURLFieldUpdateField wrapped in PatchedFieldUpdateField
func RequestURLFieldUpdateFieldAsPatchedFieldUpdateField(v *RequestURLFieldUpdateField) PatchedFieldUpdateField {
	return PatchedFieldUpdateField{
		RequestURLFieldUpdateField: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *PatchedFieldUpdateField) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into RequestBooleanFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestBooleanFieldUpdateField)
	if err == nil {
		jsonRequestBooleanFieldUpdateField, _ := json.Marshal(dst.RequestBooleanFieldUpdateField)
		if string(jsonRequestBooleanFieldUpdateField) == "{}" { // empty struct
			dst.RequestBooleanFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestBooleanFieldUpdateField = nil
	}

	// try to unmarshal data into RequestCountFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestCountFieldUpdateField)
	if err == nil {
		jsonRequestCountFieldUpdateField, _ := json.Marshal(dst.RequestCountFieldUpdateField)
		if string(jsonRequestCountFieldUpdateField) == "{}" { // empty struct
			dst.RequestCountFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestCountFieldUpdateField = nil
	}

	// try to unmarshal data into RequestCreatedOnFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestCreatedOnFieldUpdateField)
	if err == nil {
		jsonRequestCreatedOnFieldUpdateField, _ := json.Marshal(dst.RequestCreatedOnFieldUpdateField)
		if string(jsonRequestCreatedOnFieldUpdateField) == "{}" { // empty struct
			dst.RequestCreatedOnFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestCreatedOnFieldUpdateField = nil
	}

	// try to unmarshal data into RequestDateFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestDateFieldUpdateField)
	if err == nil {
		jsonRequestDateFieldUpdateField, _ := json.Marshal(dst.RequestDateFieldUpdateField)
		if string(jsonRequestDateFieldUpdateField) == "{}" { // empty struct
			dst.RequestDateFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestDateFieldUpdateField = nil
	}

	// try to unmarshal data into RequestEmailFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestEmailFieldUpdateField)
	if err == nil {
		jsonRequestEmailFieldUpdateField, _ := json.Marshal(dst.RequestEmailFieldUpdateField)
		if string(jsonRequestEmailFieldUpdateField) == "{}" { // empty struct
			dst.RequestEmailFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestEmailFieldUpdateField = nil
	}

	// try to unmarshal data into RequestFileFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestFileFieldUpdateField)
	if err == nil {
		jsonRequestFileFieldUpdateField, _ := json.Marshal(dst.RequestFileFieldUpdateField)
		if string(jsonRequestFileFieldUpdateField) == "{}" { // empty struct
			dst.RequestFileFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestFileFieldUpdateField = nil
	}

	// try to unmarshal data into RequestFormulaFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestFormulaFieldUpdateField)
	if err == nil {
		jsonRequestFormulaFieldUpdateField, _ := json.Marshal(dst.RequestFormulaFieldUpdateField)
		if string(jsonRequestFormulaFieldUpdateField) == "{}" { // empty struct
			dst.RequestFormulaFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestFormulaFieldUpdateField = nil
	}

	// try to unmarshal data into RequestLastModifiedFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestLastModifiedFieldUpdateField)
	if err == nil {
		jsonRequestLastModifiedFieldUpdateField, _ := json.Marshal(dst.RequestLastModifiedFieldUpdateField)
		if string(jsonRequestLastModifiedFieldUpdateField) == "{}" { // empty struct
			dst.RequestLastModifiedFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestLastModifiedFieldUpdateField = nil
	}

	// try to unmarshal data into RequestLinkRowFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestLinkRowFieldUpdateField)
	if err == nil {
		jsonRequestLinkRowFieldUpdateField, _ := json.Marshal(dst.RequestLinkRowFieldUpdateField)
		if string(jsonRequestLinkRowFieldUpdateField) == "{}" { // empty struct
			dst.RequestLinkRowFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestLinkRowFieldUpdateField = nil
	}

	// try to unmarshal data into RequestLongTextFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestLongTextFieldUpdateField)
	if err == nil {
		jsonRequestLongTextFieldUpdateField, _ := json.Marshal(dst.RequestLongTextFieldUpdateField)
		if string(jsonRequestLongTextFieldUpdateField) == "{}" { // empty struct
			dst.RequestLongTextFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestLongTextFieldUpdateField = nil
	}

	// try to unmarshal data into RequestLookupFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestLookupFieldUpdateField)
	if err == nil {
		jsonRequestLookupFieldUpdateField, _ := json.Marshal(dst.RequestLookupFieldUpdateField)
		if string(jsonRequestLookupFieldUpdateField) == "{}" { // empty struct
			dst.RequestLookupFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestLookupFieldUpdateField = nil
	}

	// try to unmarshal data into RequestMultipleCollaboratorsFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestMultipleCollaboratorsFieldUpdateField)
	if err == nil {
		jsonRequestMultipleCollaboratorsFieldUpdateField, _ := json.Marshal(dst.RequestMultipleCollaboratorsFieldUpdateField)
		if string(jsonRequestMultipleCollaboratorsFieldUpdateField) == "{}" { // empty struct
			dst.RequestMultipleCollaboratorsFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestMultipleCollaboratorsFieldUpdateField = nil
	}

	// try to unmarshal data into RequestMultipleSelectFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestMultipleSelectFieldUpdateField)
	if err == nil {
		jsonRequestMultipleSelectFieldUpdateField, _ := json.Marshal(dst.RequestMultipleSelectFieldUpdateField)
		if string(jsonRequestMultipleSelectFieldUpdateField) == "{}" { // empty struct
			dst.RequestMultipleSelectFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestMultipleSelectFieldUpdateField = nil
	}

	// try to unmarshal data into RequestNumberFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestNumberFieldUpdateField)
	if err == nil {
		jsonRequestNumberFieldUpdateField, _ := json.Marshal(dst.RequestNumberFieldUpdateField)
		if string(jsonRequestNumberFieldUpdateField) == "{}" { // empty struct
			dst.RequestNumberFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestNumberFieldUpdateField = nil
	}

	// try to unmarshal data into RequestPhoneNumberFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestPhoneNumberFieldUpdateField)
	if err == nil {
		jsonRequestPhoneNumberFieldUpdateField, _ := json.Marshal(dst.RequestPhoneNumberFieldUpdateField)
		if string(jsonRequestPhoneNumberFieldUpdateField) == "{}" { // empty struct
			dst.RequestPhoneNumberFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestPhoneNumberFieldUpdateField = nil
	}

	// try to unmarshal data into RequestRatingFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestRatingFieldUpdateField)
	if err == nil {
		jsonRequestRatingFieldUpdateField, _ := json.Marshal(dst.RequestRatingFieldUpdateField)
		if string(jsonRequestRatingFieldUpdateField) == "{}" { // empty struct
			dst.RequestRatingFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestRatingFieldUpdateField = nil
	}

	// try to unmarshal data into RequestRollupFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestRollupFieldUpdateField)
	if err == nil {
		jsonRequestRollupFieldUpdateField, _ := json.Marshal(dst.RequestRollupFieldUpdateField)
		if string(jsonRequestRollupFieldUpdateField) == "{}" { // empty struct
			dst.RequestRollupFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestRollupFieldUpdateField = nil
	}

	// try to unmarshal data into RequestSingleSelectFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestSingleSelectFieldUpdateField)
	if err == nil {
		jsonRequestSingleSelectFieldUpdateField, _ := json.Marshal(dst.RequestSingleSelectFieldUpdateField)
		if string(jsonRequestSingleSelectFieldUpdateField) == "{}" { // empty struct
			dst.RequestSingleSelectFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestSingleSelectFieldUpdateField = nil
	}

	// try to unmarshal data into RequestTextFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestTextFieldUpdateField)
	if err == nil {
		jsonRequestTextFieldUpdateField, _ := json.Marshal(dst.RequestTextFieldUpdateField)
		if string(jsonRequestTextFieldUpdateField) == "{}" { // empty struct
			dst.RequestTextFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestTextFieldUpdateField = nil
	}

	// try to unmarshal data into RequestURLFieldUpdateField
	err = newStrictDecoder(data).Decode(&dst.RequestURLFieldUpdateField)
	if err == nil {
		jsonRequestURLFieldUpdateField, _ := json.Marshal(dst.RequestURLFieldUpdateField)
		if string(jsonRequestURLFieldUpdateField) == "{}" { // empty struct
			dst.RequestURLFieldUpdateField = nil
		} else {
			match++
		}
	} else {
		dst.RequestURLFieldUpdateField = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.RequestBooleanFieldUpdateField = nil
		dst.RequestCountFieldUpdateField = nil
		dst.RequestCreatedOnFieldUpdateField = nil
		dst.RequestDateFieldUpdateField = nil
		dst.RequestEmailFieldUpdateField = nil
		dst.RequestFileFieldUpdateField = nil
		dst.RequestFormulaFieldUpdateField = nil
		dst.RequestLastModifiedFieldUpdateField = nil
		dst.RequestLinkRowFieldUpdateField = nil
		dst.RequestLongTextFieldUpdateField = nil
		dst.RequestLookupFieldUpdateField = nil
		dst.RequestMultipleCollaboratorsFieldUpdateField = nil
		dst.RequestMultipleSelectFieldUpdateField = nil
		dst.RequestNumberFieldUpdateField = nil
		dst.RequestPhoneNumberFieldUpdateField = nil
		dst.RequestRatingFieldUpdateField = nil
		dst.RequestRollupFieldUpdateField = nil
		dst.RequestSingleSelectFieldUpdateField = nil
		dst.RequestTextFieldUpdateField = nil
		dst.RequestURLFieldUpdateField = nil

		return fmt.Errorf("data matches more than one schema in oneOf(PatchedFieldUpdateField)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(PatchedFieldUpdateField)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src PatchedFieldUpdateField) MarshalJSON() ([]byte, error) {
	if src.RequestBooleanFieldUpdateField != nil {
		return json.Marshal(&src.RequestBooleanFieldUpdateField)
	}

	if src.RequestCountFieldUpdateField != nil {
		return json.Marshal(&src.RequestCountFieldUpdateField)
	}

	if src.RequestCreatedOnFieldUpdateField != nil {
		return json.Marshal(&src.RequestCreatedOnFieldUpdateField)
	}

	if src.RequestDateFieldUpdateField != nil {
		return json.Marshal(&src.RequestDateFieldUpdateField)
	}

	if src.RequestEmailFieldUpdateField != nil {
		return json.Marshal(&src.RequestEmailFieldUpdateField)
	}

	if src.RequestFileFieldUpdateField != nil {
		return json.Marshal(&src.RequestFileFieldUpdateField)
	}

	if src.RequestFormulaFieldUpdateField != nil {
		return json.Marshal(&src.RequestFormulaFieldUpdateField)
	}

	if src.RequestLastModifiedFieldUpdateField != nil {
		return json.Marshal(&src.RequestLastModifiedFieldUpdateField)
	}

	if src.RequestLinkRowFieldUpdateField != nil {
		return json.Marshal(&src.RequestLinkRowFieldUpdateField)
	}

	if src.RequestLongTextFieldUpdateField != nil {
		return json.Marshal(&src.RequestLongTextFieldUpdateField)
	}

	if src.RequestLookupFieldUpdateField != nil {
		return json.Marshal(&src.RequestLookupFieldUpdateField)
	}

	if src.RequestMultipleCollaboratorsFieldUpdateField != nil {
		return json.Marshal(&src.RequestMultipleCollaboratorsFieldUpdateField)
	}

	if src.RequestMultipleSelectFieldUpdateField != nil {
		return json.Marshal(&src.RequestMultipleSelectFieldUpdateField)
	}

	if src.RequestNumberFieldUpdateField != nil {
		return json.Marshal(&src.RequestNumberFieldUpdateField)
	}

	if src.RequestPhoneNumberFieldUpdateField != nil {
		return json.Marshal(&src.RequestPhoneNumberFieldUpdateField)
	}

	if src.RequestRatingFieldUpdateField != nil {
		return json.Marshal(&src.RequestRatingFieldUpdateField)
	}

	if src.RequestRollupFieldUpdateField != nil {
		return json.Marshal(&src.RequestRollupFieldUpdateField)
	}

	if src.RequestSingleSelectFieldUpdateField != nil {
		return json.Marshal(&src.RequestSingleSelectFieldUpdateField)
	}

	if src.RequestTextFieldUpdateField != nil {
		return json.Marshal(&src.RequestTextFieldUpdateField)
	}

	if src.RequestURLFieldUpdateField != nil {
		return json.Marshal(&src.RequestURLFieldUpdateField)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *PatchedFieldUpdateField) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.RequestBooleanFieldUpdateField != nil {
		return obj.RequestBooleanFieldUpdateField
	}

	if obj.RequestCountFieldUpdateField != nil {
		return obj.RequestCountFieldUpdateField
	}

	if obj.RequestCreatedOnFieldUpdateField != nil {
		return obj.RequestCreatedOnFieldUpdateField
	}

	if obj.RequestDateFieldUpdateField != nil {
		return obj.RequestDateFieldUpdateField
	}

	if obj.RequestEmailFieldUpdateField != nil {
		return obj.RequestEmailFieldUpdateField
	}

	if obj.RequestFileFieldUpdateField != nil {
		return obj.RequestFileFieldUpdateField
	}

	if obj.RequestFormulaFieldUpdateField != nil {
		return obj.RequestFormulaFieldUpdateField
	}

	if obj.RequestLastModifiedFieldUpdateField != nil {
		return obj.RequestLastModifiedFieldUpdateField
	}

	if obj.RequestLinkRowFieldUpdateField != nil {
		return obj.RequestLinkRowFieldUpdateField
	}

	if obj.RequestLongTextFieldUpdateField != nil {
		return obj.RequestLongTextFieldUpdateField
	}

	if obj.RequestLookupFieldUpdateField != nil {
		return obj.RequestLookupFieldUpdateField
	}

	if obj.RequestMultipleCollaboratorsFieldUpdateField != nil {
		return obj.RequestMultipleCollaboratorsFieldUpdateField
	}

	if obj.RequestMultipleSelectFieldUpdateField != nil {
		return obj.RequestMultipleSelectFieldUpdateField
	}

	if obj.RequestNumberFieldUpdateField != nil {
		return obj.RequestNumberFieldUpdateField
	}

	if obj.RequestPhoneNumberFieldUpdateField != nil {
		return obj.RequestPhoneNumberFieldUpdateField
	}

	if obj.RequestRatingFieldUpdateField != nil {
		return obj.RequestRatingFieldUpdateField
	}

	if obj.RequestRollupFieldUpdateField != nil {
		return obj.RequestRollupFieldUpdateField
	}

	if obj.RequestSingleSelectFieldUpdateField != nil {
		return obj.RequestSingleSelectFieldUpdateField
	}

	if obj.RequestTextFieldUpdateField != nil {
		return obj.RequestTextFieldUpdateField
	}

	if obj.RequestURLFieldUpdateField != nil {
		return obj.RequestURLFieldUpdateField
	}

	// all schemas are nil
	return nil
}

type NullablePatchedFieldUpdateField struct {
	value *PatchedFieldUpdateField
	isSet bool
}

func (v NullablePatchedFieldUpdateField) Get() *PatchedFieldUpdateField {
	return v.value
}

func (v *NullablePatchedFieldUpdateField) Set(val *PatchedFieldUpdateField) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedFieldUpdateField) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedFieldUpdateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedFieldUpdateField(val *PatchedFieldUpdateField) *NullablePatchedFieldUpdateField {
	return &NullablePatchedFieldUpdateField{value: val, isSet: true}
}

func (v NullablePatchedFieldUpdateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedFieldUpdateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


