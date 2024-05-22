# ClusterConfigurationIn

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NumNodes** | **int32** |  | 
**NumNodesMax** | **int32** |  | 
**NodeConfiguration** | [**NodeConfiguration**](NodeConfiguration.md) |  | 
**QdrantConfiguration** | Pointer to **map[string]interface{}** |  | [optional] 
**NodeSelector** | Pointer to  |  | [optional] 
**Tolerations** | Pointer to **[]map[string]string** |  | [optional] 
**ClusterAnnotations** | Pointer to **map[string]interface{}** |  | [optional] 
**AllowedIpSourceRanges** | Pointer to **[]string** |  | [optional] 

## Methods

### NewClusterConfigurationIn

`func NewClusterConfigurationIn(numNodes int32, numNodesMax int32, nodeConfiguration NodeConfiguration, ) *ClusterConfigurationIn`

NewClusterConfigurationIn instantiates a new ClusterConfigurationIn object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClusterConfigurationInWithDefaults

`func NewClusterConfigurationInWithDefaults() *ClusterConfigurationIn`

NewClusterConfigurationInWithDefaults instantiates a new ClusterConfigurationIn object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNumNodes

`func (o *ClusterConfigurationIn) GetNumNodes() int32`

GetNumNodes returns the NumNodes field if non-nil, zero value otherwise.

### GetNumNodesOk

`func (o *ClusterConfigurationIn) GetNumNodesOk() (*int32, bool)`

GetNumNodesOk returns a tuple with the NumNodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumNodes

`func (o *ClusterConfigurationIn) SetNumNodes(v int32)`

SetNumNodes sets NumNodes field to given value.


### GetNumNodesMax

`func (o *ClusterConfigurationIn) GetNumNodesMax() int32`

GetNumNodesMax returns the NumNodesMax field if non-nil, zero value otherwise.

### GetNumNodesMaxOk

`func (o *ClusterConfigurationIn) GetNumNodesMaxOk() (*int32, bool)`

GetNumNodesMaxOk returns a tuple with the NumNodesMax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumNodesMax

`func (o *ClusterConfigurationIn) SetNumNodesMax(v int32)`

SetNumNodesMax sets NumNodesMax field to given value.


### GetNodeConfiguration

`func (o *ClusterConfigurationIn) GetNodeConfiguration() NodeConfiguration`

GetNodeConfiguration returns the NodeConfiguration field if non-nil, zero value otherwise.

### GetNodeConfigurationOk

`func (o *ClusterConfigurationIn) GetNodeConfigurationOk() (*NodeConfiguration, bool)`

GetNodeConfigurationOk returns a tuple with the NodeConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeConfiguration

`func (o *ClusterConfigurationIn) SetNodeConfiguration(v NodeConfiguration)`

SetNodeConfiguration sets NodeConfiguration field to given value.


### GetQdrantConfiguration

`func (o *ClusterConfigurationIn) GetQdrantConfiguration() map[string]interface{}`

GetQdrantConfiguration returns the QdrantConfiguration field if non-nil, zero value otherwise.

### GetQdrantConfigurationOk

`func (o *ClusterConfigurationIn) GetQdrantConfigurationOk() (*map[string]interface{}, bool)`

GetQdrantConfigurationOk returns a tuple with the QdrantConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQdrantConfiguration

`func (o *ClusterConfigurationIn) SetQdrantConfiguration(v map[string]interface{})`

SetQdrantConfiguration sets QdrantConfiguration field to given value.

### HasQdrantConfiguration

`func (o *ClusterConfigurationIn) HasQdrantConfiguration() bool`

HasQdrantConfiguration returns a boolean if a field has been set.

### SetQdrantConfigurationNil

`func (o *ClusterConfigurationIn) SetQdrantConfigurationNil(b bool)`

 SetQdrantConfigurationNil sets the value for QdrantConfiguration to be an explicit nil

### UnsetQdrantConfiguration
`func (o *ClusterConfigurationIn) UnsetQdrantConfiguration()`

UnsetQdrantConfiguration ensures that no value is present for QdrantConfiguration, not even an explicit nil
### GetNodeSelector

`func (o *ClusterConfigurationIn) GetNodeSelector() map[string]string`

GetNodeSelector returns the NodeSelector field if non-nil, zero value otherwise.

### GetNodeSelectorOk

`func (o *ClusterConfigurationIn) GetNodeSelectorOk() (*map[string]string, bool)`

GetNodeSelectorOk returns a tuple with the NodeSelector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeSelector

`func (o *ClusterConfigurationIn) SetNodeSelector(v map[string]string)`

SetNodeSelector sets NodeSelector field to given value.

### HasNodeSelector

`func (o *ClusterConfigurationIn) HasNodeSelector() bool`

HasNodeSelector returns a boolean if a field has been set.

### SetNodeSelectorNil

`func (o *ClusterConfigurationIn) SetNodeSelectorNil(b bool)`

 SetNodeSelectorNil sets the value for NodeSelector to be an explicit nil

### UnsetNodeSelector
`func (o *ClusterConfigurationIn) UnsetNodeSelector()`

UnsetNodeSelector ensures that no value is present for NodeSelector, not even an explicit nil
### GetTolerations

`func (o *ClusterConfigurationIn) GetTolerations() []map[string]string`

GetTolerations returns the Tolerations field if non-nil, zero value otherwise.

### GetTolerationsOk

`func (o *ClusterConfigurationIn) GetTolerationsOk() (*[]map[string]string, bool)`

GetTolerationsOk returns a tuple with the Tolerations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTolerations

`func (o *ClusterConfigurationIn) SetTolerations(v []map[string]string)`

SetTolerations sets Tolerations field to given value.

### HasTolerations

`func (o *ClusterConfigurationIn) HasTolerations() bool`

HasTolerations returns a boolean if a field has been set.

### SetTolerationsNil

`func (o *ClusterConfigurationIn) SetTolerationsNil(b bool)`

 SetTolerationsNil sets the value for Tolerations to be an explicit nil

### UnsetTolerations
`func (o *ClusterConfigurationIn) UnsetTolerations()`

UnsetTolerations ensures that no value is present for Tolerations, not even an explicit nil
### GetClusterAnnotations

`func (o *ClusterConfigurationIn) GetClusterAnnotations() map[string]interface{}`

GetClusterAnnotations returns the ClusterAnnotations field if non-nil, zero value otherwise.

### GetClusterAnnotationsOk

`func (o *ClusterConfigurationIn) GetClusterAnnotationsOk() (*map[string]interface{}, bool)`

GetClusterAnnotationsOk returns a tuple with the ClusterAnnotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterAnnotations

`func (o *ClusterConfigurationIn) SetClusterAnnotations(v map[string]interface{})`

SetClusterAnnotations sets ClusterAnnotations field to given value.

### HasClusterAnnotations

`func (o *ClusterConfigurationIn) HasClusterAnnotations() bool`

HasClusterAnnotations returns a boolean if a field has been set.

### SetClusterAnnotationsNil

`func (o *ClusterConfigurationIn) SetClusterAnnotationsNil(b bool)`

 SetClusterAnnotationsNil sets the value for ClusterAnnotations to be an explicit nil

### UnsetClusterAnnotations
`func (o *ClusterConfigurationIn) UnsetClusterAnnotations()`

UnsetClusterAnnotations ensures that no value is present for ClusterAnnotations, not even an explicit nil
### GetAllowedIpSourceRanges

`func (o *ClusterConfigurationIn) GetAllowedIpSourceRanges() []string`

GetAllowedIpSourceRanges returns the AllowedIpSourceRanges field if non-nil, zero value otherwise.

### GetAllowedIpSourceRangesOk

`func (o *ClusterConfigurationIn) GetAllowedIpSourceRangesOk() (*[]string, bool)`

GetAllowedIpSourceRangesOk returns a tuple with the AllowedIpSourceRanges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedIpSourceRanges

`func (o *ClusterConfigurationIn) SetAllowedIpSourceRanges(v []string)`

SetAllowedIpSourceRanges sets AllowedIpSourceRanges field to given value.

### HasAllowedIpSourceRanges

`func (o *ClusterConfigurationIn) HasAllowedIpSourceRanges() bool`

HasAllowedIpSourceRanges returns a boolean if a field has been set.

### SetAllowedIpSourceRangesNil

`func (o *ClusterConfigurationIn) SetAllowedIpSourceRangesNil(b bool)`

 SetAllowedIpSourceRangesNil sets the value for AllowedIpSourceRanges to be an explicit nil

### UnsetAllowedIpSourceRanges
`func (o *ClusterConfigurationIn) UnsetAllowedIpSourceRanges()`

UnsetAllowedIpSourceRanges ensures that no value is present for AllowedIpSourceRanges, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


