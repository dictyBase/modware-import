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

// checks if the License type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &License{}

// License struct for License
type License struct {
	Id int32 `json:"id"`
	// Unique identifier of the license.
	LicenseId string `json:"license_id"`
	// Indicates if the backend deems the license valid.
	IsActive bool `json:"is_active"`
	LastCheck NullableTime `json:"last_check,omitempty"`
	// From which timestamp the license becomes active.
	ValidFrom time.Time `json:"valid_from"`
	// Until which timestamp the license is active.
	ValidThrough time.Time `json:"valid_through"`
	// The amount of free users that are currently using the license.
	FreeUsersCount int32 `json:"free_users_count"`
	// The amount of users that are currently using the license.
	SeatsTaken int32 `json:"seats_taken"`
	// The maximum amount of users that can use the license.
	Seats int32 `json:"seats"`
	// The product code that indicates what the license unlocks.
	ProductCode string `json:"product_code"`
	// The date when the license was issued. It could be that a new license is issued with the same `license_id` because it was updated. In that case, the one that has been issued last should be used.
	IssuedOn time.Time `json:"issued_on"`
	// Indicates to which email address the license has been issued.
	IssuedToEmail string `json:"issued_to_email"`
	// Indicates to whom the license has been issued.
	IssuedToName string `json:"issued_to_name"`
}

// NewLicense instantiates a new License object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLicense(id int32, licenseId string, isActive bool, validFrom time.Time, validThrough time.Time, freeUsersCount int32, seatsTaken int32, seats int32, productCode string, issuedOn time.Time, issuedToEmail string, issuedToName string) *License {
	this := License{}
	this.Id = id
	this.LicenseId = licenseId
	this.IsActive = isActive
	this.ValidFrom = validFrom
	this.ValidThrough = validThrough
	this.FreeUsersCount = freeUsersCount
	this.SeatsTaken = seatsTaken
	this.Seats = seats
	this.ProductCode = productCode
	this.IssuedOn = issuedOn
	this.IssuedToEmail = issuedToEmail
	this.IssuedToName = issuedToName
	return &this
}

// NewLicenseWithDefaults instantiates a new License object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLicenseWithDefaults() *License {
	this := License{}
	return &this
}

// GetId returns the Id field value
func (o *License) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *License) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *License) SetId(v int32) {
	o.Id = v
}

// GetLicenseId returns the LicenseId field value
func (o *License) GetLicenseId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LicenseId
}

// GetLicenseIdOk returns a tuple with the LicenseId field value
// and a boolean to check if the value has been set.
func (o *License) GetLicenseIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LicenseId, true
}

// SetLicenseId sets field value
func (o *License) SetLicenseId(v string) {
	o.LicenseId = v
}

// GetIsActive returns the IsActive field value
func (o *License) GetIsActive() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.IsActive
}

// GetIsActiveOk returns a tuple with the IsActive field value
// and a boolean to check if the value has been set.
func (o *License) GetIsActiveOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IsActive, true
}

// SetIsActive sets field value
func (o *License) SetIsActive(v bool) {
	o.IsActive = v
}

// GetLastCheck returns the LastCheck field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *License) GetLastCheck() time.Time {
	if o == nil || IsNil(o.LastCheck.Get()) {
		var ret time.Time
		return ret
	}
	return *o.LastCheck.Get()
}

// GetLastCheckOk returns a tuple with the LastCheck field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *License) GetLastCheckOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.LastCheck.Get(), o.LastCheck.IsSet()
}

// HasLastCheck returns a boolean if a field has been set.
func (o *License) HasLastCheck() bool {
	if o != nil && o.LastCheck.IsSet() {
		return true
	}

	return false
}

// SetLastCheck gets a reference to the given NullableTime and assigns it to the LastCheck field.
func (o *License) SetLastCheck(v time.Time) {
	o.LastCheck.Set(&v)
}
// SetLastCheckNil sets the value for LastCheck to be an explicit nil
func (o *License) SetLastCheckNil() {
	o.LastCheck.Set(nil)
}

// UnsetLastCheck ensures that no value is present for LastCheck, not even an explicit nil
func (o *License) UnsetLastCheck() {
	o.LastCheck.Unset()
}

// GetValidFrom returns the ValidFrom field value
func (o *License) GetValidFrom() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.ValidFrom
}

// GetValidFromOk returns a tuple with the ValidFrom field value
// and a boolean to check if the value has been set.
func (o *License) GetValidFromOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ValidFrom, true
}

// SetValidFrom sets field value
func (o *License) SetValidFrom(v time.Time) {
	o.ValidFrom = v
}

// GetValidThrough returns the ValidThrough field value
func (o *License) GetValidThrough() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.ValidThrough
}

// GetValidThroughOk returns a tuple with the ValidThrough field value
// and a boolean to check if the value has been set.
func (o *License) GetValidThroughOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ValidThrough, true
}

// SetValidThrough sets field value
func (o *License) SetValidThrough(v time.Time) {
	o.ValidThrough = v
}

// GetFreeUsersCount returns the FreeUsersCount field value
func (o *License) GetFreeUsersCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.FreeUsersCount
}

// GetFreeUsersCountOk returns a tuple with the FreeUsersCount field value
// and a boolean to check if the value has been set.
func (o *License) GetFreeUsersCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FreeUsersCount, true
}

// SetFreeUsersCount sets field value
func (o *License) SetFreeUsersCount(v int32) {
	o.FreeUsersCount = v
}

// GetSeatsTaken returns the SeatsTaken field value
func (o *License) GetSeatsTaken() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.SeatsTaken
}

// GetSeatsTakenOk returns a tuple with the SeatsTaken field value
// and a boolean to check if the value has been set.
func (o *License) GetSeatsTakenOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SeatsTaken, true
}

// SetSeatsTaken sets field value
func (o *License) SetSeatsTaken(v int32) {
	o.SeatsTaken = v
}

// GetSeats returns the Seats field value
func (o *License) GetSeats() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Seats
}

// GetSeatsOk returns a tuple with the Seats field value
// and a boolean to check if the value has been set.
func (o *License) GetSeatsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Seats, true
}

// SetSeats sets field value
func (o *License) SetSeats(v int32) {
	o.Seats = v
}

// GetProductCode returns the ProductCode field value
func (o *License) GetProductCode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProductCode
}

// GetProductCodeOk returns a tuple with the ProductCode field value
// and a boolean to check if the value has been set.
func (o *License) GetProductCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProductCode, true
}

// SetProductCode sets field value
func (o *License) SetProductCode(v string) {
	o.ProductCode = v
}

// GetIssuedOn returns the IssuedOn field value
func (o *License) GetIssuedOn() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.IssuedOn
}

// GetIssuedOnOk returns a tuple with the IssuedOn field value
// and a boolean to check if the value has been set.
func (o *License) GetIssuedOnOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IssuedOn, true
}

// SetIssuedOn sets field value
func (o *License) SetIssuedOn(v time.Time) {
	o.IssuedOn = v
}

// GetIssuedToEmail returns the IssuedToEmail field value
func (o *License) GetIssuedToEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IssuedToEmail
}

// GetIssuedToEmailOk returns a tuple with the IssuedToEmail field value
// and a boolean to check if the value has been set.
func (o *License) GetIssuedToEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IssuedToEmail, true
}

// SetIssuedToEmail sets field value
func (o *License) SetIssuedToEmail(v string) {
	o.IssuedToEmail = v
}

// GetIssuedToName returns the IssuedToName field value
func (o *License) GetIssuedToName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IssuedToName
}

// GetIssuedToNameOk returns a tuple with the IssuedToName field value
// and a boolean to check if the value has been set.
func (o *License) GetIssuedToNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IssuedToName, true
}

// SetIssuedToName sets field value
func (o *License) SetIssuedToName(v string) {
	o.IssuedToName = v
}

func (o License) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o License) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	toSerialize["license_id"] = o.LicenseId
	toSerialize["is_active"] = o.IsActive
	if o.LastCheck.IsSet() {
		toSerialize["last_check"] = o.LastCheck.Get()
	}
	toSerialize["valid_from"] = o.ValidFrom
	toSerialize["valid_through"] = o.ValidThrough
	// skip: free_users_count is readOnly
	// skip: seats_taken is readOnly
	toSerialize["seats"] = o.Seats
	toSerialize["product_code"] = o.ProductCode
	toSerialize["issued_on"] = o.IssuedOn
	toSerialize["issued_to_email"] = o.IssuedToEmail
	toSerialize["issued_to_name"] = o.IssuedToName
	return toSerialize, nil
}

type NullableLicense struct {
	value *License
	isSet bool
}

func (v NullableLicense) Get() *License {
	return v.value
}

func (v *NullableLicense) Set(val *License) {
	v.value = val
	v.isSet = true
}

func (v NullableLicense) IsSet() bool {
	return v.isSet
}

func (v *NullableLicense) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLicense(val *License) *NullableLicense {
	return &NullableLicense{value: val, isSet: true}
}

func (v NullableLicense) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLicense) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


