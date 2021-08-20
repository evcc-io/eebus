// Package ship contains models for http://docs.eebus.org/ship/xsd/v1
package ship

import (
	"encoding/json"

	"github.com/evcc-io/eebus/util"
)

// xs:element declarations

// CmiConnectionHello message container
type CmiConnectionHello struct {
	ConnectionHello ConnectionHello `json:"connectionHello"`
}

// ConnectionHello element
type ConnectionHello struct {
	Phase               ConnectionHelloPhaseType `json:"phase"`
	Waiting             *uint                    `json:"waiting,omitempty"`
	ProlongationRequest *bool                    `json:"prolongationRequest,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionHello) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionHello) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiMessageProtocolHandshake message container
type CmiMessageProtocolHandshake struct {
	MessageProtocolHandshake MessageProtocolHandshake `json:"messageProtocolHandshake"`
}

// MessageProtocolHandshake element
type MessageProtocolHandshake struct {
	HandshakeType ProtocolHandshakeTypeType  `json:"handshakeType"`
	Version       Version                    `json:"version"`
	Formats       MessageProtocolFormatsType `json:"formats"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolHandshake) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolHandshake) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiMessageProtocolHandshakeError message container
type CmiMessageProtocolHandshakeError struct {
	MessageProtocolHandshakeError MessageProtocolHandshakeError `json:"messageProtocolHandshakeError"`
}

// MessageProtocolHandshakeError element
type MessageProtocolHandshakeError struct {
	Error MessageProtocolHandshakeErrorErrorType `json:"error"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolHandshakeError) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolHandshakeError) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiConnectionPinState message container
type CmiConnectionPinState struct {
	ConnectionPinState ConnectionPinState `json:"connectionPinState"`
}

// ConnectionPinState element
type ConnectionPinState struct {
	PinState        PinStateType            `json:"pinState"`
	InputPermission *PinInputPermissionType `json:"inputPermission,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinState) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinState) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiConnectionPinInput message container
type CmiConnectionPinInput struct {
	ConnectionPinInput ConnectionPinInput `json:"connectionPinInput"`
}

// ConnectionPinInput element
type ConnectionPinInput struct {
	Pin PinValueType `json:"pin"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinInput) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinInput) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiConnectionPinError message container
type CmiConnectionPinError struct {
	ConnectionPinError ConnectionPinError `json:"connectionPinError"`
}

// ConnectionPinError element
type ConnectionPinError struct {
	Error ConnectionPinErrorErrorType `json:"error"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinError) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinError) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiData message container
type CmiData struct {
	Data Data `json:"data"`
}

// Data element
type Data struct {
	Header    HeaderType      `json:"header"`
	Payload   json.RawMessage `json:"payload"`
	Extension *ExtensionType  `json:"extension,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Data) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Data) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiConnectionClose message container
type CmiConnectionClose struct {
	ConnectionClose ConnectionClose `json:"connectionClose"`
}

// ConnectionClose element
type ConnectionClose struct {
	Phase   ConnectionClosePhaseType   `json:"phase"`
	MaxTime *uint                      `json:"maxTime,omitempty"`
	Reason  *ConnectionCloseReasonType `json:"reason,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionClose) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionClose) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiAccessMethodsRequest message container
type CmiAccessMethodsRequest struct {
	AccessMethodsRequest AccessMethodsRequest `json:"accessMethodsRequest"`
}

// AccessMethodsRequest element
type AccessMethodsRequest struct {
}

// MarshalJSON is the SHIP serialization marshaller
func (m AccessMethodsRequest) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *AccessMethodsRequest) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiAccessMethods message container
type CmiAccessMethods struct {
	AccessMethods AccessMethods `json:"accessMethods"`
}

// AccessMethods element
type AccessMethods struct {
	Id        string     `json:"id"`
	DnsSdMDns *DnsSdMDns `json:"dnsSd_mDns,omitempty"`
	Dns       *Dns       `json:"dns,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m AccessMethods) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *AccessMethods) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiVersion message container
type CmiVersion struct {
	Version Version `json:"version"`
}

// Version element
type Version struct {
	Major uint8 `json:"major"`
	Minor uint8 `json:"minor"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Version) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Version) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiDnsSdMDns message container
type CmiDnsSdMDns struct {
	DnsSdMDns DnsSdMDns `json:"dnsSd_mDns"`
}

// DnsSdMDns element
type DnsSdMDns struct {
}

// MarshalJSON is the SHIP serialization marshaller
func (m DnsSdMDns) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *DnsSdMDns) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// CmiDns message container
type CmiDns struct {
	Dns Dns `json:"dns"`
}

// Dns element
type Dns struct {
	Uri string `json:"uri"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Dns) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Dns) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// xs:complexType declarations

// ConnectionHelloType complex type
type ConnectionHelloType struct {
	Phase               ConnectionHelloPhaseType `json:"phase"`
	Waiting             *uint                    `json:"waiting"`
	ProlongationRequest *bool                    `json:"prolongationRequest"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionHelloType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionHelloType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// MessageProtocolFormatsType complex type
type MessageProtocolFormatsType struct {
	Format []MessageProtocolFormatType `json:"format"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolFormatsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolFormatsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// MessageProtocolHandshakeType complex type
type MessageProtocolHandshakeType struct {
	HandshakeType ProtocolHandshakeTypeType  `json:"handshakeType"`
	Version       Version                    `json:"version"`
	Formats       MessageProtocolFormatsType `json:"formats"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolHandshakeType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolHandshakeType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// MessageProtocolHandshakeErrorType complex type
type MessageProtocolHandshakeErrorType struct {
	Error MessageProtocolHandshakeErrorErrorType `json:"error"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolHandshakeErrorType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolHandshakeErrorType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// ConnectionPinStateType complex type
type ConnectionPinStateType struct {
	PinState        PinStateType            `json:"pinState"`
	InputPermission *PinInputPermissionType `json:"inputPermission"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinStateType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinStateType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// ConnectionPinInputType complex type
type ConnectionPinInputType struct {
	Pin PinValueType `json:"pin"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinInputType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinInputType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// ConnectionPinErrorType complex type
type ConnectionPinErrorType struct {
	Error ConnectionPinErrorErrorType `json:"error"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinErrorType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinErrorType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// HeaderType complex type
type HeaderType struct {
	ProtocolId ProtocolIdType `json:"protocolId"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m HeaderType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *HeaderType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// ExtensionType complex type
type ExtensionType struct {
	ExtensionId string `json:"extensionId"`
	Binary      *byte  `json:"binary"`
	String      string `json:"string"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ExtensionType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ExtensionType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// DataType complex type
type DataType struct {
	Header    HeaderType     `json:"header"`
	Payload   string         `json:"payload"`
	Extension *ExtensionType `json:"extension"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m DataType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *DataType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// ConnectionCloseType complex type
type ConnectionCloseType struct {
	Phase   ConnectionClosePhaseType   `json:"phase"`
	MaxTime *uint                      `json:"maxTime"`
	Reason  *ConnectionCloseReasonType `json:"reason"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionCloseType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionCloseType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// AccessMethodsRequestType complex type
type AccessMethodsRequestType struct {
}

// MarshalJSON is the SHIP serialization marshaller
func (m AccessMethodsRequestType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *AccessMethodsRequestType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// AccessMethodsType complex type
type AccessMethodsType struct {
	Id        string     `json:"id"`
	DnsSdMDns *DnsSdMDns `json:"dnsSd_mDns"`
	Dns       *Dns       `json:"dns"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m AccessMethodsType) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *AccessMethodsType) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// xs:simpleType declarations

// ConnectionHelloPhaseType type
type ConnectionHelloPhaseType string

// ConnectionHelloPhaseType constants
const (
	ConnectionHelloPhaseTypePending ConnectionHelloPhaseType = "pending"
	ConnectionHelloPhaseTypeReady   ConnectionHelloPhaseType = "ready"
	ConnectionHelloPhaseTypeAborted ConnectionHelloPhaseType = "aborted"
)

// MessageProtocolFormatType type
type MessageProtocolFormatType string

// ProtocolHandshakeTypeType type
type ProtocolHandshakeTypeType string

// ProtocolHandshakeTypeType constants
const (
	ProtocolHandshakeTypeTypeAnnouncemax ProtocolHandshakeTypeType = "announceMax"
	ProtocolHandshakeTypeTypeSelect      ProtocolHandshakeTypeType = "select"
)

// MessageProtocolHandshakeErrorErrorType type
type MessageProtocolHandshakeErrorErrorType string

// PinStateType type
type PinStateType string

// PinStateType constants
const (
	PinStateTypeRequired PinStateType = "required"
	PinStateTypeOptional PinStateType = "optional"
	PinStateTypePinok    PinStateType = "pinOk"
	PinStateTypeNone     PinStateType = "none"
)

// PinInputPermissionType type
type PinInputPermissionType string

// PinInputPermissionType constants
const (
	PinInputPermissionTypeBusy PinInputPermissionType = "busy"
	PinInputPermissionTypeOk   PinInputPermissionType = "ok"
)

// PinValueType type
type PinValueType string

// ConnectionPinErrorErrorType type
type ConnectionPinErrorErrorType string

// ProtocolIdType type
type ProtocolIdType string

// ConnectionClosePhaseType type
type ConnectionClosePhaseType string

// ConnectionClosePhaseType constants
const (
	ConnectionClosePhaseTypeAnnounce ConnectionClosePhaseType = "announce"
	ConnectionClosePhaseTypeConfirm  ConnectionClosePhaseType = "confirm"
)

// ConnectionCloseReasonType type
type ConnectionCloseReasonType string

// ConnectionCloseReasonType constants
const (
	ConnectionCloseReasonTypeUnspecific        ConnectionCloseReasonType = "unspecific"
	ConnectionCloseReasonTypeRemovedconnection ConnectionCloseReasonType = "removedConnection"
)
