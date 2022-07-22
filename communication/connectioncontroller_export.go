package communication

import (
	"errors"
	"time"
)

type ManufacturerDetails struct {
	DeviceName    string
	DeviceCode    string
	DeviceAddress string
	BrandName     string
}

type EVSEOperationStateEnumType string

const (
	EVSEOperationStateEnumTypeNormal  = "normal"
	EVSEOperationStateEnumTypeFailure = "failure"
)

type EVChargeStateEnumType string

const (
	EVChargeStateEnumTypeUnknown   = "unknown"
	EVChargeStateEnumTypeUnplugged = "unplugged"
	EVChargeStateEnumTypeError     = "error"
	EVChargeStateEnumTypePaused    = "paused"
	EVChargeStateEnumTypeActive    = "active"
	EVChargeStateEnumTypeFinished  = "finished"
)

type EVCommunicationStandardEnumType string

const (
	EVCommunicationStandardEnumTypeUnknown      EVCommunicationStandardEnumType = "unknown"
	EVCommunicationStandardEnumTypeISO151182ED1 EVCommunicationStandardEnumType = "iso15118-2ed1"
	EVCommunicationStandardEnumTypeISO151182ED2 EVCommunicationStandardEnumType = "iso15118-2ed2"
	EVCommunicationStandardEnumTypeIEC61851     EVCommunicationStandardEnumType = "iec61851"
)

type EVSEDataType struct {
	Manufacturer   ManufacturerDetails
	OperationState EVSEOperationStateEnumType
}

type EVSEClientDataType struct {
	EVSEData EVSEDataType
	EVData   EVDataType
}

type EVMeasurementsType struct {
	Timestamp     time.Time
	Current       map[uint]float64
	Power         map[uint]float64
	ChargedEnergy float64
	SoC           float64
}

type EVCurrentLimitType struct {
	Min, Max, Default float64
}

type EVPowerLimitType struct {
	Min, Max float64
}

type EVChargingStrategyEnumType string

const (
	EVChargingStrategyEnumTypeUnknown        EVChargingStrategyEnumType = "unknown"
	EVChargingStrategyEnumTypeNoDemand       EVChargingStrategyEnumType = "nodemand"
	EVChargingStrategyEnumTypeDirectCharging EVChargingStrategyEnumType = "directcharging"
	EVChargingStrategyEnumTypeTimedCharging  EVChargingStrategyEnumType = "timedcharging"
)

type EVDataType struct {
	UCSelfConsumptionAvailable     bool
	UCCoordinatedChargingAvailable bool
	UCSoCAvailable                 bool
	AsymetricChargingSupported     bool
	CommunicationStandard          EVCommunicationStandardEnumType
	OverloadProtectionActive       bool
	SelfConsumptionActive          bool
	SoCDataAvailable               bool
	ConnectedPhases                uint
	ChargingStrategy               EVChargingStrategyEnumType
	ChargingDemand                 float64
	ChargingTargetDuration         time.Duration
	Manufacturer                   ManufacturerDetails
	Identification                 string
	ChargeState                    EVChargeStateEnumType
	Limits                         map[uint]EVCurrentLimitType
	LimitsPower                    EVPowerLimitType
	Measurements                   EVMeasurementsType
}

type EVDataElementUpdateType string

const (
	EVDataElementUpdateUseCaseSelfConsumption     EVDataElementUpdateType = "usecaseselfconsumption"
	EVDataElementUpdateUseCaseSoC                 EVDataElementUpdateType = "usecasesoc"
	EVDataElementUpdateUseCaseCoordinatedCharging EVDataElementUpdateType = "usecasecoordinatedcharging"
	EVDataElementUpdateEVConnectionState          EVDataElementUpdateType = "evconnectionstate"
	EVDataElementUpdateCommunicationStandard      EVDataElementUpdateType = "communicationstandard"
	EVDataElementUpdateAsymetricChargingType      EVDataElementUpdateType = "asymetricchargingtype"
	EVDataElementUpdateEVSEOperationState         EVDataElementUpdateType = "evseoperationstate"
	EVDataElementUpdateEVChargeState              EVDataElementUpdateType = "evchargestate"
	EVDataElementUpdateConnectedPhases            EVDataElementUpdateType = "connectedphases"
	EVDataElementUpdatePowerLimits                EVDataElementUpdateType = "powerlimits"
	EVDataElementUpdateAmperageLimits             EVDataElementUpdateType = "amperagelimits"
	EVDataElementUpdateChargingStrategy           EVDataElementUpdateType = "chargingstrategy"
	EVDataElementUpdateChargingPlanRequired       EVDataElementUpdateType = "chargingplanrequired"
)

type EVChargingSlot struct {
	Duration time.Duration
	MaxValue float64 // Watts
	Pricing  float64
}

type EVChargingPlan struct {
	Duration time.Duration
	Slots    []EVChargingSlot
}

func (c *ConnectionController) GetData() (*EVSEClientDataType, error) {
	if c == nil {
		return nil, errors.New("offline")
	}

	return c.clientData, nil
}

func (c *ConnectionController) SetDataUpdateHandler(dataUpdateHandler func(EVDataElementUpdateType, *EVSEClientDataType)) {
	c.dataUpdateHandler = dataUpdateHandler
}

func (c *ConnectionController) callDataUpdateHandler(updateType EVDataElementUpdateType) {
	if c.dataUpdateHandler != nil {
		c.dataUpdateHandler(updateType, c.clientData)
	}
}
