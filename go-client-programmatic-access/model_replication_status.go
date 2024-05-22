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

// ReplicationStatus the model 'ReplicationStatus'
type ReplicationStatus string

// List of ReplicationStatus
const (
	REPLICATIONSTATUS_DOWN ReplicationStatus = "down"
	REPLICATIONSTATUS_UNDER_REPLICATED ReplicationStatus = "under_replicated"
	REPLICATIONSTATUS_HEALTHY ReplicationStatus = "healthy"
)

// All allowed values of ReplicationStatus enum
var AllowedReplicationStatusEnumValues = []ReplicationStatus{
	"down",
	"under_replicated",
	"healthy",
}

func (v *ReplicationStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ReplicationStatus(value)
	for _, existing := range AllowedReplicationStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ReplicationStatus", value)
}

// NewReplicationStatusFromValue returns a pointer to a valid ReplicationStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewReplicationStatusFromValue(v string) (*ReplicationStatus, error) {
	ev := ReplicationStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ReplicationStatus: valid values are %v", v, AllowedReplicationStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ReplicationStatus) IsValid() bool {
	for _, existing := range AllowedReplicationStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ReplicationStatus value
func (v ReplicationStatus) Ptr() *ReplicationStatus {
	return &v
}

type NullableReplicationStatus struct {
	value *ReplicationStatus
	isSet bool
}

func (v NullableReplicationStatus) Get() *ReplicationStatus {
	return v.value
}

func (v *NullableReplicationStatus) Set(val *ReplicationStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicationStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicationStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicationStatus(val *ReplicationStatus) *NullableReplicationStatus {
	return &NullableReplicationStatus{value: val, isSet: true}
}

func (v NullableReplicationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicationStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

