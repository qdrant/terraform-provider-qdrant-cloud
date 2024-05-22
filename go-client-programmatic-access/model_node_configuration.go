/*
Qdrant Cloud API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package qdrant_cloud_programmatic_access

import (
	"encoding/json"
	"fmt"
)

// checks if the NodeConfiguration type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NodeConfiguration{}

// NodeConfiguration struct for NodeConfiguration
type NodeConfiguration struct {
	PackageId string `json:"package_id"`
	Package NullablePackageOut `json:"package,omitempty"`
	ResourceConfigurations []ResourceConfiguration `json:"resource_configurations,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _NodeConfiguration NodeConfiguration

// NewNodeConfiguration instantiates a new NodeConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodeConfiguration(packageId string) *NodeConfiguration {
	this := NodeConfiguration{}
	this.PackageId = packageId
	return &this
}

// NewNodeConfigurationWithDefaults instantiates a new NodeConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodeConfigurationWithDefaults() *NodeConfiguration {
	this := NodeConfiguration{}
	return &this
}

// GetPackageId returns the PackageId field value
func (o *NodeConfiguration) GetPackageId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PackageId
}

// GetPackageIdOk returns a tuple with the PackageId field value
// and a boolean to check if the value has been set.
func (o *NodeConfiguration) GetPackageIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PackageId, true
}

// SetPackageId sets field value
func (o *NodeConfiguration) SetPackageId(v string) {
	o.PackageId = v
}

// GetPackage returns the Package field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *NodeConfiguration) GetPackage() PackageOut {
	if o == nil || IsNil(o.Package.Get()) {
		var ret PackageOut
		return ret
	}
	return *o.Package.Get()
}

// GetPackageOk returns a tuple with the Package field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NodeConfiguration) GetPackageOk() (*PackageOut, bool) {
	if o == nil {
		return nil, false
	}
	return o.Package.Get(), o.Package.IsSet()
}

// HasPackage returns a boolean if a field has been set.
func (o *NodeConfiguration) HasPackage() bool {
	if o != nil && o.Package.IsSet() {
		return true
	}

	return false
}

// SetPackage gets a reference to the given NullablePackageOut and assigns it to the Package field.
func (o *NodeConfiguration) SetPackage(v PackageOut) {
	o.Package.Set(&v)
}
// SetPackageNil sets the value for Package to be an explicit nil
func (o *NodeConfiguration) SetPackageNil() {
	o.Package.Set(nil)
}

// UnsetPackage ensures that no value is present for Package, not even an explicit nil
func (o *NodeConfiguration) UnsetPackage() {
	o.Package.Unset()
}

// GetResourceConfigurations returns the ResourceConfigurations field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *NodeConfiguration) GetResourceConfigurations() []ResourceConfiguration {
	if o == nil {
		var ret []ResourceConfiguration
		return ret
	}
	return o.ResourceConfigurations
}

// GetResourceConfigurationsOk returns a tuple with the ResourceConfigurations field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NodeConfiguration) GetResourceConfigurationsOk() ([]ResourceConfiguration, bool) {
	if o == nil || IsNil(o.ResourceConfigurations) {
		return nil, false
	}
	return o.ResourceConfigurations, true
}

// HasResourceConfigurations returns a boolean if a field has been set.
func (o *NodeConfiguration) HasResourceConfigurations() bool {
	if o != nil && !IsNil(o.ResourceConfigurations) {
		return true
	}

	return false
}

// SetResourceConfigurations gets a reference to the given []ResourceConfiguration and assigns it to the ResourceConfigurations field.
func (o *NodeConfiguration) SetResourceConfigurations(v []ResourceConfiguration) {
	o.ResourceConfigurations = v
}

func (o NodeConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NodeConfiguration) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["package_id"] = o.PackageId
	if o.Package.IsSet() {
		toSerialize["package"] = o.Package.Get()
	}
	if o.ResourceConfigurations != nil {
		toSerialize["resource_configurations"] = o.ResourceConfigurations
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *NodeConfiguration) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"package_id",
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

	varNodeConfiguration := _NodeConfiguration{}

	err = json.Unmarshal(data, &varNodeConfiguration)

	if err != nil {
		return err
	}

	*o = NodeConfiguration(varNodeConfiguration)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "package_id")
		delete(additionalProperties, "package")
		delete(additionalProperties, "resource_configurations")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNodeConfiguration struct {
	value *NodeConfiguration
	isSet bool
}

func (v NullableNodeConfiguration) Get() *NodeConfiguration {
	return v.value
}

func (v *NullableNodeConfiguration) Set(val *NodeConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableNodeConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableNodeConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodeConfiguration(val *NodeConfiguration) *NullableNodeConfiguration {
	return &NullableNodeConfiguration{value: val, isSet: true}
}

func (v NullableNodeConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodeConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

