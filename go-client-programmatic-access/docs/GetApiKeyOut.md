# GetApiKeyOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**UserId** | Pointer to **NullableString** |  | [optional] 
**Prefix** | **string** |  | 
**ClusterIdList** | Pointer to **[]interface{}** |  | [optional] 
**AccountId** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewGetApiKeyOut

`func NewGetApiKeyOut(id string, createdAt time.Time, prefix string, ) *GetApiKeyOut`

NewGetApiKeyOut instantiates a new GetApiKeyOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetApiKeyOutWithDefaults

`func NewGetApiKeyOutWithDefaults() *GetApiKeyOut`

NewGetApiKeyOutWithDefaults instantiates a new GetApiKeyOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *GetApiKeyOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GetApiKeyOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GetApiKeyOut) SetId(v string)`

SetId sets Id field to given value.


### GetCreatedAt

`func (o *GetApiKeyOut) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *GetApiKeyOut) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *GetApiKeyOut) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUserId

`func (o *GetApiKeyOut) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *GetApiKeyOut) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *GetApiKeyOut) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *GetApiKeyOut) HasUserId() bool`

HasUserId returns a boolean if a field has been set.

### SetUserIdNil

`func (o *GetApiKeyOut) SetUserIdNil(b bool)`

 SetUserIdNil sets the value for UserId to be an explicit nil

### UnsetUserId
`func (o *GetApiKeyOut) UnsetUserId()`

UnsetUserId ensures that no value is present for UserId, not even an explicit nil
### GetPrefix

`func (o *GetApiKeyOut) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *GetApiKeyOut) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *GetApiKeyOut) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.


### GetClusterIdList

`func (o *GetApiKeyOut) GetClusterIdList() []interface{}`

GetClusterIdList returns the ClusterIdList field if non-nil, zero value otherwise.

### GetClusterIdListOk

`func (o *GetApiKeyOut) GetClusterIdListOk() (*[]interface{}, bool)`

GetClusterIdListOk returns a tuple with the ClusterIdList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterIdList

`func (o *GetApiKeyOut) SetClusterIdList(v []interface{})`

SetClusterIdList sets ClusterIdList field to given value.

### HasClusterIdList

`func (o *GetApiKeyOut) HasClusterIdList() bool`

HasClusterIdList returns a boolean if a field has been set.

### SetClusterIdListNil

`func (o *GetApiKeyOut) SetClusterIdListNil(b bool)`

 SetClusterIdListNil sets the value for ClusterIdList to be an explicit nil

### UnsetClusterIdList
`func (o *GetApiKeyOut) UnsetClusterIdList()`

UnsetClusterIdList ensures that no value is present for ClusterIdList, not even an explicit nil
### GetAccountId

`func (o *GetApiKeyOut) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *GetApiKeyOut) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *GetApiKeyOut) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *GetApiKeyOut) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *GetApiKeyOut) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *GetApiKeyOut) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


