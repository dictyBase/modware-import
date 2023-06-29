/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"time"
)

// checks if the RowComment type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RowComment{}

// RowComment struct for RowComment
type RowComment struct {
	Id int32 `json:"id"`
	// The table the row this comment is for is found in. 
	TableId int32 `json:"table_id"`
	// The id of the row the comment is for.
	RowId int32 `json:"row_id"`
	// The users comment.
	Comment string `json:"comment"`
	FirstName *string `json:"first_name,omitempty"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
	// The user who made the comment.
	UserId NullableInt32 `json:"user_id"`
}

// NewRowComment instantiates a new RowComment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRowComment(id int32, tableId int32, rowId int32, comment string, createdOn time.Time, updatedOn time.Time, userId NullableInt32) *RowComment {
	this := RowComment{}
	this.Id = id
	this.TableId = tableId
	this.RowId = rowId
	this.Comment = comment
	this.CreatedOn = createdOn
	this.UpdatedOn = updatedOn
	this.UserId = userId
	return &this
}

// NewRowCommentWithDefaults instantiates a new RowComment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRowCommentWithDefaults() *RowComment {
	this := RowComment{}
	return &this
}

// GetId returns the Id field value
func (o *RowComment) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *RowComment) SetId(v int32) {
	o.Id = v
}

// GetTableId returns the TableId field value
func (o *RowComment) GetTableId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TableId
}

// GetTableIdOk returns a tuple with the TableId field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetTableIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TableId, true
}

// SetTableId sets field value
func (o *RowComment) SetTableId(v int32) {
	o.TableId = v
}

// GetRowId returns the RowId field value
func (o *RowComment) GetRowId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.RowId
}

// GetRowIdOk returns a tuple with the RowId field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetRowIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RowId, true
}

// SetRowId sets field value
func (o *RowComment) SetRowId(v int32) {
	o.RowId = v
}

// GetComment returns the Comment field value
func (o *RowComment) GetComment() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Comment
}

// GetCommentOk returns a tuple with the Comment field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetCommentOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Comment, true
}

// SetComment sets field value
func (o *RowComment) SetComment(v string) {
	o.Comment = v
}

// GetFirstName returns the FirstName field value if set, zero value otherwise.
func (o *RowComment) GetFirstName() string {
	if o == nil || IsNil(o.FirstName) {
		var ret string
		return ret
	}
	return *o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RowComment) GetFirstNameOk() (*string, bool) {
	if o == nil || IsNil(o.FirstName) {
		return nil, false
	}
	return o.FirstName, true
}

// HasFirstName returns a boolean if a field has been set.
func (o *RowComment) HasFirstName() bool {
	if o != nil && !IsNil(o.FirstName) {
		return true
	}

	return false
}

// SetFirstName gets a reference to the given string and assigns it to the FirstName field.
func (o *RowComment) SetFirstName(v string) {
	o.FirstName = &v
}

// GetCreatedOn returns the CreatedOn field value
func (o *RowComment) GetCreatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedOn
}

// GetCreatedOnOk returns a tuple with the CreatedOn field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetCreatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedOn, true
}

// SetCreatedOn sets field value
func (o *RowComment) SetCreatedOn(v time.Time) {
	o.CreatedOn = v
}

// GetUpdatedOn returns the UpdatedOn field value
func (o *RowComment) GetUpdatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.UpdatedOn
}

// GetUpdatedOnOk returns a tuple with the UpdatedOn field value
// and a boolean to check if the value has been set.
func (o *RowComment) GetUpdatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedOn, true
}

// SetUpdatedOn sets field value
func (o *RowComment) SetUpdatedOn(v time.Time) {
	o.UpdatedOn = v
}

// GetUserId returns the UserId field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *RowComment) GetUserId() int32 {
	if o == nil || o.UserId.Get() == nil {
		var ret int32
		return ret
	}

	return *o.UserId.Get()
}

// GetUserIdOk returns a tuple with the UserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RowComment) GetUserIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.UserId.Get(), o.UserId.IsSet()
}

// SetUserId sets field value
func (o *RowComment) SetUserId(v int32) {
	o.UserId.Set(&v)
}

func (o RowComment) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RowComment) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: table_id is readOnly
	toSerialize["row_id"] = o.RowId
	toSerialize["comment"] = o.Comment
	if !IsNil(o.FirstName) {
		toSerialize["first_name"] = o.FirstName
	}
	// skip: created_on is readOnly
	// skip: updated_on is readOnly
	toSerialize["user_id"] = o.UserId.Get()
	return toSerialize, nil
}

type NullableRowComment struct {
	value *RowComment
	isSet bool
}

func (v NullableRowComment) Get() *RowComment {
	return v.value
}

func (v *NullableRowComment) Set(val *RowComment) {
	v.value = val
	v.isSet = true
}

func (v NullableRowComment) IsSet() bool {
	return v.isSet
}

func (v *NullableRowComment) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRowComment(val *RowComment) *NullableRowComment {
	return &NullableRowComment{value: val, isSet: true}
}

func (v NullableRowComment) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRowComment) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


