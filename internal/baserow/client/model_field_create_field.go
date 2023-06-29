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

// FieldCreateField - struct for FieldCreateField
type FieldCreateField struct {
	BooleanFieldCreateField *BooleanFieldCreateField
	CountFieldCreateField *CountFieldCreateField
	CreatedOnFieldCreateField *CreatedOnFieldCreateField
	DateFieldCreateField *DateFieldCreateField
	EmailFieldCreateField *EmailFieldCreateField
	FileFieldCreateField *FileFieldCreateField
	FormulaFieldCreateField *FormulaFieldCreateField
	LastModifiedFieldCreateField *LastModifiedFieldCreateField
	LinkRowFieldCreateField *LinkRowFieldCreateField
	LongTextFieldCreateField *LongTextFieldCreateField
	LookupFieldCreateField *LookupFieldCreateField
	MultipleCollaboratorsFieldCreateField *MultipleCollaboratorsFieldCreateField
	MultipleSelectFieldCreateField *MultipleSelectFieldCreateField
	NumberFieldCreateField *NumberFieldCreateField
	PhoneNumberFieldCreateField *PhoneNumberFieldCreateField
	RatingFieldCreateField *RatingFieldCreateField
	RollupFieldCreateField *RollupFieldCreateField
	SingleSelectFieldCreateField *SingleSelectFieldCreateField
	TextFieldCreateField *TextFieldCreateField
	URLFieldCreateField *URLFieldCreateField
}

// BooleanFieldCreateFieldAsFieldCreateField is a convenience function that returns BooleanFieldCreateField wrapped in FieldCreateField
func BooleanFieldCreateFieldAsFieldCreateField(v *BooleanFieldCreateField) FieldCreateField {
	return FieldCreateField{
		BooleanFieldCreateField: v,
	}
}

// CountFieldCreateFieldAsFieldCreateField is a convenience function that returns CountFieldCreateField wrapped in FieldCreateField
func CountFieldCreateFieldAsFieldCreateField(v *CountFieldCreateField) FieldCreateField {
	return FieldCreateField{
		CountFieldCreateField: v,
	}
}

// CreatedOnFieldCreateFieldAsFieldCreateField is a convenience function that returns CreatedOnFieldCreateField wrapped in FieldCreateField
func CreatedOnFieldCreateFieldAsFieldCreateField(v *CreatedOnFieldCreateField) FieldCreateField {
	return FieldCreateField{
		CreatedOnFieldCreateField: v,
	}
}

// DateFieldCreateFieldAsFieldCreateField is a convenience function that returns DateFieldCreateField wrapped in FieldCreateField
func DateFieldCreateFieldAsFieldCreateField(v *DateFieldCreateField) FieldCreateField {
	return FieldCreateField{
		DateFieldCreateField: v,
	}
}

// EmailFieldCreateFieldAsFieldCreateField is a convenience function that returns EmailFieldCreateField wrapped in FieldCreateField
func EmailFieldCreateFieldAsFieldCreateField(v *EmailFieldCreateField) FieldCreateField {
	return FieldCreateField{
		EmailFieldCreateField: v,
	}
}

// FileFieldCreateFieldAsFieldCreateField is a convenience function that returns FileFieldCreateField wrapped in FieldCreateField
func FileFieldCreateFieldAsFieldCreateField(v *FileFieldCreateField) FieldCreateField {
	return FieldCreateField{
		FileFieldCreateField: v,
	}
}

// FormulaFieldCreateFieldAsFieldCreateField is a convenience function that returns FormulaFieldCreateField wrapped in FieldCreateField
func FormulaFieldCreateFieldAsFieldCreateField(v *FormulaFieldCreateField) FieldCreateField {
	return FieldCreateField{
		FormulaFieldCreateField: v,
	}
}

// LastModifiedFieldCreateFieldAsFieldCreateField is a convenience function that returns LastModifiedFieldCreateField wrapped in FieldCreateField
func LastModifiedFieldCreateFieldAsFieldCreateField(v *LastModifiedFieldCreateField) FieldCreateField {
	return FieldCreateField{
		LastModifiedFieldCreateField: v,
	}
}

// LinkRowFieldCreateFieldAsFieldCreateField is a convenience function that returns LinkRowFieldCreateField wrapped in FieldCreateField
func LinkRowFieldCreateFieldAsFieldCreateField(v *LinkRowFieldCreateField) FieldCreateField {
	return FieldCreateField{
		LinkRowFieldCreateField: v,
	}
}

// LongTextFieldCreateFieldAsFieldCreateField is a convenience function that returns LongTextFieldCreateField wrapped in FieldCreateField
func LongTextFieldCreateFieldAsFieldCreateField(v *LongTextFieldCreateField) FieldCreateField {
	return FieldCreateField{
		LongTextFieldCreateField: v,
	}
}

// LookupFieldCreateFieldAsFieldCreateField is a convenience function that returns LookupFieldCreateField wrapped in FieldCreateField
func LookupFieldCreateFieldAsFieldCreateField(v *LookupFieldCreateField) FieldCreateField {
	return FieldCreateField{
		LookupFieldCreateField: v,
	}
}

// MultipleCollaboratorsFieldCreateFieldAsFieldCreateField is a convenience function that returns MultipleCollaboratorsFieldCreateField wrapped in FieldCreateField
func MultipleCollaboratorsFieldCreateFieldAsFieldCreateField(v *MultipleCollaboratorsFieldCreateField) FieldCreateField {
	return FieldCreateField{
		MultipleCollaboratorsFieldCreateField: v,
	}
}

// MultipleSelectFieldCreateFieldAsFieldCreateField is a convenience function that returns MultipleSelectFieldCreateField wrapped in FieldCreateField
func MultipleSelectFieldCreateFieldAsFieldCreateField(v *MultipleSelectFieldCreateField) FieldCreateField {
	return FieldCreateField{
		MultipleSelectFieldCreateField: v,
	}
}

// NumberFieldCreateFieldAsFieldCreateField is a convenience function that returns NumberFieldCreateField wrapped in FieldCreateField
func NumberFieldCreateFieldAsFieldCreateField(v *NumberFieldCreateField) FieldCreateField {
	return FieldCreateField{
		NumberFieldCreateField: v,
	}
}

// PhoneNumberFieldCreateFieldAsFieldCreateField is a convenience function that returns PhoneNumberFieldCreateField wrapped in FieldCreateField
func PhoneNumberFieldCreateFieldAsFieldCreateField(v *PhoneNumberFieldCreateField) FieldCreateField {
	return FieldCreateField{
		PhoneNumberFieldCreateField: v,
	}
}

// RatingFieldCreateFieldAsFieldCreateField is a convenience function that returns RatingFieldCreateField wrapped in FieldCreateField
func RatingFieldCreateFieldAsFieldCreateField(v *RatingFieldCreateField) FieldCreateField {
	return FieldCreateField{
		RatingFieldCreateField: v,
	}
}

// RollupFieldCreateFieldAsFieldCreateField is a convenience function that returns RollupFieldCreateField wrapped in FieldCreateField
func RollupFieldCreateFieldAsFieldCreateField(v *RollupFieldCreateField) FieldCreateField {
	return FieldCreateField{
		RollupFieldCreateField: v,
	}
}

// SingleSelectFieldCreateFieldAsFieldCreateField is a convenience function that returns SingleSelectFieldCreateField wrapped in FieldCreateField
func SingleSelectFieldCreateFieldAsFieldCreateField(v *SingleSelectFieldCreateField) FieldCreateField {
	return FieldCreateField{
		SingleSelectFieldCreateField: v,
	}
}

// TextFieldCreateFieldAsFieldCreateField is a convenience function that returns TextFieldCreateField wrapped in FieldCreateField
func TextFieldCreateFieldAsFieldCreateField(v *TextFieldCreateField) FieldCreateField {
	return FieldCreateField{
		TextFieldCreateField: v,
	}
}

// URLFieldCreateFieldAsFieldCreateField is a convenience function that returns URLFieldCreateField wrapped in FieldCreateField
func URLFieldCreateFieldAsFieldCreateField(v *URLFieldCreateField) FieldCreateField {
	return FieldCreateField{
		URLFieldCreateField: v,
	}
}


// Unmarshal JSON data into one of the pointers in the struct
func (dst *FieldCreateField) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into BooleanFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.BooleanFieldCreateField)
	if err == nil {
		jsonBooleanFieldCreateField, _ := json.Marshal(dst.BooleanFieldCreateField)
		if string(jsonBooleanFieldCreateField) == "{}" { // empty struct
			dst.BooleanFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.BooleanFieldCreateField = nil
	}

	// try to unmarshal data into CountFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.CountFieldCreateField)
	if err == nil {
		jsonCountFieldCreateField, _ := json.Marshal(dst.CountFieldCreateField)
		if string(jsonCountFieldCreateField) == "{}" { // empty struct
			dst.CountFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.CountFieldCreateField = nil
	}

	// try to unmarshal data into CreatedOnFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.CreatedOnFieldCreateField)
	if err == nil {
		jsonCreatedOnFieldCreateField, _ := json.Marshal(dst.CreatedOnFieldCreateField)
		if string(jsonCreatedOnFieldCreateField) == "{}" { // empty struct
			dst.CreatedOnFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.CreatedOnFieldCreateField = nil
	}

	// try to unmarshal data into DateFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.DateFieldCreateField)
	if err == nil {
		jsonDateFieldCreateField, _ := json.Marshal(dst.DateFieldCreateField)
		if string(jsonDateFieldCreateField) == "{}" { // empty struct
			dst.DateFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.DateFieldCreateField = nil
	}

	// try to unmarshal data into EmailFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.EmailFieldCreateField)
	if err == nil {
		jsonEmailFieldCreateField, _ := json.Marshal(dst.EmailFieldCreateField)
		if string(jsonEmailFieldCreateField) == "{}" { // empty struct
			dst.EmailFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.EmailFieldCreateField = nil
	}

	// try to unmarshal data into FileFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.FileFieldCreateField)
	if err == nil {
		jsonFileFieldCreateField, _ := json.Marshal(dst.FileFieldCreateField)
		if string(jsonFileFieldCreateField) == "{}" { // empty struct
			dst.FileFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.FileFieldCreateField = nil
	}

	// try to unmarshal data into FormulaFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.FormulaFieldCreateField)
	if err == nil {
		jsonFormulaFieldCreateField, _ := json.Marshal(dst.FormulaFieldCreateField)
		if string(jsonFormulaFieldCreateField) == "{}" { // empty struct
			dst.FormulaFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.FormulaFieldCreateField = nil
	}

	// try to unmarshal data into LastModifiedFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.LastModifiedFieldCreateField)
	if err == nil {
		jsonLastModifiedFieldCreateField, _ := json.Marshal(dst.LastModifiedFieldCreateField)
		if string(jsonLastModifiedFieldCreateField) == "{}" { // empty struct
			dst.LastModifiedFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.LastModifiedFieldCreateField = nil
	}

	// try to unmarshal data into LinkRowFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.LinkRowFieldCreateField)
	if err == nil {
		jsonLinkRowFieldCreateField, _ := json.Marshal(dst.LinkRowFieldCreateField)
		if string(jsonLinkRowFieldCreateField) == "{}" { // empty struct
			dst.LinkRowFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.LinkRowFieldCreateField = nil
	}

	// try to unmarshal data into LongTextFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.LongTextFieldCreateField)
	if err == nil {
		jsonLongTextFieldCreateField, _ := json.Marshal(dst.LongTextFieldCreateField)
		if string(jsonLongTextFieldCreateField) == "{}" { // empty struct
			dst.LongTextFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.LongTextFieldCreateField = nil
	}

	// try to unmarshal data into LookupFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.LookupFieldCreateField)
	if err == nil {
		jsonLookupFieldCreateField, _ := json.Marshal(dst.LookupFieldCreateField)
		if string(jsonLookupFieldCreateField) == "{}" { // empty struct
			dst.LookupFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.LookupFieldCreateField = nil
	}

	// try to unmarshal data into MultipleCollaboratorsFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.MultipleCollaboratorsFieldCreateField)
	if err == nil {
		jsonMultipleCollaboratorsFieldCreateField, _ := json.Marshal(dst.MultipleCollaboratorsFieldCreateField)
		if string(jsonMultipleCollaboratorsFieldCreateField) == "{}" { // empty struct
			dst.MultipleCollaboratorsFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.MultipleCollaboratorsFieldCreateField = nil
	}

	// try to unmarshal data into MultipleSelectFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.MultipleSelectFieldCreateField)
	if err == nil {
		jsonMultipleSelectFieldCreateField, _ := json.Marshal(dst.MultipleSelectFieldCreateField)
		if string(jsonMultipleSelectFieldCreateField) == "{}" { // empty struct
			dst.MultipleSelectFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.MultipleSelectFieldCreateField = nil
	}

	// try to unmarshal data into NumberFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.NumberFieldCreateField)
	if err == nil {
		jsonNumberFieldCreateField, _ := json.Marshal(dst.NumberFieldCreateField)
		if string(jsonNumberFieldCreateField) == "{}" { // empty struct
			dst.NumberFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.NumberFieldCreateField = nil
	}

	// try to unmarshal data into PhoneNumberFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.PhoneNumberFieldCreateField)
	if err == nil {
		jsonPhoneNumberFieldCreateField, _ := json.Marshal(dst.PhoneNumberFieldCreateField)
		if string(jsonPhoneNumberFieldCreateField) == "{}" { // empty struct
			dst.PhoneNumberFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.PhoneNumberFieldCreateField = nil
	}

	// try to unmarshal data into RatingFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.RatingFieldCreateField)
	if err == nil {
		jsonRatingFieldCreateField, _ := json.Marshal(dst.RatingFieldCreateField)
		if string(jsonRatingFieldCreateField) == "{}" { // empty struct
			dst.RatingFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.RatingFieldCreateField = nil
	}

	// try to unmarshal data into RollupFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.RollupFieldCreateField)
	if err == nil {
		jsonRollupFieldCreateField, _ := json.Marshal(dst.RollupFieldCreateField)
		if string(jsonRollupFieldCreateField) == "{}" { // empty struct
			dst.RollupFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.RollupFieldCreateField = nil
	}

	// try to unmarshal data into SingleSelectFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.SingleSelectFieldCreateField)
	if err == nil {
		jsonSingleSelectFieldCreateField, _ := json.Marshal(dst.SingleSelectFieldCreateField)
		if string(jsonSingleSelectFieldCreateField) == "{}" { // empty struct
			dst.SingleSelectFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.SingleSelectFieldCreateField = nil
	}

	// try to unmarshal data into TextFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.TextFieldCreateField)
	if err == nil {
		jsonTextFieldCreateField, _ := json.Marshal(dst.TextFieldCreateField)
		if string(jsonTextFieldCreateField) == "{}" { // empty struct
			dst.TextFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.TextFieldCreateField = nil
	}

	// try to unmarshal data into URLFieldCreateField
	err = newStrictDecoder(data).Decode(&dst.URLFieldCreateField)
	if err == nil {
		jsonURLFieldCreateField, _ := json.Marshal(dst.URLFieldCreateField)
		if string(jsonURLFieldCreateField) == "{}" { // empty struct
			dst.URLFieldCreateField = nil
		} else {
			match++
		}
	} else {
		dst.URLFieldCreateField = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.BooleanFieldCreateField = nil
		dst.CountFieldCreateField = nil
		dst.CreatedOnFieldCreateField = nil
		dst.DateFieldCreateField = nil
		dst.EmailFieldCreateField = nil
		dst.FileFieldCreateField = nil
		dst.FormulaFieldCreateField = nil
		dst.LastModifiedFieldCreateField = nil
		dst.LinkRowFieldCreateField = nil
		dst.LongTextFieldCreateField = nil
		dst.LookupFieldCreateField = nil
		dst.MultipleCollaboratorsFieldCreateField = nil
		dst.MultipleSelectFieldCreateField = nil
		dst.NumberFieldCreateField = nil
		dst.PhoneNumberFieldCreateField = nil
		dst.RatingFieldCreateField = nil
		dst.RollupFieldCreateField = nil
		dst.SingleSelectFieldCreateField = nil
		dst.TextFieldCreateField = nil
		dst.URLFieldCreateField = nil

		return fmt.Errorf("data matches more than one schema in oneOf(FieldCreateField)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("data failed to match schemas in oneOf(FieldCreateField)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src FieldCreateField) MarshalJSON() ([]byte, error) {
	if src.BooleanFieldCreateField != nil {
		return json.Marshal(&src.BooleanFieldCreateField)
	}

	if src.CountFieldCreateField != nil {
		return json.Marshal(&src.CountFieldCreateField)
	}

	if src.CreatedOnFieldCreateField != nil {
		return json.Marshal(&src.CreatedOnFieldCreateField)
	}

	if src.DateFieldCreateField != nil {
		return json.Marshal(&src.DateFieldCreateField)
	}

	if src.EmailFieldCreateField != nil {
		return json.Marshal(&src.EmailFieldCreateField)
	}

	if src.FileFieldCreateField != nil {
		return json.Marshal(&src.FileFieldCreateField)
	}

	if src.FormulaFieldCreateField != nil {
		return json.Marshal(&src.FormulaFieldCreateField)
	}

	if src.LastModifiedFieldCreateField != nil {
		return json.Marshal(&src.LastModifiedFieldCreateField)
	}

	if src.LinkRowFieldCreateField != nil {
		return json.Marshal(&src.LinkRowFieldCreateField)
	}

	if src.LongTextFieldCreateField != nil {
		return json.Marshal(&src.LongTextFieldCreateField)
	}

	if src.LookupFieldCreateField != nil {
		return json.Marshal(&src.LookupFieldCreateField)
	}

	if src.MultipleCollaboratorsFieldCreateField != nil {
		return json.Marshal(&src.MultipleCollaboratorsFieldCreateField)
	}

	if src.MultipleSelectFieldCreateField != nil {
		return json.Marshal(&src.MultipleSelectFieldCreateField)
	}

	if src.NumberFieldCreateField != nil {
		return json.Marshal(&src.NumberFieldCreateField)
	}

	if src.PhoneNumberFieldCreateField != nil {
		return json.Marshal(&src.PhoneNumberFieldCreateField)
	}

	if src.RatingFieldCreateField != nil {
		return json.Marshal(&src.RatingFieldCreateField)
	}

	if src.RollupFieldCreateField != nil {
		return json.Marshal(&src.RollupFieldCreateField)
	}

	if src.SingleSelectFieldCreateField != nil {
		return json.Marshal(&src.SingleSelectFieldCreateField)
	}

	if src.TextFieldCreateField != nil {
		return json.Marshal(&src.TextFieldCreateField)
	}

	if src.URLFieldCreateField != nil {
		return json.Marshal(&src.URLFieldCreateField)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *FieldCreateField) GetActualInstance() (interface{}) {
	if obj == nil {
		return nil
	}
	if obj.BooleanFieldCreateField != nil {
		return obj.BooleanFieldCreateField
	}

	if obj.CountFieldCreateField != nil {
		return obj.CountFieldCreateField
	}

	if obj.CreatedOnFieldCreateField != nil {
		return obj.CreatedOnFieldCreateField
	}

	if obj.DateFieldCreateField != nil {
		return obj.DateFieldCreateField
	}

	if obj.EmailFieldCreateField != nil {
		return obj.EmailFieldCreateField
	}

	if obj.FileFieldCreateField != nil {
		return obj.FileFieldCreateField
	}

	if obj.FormulaFieldCreateField != nil {
		return obj.FormulaFieldCreateField
	}

	if obj.LastModifiedFieldCreateField != nil {
		return obj.LastModifiedFieldCreateField
	}

	if obj.LinkRowFieldCreateField != nil {
		return obj.LinkRowFieldCreateField
	}

	if obj.LongTextFieldCreateField != nil {
		return obj.LongTextFieldCreateField
	}

	if obj.LookupFieldCreateField != nil {
		return obj.LookupFieldCreateField
	}

	if obj.MultipleCollaboratorsFieldCreateField != nil {
		return obj.MultipleCollaboratorsFieldCreateField
	}

	if obj.MultipleSelectFieldCreateField != nil {
		return obj.MultipleSelectFieldCreateField
	}

	if obj.NumberFieldCreateField != nil {
		return obj.NumberFieldCreateField
	}

	if obj.PhoneNumberFieldCreateField != nil {
		return obj.PhoneNumberFieldCreateField
	}

	if obj.RatingFieldCreateField != nil {
		return obj.RatingFieldCreateField
	}

	if obj.RollupFieldCreateField != nil {
		return obj.RollupFieldCreateField
	}

	if obj.SingleSelectFieldCreateField != nil {
		return obj.SingleSelectFieldCreateField
	}

	if obj.TextFieldCreateField != nil {
		return obj.TextFieldCreateField
	}

	if obj.URLFieldCreateField != nil {
		return obj.URLFieldCreateField
	}

	// all schemas are nil
	return nil
}

type NullableFieldCreateField struct {
	value *FieldCreateField
	isSet bool
}

func (v NullableFieldCreateField) Get() *FieldCreateField {
	return v.value
}

func (v *NullableFieldCreateField) Set(val *FieldCreateField) {
	v.value = val
	v.isSet = true
}

func (v NullableFieldCreateField) IsSet() bool {
	return v.isSet
}

func (v *NullableFieldCreateField) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFieldCreateField(val *FieldCreateField) *NullableFieldCreateField {
	return &NullableFieldCreateField{value: val, isSet: true}
}

func (v NullableFieldCreateField) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFieldCreateField) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


