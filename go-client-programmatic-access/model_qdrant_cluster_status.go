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

// QdrantClusterStatus the model 'QdrantClusterStatus'
type QdrantClusterStatus string

// List of QdrantClusterStatus
const (
	QDRANTCLUSTERSTATUS_SUSPENDED QdrantClusterStatus = "suspended"
	QDRANTCLUSTERSTATUS_RUNNING QdrantClusterStatus = "running"
)

// All allowed values of QdrantClusterStatus enum
var AllowedQdrantClusterStatusEnumValues = []QdrantClusterStatus{
	"suspended",
	"running",
}

func (v *QdrantClusterStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := QdrantClusterStatus(value)
	for _, existing := range AllowedQdrantClusterStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid QdrantClusterStatus", value)
}

// NewQdrantClusterStatusFromValue returns a pointer to a valid QdrantClusterStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewQdrantClusterStatusFromValue(v string) (*QdrantClusterStatus, error) {
	ev := QdrantClusterStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for QdrantClusterStatus: valid values are %v", v, AllowedQdrantClusterStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v QdrantClusterStatus) IsValid() bool {
	for _, existing := range AllowedQdrantClusterStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to QdrantClusterStatus value
func (v QdrantClusterStatus) Ptr() *QdrantClusterStatus {
	return &v
}

type NullableQdrantClusterStatus struct {
	value *QdrantClusterStatus
	isSet bool
}

func (v NullableQdrantClusterStatus) Get() *QdrantClusterStatus {
	return v.value
}

func (v *NullableQdrantClusterStatus) Set(val *QdrantClusterStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableQdrantClusterStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableQdrantClusterStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQdrantClusterStatus(val *QdrantClusterStatus) *NullableQdrantClusterStatus {
	return &NullableQdrantClusterStatus{value: val, isSet: true}
}

func (v NullableQdrantClusterStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQdrantClusterStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
