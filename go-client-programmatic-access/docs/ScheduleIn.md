# ScheduleIn

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatorUserId** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **NullableString** |  | [optional] 
**Cron** | **string** |  | 
**Retention** | **int32** |  | 
**PrivateRegionId** | Pointer to **NullableString** |  | [optional] 
**MarkedForDeletionAt** | Pointer to **NullableTime** |  | [optional] 
**Status** | Pointer to [**NullableScheduleState**](ScheduleState.md) |  | [optional] 

## Methods

### NewScheduleIn

`func NewScheduleIn(cron string, retention int32, ) *ScheduleIn`

NewScheduleIn instantiates a new ScheduleIn object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewScheduleInWithDefaults

`func NewScheduleInWithDefaults() *ScheduleIn`

NewScheduleInWithDefaults instantiates a new ScheduleIn object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatorUserId

`func (o *ScheduleIn) GetCreatorUserId() string`

GetCreatorUserId returns the CreatorUserId field if non-nil, zero value otherwise.

### GetCreatorUserIdOk

`func (o *ScheduleIn) GetCreatorUserIdOk() (*string, bool)`

GetCreatorUserIdOk returns a tuple with the CreatorUserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatorUserId

`func (o *ScheduleIn) SetCreatorUserId(v string)`

SetCreatorUserId sets CreatorUserId field to given value.

### HasCreatorUserId

`func (o *ScheduleIn) HasCreatorUserId() bool`

HasCreatorUserId returns a boolean if a field has been set.

### SetCreatorUserIdNil

`func (o *ScheduleIn) SetCreatorUserIdNil(b bool)`

 SetCreatorUserIdNil sets the value for CreatorUserId to be an explicit nil

### UnsetCreatorUserId
`func (o *ScheduleIn) UnsetCreatorUserId()`

UnsetCreatorUserId ensures that no value is present for CreatorUserId, not even an explicit nil
### GetAccountId

`func (o *ScheduleIn) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *ScheduleIn) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *ScheduleIn) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *ScheduleIn) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *ScheduleIn) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *ScheduleIn) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
### GetCron

`func (o *ScheduleIn) GetCron() string`

GetCron returns the Cron field if non-nil, zero value otherwise.

### GetCronOk

`func (o *ScheduleIn) GetCronOk() (*string, bool)`

GetCronOk returns a tuple with the Cron field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCron

`func (o *ScheduleIn) SetCron(v string)`

SetCron sets Cron field to given value.


### GetRetention

`func (o *ScheduleIn) GetRetention() int32`

GetRetention returns the Retention field if non-nil, zero value otherwise.

### GetRetentionOk

`func (o *ScheduleIn) GetRetentionOk() (*int32, bool)`

GetRetentionOk returns a tuple with the Retention field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetention

`func (o *ScheduleIn) SetRetention(v int32)`

SetRetention sets Retention field to given value.


### GetPrivateRegionId

`func (o *ScheduleIn) GetPrivateRegionId() string`

GetPrivateRegionId returns the PrivateRegionId field if non-nil, zero value otherwise.

### GetPrivateRegionIdOk

`func (o *ScheduleIn) GetPrivateRegionIdOk() (*string, bool)`

GetPrivateRegionIdOk returns a tuple with the PrivateRegionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateRegionId

`func (o *ScheduleIn) SetPrivateRegionId(v string)`

SetPrivateRegionId sets PrivateRegionId field to given value.

### HasPrivateRegionId

`func (o *ScheduleIn) HasPrivateRegionId() bool`

HasPrivateRegionId returns a boolean if a field has been set.

### SetPrivateRegionIdNil

`func (o *ScheduleIn) SetPrivateRegionIdNil(b bool)`

 SetPrivateRegionIdNil sets the value for PrivateRegionId to be an explicit nil

### UnsetPrivateRegionId
`func (o *ScheduleIn) UnsetPrivateRegionId()`

UnsetPrivateRegionId ensures that no value is present for PrivateRegionId, not even an explicit nil
### GetMarkedForDeletionAt

`func (o *ScheduleIn) GetMarkedForDeletionAt() time.Time`

GetMarkedForDeletionAt returns the MarkedForDeletionAt field if non-nil, zero value otherwise.

### GetMarkedForDeletionAtOk

`func (o *ScheduleIn) GetMarkedForDeletionAtOk() (*time.Time, bool)`

GetMarkedForDeletionAtOk returns a tuple with the MarkedForDeletionAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMarkedForDeletionAt

`func (o *ScheduleIn) SetMarkedForDeletionAt(v time.Time)`

SetMarkedForDeletionAt sets MarkedForDeletionAt field to given value.

### HasMarkedForDeletionAt

`func (o *ScheduleIn) HasMarkedForDeletionAt() bool`

HasMarkedForDeletionAt returns a boolean if a field has been set.

### SetMarkedForDeletionAtNil

`func (o *ScheduleIn) SetMarkedForDeletionAtNil(b bool)`

 SetMarkedForDeletionAtNil sets the value for MarkedForDeletionAt to be an explicit nil

### UnsetMarkedForDeletionAt
`func (o *ScheduleIn) UnsetMarkedForDeletionAt()`

UnsetMarkedForDeletionAt ensures that no value is present for MarkedForDeletionAt, not even an explicit nil
### GetStatus

`func (o *ScheduleIn) GetStatus() ScheduleState`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ScheduleIn) GetStatusOk() (*ScheduleState, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ScheduleIn) SetStatus(v ScheduleState)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ScheduleIn) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### SetStatusNil

`func (o *ScheduleIn) SetStatusNil(b bool)`

 SetStatusNil sets the value for Status to be an explicit nil

### UnsetStatus
`func (o *ScheduleIn) UnsetStatus()`

UnsetStatus ensures that no value is present for Status, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


