package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/evcc-io/eebus/device"
	"github.com/evcc-io/eebus/device/feature"
	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
	"github.com/evcc-io/eebus/util"
	"github.com/rickb777/date/period"
)

type ConnectionController struct {
	msgNum               uint64 // 64bit values need to be defined on top of the struct to make atomic commands work on 32bit systems
	heartBeatNum         uint64 // see https://github.com/golang/go/issues/11891
	subscriptionNum      uint64
	log                  util.Logger
	conn                 ship.Conn
	localDevice          spine.Device
	remoteDevice         spine.Device // TODO multiple remote devices
	sequencesController  *SequencesController
	stopMux              sync.Mutex
	stopHeartbeatC       chan struct{}
	subscriptionEntries  []model.SubscriptionManagementEntryDataType
	specificationVersion model.SpecificationVersionType
	// EV specific data
	clientData *EVSEClientDataType
	// EVCC specific
	dataUpdateHandler func(EVDataElementUpdateType, *EVSEClientDataType)
}

func NewConnectionController(log util.Logger, conn ship.Conn, local spine.Device) *ConnectionController {
	c := &ConnectionController{
		log:                  log,
		conn:                 conn,
		specificationVersion: device.SpecificationVersion,
		localDevice:          local,
		clientData:           &EVSEClientDataType{},
		sequencesController:  NewSequencesController(log),
	}

	return c
}

func (c *ConnectionController) Boot() error {
	c.clientData.EVData.CommunicationStandard = EVCommunicationStandardEnumTypeUnknown
	c.clientData.EVData.ChargeState = EVChargeStateEnumTypeUnknown

	m := c.localDevice.Entity([]model.AddressEntityType{0}).FeatureByProps(model.FeatureTypeEnumTypeNodeManagement, model.RoleTypeSpecial)
	if f, ok := m.(*feature.NodeManagement); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeMeasurement, model.RoleTypeClient)
	if f, ok := m.(*feature.Measurement); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeElectricalConnection, model.RoleTypeClient)
	if f, ok := m.(*feature.ElectricalConnection); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeDeviceConfiguration, model.RoleTypeClient)
	if f, ok := m.(*feature.DeviceConfiguration); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeDeviceClassification, model.RoleTypeClient)
	if f, ok := m.(*feature.DeviceClassification); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeDeviceDiagnosis, model.RoleTypeClient)
	if f, ok := m.(*feature.DeviceDiagnosis); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeIdentification, model.RoleTypeClient)
	if f, ok := m.(*feature.Identification); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeClient)
	if f, ok := m.(*feature.LoadControl); ok {
		f.Delegate = c
	}

	m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeTimeSeries, model.RoleTypeClient)
	if f, ok := m.(*feature.TimeSeries); ok {
		f.Delegate = c
	}

	// m = c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeIncentiveTable, model.RoleTypeClient)
	// if f, ok := m.(*feature.IncentiveTable); ok {
	// 	f.Delegate = c
	// }

	c.sequencesController.Boot()

	go c.Run()

	err := c.requestNodeManagementDetailedDiscoveryData()
	if err != nil {
		c.log.Println("Sending DetailedDiscoveryData read request failed!")
	}

	return err
}

func (c *ConnectionController) CloseConnection(err error) {
	c.stopHeartbeat()
	_ = c.conn.Close()
}

func (c *ConnectionController) Run() {
	var err error
	for err == nil {
		if c.conn == nil || c.conn.IsConnectionClosed() {
			err = errors.New("connection closed")
			break
		}

		var data json.RawMessage
		data, err = c.conn.Read()

		if err == nil {
			var datagram model.CmiDatagramType

			if err = json.Unmarshal(data, &datagram); err != nil {
				c.log.Println("error unmarshaling datagram: ", err, string(data))
				err = nil // don't break, otherwise charing will go to max limit
				continue
			}

			if err = c.processDatagram(datagram.Datagram); err != nil {
				c.log.Println("error processing datagram: ", err)
				err = nil // don't break, otherwise charing will go to max limit
				continue
			}
		}
	}

	if err != nil {
		c.log.Println("error processing incoming message: ", err)
	}

	c.stopHeartbeat()
	_ = c.conn.Close()
}

// Feature specific

func (c *ConnectionController) UpdateDeviceConfigurationData(f *feature.DeviceConfiguration, data []feature.DeviceConfigurationDatasetDataType) {
	var comStandard EVCommunicationStandardEnumType = EVCommunicationStandardEnumTypeIEC61851
	asymtricSupport := false

	for _, item := range data {
		switch item.KeyName {
		case model.DeviceConfigurationKeyNameEnumTypeAsymmetricChargingSupported:
			if item.KeyValueType == model.DeviceConfigurationKeyValueTypeTypeBoolean {
				asymtricSupport = item.KeyValueBoolean
			}
		case model.DeviceConfigurationKeyNameEnumTypeCommunicationsStandard:
			if item.KeyValueType == model.DeviceConfigurationKeyValueTypeTypeString {
				comStandard = EVCommunicationStandardEnumType(item.KeyValueString)
				// TODO make sure only a valid value is used, otherwise set default to IEC61851
			}
		default:
		}
	}

	if c.clientData.EVData.AsymetricChargingSupported != asymtricSupport {
		c.clientData.EVData.AsymetricChargingSupported = asymtricSupport
		c.callDataUpdateHandler(EVDataElementUpdateAsymetricChargingType)
	}
	if c.clientData.EVData.CommunicationStandard != comStandard {
		c.clientData.EVData.CommunicationStandard = comStandard
		c.callDataUpdateHandler(EVDataElementUpdateCommunicationStandard)
	}

	c.log.Println("asymetric charging supported: ", asymtricSupport, ", communication standard: ", comStandard)
}

func (c *ConnectionController) UpdateDeviceDiagnosisData(f *feature.DeviceDiagnosis, rf model.FeatureAddressType, data feature.DeviceDiagnosisDataType) {
	re := c.remoteEntityForFeatureAddress(rf)
	entityType := model.EntityTypeEnumType(re.GetType())

	prevEVSEOperationState := c.clientData.EVSEData.OperationState
	prevEVChargeState := c.clientData.EVData.ChargeState

	switch entityType {
	case model.EntityTypeEnumTypeEVSE:
		c.clientData.EVSEData.OperationState = EVSEOperationStateEnumTypeNormal
		c.log.Println("operation state EVSE: ", data.OperationState)
		switch data.OperationState {
		case model.DeviceDiagnosisOperatingStateEnumTypeFailure:
			c.clientData.EVSEData.OperationState = EVSEOperationStateEnumTypeFailure
		}

	case model.EntityTypeEnumTypeEV:
		c.clientData.EVData.ChargeState = EVChargeStateEnumTypeUnknown
		c.log.Println("charge state EV: ", data.OperationState)
		switch data.OperationState {
		case model.DeviceDiagnosisOperatingStateEnumTypeNormalOperation:
			c.clientData.EVData.ChargeState = EVChargeStateEnumTypeActive
		case model.DeviceDiagnosisOperatingStateEnumTypeStandby:
			c.clientData.EVData.ChargeState = EVChargeStateEnumTypePaused
		case model.DeviceDiagnosisOperatingStateEnumTypeFailure:
			c.clientData.EVData.ChargeState = EVChargeStateEnumTypeError
		case model.DeviceDiagnosisOperatingStateEnumTypeFinished:
			c.clientData.EVData.ChargeState = EVChargeStateEnumTypeFinished
		}
	}

	if prevEVSEOperationState != c.clientData.EVSEData.OperationState {
		c.callDataUpdateHandler(EVDataElementUpdateEVSEOperationState)
	}
	if prevEVChargeState != c.clientData.EVData.ChargeState {
		c.callDataUpdateHandler(EVDataElementUpdateEVChargeState)
	}
}

// TODO make this more generic, we assume that only one electric connection exists, that only single phases values are available and more
func (c *ConnectionController) updateMeasurementData() {
	var measurementDescription []feature.MeasurementDatasetDefinitionsType
	var measurementData []feature.MeasurementDatasetDataType
	var electricalParameterDescription []feature.ElectricalConnectionParameterDescriptionDataType
	var electricalDescription []feature.ElectricalConnectionDatasetDataType
	var electricalPermittedData []feature.ElectricalConnectionPermittedDataType

	m := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeMeasurement, model.RoleTypeClient)

	if f, ok := m.(*feature.Measurement); ok {
		measurementDescription = f.GetMeasurementDescription()
		measurementData = f.GetMeasurementData()
	}

	e := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeElectricalConnection, model.RoleTypeClient)

	if f, ok := e.(*feature.ElectricalConnection); ok {
		electricalParameterDescription = f.GetElectricalConnectionDescription()
		electricalDescription = f.GetElectricalConnectionData()
		electricalPermittedData = f.GetElectricalConnectionPermittedData()
	}

	if measurementDescription == nil || measurementData == nil || electricalParameterDescription == nil || electricalDescription == nil || electricalPermittedData == nil {
		return
	}

	var measurementCurrentIds []uint
	var measurementPowerIds []uint
	var measurementChargeID uint
	var measurementSoCID uint

	for _, item := range measurementDescription {
		switch item.ScopeType {
		case model.ScopeTypeEnumTypeACCurrent:
			measurementCurrentIds = append(measurementCurrentIds, item.MeasurementId)
		case model.ScopeTypeEnumTypeACPower:
			measurementPowerIds = append(measurementPowerIds, item.MeasurementId)
		case model.ScopeTypeEnumTypeCharge:
			measurementChargeID = item.MeasurementId
		case model.ScopeTypeEnumTypeStateOfCharge:
			measurementSoCID = item.MeasurementId
		}
	}

	for _, eItem := range electricalDescription {
		if c.clientData.EVData.ConnectedPhases != eItem.ConnectedPhases {
			c.clientData.EVData.ConnectedPhases = eItem.ConnectedPhases
			c.callDataUpdateHandler(EVDataElementUpdateConnectedPhases)
		}
	}

	var measurementIdsToPhase = make(map[uint]uint)
	for _, epItem := range electricalParameterDescription {
		_, found := FindValueInSlice(measurementCurrentIds, epItem.MeasurementId)
		if found {
			measurementIdsToPhase[epItem.MeasurementId] = epItem.Phase
		}
		_, found = FindValueInSlice(measurementPowerIds, epItem.MeasurementId)
		if found {
			measurementIdsToPhase[epItem.MeasurementId] = epItem.Phase
		}
	}

	// TODO this needs to be improved (a lot)
	for _, epdItem := range electricalPermittedData {
		powerLimitsUpdated := false
		amperageLimitsUpdated := false
		for _, epItem := range electricalParameterDescription {
			if epItem.ParameterId == epdItem.ParameterId {
				if epItem.ScopeType == model.ScopeTypeEnumTypeACPowerTotal {
					if c.clientData.EVData.LimitsPower.Min != epdItem.MinValue {
						c.clientData.EVData.LimitsPower.Min = epdItem.MinValue
						powerLimitsUpdated = true
					}
					if c.clientData.EVData.LimitsPower.Max != epdItem.MaxValue {
						c.clientData.EVData.LimitsPower.Max = epdItem.MaxValue
						powerLimitsUpdated = true
					}
				} else if epItem.Phase > 0 {
					switch epItem.Phase {
					case 1:
						if c.clientData.EVData.LimitsL1.Min != epdItem.MinValue {
							c.clientData.EVData.LimitsL1.Min = epdItem.MinValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL1.Max != epdItem.MaxValue {
							c.clientData.EVData.LimitsL1.Max = epdItem.MaxValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL1.Default != epdItem.Value {
							c.clientData.EVData.LimitsL1.Default = epdItem.Value
							amperageLimitsUpdated = true
						}
					case 2:
						if c.clientData.EVData.LimitsL2.Min != epdItem.MinValue {
							c.clientData.EVData.LimitsL2.Min = epdItem.MinValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL2.Max != epdItem.MaxValue {
							c.clientData.EVData.LimitsL2.Max = epdItem.MaxValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL2.Default != epdItem.Value {
							c.clientData.EVData.LimitsL2.Default = epdItem.Value
							amperageLimitsUpdated = true
						}
					case 3:
						if c.clientData.EVData.LimitsL3.Min != epdItem.MinValue {
							c.clientData.EVData.LimitsL3.Min = epdItem.MinValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL3.Max != epdItem.MaxValue {
							c.clientData.EVData.LimitsL3.Max = epdItem.MaxValue
							amperageLimitsUpdated = true
						}
						if c.clientData.EVData.LimitsL3.Default != epdItem.Value {
							c.clientData.EVData.LimitsL3.Default = epdItem.Value
							amperageLimitsUpdated = true
						}
					}
				}
			}
		}
		if amperageLimitsUpdated {
			// Min current data should be derived from min power data
			// but as this is only properly provided via VAS the currrent
			// min values can not be trusted.
			// Min current for 3-phase should be at least 2.2A, for 1-phase 6.6A

			if c.clientData.EVData.ConnectedPhases == 1 {
				minCurrent := 6.6
				if c.clientData.EVData.LimitsL1.Min < minCurrent {
					c.clientData.EVData.LimitsL1.Min = minCurrent
				}
			} else if c.clientData.EVData.ConnectedPhases == 3 {
				minCurrent := 2.2
				if c.clientData.EVData.LimitsL1.Min < minCurrent {
					c.clientData.EVData.LimitsL1.Min = minCurrent
				}
				if c.clientData.EVData.LimitsL2.Min < minCurrent {
					c.clientData.EVData.LimitsL2.Min = minCurrent
				}
				if c.clientData.EVData.LimitsL3.Min < minCurrent {
					c.clientData.EVData.LimitsL3.Min = minCurrent
				}
			}
			c.callDataUpdateHandler(EVDataElementUpdateAmperageLimits)
		}
		if !powerLimitsUpdated {
			// Min power data is only properly provided via VAS in ISO15118-2!
			// So use the known min limits and calculate a more likely min power
			if c.clientData.EVData.ConnectedPhases == 1 {
				minPower := c.clientData.EVData.LimitsL1.Min * 230
				if c.clientData.EVData.LimitsPower.Min < minPower {
					c.clientData.EVData.LimitsPower.Min = minPower
				}
				maxPower := c.clientData.EVData.LimitsL1.Max * 230
				if c.clientData.EVData.LimitsPower.Max != maxPower {
					c.clientData.EVData.LimitsPower.Max = maxPower
				}

			} else if c.clientData.EVData.ConnectedPhases == 3 {
				minPower := c.clientData.EVData.LimitsL1.Min*230 + c.clientData.EVData.LimitsL2.Min*230 + c.clientData.EVData.LimitsL3.Min*230
				if c.clientData.EVData.LimitsPower.Min < minPower {
					c.clientData.EVData.LimitsPower.Min = minPower
				}
				maxPower := c.clientData.EVData.LimitsL1.Max*230 + c.clientData.EVData.LimitsL2.Max*230 + c.clientData.EVData.LimitsL3.Max*230
				if c.clientData.EVData.LimitsPower.Max != maxPower {
					c.clientData.EVData.LimitsPower.Max = maxPower
				}
			}
			c.callDataUpdateHandler(EVDataElementUpdatePowerLimits)
		}
	}
	c.log.Println("limits: ")
	c.log.Println("  L1 current: min ", c.clientData.EVData.LimitsL1.Min, "A, max ", c.clientData.EVData.LimitsL1.Max, "A, pause ", c.clientData.EVData.LimitsL1.Default, "A")
	c.log.Println("  L2 current: min ", c.clientData.EVData.LimitsL2.Min, "A, max ", c.clientData.EVData.LimitsL2.Max, "A, pause ", c.clientData.EVData.LimitsL2.Default, "A")
	c.log.Println("  L3 current: min ", c.clientData.EVData.LimitsL3.Min, "A, max ", c.clientData.EVData.LimitsL3.Max, "A, pause ", c.clientData.EVData.LimitsL3.Default, "A")
	c.log.Println("       Power: min ", c.clientData.EVData.LimitsPower.Min, "W, max ", c.clientData.EVData.LimitsPower.Max, "W")

	for _, item := range measurementData {
		if item.MeasurementId == measurementChargeID {
			c.clientData.EVData.Measurements.ChargedEnergy = item.Value
		}
		if item.MeasurementId == measurementSoCID {
			c.clientData.EVData.SoCDataAvailable = true
			c.clientData.EVData.Measurements.SoC = item.Value
		}
		_, found := FindValueInSlice(measurementCurrentIds, item.MeasurementId)
		if found {
			c.clientData.EVData.Measurements.Timestamp = item.Timestamp
			switch measurementIdsToPhase[item.MeasurementId] {
			case 1:
				c.clientData.EVData.Measurements.CurrentL1 = item.Value
			case 2:
				c.clientData.EVData.Measurements.CurrentL2 = item.Value
			case 3:
				c.clientData.EVData.Measurements.CurrentL3 = item.Value
			default:
			}
		}
		_, found = FindValueInSlice(measurementPowerIds, item.MeasurementId)
		if found {
			c.clientData.EVData.Measurements.Timestamp = item.Timestamp
			switch measurementIdsToPhase[item.MeasurementId] {
			case 1:
				// in case we didn't receive power measurements, use current measurements
				if item.Value == 0 && c.clientData.EVData.Measurements.CurrentL1 != 0 {
					c.log.Println("L1 power fallback")
					item.Value = c.clientData.EVData.Measurements.CurrentL1 * 230
				}
				c.clientData.EVData.Measurements.PowerL1 = item.Value
			case 2:
				// in case we didn't receive power measurements, use current measurements
				if item.Value == 0 && c.clientData.EVData.Measurements.CurrentL2 != 0 {
					c.log.Println("L2 power fallback")
					item.Value = c.clientData.EVData.Measurements.CurrentL2 * 230
				}
				c.clientData.EVData.Measurements.PowerL2 = item.Value
			case 3:
				// in case we didn't receive power measurements, use current measurements
				if item.Value == 0 && c.clientData.EVData.Measurements.CurrentL3 != 0 {
					c.log.Println("L3 power fallback")
					item.Value = c.clientData.EVData.Measurements.CurrentL3 * 230
				}
				c.clientData.EVData.Measurements.PowerL3 = item.Value
			default:
			}
		} else {
			c.clientData.EVData.Measurements.Timestamp = item.Timestamp
			item.Value = 0
			switch measurementIdsToPhase[item.MeasurementId] {
			case 1:
				// in case we didn't receive power measurements, use current measurements
				if c.clientData.EVData.Measurements.CurrentL1 != 0 {
					item.Value = c.clientData.EVData.Measurements.CurrentL1 * 230
				}
				c.clientData.EVData.Measurements.PowerL1 = item.Value
			case 2:
				// in case we didn't receive power measurements, use current measurements
				if c.clientData.EVData.Measurements.CurrentL2 != 0 {
					item.Value = c.clientData.EVData.Measurements.CurrentL2 * 230
				}
				c.clientData.EVData.Measurements.PowerL2 = item.Value
			case 3:
				// in case we didn't receive power measurements, use current measurements
				if c.clientData.EVData.Measurements.CurrentL3 != 0 {
					item.Value = c.clientData.EVData.Measurements.CurrentL3 * 230
				}
				c.clientData.EVData.Measurements.PowerL3 = item.Value
			default:
			}
		}
	}

	if len(measurementPowerIds) == 0 {
		// we did not receive any Power measurements, so calculate them
		c.clientData.EVData.Measurements.PowerL1 = c.clientData.EVData.Measurements.CurrentL1 * 230
		c.clientData.EVData.Measurements.PowerL2 = c.clientData.EVData.Measurements.CurrentL2 * 230
		c.clientData.EVData.Measurements.PowerL3 = c.clientData.EVData.Measurements.CurrentL3 * 230
	}

	c.log.Println("phases: ", c.clientData.EVData.ConnectedPhases, ", charged energy: ", c.clientData.EVData.Measurements.ChargedEnergy, "Wh")
	c.log.Println("current current: L1 ", c.clientData.EVData.Measurements.CurrentL1, "A, L2 ", c.clientData.EVData.Measurements.CurrentL2, "A L3 ", c.clientData.EVData.Measurements.CurrentL3, "A")
	c.log.Println("current power: L1 ", c.clientData.EVData.Measurements.PowerL1, "W, L2 ", c.clientData.EVData.Measurements.PowerL2, "W L3 ", c.clientData.EVData.Measurements.PowerL3, "W")

	// TODO REMOVE THIS, this is for testing only!!!
	// currentsPerPhase := []float64{0, 0, 0}
	// c.WriteCurrentLimitData(currentsPerPhase)
}

func (c *ConnectionController) UpdateElectricalConnectionData(f *feature.ElectricalConnection) {
	c.updateMeasurementData()
}

func (c *ConnectionController) UpdateMeasurementData(f *feature.Measurement) {
	c.updateMeasurementData()
}

func (c *ConnectionController) UpdateIdentificationData(f *feature.Identification, data []feature.IdentificationDatasetDataType) {
	for _, item := range data {
		// in case of the EVSE-EV-communication dropping back from ISO to IEC, this value will be updated
		// with an empty string, so the previous identifier is lost
		if len(item.IdentificationValue) > 0 {
			c.clientData.EVData.Identification = item.IdentificationValue
			c.log.Println("EV Identification: ", item.IdentificationValue)
		}
	}
}

func (c *ConnectionController) UpdateUseCaseSupportData(f *feature.NodeManagement, useCasename model.UseCaseNameType, available bool) {
	switch useCasename {
	case model.UseCaseNameType(model.UseCaseNameEnumTypeEVStateOfCharge):
		c.clientData.EVData.UCSoCAvailable = available
		c.log.Println("SoC support: ", available)
		c.callDataUpdateHandler(EVDataElementUpdateUseCaseSoC)

	case model.UseCaseNameType(model.UseCaseNameEnumTypeOptimizationOfSelfConsumptionDuringEVCharging):
		c.clientData.EVData.UCSelfConsumptionAvailable = available
		c.log.Println("Self consumption support: ", available)
		c.callDataUpdateHandler(EVDataElementUpdateUseCaseSelfConsumption)

	case model.UseCaseNameType(model.UseCaseNameEnumTypeCoordinatedEVCharging):
		c.clientData.EVData.UCCoordinatedChargingAvailable = available
		c.log.Println("Coordinated charging support: ", available)
		c.callDataUpdateHandler(EVDataElementUpdateUseCaseCoordinatedCharging)
	}
}

func (c *ConnectionController) remoteEntityForFeatureAddress(rf model.FeatureAddressType) spine.Entity {
	re := c.remoteDevice.GetEntities()

	for _, entity := range re {
		if reflect.DeepEqual(entity.GetAddress(), rf.Entity) {
			return entity
		}
	}

	return nil
}

func (c *ConnectionController) UpdateDeviceClassificationData(f *feature.DeviceClassification, rf model.FeatureAddressType, data model.DeviceClassificationManufacturerDataType) {
	re := c.remoteEntityForFeatureAddress(rf)
	entityType := model.EntityTypeEnumType(re.GetType())

	if entityType == model.EntityTypeEnumTypeEVSE {
		if data.BrandName != nil {
			c.clientData.EVSEData.Manufacturer.BrandName = string(*data.BrandName)
		}
		if data.DeviceCode != nil {
			c.clientData.EVSEData.Manufacturer.DeviceCode = string(*data.DeviceCode)
		}
		if data.DeviceName != nil {
			c.clientData.EVSEData.Manufacturer.DeviceName = string(*data.DeviceName)
		}
	} else if entityType == model.EntityTypeEnumTypeEV {
		if data.BrandName != nil {
			c.clientData.EVData.Manufacturer.BrandName = string(*data.BrandName)
		}
		if data.DeviceCode != nil {
			c.clientData.EVData.Manufacturer.DeviceCode = string(*data.DeviceCode)
		}
		if data.DeviceName != nil {
			c.clientData.EVData.Manufacturer.DeviceName = string(*data.DeviceName)
		}
	}
}

func (c *ConnectionController) UpdateLoadControlLimitData(f *feature.LoadControl) {
	limitDescriptionData := f.GetLoadControlLimitDescriptionData()
	limitData := f.GetLoadControlLimitData()

	if limitDescriptionData == nil || limitData == nil {
		return
	}

	for _, item := range limitDescriptionData {
		for _, dataItem := range limitData {
			if dataItem.LimitId == item.LimitId {
				switch item.ScopeType {
				case model.ScopeTypeEnumTypeOverloadProtection:
					c.clientData.EVData.OverloadProtectionActive = dataItem.IsLimitActive
				case model.ScopeTypeEnumTypeSelfConsumption:
					c.clientData.EVData.SelfConsumptionActive = dataItem.IsLimitActive

				}
			}
		}
	}
}

func (c *ConnectionController) UpdateTimeSeriesDescriptionData(f *feature.TimeSeries) {
	timeSeriesDescriptionData := f.GetTimeSeriesDescriptionData()

	for _, item := range timeSeriesDescriptionData {
		// TODO: add processing
		switch item.TimeSeriesType {
		case model.TimeSeriesTypeEnumTypeConstraints:
			if item.UpdateRequired {
				// we need to send a response with a plan (within 20s or something like that)
				c.callDataUpdateHandler(EVDataElementUpdateChargingPlanRequired)
			}
			return
		case model.TimeSeriesTypeEnumTypePlan:
			return
		case model.TimeSeriesTypeEnumTypeSingleDemand:
			return
		}
	}
}

func (c *ConnectionController) UpdateTimeSeriesData(f *feature.TimeSeries, timeSeriesData feature.TimeSeriesDatasetType) {
	timeSeriesDescriptionData := f.GetTimeSeriesDescriptionData()

	c.clientData.EVData.ChargingStrategy = EVChargingStrategyEnumTypeUnknown

	if timeSeriesDescriptionData == nil {
		return
	}

	timeSeriesType, err := f.GetTimeSeriesTypeForId(timeSeriesData.TimeSeriesId)
	if err != nil {
		c.log.Printf("Error getting Time Series Type for ID %d: %s\n", timeSeriesData.TimeSeriesId, err)
		return
	}

	// TODO: add processing
	switch timeSeriesType {
	case model.TimeSeriesTypeEnumTypeConstraints:
		return
	case model.TimeSeriesTypeEnumTypePlan:
		output := "EV informed about its charging plan:\n"
		if timeSeriesData.TimePeriod == nil {
			c.log.Printf("The time series plan is empty %d: %s\n", timeSeriesData.TimeSeriesId, err)
			return
		}
		if timeSeriesData.TimePeriod.StartTime != nil {
			output += fmt.Sprintf("\tStartTime: %s\n", *timeSeriesData.TimePeriod.StartTime)
		}
		if timeSeriesData.TimePeriod.EndTime != nil {
			output += fmt.Sprintf("\tEndTime: %s\n", *timeSeriesData.TimePeriod.EndTime)
		}
		for _, slot := range timeSeriesData.TimeSeriesSlots {
			output += fmt.Sprintf("\t%.1f: %s\n", slot.MaxValue.GetValue(), *slot.Duration)
		}
		c.log.Println(output)
		return
	case model.TimeSeriesTypeEnumTypeSingleDemand:
		if timeSeriesData.TimeSeriesSlots == nil {
			c.log.Printf("The time series slots are empty %d: %s\n", timeSeriesData.TimeSeriesId, err)
			return
		}
		demand := timeSeriesData.TimeSeriesSlots[0].Value.GetValue()
		c.clientData.EVData.ChargingDemand = demand / 1000 // return kWh
		c.clientData.EVData.ChargingTargetDuration = time.Duration(24) * time.Hour

		if demand > 0 {
			// if demand is > 0 and duration is not existing, the EV is not charging via a timer
			// but either via direct charging enabled or charging to minimum SoC using a profile
			if timeSeriesData.TimeSeriesSlots[0].Duration == nil {
				c.clientData.EVData.ChargingStrategy = EVChargingStrategyEnumTypeDirectCharging
				c.log.Printf("EV is charging via direct charging. Demand: %.1f kWh\n", (demand / 1000))
			} else {
				p, err := period.Parse(*timeSeriesData.TimeSeriesSlots[0].Duration)
				if err != nil {
					c.log.Printf("Error parsing duration: %s\n", err)
					return
				}
				c.clientData.EVData.ChargingTargetDuration, _ = p.Duration()
				c.clientData.EVData.ChargingStrategy = EVChargingStrategyEnumTypeTimedCharging
				c.log.Printf("EV is charging via timed charging. Demand: %.1f kWh, Duration: %s\n", (demand / 1000), c.clientData.EVData.ChargingTargetDuration.String())
			}
		} else {
			c.clientData.EVData.ChargingStrategy = EVChargingStrategyEnumTypeNoDemand
			c.log.Println("EV not reporting a demand")
		}
		c.callDataUpdateHandler(EVDataElementUpdateChargingStrategy)
		return
	}
}

func (c *ConnectionController) UpdateIncentiveConstraintsData(f *feature.IncentiveTable) {
	// we received tariff, tiers, boundaries, incentives and slotcount limits

	// now we need to reply with the incentiveTableDescription

	if c.remoteDevice == nil {
		// errors.New("charger is not connected")
		return
	}

	evEntity := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))
	if evEntity == nil {
		// errors.New("no ev connected")
		return
	}

	rf := evEntity.FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeServer)

	ctx := c.context(nil)

	lf := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeIncentiveTable, model.RoleTypeClient)
	if l, ok := lf.(*feature.IncentiveTable); ok {
		_ = l.WriteDescriptionData(ctx, rf)
	}
}

func (c *ConnectionController) WriteChargingPlan(chargingPlan EVChargingPlan) error {
	if c.remoteDevice == nil {
		return errors.New("charger is not connected")
	}

	evEntity := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))
	if evEntity == nil {
		return errors.New("no ev connected")
	}

	ctx := c.context(nil)

	// at first we need to send the power limits plan
	lf := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeTimeSeries, model.RoleTypeClient)
	l, ok := lf.(*feature.TimeSeries)

	if !ok {
		return errors.New("timeseries feature is not available on local device")
	}

	rf := evEntity.FeatureByProps(model.FeatureTypeEnumTypeTimeSeries, model.RoleTypeServer)

	timeSeriesSlots := []feature.TimeSeriesChargingSlot{}
	for _, slot := range chargingPlan.Slots {
		timeSeriesSlots = append(timeSeriesSlots, feature.TimeSeriesChargingSlot{
			MaxValue: slot.MaxValue,
			Duration: slot.Duration,
		})
	}

	timeSeriesChargingPlan := feature.TimeSeriesChargingPlan{
		Duration: chargingPlan.Duration,
		Slots:    timeSeriesSlots,
	}

	output := "Sending time Series slots:\n"
	for index, slot := range timeSeriesSlots {
		if index > 0 && index%4 == 0 {
			output += "\n"
		}
		output += fmt.Sprintf("\t%.1f: %s", slot.MaxValue, slot.Duration.String())
		if (index + 1) == len(timeSeriesSlots) {
			output += "\n"
		}
	}
	c.log.Println(output)

	if err := l.WriteTimeSeriesPlanData(ctx, rf, timeSeriesChargingPlan); err != nil {
		c.log.Println("error sending loadcontrol limits ", err)
		return err
	}

	// now we need to send the incentive plan
	lf2 := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeIncentiveTable, model.RoleTypeClient)
	l2, ok := lf2.(*feature.IncentiveTable)

	if !ok {
		return errors.New("incentivetable feature is not available on local device")
	}

	rf2 := evEntity.FeatureByProps(model.FeatureTypeEnumTypeIncentiveTable, model.RoleTypeServer)

	incentiveSlots := []feature.IncentiveChargingSlot{}
	for _, slot := range chargingPlan.Slots {
		incentiveSlots = append(incentiveSlots, feature.IncentiveChargingSlot{
			Pricing:  slot.Pricing,
			Duration: slot.Duration,
		})
	}

	incentiveChargingPlan := feature.IncentiveChargingPlan{
		Duration: chargingPlan.Duration,
		Slots:    incentiveSlots,
	}

	output = "Sending incentive slots:\n"
	for index, slot := range incentiveSlots {
		if index > 0 && index%4 == 0 {
			output += "\n"
		}
		output += fmt.Sprintf("\t%.3f: %s", slot.Pricing, slot.Duration.String())
		if (index + 1) == len(incentiveSlots) {
			output += "\n"
		}
	}
	c.log.Println(output)

	if err := l2.WriteIncentiveTablePlanData(ctx, rf2, incentiveChargingPlan); err != nil {
		c.log.Println("error sending loadcontrol limits ", err)
		return err
	}

	return nil
}

// TODO error handling and returning
func (c *ConnectionController) WriteCurrentLimitData(overloadProtectionCurrentsPerPhase []float64, selfConsumptionCurrentsPerPhase []float64, evData EVDataType) error {
	var electricalParameterDescription []feature.ElectricalConnectionParameterDescriptionDataType
	var measurementDescription []feature.MeasurementDatasetDefinitionsType
	var limitDescription []feature.LoadControlLimitDescriptionDataType

	e := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeElectricalConnection, model.RoleTypeClient)

	if f, ok := e.(*feature.ElectricalConnection); ok {
		electricalParameterDescription = f.GetElectricalConnectionDescription()
	}

	m := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeMeasurement, model.RoleTypeClient)

	if f, ok := m.(*feature.Measurement); ok {
		measurementDescription = f.GetMeasurementDescription()
	}

	lf := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM)).FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeClient)
	l, ok := lf.(*feature.LoadControl)

	if !ok {
		return errors.New("loadcontrol feature is not available on local device")
	}

	if c.remoteDevice == nil {
		return errors.New("charger is not connected")
	}

	evEntity := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))
	if evEntity == nil {
		return errors.New("no ev connected")
	}

	rf := evEntity.FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeServer)

	limitDescription = l.GetLoadControlLimitDescriptionData()

	if electricalParameterDescription == nil || measurementDescription == nil || limitDescription == nil {
		return errors.New("no eletrical paramaters, measurements, or limits available yet")
	}

	var measurementCurrentIds []uint

	for _, item := range measurementDescription {
		if item.ScopeType == model.ScopeTypeEnumTypeACCurrent {
			measurementCurrentIds = append(measurementCurrentIds, item.MeasurementId)
		}
	}

	var measurementIdsToPhase = make(map[uint]uint)
	for _, epItem := range electricalParameterDescription {
		_, found := FindValueInSlice(measurementCurrentIds, epItem.MeasurementId)
		if found {
			measurementIdsToPhase[epItem.MeasurementId] = epItem.Phase
		}
	}

	var limitItems []feature.LoadControlLimitDatasetType

	for scopeTypes := 0; scopeTypes < 2; scopeTypes++ {
		currentsPerPhase := overloadProtectionCurrentsPerPhase
		if scopeTypes == 1 {
			currentsPerPhase = selfConsumptionCurrentsPerPhase
		}

		for index, current := range currentsPerPhase {
			if index < int(evData.ConnectedPhases) {
				measurementId := measurementCurrentIds[index]
				currentValue := current

				var limitId model.LoadControlLimitIdType = 0

				for _, item := range limitDescription {
					if item.MeasurementId == measurementId {
						if (scopeTypes == 0 && item.ScopeType == model.ScopeTypeEnumTypeOverloadProtection) ||
							(scopeTypes == 1 && item.ScopeType == model.ScopeTypeEnumTypeSelfConsumption) {
							limitId = model.LoadControlLimitIdType(item.LimitId)
						}
					}
				}

				switch index {
				case 0:
					if currentValue < evData.LimitsL1.Min {
						currentValue = evData.LimitsL1.Default
					}
					if currentValue > evData.LimitsL1.Max {
						currentValue = evData.LimitsL1.Max
					}
				case 1:
					if currentValue < evData.LimitsL2.Min {
						currentValue = evData.LimitsL2.Default
					}
					if currentValue > evData.LimitsL2.Max {
						currentValue = evData.LimitsL2.Max
					}
				case 2:
					if currentValue < evData.LimitsL3.Min {
						currentValue = evData.LimitsL3.Default
					}
					if currentValue > evData.LimitsL3.Max {
						currentValue = evData.LimitsL3.Max
					}
				}

				if limitId > 0 {
					newItem := feature.LoadControlLimitDatasetType{
						LimitId: uint(limitId),
						Value:   currentValue,
					}
					limitItems = append(limitItems, newItem)
				}
			}
		}
	}

	if limitItems == nil {
		return errors.New("no limits available")
	}

	sort.Slice(limitItems, func(i, j int) bool {
		return limitItems[i].LimitId < limitItems[j].LimitId
	})

	ctx := c.context(nil)

	if err := l.WriteLoadControlLimitListData(ctx, rf, limitItems); err != nil {
		c.log.Println("error sending loadcontrol limits ", err)
		return err
	}

	return nil
}

func FindValueInSlice(slice []uint, value uint) (int, bool) {
	for i, item := range slice {
		if item == value {
			return i, true
		}
	}
	return -1, false
}
