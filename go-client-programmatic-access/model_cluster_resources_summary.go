/*
Qdrant Cloud API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package qdrant_cloud_programmatic_access

import (
	"encoding/json"
)

// checks if the ClusterResourcesSummary type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterResourcesSummary{}

// ClusterResourcesSummary struct for ClusterResourcesSummary
type ClusterResourcesSummary struct {
	Disk NullableClusterResources `json:"disk,omitempty"`
	Ram NullableClusterResources `json:"ram,omitempty"`
	Cpu NullableClusterResources `json:"cpu,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ClusterResourcesSummary ClusterResourcesSummary

// NewClusterResourcesSummary instantiates a new ClusterResourcesSummary object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterResourcesSummary() *ClusterResourcesSummary {
	this := ClusterResourcesSummary{}
	return &this
}

// NewClusterResourcesSummaryWithDefaults instantiates a new ClusterResourcesSummary object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterResourcesSummaryWithDefaults() *ClusterResourcesSummary {
	this := ClusterResourcesSummary{}
	return &this
}

// GetDisk returns the Disk field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterResourcesSummary) GetDisk() ClusterResources {
	if o == nil || IsNil(o.Disk.Get()) {
		var ret ClusterResources
		return ret
	}
	return *o.Disk.Get()
}

// GetDiskOk returns a tuple with the Disk field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterResourcesSummary) GetDiskOk() (*ClusterResources, bool) {
	if o == nil {
		return nil, false
	}
	return o.Disk.Get(), o.Disk.IsSet()
}

// HasDisk returns a boolean if a field has been set.
func (o *ClusterResourcesSummary) HasDisk() bool {
	if o != nil && o.Disk.IsSet() {
		return true
	}

	return false
}

// SetDisk gets a reference to the given NullableClusterResources and assigns it to the Disk field.
func (o *ClusterResourcesSummary) SetDisk(v ClusterResources) {
	o.Disk.Set(&v)
}
// SetDiskNil sets the value for Disk to be an explicit nil
func (o *ClusterResourcesSummary) SetDiskNil() {
	o.Disk.Set(nil)
}

// UnsetDisk ensures that no value is present for Disk, not even an explicit nil
func (o *ClusterResourcesSummary) UnsetDisk() {
	o.Disk.Unset()
}

// GetRam returns the Ram field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterResourcesSummary) GetRam() ClusterResources {
	if o == nil || IsNil(o.Ram.Get()) {
		var ret ClusterResources
		return ret
	}
	return *o.Ram.Get()
}

// GetRamOk returns a tuple with the Ram field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterResourcesSummary) GetRamOk() (*ClusterResources, bool) {
	if o == nil {
		return nil, false
	}
	return o.Ram.Get(), o.Ram.IsSet()
}

// HasRam returns a boolean if a field has been set.
func (o *ClusterResourcesSummary) HasRam() bool {
	if o != nil && o.Ram.IsSet() {
		return true
	}

	return false
}

// SetRam gets a reference to the given NullableClusterResources and assigns it to the Ram field.
func (o *ClusterResourcesSummary) SetRam(v ClusterResources) {
	o.Ram.Set(&v)
}
// SetRamNil sets the value for Ram to be an explicit nil
func (o *ClusterResourcesSummary) SetRamNil() {
	o.Ram.Set(nil)
}

// UnsetRam ensures that no value is present for Ram, not even an explicit nil
func (o *ClusterResourcesSummary) UnsetRam() {
	o.Ram.Unset()
}

// GetCpu returns the Cpu field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterResourcesSummary) GetCpu() ClusterResources {
	if o == nil || IsNil(o.Cpu.Get()) {
		var ret ClusterResources
		return ret
	}
	return *o.Cpu.Get()
}

// GetCpuOk returns a tuple with the Cpu field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterResourcesSummary) GetCpuOk() (*ClusterResources, bool) {
	if o == nil {
		return nil, false
	}
	return o.Cpu.Get(), o.Cpu.IsSet()
}

// HasCpu returns a boolean if a field has been set.
func (o *ClusterResourcesSummary) HasCpu() bool {
	if o != nil && o.Cpu.IsSet() {
		return true
	}

	return false
}

// SetCpu gets a reference to the given NullableClusterResources and assigns it to the Cpu field.
func (o *ClusterResourcesSummary) SetCpu(v ClusterResources) {
	o.Cpu.Set(&v)
}
// SetCpuNil sets the value for Cpu to be an explicit nil
func (o *ClusterResourcesSummary) SetCpuNil() {
	o.Cpu.Set(nil)
}

// UnsetCpu ensures that no value is present for Cpu, not even an explicit nil
func (o *ClusterResourcesSummary) UnsetCpu() {
	o.Cpu.Unset()
}

func (o ClusterResourcesSummary) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClusterResourcesSummary) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Disk.IsSet() {
		toSerialize["disk"] = o.Disk.Get()
	}
	if o.Ram.IsSet() {
		toSerialize["ram"] = o.Ram.Get()
	}
	if o.Cpu.IsSet() {
		toSerialize["cpu"] = o.Cpu.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ClusterResourcesSummary) UnmarshalJSON(data []byte) (err error) {
	varClusterResourcesSummary := _ClusterResourcesSummary{}

	err = json.Unmarshal(data, &varClusterResourcesSummary)

	if err != nil {
		return err
	}

	*o = ClusterResourcesSummary(varClusterResourcesSummary)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "disk")
		delete(additionalProperties, "ram")
		delete(additionalProperties, "cpu")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableClusterResourcesSummary struct {
	value *ClusterResourcesSummary
	isSet bool
}

func (v NullableClusterResourcesSummary) Get() *ClusterResourcesSummary {
	return v.value
}

func (v *NullableClusterResourcesSummary) Set(val *ClusterResourcesSummary) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterResourcesSummary) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterResourcesSummary) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterResourcesSummary(val *ClusterResourcesSummary) *NullableClusterResourcesSummary {
	return &NullableClusterResourcesSummary{value: val, isSet: true}
}

func (v NullableClusterResourcesSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterResourcesSummary) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


