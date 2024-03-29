/*
Baserow API spec

For more information about our REST API, please visit [this page](https://baserow.io/docs/apis%2Frest-api).  For more information about our deprecation policy, please visit [this page](https://baserow.io/docs/apis%2Fdeprecations).

API version: 1.18.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the GitHubAuthProviderModelAuthProvider type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GitHubAuthProviderModelAuthProvider{}

// GitHubAuthProviderModelAuthProvider struct for GitHubAuthProviderModelAuthProvider
type GitHubAuthProviderModelAuthProvider struct {
	Id int32 `json:"id"`
	// The type of the related field.
	Type string `json:"type"`
	Domain NullableString `json:"domain,omitempty"`
	Enabled *bool `json:"enabled,omitempty"`
	Name string `json:"name"`
	// App ID, or consumer key
	ClientId string `json:"client_id"`
	// API secret, client secret, or consumer secret
	Secret string `json:"secret"`
}

// NewGitHubAuthProviderModelAuthProvider instantiates a new GitHubAuthProviderModelAuthProvider object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGitHubAuthProviderModelAuthProvider(id int32, type_ string, name string, clientId string, secret string) *GitHubAuthProviderModelAuthProvider {
	this := GitHubAuthProviderModelAuthProvider{}
	this.Id = id
	this.Type = type_
	this.Name = name
	this.ClientId = clientId
	this.Secret = secret
	return &this
}

// NewGitHubAuthProviderModelAuthProviderWithDefaults instantiates a new GitHubAuthProviderModelAuthProvider object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGitHubAuthProviderModelAuthProviderWithDefaults() *GitHubAuthProviderModelAuthProvider {
	this := GitHubAuthProviderModelAuthProvider{}
	return &this
}

// GetId returns the Id field value
func (o *GitHubAuthProviderModelAuthProvider) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GitHubAuthProviderModelAuthProvider) SetId(v int32) {
	o.Id = v
}

// GetType returns the Type field value
func (o *GitHubAuthProviderModelAuthProvider) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GitHubAuthProviderModelAuthProvider) SetType(v string) {
	o.Type = v
}

// GetDomain returns the Domain field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *GitHubAuthProviderModelAuthProvider) GetDomain() string {
	if o == nil || IsNil(o.Domain.Get()) {
		var ret string
		return ret
	}
	return *o.Domain.Get()
}

// GetDomainOk returns a tuple with the Domain field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GitHubAuthProviderModelAuthProvider) GetDomainOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Domain.Get(), o.Domain.IsSet()
}

// HasDomain returns a boolean if a field has been set.
func (o *GitHubAuthProviderModelAuthProvider) HasDomain() bool {
	if o != nil && o.Domain.IsSet() {
		return true
	}

	return false
}

// SetDomain gets a reference to the given NullableString and assigns it to the Domain field.
func (o *GitHubAuthProviderModelAuthProvider) SetDomain(v string) {
	o.Domain.Set(&v)
}
// SetDomainNil sets the value for Domain to be an explicit nil
func (o *GitHubAuthProviderModelAuthProvider) SetDomainNil() {
	o.Domain.Set(nil)
}

// UnsetDomain ensures that no value is present for Domain, not even an explicit nil
func (o *GitHubAuthProviderModelAuthProvider) UnsetDomain() {
	o.Domain.Unset()
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *GitHubAuthProviderModelAuthProvider) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *GitHubAuthProviderModelAuthProvider) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *GitHubAuthProviderModelAuthProvider) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetName returns the Name field value
func (o *GitHubAuthProviderModelAuthProvider) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *GitHubAuthProviderModelAuthProvider) SetName(v string) {
	o.Name = v
}

// GetClientId returns the ClientId field value
func (o *GitHubAuthProviderModelAuthProvider) GetClientId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetClientIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClientId, true
}

// SetClientId sets field value
func (o *GitHubAuthProviderModelAuthProvider) SetClientId(v string) {
	o.ClientId = v
}

// GetSecret returns the Secret field value
func (o *GitHubAuthProviderModelAuthProvider) GetSecret() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Secret
}

// GetSecretOk returns a tuple with the Secret field value
// and a boolean to check if the value has been set.
func (o *GitHubAuthProviderModelAuthProvider) GetSecretOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Secret, true
}

// SetSecret sets field value
func (o *GitHubAuthProviderModelAuthProvider) SetSecret(v string) {
	o.Secret = v
}

func (o GitHubAuthProviderModelAuthProvider) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GitHubAuthProviderModelAuthProvider) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	// skip: type is readOnly
	if o.Domain.IsSet() {
		toSerialize["domain"] = o.Domain.Get()
	}
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	toSerialize["name"] = o.Name
	toSerialize["client_id"] = o.ClientId
	toSerialize["secret"] = o.Secret
	return toSerialize, nil
}

type NullableGitHubAuthProviderModelAuthProvider struct {
	value *GitHubAuthProviderModelAuthProvider
	isSet bool
}

func (v NullableGitHubAuthProviderModelAuthProvider) Get() *GitHubAuthProviderModelAuthProvider {
	return v.value
}

func (v *NullableGitHubAuthProviderModelAuthProvider) Set(val *GitHubAuthProviderModelAuthProvider) {
	v.value = val
	v.isSet = true
}

func (v NullableGitHubAuthProviderModelAuthProvider) IsSet() bool {
	return v.isSet
}

func (v *NullableGitHubAuthProviderModelAuthProvider) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGitHubAuthProviderModelAuthProvider(val *GitHubAuthProviderModelAuthProvider) *NullableGitHubAuthProviderModelAuthProvider {
	return &NullableGitHubAuthProviderModelAuthProvider{value: val, isSet: true}
}

func (v NullableGitHubAuthProviderModelAuthProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGitHubAuthProviderModelAuthProvider) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


