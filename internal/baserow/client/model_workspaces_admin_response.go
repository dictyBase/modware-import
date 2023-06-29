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

// checks if the WorkspacesAdminResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkspacesAdminResponse{}

// WorkspacesAdminResponse struct for WorkspacesAdminResponse
type WorkspacesAdminResponse struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
	Users []WorkspaceAdminUsers `json:"users"`
	ApplicationCount int32 `json:"application_count"`
	RowCount int32 `json:"row_count"`
	StorageUsage NullableInt32 `json:"storage_usage,omitempty"`
	SeatsTaken int32 `json:"seats_taken"`
	FreeUsers int32 `json:"free_users"`
	CreatedOn time.Time `json:"created_on"`
}

// NewWorkspacesAdminResponse instantiates a new WorkspacesAdminResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkspacesAdminResponse(id int32, name string, users []WorkspaceAdminUsers, applicationCount int32, rowCount int32, seatsTaken int32, freeUsers int32, createdOn time.Time) *WorkspacesAdminResponse {
	this := WorkspacesAdminResponse{}
	this.Id = id
	this.Name = name
	this.Users = users
	this.ApplicationCount = applicationCount
	this.RowCount = rowCount
	this.SeatsTaken = seatsTaken
	this.FreeUsers = freeUsers
	this.CreatedOn = createdOn
	return &this
}

// NewWorkspacesAdminResponseWithDefaults instantiates a new WorkspacesAdminResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkspacesAdminResponseWithDefaults() *WorkspacesAdminResponse {
	this := WorkspacesAdminResponse{}
	return &this
}

// GetId returns the Id field value
func (o *WorkspacesAdminResponse) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkspacesAdminResponse) SetId(v int32) {
	o.Id = v
}

// GetName returns the Name field value
func (o *WorkspacesAdminResponse) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WorkspacesAdminResponse) SetName(v string) {
	o.Name = v
}

// GetUsers returns the Users field value
func (o *WorkspacesAdminResponse) GetUsers() []WorkspaceAdminUsers {
	if o == nil {
		var ret []WorkspaceAdminUsers
		return ret
	}

	return o.Users
}

// GetUsersOk returns a tuple with the Users field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetUsersOk() ([]WorkspaceAdminUsers, bool) {
	if o == nil {
		return nil, false
	}
	return o.Users, true
}

// SetUsers sets field value
func (o *WorkspacesAdminResponse) SetUsers(v []WorkspaceAdminUsers) {
	o.Users = v
}

// GetApplicationCount returns the ApplicationCount field value
func (o *WorkspacesAdminResponse) GetApplicationCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ApplicationCount
}

// GetApplicationCountOk returns a tuple with the ApplicationCount field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetApplicationCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApplicationCount, true
}

// SetApplicationCount sets field value
func (o *WorkspacesAdminResponse) SetApplicationCount(v int32) {
	o.ApplicationCount = v
}

// GetRowCount returns the RowCount field value
func (o *WorkspacesAdminResponse) GetRowCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.RowCount
}

// GetRowCountOk returns a tuple with the RowCount field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetRowCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RowCount, true
}

// SetRowCount sets field value
func (o *WorkspacesAdminResponse) SetRowCount(v int32) {
	o.RowCount = v
}

// GetStorageUsage returns the StorageUsage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WorkspacesAdminResponse) GetStorageUsage() int32 {
	if o == nil || IsNil(o.StorageUsage.Get()) {
		var ret int32
		return ret
	}
	return *o.StorageUsage.Get()
}

// GetStorageUsageOk returns a tuple with the StorageUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WorkspacesAdminResponse) GetStorageUsageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.StorageUsage.Get(), o.StorageUsage.IsSet()
}

// HasStorageUsage returns a boolean if a field has been set.
func (o *WorkspacesAdminResponse) HasStorageUsage() bool {
	if o != nil && o.StorageUsage.IsSet() {
		return true
	}

	return false
}

// SetStorageUsage gets a reference to the given NullableInt32 and assigns it to the StorageUsage field.
func (o *WorkspacesAdminResponse) SetStorageUsage(v int32) {
	o.StorageUsage.Set(&v)
}
// SetStorageUsageNil sets the value for StorageUsage to be an explicit nil
func (o *WorkspacesAdminResponse) SetStorageUsageNil() {
	o.StorageUsage.Set(nil)
}

// UnsetStorageUsage ensures that no value is present for StorageUsage, not even an explicit nil
func (o *WorkspacesAdminResponse) UnsetStorageUsage() {
	o.StorageUsage.Unset()
}

// GetSeatsTaken returns the SeatsTaken field value
func (o *WorkspacesAdminResponse) GetSeatsTaken() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.SeatsTaken
}

// GetSeatsTakenOk returns a tuple with the SeatsTaken field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetSeatsTakenOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SeatsTaken, true
}

// SetSeatsTaken sets field value
func (o *WorkspacesAdminResponse) SetSeatsTaken(v int32) {
	o.SeatsTaken = v
}

// GetFreeUsers returns the FreeUsers field value
func (o *WorkspacesAdminResponse) GetFreeUsers() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.FreeUsers
}

// GetFreeUsersOk returns a tuple with the FreeUsers field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetFreeUsersOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FreeUsers, true
}

// SetFreeUsers sets field value
func (o *WorkspacesAdminResponse) SetFreeUsers(v int32) {
	o.FreeUsers = v
}

// GetCreatedOn returns the CreatedOn field value
func (o *WorkspacesAdminResponse) GetCreatedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedOn
}

// GetCreatedOnOk returns a tuple with the CreatedOn field value
// and a boolean to check if the value has been set.
func (o *WorkspacesAdminResponse) GetCreatedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedOn, true
}

// SetCreatedOn sets field value
func (o *WorkspacesAdminResponse) SetCreatedOn(v time.Time) {
	o.CreatedOn = v
}

func (o WorkspacesAdminResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkspacesAdminResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["name"] = o.Name
	toSerialize["users"] = o.Users
	toSerialize["application_count"] = o.ApplicationCount
	// skip: row_count is readOnly
	if o.StorageUsage.IsSet() {
		toSerialize["storage_usage"] = o.StorageUsage.Get()
	}
	toSerialize["seats_taken"] = o.SeatsTaken
	// skip: free_users is readOnly
	// skip: created_on is readOnly
	return toSerialize, nil
}

type NullableWorkspacesAdminResponse struct {
	value *WorkspacesAdminResponse
	isSet bool
}

func (v NullableWorkspacesAdminResponse) Get() *WorkspacesAdminResponse {
	return v.value
}

func (v *NullableWorkspacesAdminResponse) Set(val *WorkspacesAdminResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkspacesAdminResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkspacesAdminResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkspacesAdminResponse(val *WorkspacesAdminResponse) *NullableWorkspacesAdminResponse {
	return &NullableWorkspacesAdminResponse{value: val, isSet: true}
}

func (v NullableWorkspacesAdminResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkspacesAdminResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


