/*
Qdrant Cloud API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package qdrant_cloud_programmatic_access

import (
	"encoding/json"
	"time"
	"fmt"
)

// checks if the CreateApiKeyOut type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateApiKeyOut{}

// CreateApiKeyOut struct for CreateApiKeyOut
type CreateApiKeyOut struct {
	Id string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserId NullableString `json:"user_id,omitempty"`
	Prefix string `json:"prefix"`
	ClusterIdList []interface{} `json:"cluster_id_list,omitempty"`
	AccountId NullableString `json:"account_id,omitempty"`
	Token string `json:"token"`
	AdditionalProperties map[string]interface{}
}

type _CreateApiKeyOut CreateApiKeyOut

// NewCreateApiKeyOut instantiates a new CreateApiKeyOut object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateApiKeyOut(id string, createdAt time.Time, prefix string, token string) *CreateApiKeyOut {
	this := CreateApiKeyOut{}
	this.Id = id
	this.CreatedAt = createdAt
	this.Prefix = prefix
	this.Token = token
	return &this
}

// NewCreateApiKeyOutWithDefaults instantiates a new CreateApiKeyOut object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateApiKeyOutWithDefaults() *CreateApiKeyOut {
	this := CreateApiKeyOut{}
	return &this
}

// GetId returns the Id field value
func (o *CreateApiKeyOut) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *CreateApiKeyOut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *CreateApiKeyOut) SetId(v string) {
	o.Id = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *CreateApiKeyOut) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *CreateApiKeyOut) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *CreateApiKeyOut) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetUserId returns the UserId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CreateApiKeyOut) GetUserId() string {
	if o == nil || IsNil(o.UserId.Get()) {
		var ret string
		return ret
	}
	return *o.UserId.Get()
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateApiKeyOut) GetUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.UserId.Get(), o.UserId.IsSet()
}

// HasUserId returns a boolean if a field has been set.
func (o *CreateApiKeyOut) HasUserId() bool {
	if o != nil && o.UserId.IsSet() {
		return true
	}

	return false
}

// SetUserId gets a reference to the given NullableString and assigns it to the UserId field.
func (o *CreateApiKeyOut) SetUserId(v string) {
	o.UserId.Set(&v)
}
// SetUserIdNil sets the value for UserId to be an explicit nil
func (o *CreateApiKeyOut) SetUserIdNil() {
	o.UserId.Set(nil)
}

// UnsetUserId ensures that no value is present for UserId, not even an explicit nil
func (o *CreateApiKeyOut) UnsetUserId() {
	o.UserId.Unset()
}

// GetPrefix returns the Prefix field value
func (o *CreateApiKeyOut) GetPrefix() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value
// and a boolean to check if the value has been set.
func (o *CreateApiKeyOut) GetPrefixOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Prefix, true
}

// SetPrefix sets field value
func (o *CreateApiKeyOut) SetPrefix(v string) {
	o.Prefix = v
}

// GetClusterIdList returns the ClusterIdList field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CreateApiKeyOut) GetClusterIdList() []interface{} {
	if o == nil {
		var ret []interface{}
		return ret
	}
	return o.ClusterIdList
}

// GetClusterIdListOk returns a tuple with the ClusterIdList field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateApiKeyOut) GetClusterIdListOk() ([]interface{}, bool) {
	if o == nil || IsNil(o.ClusterIdList) {
		return nil, false
	}
	return o.ClusterIdList, true
}

// HasClusterIdList returns a boolean if a field has been set.
func (o *CreateApiKeyOut) HasClusterIdList() bool {
	if o != nil && !IsNil(o.ClusterIdList) {
		return true
	}

	return false
}

// SetClusterIdList gets a reference to the given []interface{} and assigns it to the ClusterIdList field.
func (o *CreateApiKeyOut) SetClusterIdList(v []interface{}) {
	o.ClusterIdList = v
}

// GetAccountId returns the AccountId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *CreateApiKeyOut) GetAccountId() string {
	if o == nil || IsNil(o.AccountId.Get()) {
		var ret string
		return ret
	}
	return *o.AccountId.Get()
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateApiKeyOut) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AccountId.Get(), o.AccountId.IsSet()
}

// HasAccountId returns a boolean if a field has been set.
func (o *CreateApiKeyOut) HasAccountId() bool {
	if o != nil && o.AccountId.IsSet() {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given NullableString and assigns it to the AccountId field.
func (o *CreateApiKeyOut) SetAccountId(v string) {
	o.AccountId.Set(&v)
}
// SetAccountIdNil sets the value for AccountId to be an explicit nil
func (o *CreateApiKeyOut) SetAccountIdNil() {
	o.AccountId.Set(nil)
}

// UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
func (o *CreateApiKeyOut) UnsetAccountId() {
	o.AccountId.Unset()
}

// GetToken returns the Token field value
func (o *CreateApiKeyOut) GetToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Token
}

// GetTokenOk returns a tuple with the Token field value
// and a boolean to check if the value has been set.
func (o *CreateApiKeyOut) GetTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Token, true
}

// SetToken sets field value
func (o *CreateApiKeyOut) SetToken(v string) {
	o.Token = v
}

func (o CreateApiKeyOut) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateApiKeyOut) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["created_at"] = o.CreatedAt
	if o.UserId.IsSet() {
		toSerialize["user_id"] = o.UserId.Get()
	}
	toSerialize["prefix"] = o.Prefix
	if o.ClusterIdList != nil {
		toSerialize["cluster_id_list"] = o.ClusterIdList
	}
	if o.AccountId.IsSet() {
		toSerialize["account_id"] = o.AccountId.Get()
	}
	toSerialize["token"] = o.Token

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CreateApiKeyOut) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"created_at",
		"prefix",
		"token",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varCreateApiKeyOut := _CreateApiKeyOut{}

	err = json.Unmarshal(data, &varCreateApiKeyOut)

	if err != nil {
		return err
	}

	*o = CreateApiKeyOut(varCreateApiKeyOut)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "user_id")
		delete(additionalProperties, "prefix")
		delete(additionalProperties, "cluster_id_list")
		delete(additionalProperties, "account_id")
		delete(additionalProperties, "token")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCreateApiKeyOut struct {
	value *CreateApiKeyOut
	isSet bool
}

func (v NullableCreateApiKeyOut) Get() *CreateApiKeyOut {
	return v.value
}

func (v *NullableCreateApiKeyOut) Set(val *CreateApiKeyOut) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateApiKeyOut) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateApiKeyOut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateApiKeyOut(val *CreateApiKeyOut) *NullableCreateApiKeyOut {
	return &NullableCreateApiKeyOut{value: val, isSet: true}
}

func (v NullableCreateApiKeyOut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateApiKeyOut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

