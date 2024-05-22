# QdrantClusterState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Version** | Pointer to **NullableString** |  | [optional] 
**ReplicationStatus** | [**ReplicationStatus**](ReplicationStatus.md) |  | 
**Replicas** | **int32** |  | 
**AvailableReplicas** | Pointer to **NullableInt32** |  | [optional] 
**RestartedAt** | Pointer to **NullableTime** |  | [optional] 
**Nodes** | **interface{}** |  | 
**Phase** | Pointer to [**NullableClusterState**](ClusterState.md) |  | [optional] 
**Reason** | Pointer to **NullableString** |  | [optional] 
**Endpoint** | Pointer to **NullableString** |  | [optional] 
**Current** | Pointer to [**QdrantClusterStatus**](QdrantClusterStatus.md) | Whether the cluster is running. | [optional] [default to QDRANTCLUSTERSTATUS_RUNNING]

## Methods

### NewQdrantClusterState

`func NewQdrantClusterState(id string, replicationStatus ReplicationStatus, replicas int32, nodes interface{}, ) *QdrantClusterState`

NewQdrantClusterState instantiates a new QdrantClusterState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQdrantClusterStateWithDefaults

`func NewQdrantClusterStateWithDefaults() *QdrantClusterState`

NewQdrantClusterStateWithDefaults instantiates a new QdrantClusterState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *QdrantClusterState) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *QdrantClusterState) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *QdrantClusterState) SetId(v string)`

SetId sets Id field to given value.


### GetVersion

`func (o *QdrantClusterState) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *QdrantClusterState) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *QdrantClusterState) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *QdrantClusterState) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### SetVersionNil

`func (o *QdrantClusterState) SetVersionNil(b bool)`

 SetVersionNil sets the value for Version to be an explicit nil

### UnsetVersion
`func (o *QdrantClusterState) UnsetVersion()`

UnsetVersion ensures that no value is present for Version, not even an explicit nil
### GetReplicationStatus

`func (o *QdrantClusterState) GetReplicationStatus() ReplicationStatus`

GetReplicationStatus returns the ReplicationStatus field if non-nil, zero value otherwise.

### GetReplicationStatusOk

`func (o *QdrantClusterState) GetReplicationStatusOk() (*ReplicationStatus, bool)`

GetReplicationStatusOk returns a tuple with the ReplicationStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicationStatus

`func (o *QdrantClusterState) SetReplicationStatus(v ReplicationStatus)`

SetReplicationStatus sets ReplicationStatus field to given value.


### GetReplicas

`func (o *QdrantClusterState) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *QdrantClusterState) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *QdrantClusterState) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.


### GetAvailableReplicas

`func (o *QdrantClusterState) GetAvailableReplicas() int32`

GetAvailableReplicas returns the AvailableReplicas field if non-nil, zero value otherwise.

### GetAvailableReplicasOk

`func (o *QdrantClusterState) GetAvailableReplicasOk() (*int32, bool)`

GetAvailableReplicasOk returns a tuple with the AvailableReplicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableReplicas

`func (o *QdrantClusterState) SetAvailableReplicas(v int32)`

SetAvailableReplicas sets AvailableReplicas field to given value.

### HasAvailableReplicas

`func (o *QdrantClusterState) HasAvailableReplicas() bool`

HasAvailableReplicas returns a boolean if a field has been set.

### SetAvailableReplicasNil

`func (o *QdrantClusterState) SetAvailableReplicasNil(b bool)`

 SetAvailableReplicasNil sets the value for AvailableReplicas to be an explicit nil

### UnsetAvailableReplicas
`func (o *QdrantClusterState) UnsetAvailableReplicas()`

UnsetAvailableReplicas ensures that no value is present for AvailableReplicas, not even an explicit nil
### GetRestartedAt

`func (o *QdrantClusterState) GetRestartedAt() time.Time`

GetRestartedAt returns the RestartedAt field if non-nil, zero value otherwise.

### GetRestartedAtOk

`func (o *QdrantClusterState) GetRestartedAtOk() (*time.Time, bool)`

GetRestartedAtOk returns a tuple with the RestartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestartedAt

`func (o *QdrantClusterState) SetRestartedAt(v time.Time)`

SetRestartedAt sets RestartedAt field to given value.

### HasRestartedAt

`func (o *QdrantClusterState) HasRestartedAt() bool`

HasRestartedAt returns a boolean if a field has been set.

### SetRestartedAtNil

`func (o *QdrantClusterState) SetRestartedAtNil(b bool)`

 SetRestartedAtNil sets the value for RestartedAt to be an explicit nil

### UnsetRestartedAt
`func (o *QdrantClusterState) UnsetRestartedAt()`

UnsetRestartedAt ensures that no value is present for RestartedAt, not even an explicit nil
### GetNodes

`func (o *QdrantClusterState) GetNodes() interface{}`

GetNodes returns the Nodes field if non-nil, zero value otherwise.

### GetNodesOk

`func (o *QdrantClusterState) GetNodesOk() (*interface{}, bool)`

GetNodesOk returns a tuple with the Nodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodes

`func (o *QdrantClusterState) SetNodes(v interface{})`

SetNodes sets Nodes field to given value.


### SetNodesNil

`func (o *QdrantClusterState) SetNodesNil(b bool)`

 SetNodesNil sets the value for Nodes to be an explicit nil

### UnsetNodes
`func (o *QdrantClusterState) UnsetNodes()`

UnsetNodes ensures that no value is present for Nodes, not even an explicit nil
### GetPhase

`func (o *QdrantClusterState) GetPhase() ClusterState`

GetPhase returns the Phase field if non-nil, zero value otherwise.

### GetPhaseOk

`func (o *QdrantClusterState) GetPhaseOk() (*ClusterState, bool)`

GetPhaseOk returns a tuple with the Phase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhase

`func (o *QdrantClusterState) SetPhase(v ClusterState)`

SetPhase sets Phase field to given value.

### HasPhase

`func (o *QdrantClusterState) HasPhase() bool`

HasPhase returns a boolean if a field has been set.

### SetPhaseNil

`func (o *QdrantClusterState) SetPhaseNil(b bool)`

 SetPhaseNil sets the value for Phase to be an explicit nil

### UnsetPhase
`func (o *QdrantClusterState) UnsetPhase()`

UnsetPhase ensures that no value is present for Phase, not even an explicit nil
### GetReason

`func (o *QdrantClusterState) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *QdrantClusterState) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *QdrantClusterState) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *QdrantClusterState) HasReason() bool`

HasReason returns a boolean if a field has been set.

### SetReasonNil

`func (o *QdrantClusterState) SetReasonNil(b bool)`

 SetReasonNil sets the value for Reason to be an explicit nil

### UnsetReason
`func (o *QdrantClusterState) UnsetReason()`

UnsetReason ensures that no value is present for Reason, not even an explicit nil
### GetEndpoint

`func (o *QdrantClusterState) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *QdrantClusterState) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *QdrantClusterState) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *QdrantClusterState) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### SetEndpointNil

`func (o *QdrantClusterState) SetEndpointNil(b bool)`

 SetEndpointNil sets the value for Endpoint to be an explicit nil

### UnsetEndpoint
`func (o *QdrantClusterState) UnsetEndpoint()`

UnsetEndpoint ensures that no value is present for Endpoint, not even an explicit nil
### GetCurrent

`func (o *QdrantClusterState) GetCurrent() QdrantClusterStatus`

GetCurrent returns the Current field if non-nil, zero value otherwise.

### GetCurrentOk

`func (o *QdrantClusterState) GetCurrentOk() (*QdrantClusterStatus, bool)`

GetCurrentOk returns a tuple with the Current field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrent

`func (o *QdrantClusterState) SetCurrent(v QdrantClusterStatus)`

SetCurrent sets Current field to given value.

### HasCurrent

`func (o *QdrantClusterState) HasCurrent() bool`

HasCurrent returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


