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

// ExportCharsetEnum * `utf-8` - utf-8 * `iso-8859-6` - iso-8859-6 * `windows-1256` - windows-1256 * `iso-8859-4` - iso-8859-4 * `windows-1257` - windows-1257 * `iso-8859-14` - iso-8859-14 * `iso-8859-2` - iso-8859-2 * `windows-1250` - windows-1250 * `gbk` - gbk * `gb18030` - gb18030 * `big5` - big5 * `koi8-r` - koi8-r * `koi8-u` - koi8-u * `iso-8859-5` - iso-8859-5 * `windows-1251` - windows-1251 * `x-mac-cyrillic` - mac-cyrillic * `iso-8859-7` - iso-8859-7 * `windows-1253` - windows-1253 * `iso-8859-8` - iso-8859-8 * `windows-1255` - windows-1255 * `euc-jp` - euc-jp * `iso-2022-jp` - iso-2022-jp * `shift-jis` - shift-jis * `euc-kr` - euc-kr * `macintosh` - macintosh * `iso-8859-10` - iso-8859-10 * `iso-8859-16` - iso-8859-16 * `windows-874` - cp874 * `windows-1254` - windows-1254 * `windows-1258` - windows-1258 * `iso-8859-1` - iso-8859-1 * `windows-1252` - windows-1252 * `iso-8859-3` - iso-8859-3
type ExportCharsetEnum string

// List of ExportCharsetEnum
const (
	UTF_8 ExportCharsetEnum = "utf-8"
	ISO_8859_6 ExportCharsetEnum = "iso-8859-6"
	WINDOWS_1256 ExportCharsetEnum = "windows-1256"
	ISO_8859_4 ExportCharsetEnum = "iso-8859-4"
	WINDOWS_1257 ExportCharsetEnum = "windows-1257"
	ISO_8859_14 ExportCharsetEnum = "iso-8859-14"
	ISO_8859_2 ExportCharsetEnum = "iso-8859-2"
	WINDOWS_1250 ExportCharsetEnum = "windows-1250"
	GBK ExportCharsetEnum = "gbk"
	GB18030 ExportCharsetEnum = "gb18030"
	BIG5 ExportCharsetEnum = "big5"
	KOI8_R ExportCharsetEnum = "koi8-r"
	KOI8_U ExportCharsetEnum = "koi8-u"
	ISO_8859_5 ExportCharsetEnum = "iso-8859-5"
	WINDOWS_1251 ExportCharsetEnum = "windows-1251"
	X_MAC_CYRILLIC ExportCharsetEnum = "x-mac-cyrillic"
	ISO_8859_7 ExportCharsetEnum = "iso-8859-7"
	WINDOWS_1253 ExportCharsetEnum = "windows-1253"
	ISO_8859_8 ExportCharsetEnum = "iso-8859-8"
	WINDOWS_1255 ExportCharsetEnum = "windows-1255"
	EUC_JP ExportCharsetEnum = "euc-jp"
	ISO_2022_JP ExportCharsetEnum = "iso-2022-jp"
	SHIFT_JIS ExportCharsetEnum = "shift-jis"
	EUC_KR ExportCharsetEnum = "euc-kr"
	MACINTOSH ExportCharsetEnum = "macintosh"
	ISO_8859_10 ExportCharsetEnum = "iso-8859-10"
	ISO_8859_16 ExportCharsetEnum = "iso-8859-16"
	WINDOWS_874 ExportCharsetEnum = "windows-874"
	WINDOWS_1254 ExportCharsetEnum = "windows-1254"
	WINDOWS_1258 ExportCharsetEnum = "windows-1258"
	ISO_8859_1 ExportCharsetEnum = "iso-8859-1"
	WINDOWS_1252 ExportCharsetEnum = "windows-1252"
	ISO_8859_3 ExportCharsetEnum = "iso-8859-3"
)

// All allowed values of ExportCharsetEnum enum
var AllowedExportCharsetEnumEnumValues = []ExportCharsetEnum{
	"utf-8",
	"iso-8859-6",
	"windows-1256",
	"iso-8859-4",
	"windows-1257",
	"iso-8859-14",
	"iso-8859-2",
	"windows-1250",
	"gbk",
	"gb18030",
	"big5",
	"koi8-r",
	"koi8-u",
	"iso-8859-5",
	"windows-1251",
	"x-mac-cyrillic",
	"iso-8859-7",
	"windows-1253",
	"iso-8859-8",
	"windows-1255",
	"euc-jp",
	"iso-2022-jp",
	"shift-jis",
	"euc-kr",
	"macintosh",
	"iso-8859-10",
	"iso-8859-16",
	"windows-874",
	"windows-1254",
	"windows-1258",
	"iso-8859-1",
	"windows-1252",
	"iso-8859-3",
}

func (v *ExportCharsetEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ExportCharsetEnum(value)
	for _, existing := range AllowedExportCharsetEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ExportCharsetEnum", value)
}

// NewExportCharsetEnumFromValue returns a pointer to a valid ExportCharsetEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewExportCharsetEnumFromValue(v string) (*ExportCharsetEnum, error) {
	ev := ExportCharsetEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ExportCharsetEnum: valid values are %v", v, AllowedExportCharsetEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ExportCharsetEnum) IsValid() bool {
	for _, existing := range AllowedExportCharsetEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ExportCharsetEnum value
func (v ExportCharsetEnum) Ptr() *ExportCharsetEnum {
	return &v
}

type NullableExportCharsetEnum struct {
	value *ExportCharsetEnum
	isSet bool
}

func (v NullableExportCharsetEnum) Get() *ExportCharsetEnum {
	return v.value
}

func (v *NullableExportCharsetEnum) Set(val *ExportCharsetEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableExportCharsetEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableExportCharsetEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportCharsetEnum(val *ExportCharsetEnum) *NullableExportCharsetEnum {
	return &NullableExportCharsetEnum{value: val, isSet: true}
}

func (v NullableExportCharsetEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportCharsetEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

