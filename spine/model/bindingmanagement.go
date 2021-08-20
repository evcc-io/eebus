// Package ns_p contains models for http://docs.eebus.org/spine/xsd/v1
package model

// Code generated by github.com/andig/xsd2go. DO NOT EDIT.

import "github.com/evcc-io/eebus/util"

// BindingIdType type
type BindingIdType uint

// BindingManagementEntryDataType complex type
type BindingManagementEntryDataType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
	Label         *LabelType          `json:"label,omitempty"`
	Description   *DescriptionType    `json:"description,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementEntryDataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementEntryDataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementEntryDataElementsType complex type
type BindingManagementEntryDataElementsType struct {
	BindingId     *ElementTagType             `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressElementsType `json:"serverAddress,omitempty"`
	Label         *ElementTagType             `json:"label,omitempty"`
	Description   *ElementTagType             `json:"description,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementEntryDataElementsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementEntryDataElementsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementEntryListDataType complex type
type BindingManagementEntryListDataType struct {
	BindingManagementEntryData []BindingManagementEntryDataType `json:"bindingManagementEntryData,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementEntryListDataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementEntryListDataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementEntryListDataSelectorsType complex type
type BindingManagementEntryListDataSelectorsType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementEntryListDataSelectorsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementEntryListDataSelectorsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementRequestCallType complex type
type BindingManagementRequestCallType struct {
	ClientAddress     *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress     *FeatureAddressType `json:"serverAddress,omitempty"`
	ServerFeatureType *FeatureTypeType    `json:"serverFeatureType,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementRequestCallType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementRequestCallType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementRequestCallElementsType complex type
type BindingManagementRequestCallElementsType struct {
	ClientAddress     *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress     *FeatureAddressElementsType `json:"serverAddress,omitempty"`
	ServerFeatureType *ElementTagType             `json:"serverFeatureType,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementRequestCallElementsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementRequestCallElementsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementDeleteCallType complex type
type BindingManagementDeleteCallType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementDeleteCallType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementDeleteCallType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// BindingManagementDeleteCallElementsType complex type
type BindingManagementDeleteCallElementsType struct {
	BindingId     *ElementTagType             `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressElementsType `json:"serverAddress,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m BindingManagementDeleteCallElementsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *BindingManagementDeleteCallElementsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
