/*
Qdrant Cloud API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package qdrant_cloud_programmatic_access

import (
	"encoding/json"
	"fmt"
)

// Currency the model 'Currency'
type Currency string

// List of Currency
const (
	CURRENCY_USD Currency = "usd"
	CURRENCY_EUR Currency = "eur"
)

// All allowed values of Currency enum
var AllowedCurrencyEnumValues = []Currency{
	"usd",
	"eur",
}

func (v *Currency) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Currency(value)
	for _, existing := range AllowedCurrencyEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Currency", value)
}

// NewCurrencyFromValue returns a pointer to a valid Currency
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewCurrencyFromValue(v string) (*Currency, error) {
	ev := Currency(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Currency: valid values are %v", v, AllowedCurrencyEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Currency) IsValid() bool {
	for _, existing := range AllowedCurrencyEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Currency value
func (v Currency) Ptr() *Currency {
	return &v
}

type NullableCurrency struct {
	value *Currency
	isSet bool
}

func (v NullableCurrency) Get() *Currency {
	return v.value
}

func (v *NullableCurrency) Set(val *Currency) {
	v.value = val
	v.isSet = true
}

func (v NullableCurrency) IsSet() bool {
	return v.isSet
}

func (v *NullableCurrency) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCurrency(val *Currency) *NullableCurrency {
	return &NullableCurrency{value: val, isSet: true}
}

func (v NullableCurrency) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCurrency) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
