# ClusterOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**OwnerId** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **NullableString** |  | [optional] 
**Name** | **string** |  | 
**CloudProvider** | **string** |  | 
**CloudRegion** | **string** |  | 
**CloudRegionAz** | Pointer to **NullableString** |  | [optional] 
**CloudRegionSetup** | Pointer to **NullableString** |  | [optional] 
**CurrentConfigurationId** | **string** |  | 
**PrivateRegionId** | Pointer to **NullableString** |  | [optional] 
**EncryptionKeyId** | Pointer to **NullableString** |  | [optional] 
**MarkedForDeletionAt** | Pointer to **NullableTime** |  | [optional] 
**Version** | Pointer to **NullableString** |  | [optional] 
**Url** | **string** |  | 
**State** | Pointer to [**NullableQdrantClusterState**](QdrantClusterState.md) |  | [optional] 
**Configuration** | Pointer to [**NullableClusterConfigurationOut**](ClusterConfigurationOut.md) |  | [optional] 
**Resources** | Pointer to [**NullableClusterResourcesSummary**](ClusterResourcesSummary.md) |  | [optional] 
**TotalExtraDisk** | Pointer to **int32** |  | [optional] [default to 0]

## Methods

### NewClusterOut

`func NewClusterOut(id string, createdAt time.Time, name string, cloudProvider string, cloudRegion string, currentConfigurationId string, url string, ) *ClusterOut`

NewClusterOut instantiates a new ClusterOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClusterOutWithDefaults

`func NewClusterOutWithDefaults() *ClusterOut`

NewClusterOutWithDefaults instantiates a new ClusterOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ClusterOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClusterOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClusterOut) SetId(v string)`

SetId sets Id field to given value.


### GetCreatedAt

`func (o *ClusterOut) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ClusterOut) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ClusterOut) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetOwnerId

`func (o *ClusterOut) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *ClusterOut) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *ClusterOut) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.

### HasOwnerId

`func (o *ClusterOut) HasOwnerId() bool`

HasOwnerId returns a boolean if a field has been set.

### SetOwnerIdNil

`func (o *ClusterOut) SetOwnerIdNil(b bool)`

 SetOwnerIdNil sets the value for OwnerId to be an explicit nil

### UnsetOwnerId
`func (o *ClusterOut) UnsetOwnerId()`

UnsetOwnerId ensures that no value is present for OwnerId, not even an explicit nil
### GetAccountId

`func (o *ClusterOut) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *ClusterOut) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *ClusterOut) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *ClusterOut) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *ClusterOut) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *ClusterOut) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
### GetName

`func (o *ClusterOut) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClusterOut) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClusterOut) SetName(v string)`

SetName sets Name field to given value.


### GetCloudProvider

`func (o *ClusterOut) GetCloudProvider() string`

GetCloudProvider returns the CloudProvider field if non-nil, zero value otherwise.

### GetCloudProviderOk

`func (o *ClusterOut) GetCloudProviderOk() (*string, bool)`

GetCloudProviderOk returns a tuple with the CloudProvider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudProvider

`func (o *ClusterOut) SetCloudProvider(v string)`

SetCloudProvider sets CloudProvider field to given value.


### GetCloudRegion

`func (o *ClusterOut) GetCloudRegion() string`

GetCloudRegion returns the CloudRegion field if non-nil, zero value otherwise.

### GetCloudRegionOk

`func (o *ClusterOut) GetCloudRegionOk() (*string, bool)`

GetCloudRegionOk returns a tuple with the CloudRegion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegion

`func (o *ClusterOut) SetCloudRegion(v string)`

SetCloudRegion sets CloudRegion field to given value.


### GetCloudRegionAz

`func (o *ClusterOut) GetCloudRegionAz() string`

GetCloudRegionAz returns the CloudRegionAz field if non-nil, zero value otherwise.

### GetCloudRegionAzOk

`func (o *ClusterOut) GetCloudRegionAzOk() (*string, bool)`

GetCloudRegionAzOk returns a tuple with the CloudRegionAz field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegionAz

`func (o *ClusterOut) SetCloudRegionAz(v string)`

SetCloudRegionAz sets CloudRegionAz field to given value.

### HasCloudRegionAz

`func (o *ClusterOut) HasCloudRegionAz() bool`

HasCloudRegionAz returns a boolean if a field has been set.

### SetCloudRegionAzNil

`func (o *ClusterOut) SetCloudRegionAzNil(b bool)`

 SetCloudRegionAzNil sets the value for CloudRegionAz to be an explicit nil

### UnsetCloudRegionAz
`func (o *ClusterOut) UnsetCloudRegionAz()`

UnsetCloudRegionAz ensures that no value is present for CloudRegionAz, not even an explicit nil
### GetCloudRegionSetup

`func (o *ClusterOut) GetCloudRegionSetup() string`

GetCloudRegionSetup returns the CloudRegionSetup field if non-nil, zero value otherwise.

### GetCloudRegionSetupOk

`func (o *ClusterOut) GetCloudRegionSetupOk() (*string, bool)`

GetCloudRegionSetupOk returns a tuple with the CloudRegionSetup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegionSetup

`func (o *ClusterOut) SetCloudRegionSetup(v string)`

SetCloudRegionSetup sets CloudRegionSetup field to given value.

### HasCloudRegionSetup

`func (o *ClusterOut) HasCloudRegionSetup() bool`

HasCloudRegionSetup returns a boolean if a field has been set.

### SetCloudRegionSetupNil

`func (o *ClusterOut) SetCloudRegionSetupNil(b bool)`

 SetCloudRegionSetupNil sets the value for CloudRegionSetup to be an explicit nil

### UnsetCloudRegionSetup
`func (o *ClusterOut) UnsetCloudRegionSetup()`

UnsetCloudRegionSetup ensures that no value is present for CloudRegionSetup, not even an explicit nil
### GetCurrentConfigurationId

`func (o *ClusterOut) GetCurrentConfigurationId() string`

GetCurrentConfigurationId returns the CurrentConfigurationId field if non-nil, zero value otherwise.

### GetCurrentConfigurationIdOk

`func (o *ClusterOut) GetCurrentConfigurationIdOk() (*string, bool)`

GetCurrentConfigurationIdOk returns a tuple with the CurrentConfigurationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentConfigurationId

`func (o *ClusterOut) SetCurrentConfigurationId(v string)`

SetCurrentConfigurationId sets CurrentConfigurationId field to given value.


### GetPrivateRegionId

`func (o *ClusterOut) GetPrivateRegionId() string`

GetPrivateRegionId returns the PrivateRegionId field if non-nil, zero value otherwise.

### GetPrivateRegionIdOk

`func (o *ClusterOut) GetPrivateRegionIdOk() (*string, bool)`

GetPrivateRegionIdOk returns a tuple with the PrivateRegionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateRegionId

`func (o *ClusterOut) SetPrivateRegionId(v string)`

SetPrivateRegionId sets PrivateRegionId field to given value.

### HasPrivateRegionId

`func (o *ClusterOut) HasPrivateRegionId() bool`

HasPrivateRegionId returns a boolean if a field has been set.

### SetPrivateRegionIdNil

`func (o *ClusterOut) SetPrivateRegionIdNil(b bool)`

 SetPrivateRegionIdNil sets the value for PrivateRegionId to be an explicit nil

### UnsetPrivateRegionId
`func (o *ClusterOut) UnsetPrivateRegionId()`

UnsetPrivateRegionId ensures that no value is present for PrivateRegionId, not even an explicit nil
### GetEncryptionKeyId

`func (o *ClusterOut) GetEncryptionKeyId() string`

GetEncryptionKeyId returns the EncryptionKeyId field if non-nil, zero value otherwise.

### GetEncryptionKeyIdOk

`func (o *ClusterOut) GetEncryptionKeyIdOk() (*string, bool)`

GetEncryptionKeyIdOk returns a tuple with the EncryptionKeyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionKeyId

`func (o *ClusterOut) SetEncryptionKeyId(v string)`

SetEncryptionKeyId sets EncryptionKeyId field to given value.

### HasEncryptionKeyId

`func (o *ClusterOut) HasEncryptionKeyId() bool`

HasEncryptionKeyId returns a boolean if a field has been set.

### SetEncryptionKeyIdNil

`func (o *ClusterOut) SetEncryptionKeyIdNil(b bool)`

 SetEncryptionKeyIdNil sets the value for EncryptionKeyId to be an explicit nil

### UnsetEncryptionKeyId
`func (o *ClusterOut) UnsetEncryptionKeyId()`

UnsetEncryptionKeyId ensures that no value is present for EncryptionKeyId, not even an explicit nil
### GetMarkedForDeletionAt

`func (o *ClusterOut) GetMarkedForDeletionAt() time.Time`

GetMarkedForDeletionAt returns the MarkedForDeletionAt field if non-nil, zero value otherwise.

### GetMarkedForDeletionAtOk

`func (o *ClusterOut) GetMarkedForDeletionAtOk() (*time.Time, bool)`

GetMarkedForDeletionAtOk returns a tuple with the MarkedForDeletionAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMarkedForDeletionAt

`func (o *ClusterOut) SetMarkedForDeletionAt(v time.Time)`

SetMarkedForDeletionAt sets MarkedForDeletionAt field to given value.

### HasMarkedForDeletionAt

`func (o *ClusterOut) HasMarkedForDeletionAt() bool`

HasMarkedForDeletionAt returns a boolean if a field has been set.

### SetMarkedForDeletionAtNil

`func (o *ClusterOut) SetMarkedForDeletionAtNil(b bool)`

 SetMarkedForDeletionAtNil sets the value for MarkedForDeletionAt to be an explicit nil

### UnsetMarkedForDeletionAt
`func (o *ClusterOut) UnsetMarkedForDeletionAt()`

UnsetMarkedForDeletionAt ensures that no value is present for MarkedForDeletionAt, not even an explicit nil
### GetVersion

`func (o *ClusterOut) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ClusterOut) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ClusterOut) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ClusterOut) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### SetVersionNil

`func (o *ClusterOut) SetVersionNil(b bool)`

 SetVersionNil sets the value for Version to be an explicit nil

### UnsetVersion
`func (o *ClusterOut) UnsetVersion()`

UnsetVersion ensures that no value is present for Version, not even an explicit nil
### GetUrl

`func (o *ClusterOut) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ClusterOut) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ClusterOut) SetUrl(v string)`

SetUrl sets Url field to given value.


### GetState

`func (o *ClusterOut) GetState() QdrantClusterState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ClusterOut) GetStateOk() (*QdrantClusterState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ClusterOut) SetState(v QdrantClusterState)`

SetState sets State field to given value.

### HasState

`func (o *ClusterOut) HasState() bool`

HasState returns a boolean if a field has been set.

### SetStateNil

`func (o *ClusterOut) SetStateNil(b bool)`

 SetStateNil sets the value for State to be an explicit nil

### UnsetState
`func (o *ClusterOut) UnsetState()`

UnsetState ensures that no value is present for State, not even an explicit nil
### GetConfiguration

`func (o *ClusterOut) GetConfiguration() ClusterConfigurationOut`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *ClusterOut) GetConfigurationOk() (*ClusterConfigurationOut, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *ClusterOut) SetConfiguration(v ClusterConfigurationOut)`

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *ClusterOut) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.

### SetConfigurationNil

`func (o *ClusterOut) SetConfigurationNil(b bool)`

 SetConfigurationNil sets the value for Configuration to be an explicit nil

### UnsetConfiguration
`func (o *ClusterOut) UnsetConfiguration()`

UnsetConfiguration ensures that no value is present for Configuration, not even an explicit nil
### GetResources

`func (o *ClusterOut) GetResources() ClusterResourcesSummary`

GetResources returns the Resources field if non-nil, zero value otherwise.

### GetResourcesOk

`func (o *ClusterOut) GetResourcesOk() (*ClusterResourcesSummary, bool)`

GetResourcesOk returns a tuple with the Resources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResources

`func (o *ClusterOut) SetResources(v ClusterResourcesSummary)`

SetResources sets Resources field to given value.

### HasResources

`func (o *ClusterOut) HasResources() bool`

HasResources returns a boolean if a field has been set.

### SetResourcesNil

`func (o *ClusterOut) SetResourcesNil(b bool)`

 SetResourcesNil sets the value for Resources to be an explicit nil

### UnsetResources
`func (o *ClusterOut) UnsetResources()`

UnsetResources ensures that no value is present for Resources, not even an explicit nil
### GetTotalExtraDisk

`func (o *ClusterOut) GetTotalExtraDisk() int32`

GetTotalExtraDisk returns the TotalExtraDisk field if non-nil, zero value otherwise.

### GetTotalExtraDiskOk

`func (o *ClusterOut) GetTotalExtraDiskOk() (*int32, bool)`

GetTotalExtraDiskOk returns a tuple with the TotalExtraDisk field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalExtraDisk

`func (o *ClusterOut) SetTotalExtraDisk(v int32)`

SetTotalExtraDisk sets TotalExtraDisk field to given value.

### HasTotalExtraDisk

`func (o *ClusterOut) HasTotalExtraDisk() bool`

HasTotalExtraDisk returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


