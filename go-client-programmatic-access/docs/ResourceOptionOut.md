# ResourceOptionOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**ResourceType** | [**ResourceType**](ResourceType.md) |  | 
**Status** | [**BookingStatus**](BookingStatus.md) |  | 
**Name** | Pointer to **NullableString** |  | [optional] 
**ResourceUnit** | **string** |  | 
**Currency** | [**Currency**](Currency.md) |  | 
**UnitIntPricePerHour** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerDay** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerMonth** | Pointer to **NullableInt32** |  | [optional] 
**UnitIntPricePerYear** | Pointer to **NullableInt32** |  | [optional] 

## Methods

### NewResourceOptionOut

`func NewResourceOptionOut(id string, resourceType ResourceType, status BookingStatus, resourceUnit string, currency Currency, ) *ResourceOptionOut`

NewResourceOptionOut instantiates a new ResourceOptionOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewResourceOptionOutWithDefaults

`func NewResourceOptionOutWithDefaults() *ResourceOptionOut`

NewResourceOptionOutWithDefaults instantiates a new ResourceOptionOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ResourceOptionOut) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ResourceOptionOut) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ResourceOptionOut) SetId(v string)`

SetId sets Id field to given value.


### GetResourceType

`func (o *ResourceOptionOut) GetResourceType() ResourceType`

GetResourceType returns the ResourceType field if non-nil, zero value otherwise.

### GetResourceTypeOk

`func (o *ResourceOptionOut) GetResourceTypeOk() (*ResourceType, bool)`

GetResourceTypeOk returns a tuple with the ResourceType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceType

`func (o *ResourceOptionOut) SetResourceType(v ResourceType)`

SetResourceType sets ResourceType field to given value.


### GetStatus

`func (o *ResourceOptionOut) GetStatus() BookingStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ResourceOptionOut) GetStatusOk() (*BookingStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ResourceOptionOut) SetStatus(v BookingStatus)`

SetStatus sets Status field to given value.


### GetName

`func (o *ResourceOptionOut) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ResourceOptionOut) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ResourceOptionOut) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ResourceOptionOut) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *ResourceOptionOut) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *ResourceOptionOut) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetResourceUnit

`func (o *ResourceOptionOut) GetResourceUnit() string`

GetResourceUnit returns the ResourceUnit field if non-nil, zero value otherwise.

### GetResourceUnitOk

`func (o *ResourceOptionOut) GetResourceUnitOk() (*string, bool)`

GetResourceUnitOk returns a tuple with the ResourceUnit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceUnit

`func (o *ResourceOptionOut) SetResourceUnit(v string)`

SetResourceUnit sets ResourceUnit field to given value.


### GetCurrency

`func (o *ResourceOptionOut) GetCurrency() Currency`

GetCurrency returns the Currency field if non-nil, zero value otherwise.

### GetCurrencyOk

`func (o *ResourceOptionOut) GetCurrencyOk() (*Currency, bool)`

GetCurrencyOk returns a tuple with the Currency field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrency

`func (o *ResourceOptionOut) SetCurrency(v Currency)`

SetCurrency sets Currency field to given value.


### GetUnitIntPricePerHour

`func (o *ResourceOptionOut) GetUnitIntPricePerHour() int32`

GetUnitIntPricePerHour returns the UnitIntPricePerHour field if non-nil, zero value otherwise.

### GetUnitIntPricePerHourOk

`func (o *ResourceOptionOut) GetUnitIntPricePerHourOk() (*int32, bool)`

GetUnitIntPricePerHourOk returns a tuple with the UnitIntPricePerHour field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerHour

`func (o *ResourceOptionOut) SetUnitIntPricePerHour(v int32)`

SetUnitIntPricePerHour sets UnitIntPricePerHour field to given value.

### HasUnitIntPricePerHour

`func (o *ResourceOptionOut) HasUnitIntPricePerHour() bool`

HasUnitIntPricePerHour returns a boolean if a field has been set.

### SetUnitIntPricePerHourNil

`func (o *ResourceOptionOut) SetUnitIntPricePerHourNil(b bool)`

 SetUnitIntPricePerHourNil sets the value for UnitIntPricePerHour to be an explicit nil

### UnsetUnitIntPricePerHour
`func (o *ResourceOptionOut) UnsetUnitIntPricePerHour()`

UnsetUnitIntPricePerHour ensures that no value is present for UnitIntPricePerHour, not even an explicit nil
### GetUnitIntPricePerDay

`func (o *ResourceOptionOut) GetUnitIntPricePerDay() int32`

GetUnitIntPricePerDay returns the UnitIntPricePerDay field if non-nil, zero value otherwise.

### GetUnitIntPricePerDayOk

`func (o *ResourceOptionOut) GetUnitIntPricePerDayOk() (*int32, bool)`

GetUnitIntPricePerDayOk returns a tuple with the UnitIntPricePerDay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerDay

`func (o *ResourceOptionOut) SetUnitIntPricePerDay(v int32)`

SetUnitIntPricePerDay sets UnitIntPricePerDay field to given value.

### HasUnitIntPricePerDay

`func (o *ResourceOptionOut) HasUnitIntPricePerDay() bool`

HasUnitIntPricePerDay returns a boolean if a field has been set.

### SetUnitIntPricePerDayNil

`func (o *ResourceOptionOut) SetUnitIntPricePerDayNil(b bool)`

 SetUnitIntPricePerDayNil sets the value for UnitIntPricePerDay to be an explicit nil

### UnsetUnitIntPricePerDay
`func (o *ResourceOptionOut) UnsetUnitIntPricePerDay()`

UnsetUnitIntPricePerDay ensures that no value is present for UnitIntPricePerDay, not even an explicit nil
### GetUnitIntPricePerMonth

`func (o *ResourceOptionOut) GetUnitIntPricePerMonth() int32`

GetUnitIntPricePerMonth returns the UnitIntPricePerMonth field if non-nil, zero value otherwise.

### GetUnitIntPricePerMonthOk

`func (o *ResourceOptionOut) GetUnitIntPricePerMonthOk() (*int32, bool)`

GetUnitIntPricePerMonthOk returns a tuple with the UnitIntPricePerMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerMonth

`func (o *ResourceOptionOut) SetUnitIntPricePerMonth(v int32)`

SetUnitIntPricePerMonth sets UnitIntPricePerMonth field to given value.

### HasUnitIntPricePerMonth

`func (o *ResourceOptionOut) HasUnitIntPricePerMonth() bool`

HasUnitIntPricePerMonth returns a boolean if a field has been set.

### SetUnitIntPricePerMonthNil

`func (o *ResourceOptionOut) SetUnitIntPricePerMonthNil(b bool)`

 SetUnitIntPricePerMonthNil sets the value for UnitIntPricePerMonth to be an explicit nil

### UnsetUnitIntPricePerMonth
`func (o *ResourceOptionOut) UnsetUnitIntPricePerMonth()`

UnsetUnitIntPricePerMonth ensures that no value is present for UnitIntPricePerMonth, not even an explicit nil
### GetUnitIntPricePerYear

`func (o *ResourceOptionOut) GetUnitIntPricePerYear() int32`

GetUnitIntPricePerYear returns the UnitIntPricePerYear field if non-nil, zero value otherwise.

### GetUnitIntPricePerYearOk

`func (o *ResourceOptionOut) GetUnitIntPricePerYearOk() (*int32, bool)`

GetUnitIntPricePerYearOk returns a tuple with the UnitIntPricePerYear field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnitIntPricePerYear

`func (o *ResourceOptionOut) SetUnitIntPricePerYear(v int32)`

SetUnitIntPricePerYear sets UnitIntPricePerYear field to given value.

### HasUnitIntPricePerYear

`func (o *ResourceOptionOut) HasUnitIntPricePerYear() bool`

HasUnitIntPricePerYear returns a boolean if a field has been set.

### SetUnitIntPricePerYearNil

`func (o *ResourceOptionOut) SetUnitIntPricePerYearNil(b bool)`

 SetUnitIntPricePerYearNil sets the value for UnitIntPricePerYear to be an explicit nil

### UnsetUnitIntPricePerYear
`func (o *ResourceOptionOut) UnsetUnitIntPricePerYear()`

UnsetUnitIntPricePerYear ensures that no value is present for UnitIntPricePerYear, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


