# QdrantClusterStateNodesValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**StartedAt** | Pointer to **time.Time** |  | [optional] 
**State** | Pointer to  |  | [optional] [default to {}]
**Version** | Pointer to **string** |  | [optional] [default to ""]
**Endpoint** | Pointer to **string** |  | [optional] [default to ""]

## Methods

### NewQdrantClusterStateNodesValue

`func NewQdrantClusterStateNodesValue(name string, ) *QdrantClusterStateNodesValue`

NewQdrantClusterStateNodesValue instantiates a new QdrantClusterStateNodesValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQdrantClusterStateNodesValueWithDefaults

`func NewQdrantClusterStateNodesValueWithDefaults() *QdrantClusterStateNodesValue`

NewQdrantClusterStateNodesValueWithDefaults instantiates a new QdrantClusterStateNodesValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *QdrantClusterStateNodesValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *QdrantClusterStateNodesValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *QdrantClusterStateNodesValue) SetName(v string)`

SetName sets Name field to given value.


### GetStartedAt

`func (o *QdrantClusterStateNodesValue) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *QdrantClusterStateNodesValue) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *QdrantClusterStateNodesValue) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *QdrantClusterStateNodesValue) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetState

`func (o *QdrantClusterStateNodesValue) GetState() map[string]interface{}`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *QdrantClusterStateNodesValue) GetStateOk() (*map[string]interface{}, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *QdrantClusterStateNodesValue) SetState(v map[string]interface{})`

SetState sets State field to given value.

### HasState

`func (o *QdrantClusterStateNodesValue) HasState() bool`

HasState returns a boolean if a field has been set.

### GetVersion

`func (o *QdrantClusterStateNodesValue) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *QdrantClusterStateNodesValue) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *QdrantClusterStateNodesValue) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *QdrantClusterStateNodesValue) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetEndpoint

`func (o *QdrantClusterStateNodesValue) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *QdrantClusterStateNodesValue) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *QdrantClusterStateNodesValue) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *QdrantClusterStateNodesValue) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


