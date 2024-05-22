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

// checks if the QdrantClusterNodeState type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &QdrantClusterNodeState{}

// QdrantClusterNodeState struct for QdrantClusterNodeState
type QdrantClusterNodeState struct {
	Name string `json:"name"`
	StartedAt NullableTime `json:"started_at,omitempty"`
	State map[string]string `json:"state,omitempty"`
	Version *string `json:"version,omitempty"`
	Endpoint *string `json:"endpoint,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _QdrantClusterNodeState QdrantClusterNodeState

// NewQdrantClusterNodeState instantiates a new QdrantClusterNodeState object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQdrantClusterNodeState(name string) *QdrantClusterNodeState {
	this := QdrantClusterNodeState{}
	this.Name = name
	var version string = ""
	this.Version = &version
	var endpoint string = ""
	this.Endpoint = &endpoint
	return &this
}

// NewQdrantClusterNodeStateWithDefaults instantiates a new QdrantClusterNodeState object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQdrantClusterNodeStateWithDefaults() *QdrantClusterNodeState {
	this := QdrantClusterNodeState{}
	var version string = ""
	this.Version = &version
	var endpoint string = ""
	this.Endpoint = &endpoint
	return &this
}

// GetName returns the Name field value
func (o *QdrantClusterNodeState) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *QdrantClusterNodeState) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *QdrantClusterNodeState) SetName(v string) {
	o.Name = v
}

// GetStartedAt returns the StartedAt field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *QdrantClusterNodeState) GetStartedAt() time.Time {
	if o == nil || IsNil(o.StartedAt.Get()) {
		var ret time.Time
		return ret
	}
	return *o.StartedAt.Get()
}

// GetStartedAtOk returns a tuple with the StartedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *QdrantClusterNodeState) GetStartedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.StartedAt.Get(), o.StartedAt.IsSet()
}

// HasStartedAt returns a boolean if a field has been set.
func (o *QdrantClusterNodeState) HasStartedAt() bool {
	if o != nil && o.StartedAt.IsSet() {
		return true
	}

	return false
}

// SetStartedAt gets a reference to the given NullableTime and assigns it to the StartedAt field.
func (o *QdrantClusterNodeState) SetStartedAt(v time.Time) {
	o.StartedAt.Set(&v)
}
// SetStartedAtNil sets the value for StartedAt to be an explicit nil
func (o *QdrantClusterNodeState) SetStartedAtNil() {
	o.StartedAt.Set(nil)
}

// UnsetStartedAt ensures that no value is present for StartedAt, not even an explicit nil
func (o *QdrantClusterNodeState) UnsetStartedAt() {
	o.StartedAt.Unset()
}

// GetState returns the State field value if set, zero value otherwise.
func (o *QdrantClusterNodeState) GetState() map[string]string {
	if o == nil || IsNil(o.State) {
		var ret map[string]string
		return ret
	}
	return o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QdrantClusterNodeState) GetStateOk() (map[string]string, bool) {
	if o == nil || IsNil(o.State) {
		return map[string]string{}, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *QdrantClusterNodeState) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given map[string]string and assigns it to the State field.
func (o *QdrantClusterNodeState) SetState(v map[string]string) {
	o.State = v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *QdrantClusterNodeState) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QdrantClusterNodeState) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *QdrantClusterNodeState) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *QdrantClusterNodeState) SetVersion(v string) {
	o.Version = &v
}

// GetEndpoint returns the Endpoint field value if set, zero value otherwise.
func (o *QdrantClusterNodeState) GetEndpoint() string {
	if o == nil || IsNil(o.Endpoint) {
		var ret string
		return ret
	}
	return *o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QdrantClusterNodeState) GetEndpointOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint) {
		return nil, false
	}
	return o.Endpoint, true
}

// HasEndpoint returns a boolean if a field has been set.
func (o *QdrantClusterNodeState) HasEndpoint() bool {
	if o != nil && !IsNil(o.Endpoint) {
		return true
	}

	return false
}

// SetEndpoint gets a reference to the given string and assigns it to the Endpoint field.
func (o *QdrantClusterNodeState) SetEndpoint(v string) {
	o.Endpoint = &v
}

func (o QdrantClusterNodeState) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o QdrantClusterNodeState) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if o.StartedAt.IsSet() {
		toSerialize["started_at"] = o.StartedAt.Get()
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.Endpoint) {
		toSerialize["endpoint"] = o.Endpoint
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *QdrantClusterNodeState) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
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

	varQdrantClusterNodeState := _QdrantClusterNodeState{}

	err = json.Unmarshal(data, &varQdrantClusterNodeState)

	if err != nil {
		return err
	}

	*o = QdrantClusterNodeState(varQdrantClusterNodeState)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "started_at")
		delete(additionalProperties, "state")
		delete(additionalProperties, "version")
		delete(additionalProperties, "endpoint")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableQdrantClusterNodeState struct {
	value *QdrantClusterNodeState
	isSet bool
}

func (v NullableQdrantClusterNodeState) Get() *QdrantClusterNodeState {
	return v.value
}

func (v *NullableQdrantClusterNodeState) Set(val *QdrantClusterNodeState) {
	v.value = val
	v.isSet = true
}

func (v NullableQdrantClusterNodeState) IsSet() bool {
	return v.isSet
}

func (v *NullableQdrantClusterNodeState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQdrantClusterNodeState(val *QdrantClusterNodeState) *NullableQdrantClusterNodeState {
	return &NullableQdrantClusterNodeState{value: val, isSet: true}
}

func (v NullableQdrantClusterNodeState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQdrantClusterNodeState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


