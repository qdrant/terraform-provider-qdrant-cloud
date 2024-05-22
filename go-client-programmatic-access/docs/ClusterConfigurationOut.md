# ClusterConfigurationOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**NumNodes** | **int32** |  | 
**NumNodesMax** | **int32** |  | 
**ClusterId** | **string** |  | 
**NodeConfiguration** | [**NodeConfiguration**](NodeConfiguration.md) |  | 
**QdrantConfiguration** | Pointer to **map[string]interface{}** |  | [optional] 
**NodeSelector** | Pointer to  |  | [optional] 
**Tolerations** | Pointer to **[]map[string]string** |  | [optional] 
**ClusterAnnotations** | Pointer to **map[string]interface{}** |  | [optional] 
**AllowedIpSourceRanges** | Pointer to **[]string** |  | [optional] 

## Methods

### NewClusterConfigurationOut

`func NewClusterConfigurationOut(id string, createdAt time.Time, numNodes int32, numNodesMax int32, clusterId string, nodeConfiguration NodeConfiguration, ) *ClusterConfigurationOut`

NewClusterConfigurationOut instantiates a new ClusterConfigurationOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClusterConfigurationOutWithDefaults

`func NewClusterConfigurationOutWithDefaults() *ClusterConfigurationOut`

NewClusterConfigurationOutWithDefaults instantiates a new ClusterConfigurationOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ClusterConfigurationOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClusterConfigurationOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClusterConfigurationOut) SetId(v string)`

SetId sets Id field to given value.


### GetCreatedAt

`func (o *ClusterConfigurationOut) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ClusterConfigurationOut) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ClusterConfigurationOut) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetNumNodes

`func (o *ClusterConfigurationOut) GetNumNodes() int32`

GetNumNodes returns the NumNodes field if non-nil, zero value otherwise.

### GetNumNodesOk

`func (o *ClusterConfigurationOut) GetNumNodesOk() (*int32, bool)`

GetNumNodesOk returns a tuple with the NumNodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumNodes

`func (o *ClusterConfigurationOut) SetNumNodes(v int32)`

SetNumNodes sets NumNodes field to given value.


### GetNumNodesMax

`func (o *ClusterConfigurationOut) GetNumNodesMax() int32`

GetNumNodesMax returns the NumNodesMax field if non-nil, zero value otherwise.

### GetNumNodesMaxOk

`func (o *ClusterConfigurationOut) GetNumNodesMaxOk() (*int32, bool)`

GetNumNodesMaxOk returns a tuple with the NumNodesMax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumNodesMax

`func (o *ClusterConfigurationOut) SetNumNodesMax(v int32)`

SetNumNodesMax sets NumNodesMax field to given value.


### GetClusterId

`func (o *ClusterConfigurationOut) GetClusterId() string`

GetClusterId returns the ClusterId field if non-nil, zero value otherwise.

### GetClusterIdOk

`func (o *ClusterConfigurationOut) GetClusterIdOk() (*string, bool)`

GetClusterIdOk returns a tuple with the ClusterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterId

`func (o *ClusterConfigurationOut) SetClusterId(v string)`

SetClusterId sets ClusterId field to given value.


### GetNodeConfiguration

`func (o *ClusterConfigurationOut) GetNodeConfiguration() NodeConfiguration`

GetNodeConfiguration returns the NodeConfiguration field if non-nil, zero value otherwise.

### GetNodeConfigurationOk

`func (o *ClusterConfigurationOut) GetNodeConfigurationOk() (*NodeConfiguration, bool)`

GetNodeConfigurationOk returns a tuple with the NodeConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeConfiguration

`func (o *ClusterConfigurationOut) SetNodeConfiguration(v NodeConfiguration)`

SetNodeConfiguration sets NodeConfiguration field to given value.


### GetQdrantConfiguration

`func (o *ClusterConfigurationOut) GetQdrantConfiguration() map[string]interface{}`

GetQdrantConfiguration returns the QdrantConfiguration field if non-nil, zero value otherwise.

### GetQdrantConfigurationOk

`func (o *ClusterConfigurationOut) GetQdrantConfigurationOk() (*map[string]interface{}, bool)`

GetQdrantConfigurationOk returns a tuple with the QdrantConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQdrantConfiguration

`func (o *ClusterConfigurationOut) SetQdrantConfiguration(v map[string]interface{})`

SetQdrantConfiguration sets QdrantConfiguration field to given value.

### HasQdrantConfiguration

`func (o *ClusterConfigurationOut) HasQdrantConfiguration() bool`

HasQdrantConfiguration returns a boolean if a field has been set.

### SetQdrantConfigurationNil

`func (o *ClusterConfigurationOut) SetQdrantConfigurationNil(b bool)`

 SetQdrantConfigurationNil sets the value for QdrantConfiguration to be an explicit nil

### UnsetQdrantConfiguration
`func (o *ClusterConfigurationOut) UnsetQdrantConfiguration()`

UnsetQdrantConfiguration ensures that no value is present for QdrantConfiguration, not even an explicit nil
### GetNodeSelector

`func (o *ClusterConfigurationOut) GetNodeSelector() map[string]string`

GetNodeSelector returns the NodeSelector field if non-nil, zero value otherwise.

### GetNodeSelectorOk

`func (o *ClusterConfigurationOut) GetNodeSelectorOk() (*map[string]string, bool)`

GetNodeSelectorOk returns a tuple with the NodeSelector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeSelector

`func (o *ClusterConfigurationOut) SetNodeSelector(v map[string]string)`

SetNodeSelector sets NodeSelector field to given value.

### HasNodeSelector

`func (o *ClusterConfigurationOut) HasNodeSelector() bool`

HasNodeSelector returns a boolean if a field has been set.

### SetNodeSelectorNil

`func (o *ClusterConfigurationOut) SetNodeSelectorNil(b bool)`

 SetNodeSelectorNil sets the value for NodeSelector to be an explicit nil

### UnsetNodeSelector
`func (o *ClusterConfigurationOut) UnsetNodeSelector()`

UnsetNodeSelector ensures that no value is present for NodeSelector, not even an explicit nil
### GetTolerations

`func (o *ClusterConfigurationOut) GetTolerations() []map[string]string`

GetTolerations returns the Tolerations field if non-nil, zero value otherwise.

### GetTolerationsOk

`func (o *ClusterConfigurationOut) GetTolerationsOk() (*[]map[string]string, bool)`

GetTolerationsOk returns a tuple with the Tolerations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTolerations

`func (o *ClusterConfigurationOut) SetTolerations(v []map[string]string)`

SetTolerations sets Tolerations field to given value.

### HasTolerations

`func (o *ClusterConfigurationOut) HasTolerations() bool`

HasTolerations returns a boolean if a field has been set.

### SetTolerationsNil

`func (o *ClusterConfigurationOut) SetTolerationsNil(b bool)`

 SetTolerationsNil sets the value for Tolerations to be an explicit nil

### UnsetTolerations
`func (o *ClusterConfigurationOut) UnsetTolerations()`

UnsetTolerations ensures that no value is present for Tolerations, not even an explicit nil
### GetClusterAnnotations

`func (o *ClusterConfigurationOut) GetClusterAnnotations() map[string]interface{}`

GetClusterAnnotations returns the ClusterAnnotations field if non-nil, zero value otherwise.

### GetClusterAnnotationsOk

`func (o *ClusterConfigurationOut) GetClusterAnnotationsOk() (*map[string]interface{}, bool)`

GetClusterAnnotationsOk returns a tuple with the ClusterAnnotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterAnnotations

`func (o *ClusterConfigurationOut) SetClusterAnnotations(v map[string]interface{})`

SetClusterAnnotations sets ClusterAnnotations field to given value.

### HasClusterAnnotations

`func (o *ClusterConfigurationOut) HasClusterAnnotations() bool`

HasClusterAnnotations returns a boolean if a field has been set.

### SetClusterAnnotationsNil

`func (o *ClusterConfigurationOut) SetClusterAnnotationsNil(b bool)`

 SetClusterAnnotationsNil sets the value for ClusterAnnotations to be an explicit nil

### UnsetClusterAnnotations
`func (o *ClusterConfigurationOut) UnsetClusterAnnotations()`

UnsetClusterAnnotations ensures that no value is present for ClusterAnnotations, not even an explicit nil
### GetAllowedIpSourceRanges

`func (o *ClusterConfigurationOut) GetAllowedIpSourceRanges() []string`

GetAllowedIpSourceRanges returns the AllowedIpSourceRanges field if non-nil, zero value otherwise.

### GetAllowedIpSourceRangesOk

`func (o *ClusterConfigurationOut) GetAllowedIpSourceRangesOk() (*[]string, bool)`

GetAllowedIpSourceRangesOk returns a tuple with the AllowedIpSourceRanges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedIpSourceRanges

`func (o *ClusterConfigurationOut) SetAllowedIpSourceRanges(v []string)`

SetAllowedIpSourceRanges sets AllowedIpSourceRanges field to given value.

### HasAllowedIpSourceRanges

`func (o *ClusterConfigurationOut) HasAllowedIpSourceRanges() bool`

HasAllowedIpSourceRanges returns a boolean if a field has been set.

### SetAllowedIpSourceRangesNil

`func (o *ClusterConfigurationOut) SetAllowedIpSourceRangesNil(b bool)`

 SetAllowedIpSourceRangesNil sets the value for AllowedIpSourceRanges to be an explicit nil

### UnsetAllowedIpSourceRanges
`func (o *ClusterConfigurationOut) UnsetAllowedIpSourceRanges()`

UnsetAllowedIpSourceRanges ensures that no value is present for AllowedIpSourceRanges, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


