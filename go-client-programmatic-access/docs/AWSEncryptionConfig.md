# AWSEncryptionConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Managed** | Pointer to **bool** |  | [optional] [default to true]
**EncryptionKeyId** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewAWSEncryptionConfig

`func NewAWSEncryptionConfig() *AWSEncryptionConfig`

NewAWSEncryptionConfig instantiates a new AWSEncryptionConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAWSEncryptionConfigWithDefaults

`func NewAWSEncryptionConfigWithDefaults() *AWSEncryptionConfig`

NewAWSEncryptionConfigWithDefaults instantiates a new AWSEncryptionConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetManaged

`func (o *AWSEncryptionConfig) GetManaged() bool`

GetManaged returns the Managed field if non-nil, zero value otherwise.

### GetManagedOk

`func (o *AWSEncryptionConfig) GetManagedOk() (*bool, bool)`

GetManagedOk returns a tuple with the Managed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManaged

`func (o *AWSEncryptionConfig) SetManaged(v bool)`

SetManaged sets Managed field to given value.

### HasManaged

`func (o *AWSEncryptionConfig) HasManaged() bool`

HasManaged returns a boolean if a field has been set.

### GetEncryptionKeyId

`func (o *AWSEncryptionConfig) GetEncryptionKeyId() string`

GetEncryptionKeyId returns the EncryptionKeyId field if non-nil, zero value otherwise.

### GetEncryptionKeyIdOk

`func (o *AWSEncryptionConfig) GetEncryptionKeyIdOk() (*string, bool)`

GetEncryptionKeyIdOk returns a tuple with the EncryptionKeyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionKeyId

`func (o *AWSEncryptionConfig) SetEncryptionKeyId(v string)`

SetEncryptionKeyId sets EncryptionKeyId field to given value.

### HasEncryptionKeyId

`func (o *AWSEncryptionConfig) HasEncryptionKeyId() bool`

HasEncryptionKeyId returns a boolean if a field has been set.

### SetEncryptionKeyIdNil

`func (o *AWSEncryptionConfig) SetEncryptionKeyIdNil(b bool)`

 SetEncryptionKeyIdNil sets the value for EncryptionKeyId to be an explicit nil

### UnsetEncryptionKeyId
`func (o *AWSEncryptionConfig) UnsetEncryptionKeyId()`

UnsetEncryptionKeyId ensures that no value is present for EncryptionKeyId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


