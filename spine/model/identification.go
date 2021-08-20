// Package ns_p contains models for http://docs.eebus.org/spine/xsd/v1
package model

// Code generated by github.com/andig/xsd2go. DO NOT EDIT.

import "github.com/evcc-io/eebus/util"

// IdentificationIdType type
type IdentificationIdType uint

// IdentificationTypeType type
type IdentificationTypeType IdentificationTypeEnumType

// IdentificationTypeEnumType type
type IdentificationTypeEnumType string

// IdentificationTypeEnumType constants
const (
	IdentificationTypeEnumTypeEui48       IdentificationTypeEnumType = "eui48"
	IdentificationTypeEnumTypeEui64       IdentificationTypeEnumType = "eui64"
	IdentificationTypeEnumTypeUserrfidtag IdentificationTypeEnumType = "userRfidTag"
)

// IdentificationValueType type
type IdentificationValueType string

// IdentificationDataType complex type
type IdentificationDataType struct {
	IdentificationId    *IdentificationIdType    `json:"identificationId,omitempty"`
	IdentificationType  *IdentificationTypeType  `json:"identificationType,omitempty"`
	IdentificationValue *IdentificationValueType `json:"identificationValue,omitempty"`
	Authorized          *bool                    `json:"authorized,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m IdentificationDataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *IdentificationDataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// IdentificationDataElementsType complex type
type IdentificationDataElementsType struct {
	IdentificationId    *ElementTagType `json:"identificationId,omitempty"`
	IdentificationType  *ElementTagType `json:"identificationType,omitempty"`
	IdentificationValue *ElementTagType `json:"identificationValue,omitempty"`
	Authorized          *ElementTagType `json:"authorized,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m IdentificationDataElementsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *IdentificationDataElementsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// IdentificationListDataType complex type
type IdentificationListDataType struct {
	IdentificationData []IdentificationDataType `json:"identificationData,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m IdentificationListDataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *IdentificationListDataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// IdentificationListDataSelectorsType complex type
type IdentificationListDataSelectorsType struct {
	IdentificationId   *IdentificationIdType   `json:"identificationId,omitempty"`
	IdentificationType *IdentificationTypeType `json:"identificationType,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m IdentificationListDataSelectorsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *IdentificationListDataSelectorsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
