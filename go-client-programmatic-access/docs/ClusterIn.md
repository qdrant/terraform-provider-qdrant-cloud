# ClusterIn

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**OwnerId** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **NullableString** |  | [optional] 
**Name** | **string** |  | 
**CloudProvider** | **string** |  | 
**CloudRegion** | **string** |  | 
**CloudRegionAz** | Pointer to **NullableString** |  | [optional] 
**CloudRegionSetup** | Pointer to **NullableString** |  | [optional] 
**PrivateRegionId** | Pointer to **NullableString** |  | [optional] 
**Version** | Pointer to **NullableString** |  | [optional] 
**Configuration** | [**ClusterConfigurationIn**](ClusterConfigurationIn.md) |  | 
**Schedule** | Pointer to [**NullableScheduleIn**](ScheduleIn.md) |  | [optional] 
**EncryptionConfig** | Pointer to [**NullableEncryptionConfigIn**](EncryptionConfigIn.md) |  | [optional] 

## Methods

### NewClusterIn

`func NewClusterIn(name string, cloudProvider string, cloudRegion string, configuration ClusterConfigurationIn, ) *ClusterIn`

NewClusterIn instantiates a new ClusterIn object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClusterInWithDefaults

`func NewClusterInWithDefaults() *ClusterIn`

NewClusterInWithDefaults instantiates a new ClusterIn object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOwnerId

`func (o *ClusterIn) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *ClusterIn) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *ClusterIn) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.

### HasOwnerId

`func (o *ClusterIn) HasOwnerId() bool`

HasOwnerId returns a boolean if a field has been set.

### SetOwnerIdNil

`func (o *ClusterIn) SetOwnerIdNil(b bool)`

 SetOwnerIdNil sets the value for OwnerId to be an explicit nil

### UnsetOwnerId
`func (o *ClusterIn) UnsetOwnerId()`

UnsetOwnerId ensures that no value is present for OwnerId, not even an explicit nil
### GetAccountId

`func (o *ClusterIn) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *ClusterIn) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *ClusterIn) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *ClusterIn) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *ClusterIn) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *ClusterIn) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
### GetName

`func (o *ClusterIn) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClusterIn) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClusterIn) SetName(v string)`

SetName sets Name field to given value.


### GetCloudProvider

`func (o *ClusterIn) GetCloudProvider() string`

GetCloudProvider returns the CloudProvider field if non-nil, zero value otherwise.

### GetCloudProviderOk

`func (o *ClusterIn) GetCloudProviderOk() (*string, bool)`

GetCloudProviderOk returns a tuple with the CloudProvider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudProvider

`func (o *ClusterIn) SetCloudProvider(v string)`

SetCloudProvider sets CloudProvider field to given value.


### GetCloudRegion

`func (o *ClusterIn) GetCloudRegion() string`

GetCloudRegion returns the CloudRegion field if non-nil, zero value otherwise.

### GetCloudRegionOk

`func (o *ClusterIn) GetCloudRegionOk() (*string, bool)`

GetCloudRegionOk returns a tuple with the CloudRegion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegion

`func (o *ClusterIn) SetCloudRegion(v string)`

SetCloudRegion sets CloudRegion field to given value.


### GetCloudRegionAz

`func (o *ClusterIn) GetCloudRegionAz() string`

GetCloudRegionAz returns the CloudRegionAz field if non-nil, zero value otherwise.

### GetCloudRegionAzOk

`func (o *ClusterIn) GetCloudRegionAzOk() (*string, bool)`

GetCloudRegionAzOk returns a tuple with the CloudRegionAz field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegionAz

`func (o *ClusterIn) SetCloudRegionAz(v string)`

SetCloudRegionAz sets CloudRegionAz field to given value.

### HasCloudRegionAz

`func (o *ClusterIn) HasCloudRegionAz() bool`

HasCloudRegionAz returns a boolean if a field has been set.

### SetCloudRegionAzNil

`func (o *ClusterIn) SetCloudRegionAzNil(b bool)`

 SetCloudRegionAzNil sets the value for CloudRegionAz to be an explicit nil

### UnsetCloudRegionAz
`func (o *ClusterIn) UnsetCloudRegionAz()`

UnsetCloudRegionAz ensures that no value is present for CloudRegionAz, not even an explicit nil
### GetCloudRegionSetup

`func (o *ClusterIn) GetCloudRegionSetup() string`

GetCloudRegionSetup returns the CloudRegionSetup field if non-nil, zero value otherwise.

### GetCloudRegionSetupOk

`func (o *ClusterIn) GetCloudRegionSetupOk() (*string, bool)`

GetCloudRegionSetupOk returns a tuple with the CloudRegionSetup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCloudRegionSetup

`func (o *ClusterIn) SetCloudRegionSetup(v string)`

SetCloudRegionSetup sets CloudRegionSetup field to given value.

### HasCloudRegionSetup

`func (o *ClusterIn) HasCloudRegionSetup() bool`

HasCloudRegionSetup returns a boolean if a field has been set.

### SetCloudRegionSetupNil

`func (o *ClusterIn) SetCloudRegionSetupNil(b bool)`

 SetCloudRegionSetupNil sets the value for CloudRegionSetup to be an explicit nil

### UnsetCloudRegionSetup
`func (o *ClusterIn) UnsetCloudRegionSetup()`

UnsetCloudRegionSetup ensures that no value is present for CloudRegionSetup, not even an explicit nil
### GetPrivateRegionId

`func (o *ClusterIn) GetPrivateRegionId() string`

GetPrivateRegionId returns the PrivateRegionId field if non-nil, zero value otherwise.

### GetPrivateRegionIdOk

`func (o *ClusterIn) GetPrivateRegionIdOk() (*string, bool)`

GetPrivateRegionIdOk returns a tuple with the PrivateRegionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateRegionId

`func (o *ClusterIn) SetPrivateRegionId(v string)`

SetPrivateRegionId sets PrivateRegionId field to given value.

### HasPrivateRegionId

`func (o *ClusterIn) HasPrivateRegionId() bool`

HasPrivateRegionId returns a boolean if a field has been set.

### SetPrivateRegionIdNil

`func (o *ClusterIn) SetPrivateRegionIdNil(b bool)`

 SetPrivateRegionIdNil sets the value for PrivateRegionId to be an explicit nil

### UnsetPrivateRegionId
`func (o *ClusterIn) UnsetPrivateRegionId()`

UnsetPrivateRegionId ensures that no value is present for PrivateRegionId, not even an explicit nil
### GetVersion

`func (o *ClusterIn) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ClusterIn) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ClusterIn) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ClusterIn) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### SetVersionNil

`func (o *ClusterIn) SetVersionNil(b bool)`

 SetVersionNil sets the value for Version to be an explicit nil

### UnsetVersion
`func (o *ClusterIn) UnsetVersion()`

UnsetVersion ensures that no value is present for Version, not even an explicit nil
### GetConfiguration

`func (o *ClusterIn) GetConfiguration() ClusterConfigurationIn`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *ClusterIn) GetConfigurationOk() (*ClusterConfigurationIn, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *ClusterIn) SetConfiguration(v ClusterConfigurationIn)`

SetConfiguration sets Configuration field to given value.


### GetSchedule

`func (o *ClusterIn) GetSchedule() ScheduleIn`

GetSchedule returns the Schedule field if non-nil, zero value otherwise.

### GetScheduleOk

`func (o *ClusterIn) GetScheduleOk() (*ScheduleIn, bool)`

GetScheduleOk returns a tuple with the Schedule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchedule

`func (o *ClusterIn) SetSchedule(v ScheduleIn)`

SetSchedule sets Schedule field to given value.

### HasSchedule

`func (o *ClusterIn) HasSchedule() bool`

HasSchedule returns a boolean if a field has been set.

### SetScheduleNil

`func (o *ClusterIn) SetScheduleNil(b bool)`

 SetScheduleNil sets the value for Schedule to be an explicit nil

### UnsetSchedule
`func (o *ClusterIn) UnsetSchedule()`

UnsetSchedule ensures that no value is present for Schedule, not even an explicit nil
### GetEncryptionConfig

`func (o *ClusterIn) GetEncryptionConfig() EncryptionConfigIn`

GetEncryptionConfig returns the EncryptionConfig field if non-nil, zero value otherwise.

### GetEncryptionConfigOk

`func (o *ClusterIn) GetEncryptionConfigOk() (*EncryptionConfigIn, bool)`

GetEncryptionConfigOk returns a tuple with the EncryptionConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionConfig

`func (o *ClusterIn) SetEncryptionConfig(v EncryptionConfigIn)`

SetEncryptionConfig sets EncryptionConfig field to given value.

### HasEncryptionConfig

`func (o *ClusterIn) HasEncryptionConfig() bool`

HasEncryptionConfig returns a boolean if a field has been set.

### SetEncryptionConfigNil

`func (o *ClusterIn) SetEncryptionConfigNil(b bool)`

 SetEncryptionConfigNil sets the value for EncryptionConfig to be an explicit nil

### UnsetEncryptionConfig
`func (o *ClusterIn) UnsetEncryptionConfig()`

UnsetEncryptionConfig ensures that no value is present for EncryptionConfig, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


