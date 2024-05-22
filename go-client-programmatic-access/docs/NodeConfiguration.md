# NodeConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PackageId** | **string** |  | 
**Package** | Pointer to [**NullablePackageOut**](PackageOut.md) |  | [optional] 
**ResourceConfigurations** | Pointer to [**[]ResourceConfiguration**](ResourceConfiguration.md) |  | [optional] 

## Methods

### NewNodeConfiguration

`func NewNodeConfiguration(packageId string, ) *NodeConfiguration`

NewNodeConfiguration instantiates a new NodeConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeConfigurationWithDefaults

`func NewNodeConfigurationWithDefaults() *NodeConfiguration`

NewNodeConfigurationWithDefaults instantiates a new NodeConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPackageId

`func (o *NodeConfiguration) GetPackageId() string`

GetPackageId returns the PackageId field if non-nil, zero value otherwise.

### GetPackageIdOk

`func (o *NodeConfiguration) GetPackageIdOk() (*string, bool)`

GetPackageIdOk returns a tuple with the PackageId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageId

`func (o *NodeConfiguration) SetPackageId(v string)`

SetPackageId sets PackageId field to given value.


### GetPackage

`func (o *NodeConfiguration) GetPackage() PackageOut`

GetPackage returns the Package field if non-nil, zero value otherwise.

### GetPackageOk

`func (o *NodeConfiguration) GetPackageOk() (*PackageOut, bool)`

GetPackageOk returns a tuple with the Package field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackage

`func (o *NodeConfiguration) SetPackage(v PackageOut)`

SetPackage sets Package field to given value.

### HasPackage

`func (o *NodeConfiguration) HasPackage() bool`

HasPackage returns a boolean if a field has been set.

### SetPackageNil

`func (o *NodeConfiguration) SetPackageNil(b bool)`

 SetPackageNil sets the value for Package to be an explicit nil

### UnsetPackage
`func (o *NodeConfiguration) UnsetPackage()`

UnsetPackage ensures that no value is present for Package, not even an explicit nil
### GetResourceConfigurations

`func (o *NodeConfiguration) GetResourceConfigurations() []ResourceConfiguration`

GetResourceConfigurations returns the ResourceConfigurations field if non-nil, zero value otherwise.

### GetResourceConfigurationsOk

`func (o *NodeConfiguration) GetResourceConfigurationsOk() (*[]ResourceConfiguration, bool)`

GetResourceConfigurationsOk returns a tuple with the ResourceConfigurations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceConfigurations

`func (o *NodeConfiguration) SetResourceConfigurations(v []ResourceConfiguration)`

SetResourceConfigurations sets ResourceConfigurations field to given value.

### HasResourceConfigurations

`func (o *NodeConfiguration) HasResourceConfigurations() bool`

HasResourceConfigurations returns a boolean if a field has been set.

### SetResourceConfigurationsNil

`func (o *NodeConfiguration) SetResourceConfigurationsNil(b bool)`

 SetResourceConfigurationsNil sets the value for ResourceConfigurations to be an explicit nil

### UnsetResourceConfigurations
`func (o *NodeConfiguration) UnsetResourceConfigurations()`

UnsetResourceConfigurations ensures that no value is present for ResourceConfigurations, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


