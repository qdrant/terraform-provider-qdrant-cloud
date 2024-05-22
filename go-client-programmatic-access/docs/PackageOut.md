# PackageOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **NullableString** |  | [optional] 
**ResourceConfiguration** | [**[]ResourceConfiguration**](ResourceConfiguration.md) |  | 
**Name** | **string** |  | 
**Status** | [**BookingStatus**](BookingStatus.md) |  | 
**Currency** | [**Currency**](Currency.md) |  | 
**UnitIntPricePerHour** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerDay** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerMonth** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerYear** | Pointer to **NullableInt32** |  | [optional] 
**RegionalMappingId** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewPackageOut

`func NewPackageOut(resourceConfiguration []ResourceConfiguration, name string, status BookingStatus, currency Currency, ) *PackageOut`

NewPackageOut instantiates a new PackageOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPackageOutWithDefaults

`func NewPackageOutWithDefaults() *PackageOut`

NewPackageOutWithDefaults instantiates a new PackageOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PackageOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PackageOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PackageOut) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *PackageOut) HasId() bool`

HasId returns a boolean if a field has been set.

### SetIdNil

`func (o *PackageOut) SetIdNil(b bool)`

 SetIdNil sets the value for Id to be an explicit nil

### UnsetId
`func (o *PackageOut) UnsetId()`

UnsetId ensures that no value is present for Id, not even an explicit nil
### GetResourceConfiguration

`func (o *PackageOut) GetResourceConfiguration() []ResourceConfiguration`

GetResourceConfiguration returns the ResourceConfiguration field if non-nil, zero value otherwise.

### GetResourceConfigurationOk

`func (o *PackageOut) GetResourceConfigurationOk() (*[]ResourceConfiguration, bool)`

GetResourceConfigurationOk returns a tuple with the ResourceConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceConfiguration

`func (o *PackageOut) SetResourceConfiguration(v []ResourceConfiguration)`

SetResourceConfiguration sets ResourceConfiguration field to given value.


### GetName

`func (o *PackageOut) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PackageOut) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PackageOut) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *PackageOut) GetStatus() BookingStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PackageOut) GetStatusOk() (*BookingStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PackageOut) SetStatus(v BookingStatus)`

SetStatus sets Status field to given value.


### GetCurrency

`func (o *PackageOut) GetCurrency() Currency`

GetCurrency returns the Currency field if non-nil, zero value otherwise.

### GetCurrencyOk

`func (o *PackageOut) GetCurrencyOk() (*Currency, bool)`

GetCurrencyOk returns a tuple with the Currency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrency

`func (o *PackageOut) SetCurrency(v Currency)`

SetCurrency sets Currency field to given value.


### GetUnitIntPricePerHour

`func (o *PackageOut) GetUnitIntPricePerHour() int32`

GetUnitIntPricePerHour returns the UnitIntPricePerHour field if non-nil, zero value otherwise.

### GetUnitIntPricePerHourOk

`func (o *PackageOut) GetUnitIntPricePerHourOk() (*int32, bool)`

GetUnitIntPricePerHourOk returns a tuple with the UnitIntPricePerHour field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerHour

`func (o *PackageOut) SetUnitIntPricePerHour(v int32)`

SetUnitIntPricePerHour sets UnitIntPricePerHour field to given value.

### HasUnitIntPricePerHour

`func (o *PackageOut) HasUnitIntPricePerHour() bool`

HasUnitIntPricePerHour returns a boolean if a field has been set.

### SetUnitIntPricePerHourNil

`func (o *PackageOut) SetUnitIntPricePerHourNil(b bool)`

 SetUnitIntPricePerHourNil sets the value for UnitIntPricePerHour to be an explicit nil

### UnsetUnitIntPricePerHour
`func (o *PackageOut) UnsetUnitIntPricePerHour()`

UnsetUnitIntPricePerHour ensures that no value is present for UnitIntPricePerHour, not even an explicit nil
### GetUnitIntPricePerDay

`func (o *PackageOut) GetUnitIntPricePerDay() int32`

GetUnitIntPricePerDay returns the UnitIntPricePerDay field if non-nil, zero value otherwise.

### GetUnitIntPricePerDayOk

`func (o *PackageOut) GetUnitIntPricePerDayOk() (*int32, bool)`

GetUnitIntPricePerDayOk returns a tuple with the UnitIntPricePerDay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerDay

`func (o *PackageOut) SetUnitIntPricePerDay(v int32)`

SetUnitIntPricePerDay sets UnitIntPricePerDay field to given value.

### HasUnitIntPricePerDay

`func (o *PackageOut) HasUnitIntPricePerDay() bool`

HasUnitIntPricePerDay returns a boolean if a field has been set.

### SetUnitIntPricePerDayNil

`func (o *PackageOut) SetUnitIntPricePerDayNil(b bool)`

 SetUnitIntPricePerDayNil sets the value for UnitIntPricePerDay to be an explicit nil

### UnsetUnitIntPricePerDay
`func (o *PackageOut) UnsetUnitIntPricePerDay()`

UnsetUnitIntPricePerDay ensures that no value is present for UnitIntPricePerDay, not even an explicit nil
### GetUnitIntPricePerMonth

`func (o *PackageOut) GetUnitIntPricePerMonth() int32`

GetUnitIntPricePerMonth returns the UnitIntPricePerMonth field if non-nil, zero value otherwise.

### GetUnitIntPricePerMonthOk

`func (o *PackageOut) GetUnitIntPricePerMonthOk() (*int32, bool)`

GetUnitIntPricePerMonthOk returns a tuple with the UnitIntPricePerMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerMonth

`func (o *PackageOut) SetUnitIntPricePerMonth(v int32)`

SetUnitIntPricePerMonth sets UnitIntPricePerMonth field to given value.

### HasUnitIntPricePerMonth

`func (o *PackageOut) HasUnitIntPricePerMonth() bool`

HasUnitIntPricePerMonth returns a boolean if a field has been set.

### SetUnitIntPricePerMonthNil

`func (o *PackageOut) SetUnitIntPricePerMonthNil(b bool)`

 SetUnitIntPricePerMonthNil sets the value for UnitIntPricePerMonth to be an explicit nil

### UnsetUnitIntPricePerMonth
`func (o *PackageOut) UnsetUnitIntPricePerMonth()`

UnsetUnitIntPricePerMonth ensures that no value is present for UnitIntPricePerMonth, not even an explicit nil
### GetUnitIntPricePerYear

`func (o *PackageOut) GetUnitIntPricePerYear() int32`

GetUnitIntPricePerYear returns the UnitIntPricePerYear field if non-nil, zero value otherwise.

### GetUnitIntPricePerYearOk

`func (o *PackageOut) GetUnitIntPricePerYearOk() (*int32, bool)`

GetUnitIntPricePerYearOk returns a tuple with the UnitIntPricePerYear field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerYear

`func (o *PackageOut) SetUnitIntPricePerYear(v int32)`

SetUnitIntPricePerYear sets UnitIntPricePerYear field to given value.

### HasUnitIntPricePerYear

`func (o *PackageOut) HasUnitIntPricePerYear() bool`

HasUnitIntPricePerYear returns a boolean if a field has been set.

### SetUnitIntPricePerYearNil

`func (o *PackageOut) SetUnitIntPricePerYearNil(b bool)`

 SetUnitIntPricePerYearNil sets the value for UnitIntPricePerYear to be an explicit nil

### UnsetUnitIntPricePerYear
`func (o *PackageOut) UnsetUnitIntPricePerYear()`

UnsetUnitIntPricePerYear ensures that no value is present for UnitIntPricePerYear, not even an explicit nil
### GetRegionalMappingId

`func (o *PackageOut) GetRegionalMappingId() string`

GetRegionalMappingId returns the RegionalMappingId field if non-nil, zero value otherwise.

### GetRegionalMappingIdOk

`func (o *PackageOut) GetRegionalMappingIdOk() (*string, bool)`

GetRegionalMappingIdOk returns a tuple with the RegionalMappingId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionalMappingId

`func (o *PackageOut) SetRegionalMappingId(v string)`

SetRegionalMappingId sets RegionalMappingId field to given value.

### HasRegionalMappingId

`func (o *PackageOut) HasRegionalMappingId() bool`

HasRegionalMappingId returns a boolean if a field has been set.

### SetRegionalMappingIdNil

`func (o *PackageOut) SetRegionalMappingIdNil(b bool)`

 SetRegionalMappingIdNil sets the value for RegionalMappingId to be an explicit nil

### UnsetRegionalMappingId
`func (o *PackageOut) UnsetRegionalMappingId()`

UnsetRegionalMappingId ensures that no value is present for RegionalMappingId, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


