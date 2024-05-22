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

// checks if the ClusterConfigurationOut type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterConfigurationOut{}

// ClusterConfigurationOut struct for ClusterConfigurationOut
type ClusterConfigurationOut struct {
	Id string `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	NumNodes int32 `json:"num_nodes"`
	NumNodesMax int32 `json:"num_nodes_max"`
	ClusterId string `json:"cluster_id"`
	NodeConfiguration NodeConfiguration `json:"node_configuration"`
	QdrantConfiguration map[string]interface{} `json:"qdrant_configuration,omitempty"`
	NodeSelector map[string]string `json:"node_selector,omitempty"`
	Tolerations []map[string]string `json:"tolerations,omitempty"`
	ClusterAnnotations map[string]interface{} `json:"cluster_annotations,omitempty"`
	AllowedIpSourceRanges []string `json:"allowed_ip_source_ranges,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ClusterConfigurationOut ClusterConfigurationOut

// NewClusterConfigurationOut instantiates a new ClusterConfigurationOut object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterConfigurationOut(id string, createdAt time.Time, numNodes int32, numNodesMax int32, clusterId string, nodeConfiguration NodeConfiguration) *ClusterConfigurationOut {
	this := ClusterConfigurationOut{}
	this.Id = id
	this.CreatedAt = createdAt
	this.NumNodes = numNodes
	this.NumNodesMax = numNodesMax
	this.ClusterId = clusterId
	this.NodeConfiguration = nodeConfiguration
	return &this
}

// NewClusterConfigurationOutWithDefaults instantiates a new ClusterConfigurationOut object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterConfigurationOutWithDefaults() *ClusterConfigurationOut {
	this := ClusterConfigurationOut{}
	return &this
}

// GetId returns the Id field value
func (o *ClusterConfigurationOut) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *ClusterConfigurationOut) SetId(v string) {
	o.Id = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *ClusterConfigurationOut) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *ClusterConfigurationOut) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetNumNodes returns the NumNodes field value
func (o *ClusterConfigurationOut) GetNumNodes() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NumNodes
}

// GetNumNodesOk returns a tuple with the NumNodes field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetNumNodesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumNodes, true
}

// SetNumNodes sets field value
func (o *ClusterConfigurationOut) SetNumNodes(v int32) {
	o.NumNodes = v
}

// GetNumNodesMax returns the NumNodesMax field value
func (o *ClusterConfigurationOut) GetNumNodesMax() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NumNodesMax
}

// GetNumNodesMaxOk returns a tuple with the NumNodesMax field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetNumNodesMaxOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumNodesMax, true
}

// SetNumNodesMax sets field value
func (o *ClusterConfigurationOut) SetNumNodesMax(v int32) {
	o.NumNodesMax = v
}

// GetClusterId returns the ClusterId field value
func (o *ClusterConfigurationOut) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *ClusterConfigurationOut) SetClusterId(v string) {
	o.ClusterId = v
}

// GetNodeConfiguration returns the NodeConfiguration field value
func (o *ClusterConfigurationOut) GetNodeConfiguration() NodeConfiguration {
	if o == nil {
		var ret NodeConfiguration
		return ret
	}

	return o.NodeConfiguration
}

// GetNodeConfigurationOk returns a tuple with the NodeConfiguration field value
// and a boolean to check if the value has been set.
func (o *ClusterConfigurationOut) GetNodeConfigurationOk() (*NodeConfiguration, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeConfiguration, true
}

// SetNodeConfiguration sets field value
func (o *ClusterConfigurationOut) SetNodeConfiguration(v NodeConfiguration) {
	o.NodeConfiguration = v
}

// GetQdrantConfiguration returns the QdrantConfiguration field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterConfigurationOut) GetQdrantConfiguration() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.QdrantConfiguration
}

// GetQdrantConfigurationOk returns a tuple with the QdrantConfiguration field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterConfigurationOut) GetQdrantConfigurationOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.QdrantConfiguration) {
		return map[string]interface{}{}, false
	}
	return o.QdrantConfiguration, true
}

// HasQdrantConfiguration returns a boolean if a field has been set.
func (o *ClusterConfigurationOut) HasQdrantConfiguration() bool {
	if o != nil && !IsNil(o.QdrantConfiguration) {
		return true
	}

	return false
}

// SetQdrantConfiguration gets a reference to the given map[string]interface{} and assigns it to the QdrantConfiguration field.
func (o *ClusterConfigurationOut) SetQdrantConfiguration(v map[string]interface{}) {
	o.QdrantConfiguration = v
}

// GetNodeSelector returns the NodeSelector field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterConfigurationOut) GetNodeSelector() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}
	return o.NodeSelector
}

// GetNodeSelectorOk returns a tuple with the NodeSelector field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterConfigurationOut) GetNodeSelectorOk() (map[string]string, bool) {
	if o == nil || IsNil(o.NodeSelector) {
		return map[string]string{}, false
	}
	return o.NodeSelector, true
}

// HasNodeSelector returns a boolean if a field has been set.
func (o *ClusterConfigurationOut) HasNodeSelector() bool {
	if o != nil && !IsNil(o.NodeSelector) {
		return true
	}

	return false
}

// SetNodeSelector gets a reference to the given map[string]string and assigns it to the NodeSelector field.
func (o *ClusterConfigurationOut) SetNodeSelector(v map[string]string) {
	o.NodeSelector = v
}

// GetTolerations returns the Tolerations field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterConfigurationOut) GetTolerations() []map[string]string {
	if o == nil {
		var ret []map[string]string
		return ret
	}
	return o.Tolerations
}

// GetTolerationsOk returns a tuple with the Tolerations field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterConfigurationOut) GetTolerationsOk() ([]map[string]string, bool) {
	if o == nil || IsNil(o.Tolerations) {
		return nil, false
	}
	return o.Tolerations, true
}

// HasTolerations returns a boolean if a field has been set.
func (o *ClusterConfigurationOut) HasTolerations() bool {
	if o != nil && !IsNil(o.Tolerations) {
		return true
	}

	return false
}

// SetTolerations gets a reference to the given []map[string]string and assigns it to the Tolerations field.
func (o *ClusterConfigurationOut) SetTolerations(v []map[string]string) {
	o.Tolerations = v
}

// GetClusterAnnotations returns the ClusterAnnotations field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterConfigurationOut) GetClusterAnnotations() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.ClusterAnnotations
}

// GetClusterAnnotationsOk returns a tuple with the ClusterAnnotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterConfigurationOut) GetClusterAnnotationsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.ClusterAnnotations) {
		return map[string]interface{}{}, false
	}
	return o.ClusterAnnotations, true
}

// HasClusterAnnotations returns a boolean if a field has been set.
func (o *ClusterConfigurationOut) HasClusterAnnotations() bool {
	if o != nil && !IsNil(o.ClusterAnnotations) {
		return true
	}

	return false
}

// SetClusterAnnotations gets a reference to the given map[string]interface{} and assigns it to the ClusterAnnotations field.
func (o *ClusterConfigurationOut) SetClusterAnnotations(v map[string]interface{}) {
	o.ClusterAnnotations = v
}

// GetAllowedIpSourceRanges returns the AllowedIpSourceRanges field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ClusterConfigurationOut) GetAllowedIpSourceRanges() []string {
	if o == nil {
		var ret []string
		return ret
	}
	return o.AllowedIpSourceRanges
}

// GetAllowedIpSourceRangesOk returns a tuple with the AllowedIpSourceRanges field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterConfigurationOut) GetAllowedIpSourceRangesOk() ([]string, bool) {
	if o == nil || IsNil(o.AllowedIpSourceRanges) {
		return nil, false
	}
	return o.AllowedIpSourceRanges, true
}

// HasAllowedIpSourceRanges returns a boolean if a field has been set.
func (o *ClusterConfigurationOut) HasAllowedIpSourceRanges() bool {
	if o != nil && !IsNil(o.AllowedIpSourceRanges) {
		return true
	}

	return false
}

// SetAllowedIpSourceRanges gets a reference to the given []string and assigns it to the AllowedIpSourceRanges field.
func (o *ClusterConfigurationOut) SetAllowedIpSourceRanges(v []string) {
	o.AllowedIpSourceRanges = v
}

func (o ClusterConfigurationOut) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClusterConfigurationOut) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["created_at"] = o.CreatedAt
	toSerialize["num_nodes"] = o.NumNodes
	toSerialize["num_nodes_max"] = o.NumNodesMax
	toSerialize["cluster_id"] = o.ClusterId
	toSerialize["node_configuration"] = o.NodeConfiguration
	if o.QdrantConfiguration != nil {
		toSerialize["qdrant_configuration"] = o.QdrantConfiguration
	}
	if o.NodeSelector != nil {
		toSerialize["node_selector"] = o.NodeSelector
	}
	if o.Tolerations != nil {
		toSerialize["tolerations"] = o.Tolerations
	}
	if o.ClusterAnnotations != nil {
		toSerialize["cluster_annotations"] = o.ClusterAnnotations
	}
	if o.AllowedIpSourceRanges != nil {
		toSerialize["allowed_ip_source_ranges"] = o.AllowedIpSourceRanges
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ClusterConfigurationOut) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"created_at",
		"num_nodes",
		"num_nodes_max",
		"cluster_id",
		"node_configuration",
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

	varClusterConfigurationOut := _ClusterConfigurationOut{}

	err = json.Unmarshal(data, &varClusterConfigurationOut)

	if err != nil {
		return err
	}

	*o = ClusterConfigurationOut(varClusterConfigurationOut)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "num_nodes")
		delete(additionalProperties, "num_nodes_max")
		delete(additionalProperties, "cluster_id")
		delete(additionalProperties, "node_configuration")
		delete(additionalProperties, "qdrant_configuration")
		delete(additionalProperties, "node_selector")
		delete(additionalProperties, "tolerations")
		delete(additionalProperties, "cluster_annotations")
		delete(additionalProperties, "allowed_ip_source_ranges")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableClusterConfigurationOut struct {
	value *ClusterConfigurationOut
	isSet bool
}

func (v NullableClusterConfigurationOut) Get() *ClusterConfigurationOut {
	return v.value
}

func (v *NullableClusterConfigurationOut) Set(val *ClusterConfigurationOut) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterConfigurationOut) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterConfigurationOut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterConfigurationOut(val *ClusterConfigurationOut) *NullableClusterConfigurationOut {
	return &NullableClusterConfigurationOut{value: val, isSet: true}
}

func (v NullableClusterConfigurationOut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterConfigurationOut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

