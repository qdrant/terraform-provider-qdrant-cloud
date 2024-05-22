# QdrantClusterNodeState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**StartedAt** | Pointer to **NullableTime** |  | [optional] 
**State** | Pointer to  |  | [optional] [default to {}]
**Version** | Pointer to **string** |  | [optional] [default to ""]
**Endpoint** | Pointer to **string** |  | [optional] [default to ""]

## Methods

### NewQdrantClusterNodeState

`func NewQdrantClusterNodeState(name string, ) *QdrantClusterNodeState`

NewQdrantClusterNodeState instantiates a new QdrantClusterNodeState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQdrantClusterNodeStateWithDefaults

`func NewQdrantClusterNodeStateWithDefaults() *QdrantClusterNodeState`

NewQdrantClusterNodeStateWithDefaults instantiates a new QdrantClusterNodeState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *QdrantClusterNodeState) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *QdrantClusterNodeState) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *QdrantClusterNodeState) SetName(v string)`

SetName sets Name field to given value.


### GetStartedAt

`func (o *QdrantClusterNodeState) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *QdrantClusterNodeState) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *QdrantClusterNodeState) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *QdrantClusterNodeState) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### SetStartedAtNil

`func (o *QdrantClusterNodeState) SetStartedAtNil(b bool)`

 SetStartedAtNil sets the value for StartedAt to be an explicit nil

### UnsetStartedAt
`func (o *QdrantClusterNodeState) UnsetStartedAt()`

UnsetStartedAt ensures that no value is present for StartedAt, not even an explicit nil
### GetState

`func (o *QdrantClusterNodeState) GetState() map[string]string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *QdrantClusterNodeState) GetStateOk() (*map[string]string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *QdrantClusterNodeState) SetState(v map[string]string)`

SetState sets State field to given value.

### HasState

`func (o *QdrantClusterNodeState) HasState() bool`

HasState returns a boolean if a field has been set.

### GetVersion

`func (o *QdrantClusterNodeState) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *QdrantClusterNodeState) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *QdrantClusterNodeState) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *QdrantClusterNodeState) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetEndpoint

`func (o *QdrantClusterNodeState) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *QdrantClusterNodeState) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *QdrantClusterNodeState) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *QdrantClusterNodeState) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


