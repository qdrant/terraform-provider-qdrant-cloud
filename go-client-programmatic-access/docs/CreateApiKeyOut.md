# CreateApiKeyOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**CreatedAt** | **time.Time** |  | 
**UserId** | Pointer to **NullableString** |  | [optional] 
**Prefix** | **string** |  | 
**ClusterIdList** | Pointer to **[]interface{}** |  | [optional] 
**AccountId** | Pointer to **NullableString** |  | [optional] 
**Token** | **string** |  | 

## Methods

### NewCreateApiKeyOut

`func NewCreateApiKeyOut(id string, createdAt time.Time, prefix string, token string, ) *CreateApiKeyOut`

NewCreateApiKeyOut instantiates a new CreateApiKeyOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateApiKeyOutWithDefaults

`func NewCreateApiKeyOutWithDefaults() *CreateApiKeyOut`

NewCreateApiKeyOutWithDefaults instantiates a new CreateApiKeyOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CreateApiKeyOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreateApiKeyOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreateApiKeyOut) SetId(v string)`

SetId sets Id field to given value.


### GetCreatedAt

`func (o *CreateApiKeyOut) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CreateApiKeyOut) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CreateApiKeyOut) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUserId

`func (o *CreateApiKeyOut) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *CreateApiKeyOut) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *CreateApiKeyOut) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *CreateApiKeyOut) HasUserId() bool`

HasUserId returns a boolean if a field has been set.

### SetUserIdNil

`func (o *CreateApiKeyOut) SetUserIdNil(b bool)`

 SetUserIdNil sets the value for UserId to be an explicit nil

### UnsetUserId
`func (o *CreateApiKeyOut) UnsetUserId()`

UnsetUserId ensures that no value is present for UserId, not even an explicit nil
### GetPrefix

`func (o *CreateApiKeyOut) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *CreateApiKeyOut) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *CreateApiKeyOut) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.


### GetClusterIdList

`func (o *CreateApiKeyOut) GetClusterIdList() []interface{}`

GetClusterIdList returns the ClusterIdList field if non-nil, zero value otherwise.

### GetClusterIdListOk

`func (o *CreateApiKeyOut) GetClusterIdListOk() (*[]interface{}, bool)`

GetClusterIdListOk returns a tuple with the ClusterIdList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterIdList

`func (o *CreateApiKeyOut) SetClusterIdList(v []interface{})`

SetClusterIdList sets ClusterIdList field to given value.

### HasClusterIdList

`func (o *CreateApiKeyOut) HasClusterIdList() bool`

HasClusterIdList returns a boolean if a field has been set.

### SetClusterIdListNil

`func (o *CreateApiKeyOut) SetClusterIdListNil(b bool)`

 SetClusterIdListNil sets the value for ClusterIdList to be an explicit nil

### UnsetClusterIdList
`func (o *CreateApiKeyOut) UnsetClusterIdList()`

UnsetClusterIdList ensures that no value is present for ClusterIdList, not even an explicit nil
### GetAccountId

`func (o *CreateApiKeyOut) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *CreateApiKeyOut) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *CreateApiKeyOut) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *CreateApiKeyOut) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### SetAccountIdNil

`func (o *CreateApiKeyOut) SetAccountIdNil(b bool)`

 SetAccountIdNil sets the value for AccountId to be an explicit nil

### UnsetAccountId
`func (o *CreateApiKeyOut) UnsetAccountId()`

UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
### GetToken

`func (o *CreateApiKeyOut) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *CreateApiKeyOut) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *CreateApiKeyOut) SetToken(v string)`

SetToken sets Token field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


