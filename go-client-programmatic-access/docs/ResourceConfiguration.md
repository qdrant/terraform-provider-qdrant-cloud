# ResourceConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceOptionId** | **string** |  | 
**ResourceOption** | Pointer to [**NullableResourceOptionOut**](ResourceOptionOut.md) |  | [optional] 
**Amount** | **int32** |  | 

## Methods

### NewResourceConfiguration

`func NewResourceConfiguration(resourceOptionId string, amount int32, ) *ResourceConfiguration`

NewResourceConfiguration instantiates a new ResourceConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResourceConfigurationWithDefaults

`func NewResourceConfigurationWithDefaults() *ResourceConfiguration`

NewResourceConfigurationWithDefaults instantiates a new ResourceConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetResourceOptionId

`func (o *ResourceConfiguration) GetResourceOptionId() string`

GetResourceOptionId returns the ResourceOptionId field if non-nil, zero value otherwise.

### GetResourceOptionIdOk

`func (o *ResourceConfiguration) GetResourceOptionIdOk() (*string, bool)`

GetResourceOptionIdOk returns a tuple with the ResourceOptionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceOptionId

`func (o *ResourceConfiguration) SetResourceOptionId(v string)`

SetResourceOptionId sets ResourceOptionId field to given value.


### GetResourceOption

`func (o *ResourceConfiguration) GetResourceOption() ResourceOptionOut`

GetResourceOption returns the ResourceOption field if non-nil, zero value otherwise.

### GetResourceOptionOk

`func (o *ResourceConfiguration) GetResourceOptionOk() (*ResourceOptionOut, bool)`

GetResourceOptionOk returns a tuple with the ResourceOption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceOption

`func (o *ResourceConfiguration) SetResourceOption(v ResourceOptionOut)`

SetResourceOption sets ResourceOption field to given value.

### HasResourceOption

`func (o *ResourceConfiguration) HasResourceOption() bool`

HasResourceOption returns a boolean if a field has been set.

### SetResourceOptionNil

`func (o *ResourceConfiguration) SetResourceOptionNil(b bool)`

 SetResourceOptionNil sets the value for ResourceOption to be an explicit nil

### UnsetResourceOption
`func (o *ResourceConfiguration) UnsetResourceOption()`

UnsetResourceOption ensures that no value is present for ResourceOption, not even an explicit nil
### GetAmount

`func (o *ResourceConfiguration) GetAmount() int32`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *ResourceConfiguration) GetAmountOk() (*int32, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *ResourceConfiguration) SetAmount(v int32)`

SetAmount sets Amount field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


