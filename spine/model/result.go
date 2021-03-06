// Package ns_p contains models for http://docs.eebus.org/spine/xsd/v1
package model

// Code generated by github.com/andig/xsd2go. DO NOT EDIT.

import "github.com/evcc-io/eebus/util"

// ErrorNumberType type
type ErrorNumberType uint

// ResultDataType complex type
type ResultDataType struct {
	ErrorNumber *ErrorNumberType `json:"errorNumber,omitempty"`
	Description *DescriptionType `json:"description,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ResultDataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ResultDataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
