package communication

import (
	"errors"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
	"github.com/evcc-io/eebus/util"
)

const (
	SequenceEnumTypeEVDeviceClassification = "EVDeviceClassification"
	SequenceEnumTypeEVIdentification       = "EVIdentification"
	SequenceEnumTypeEVDeviceConfiguration  = "EVDeviceConfiguration"
	SequenceEnumTypeEVMeasurement          = "EVMeasurement"
	SequenceEnumTypeEVLoadControl          = "EVLoadControl"
	SequenceEnumTypeEVCoordinatedCharging  = "EVCoordinatedCharging"
)

type SequenceElement struct {
	featureType   model.FeatureTypeEnumType
	functionType  model.FunctionEnumType
	cmdClassifier model.CmdClassifierType
	msgCounter    model.MsgCounterType
}

type SequenceFlow struct {
	currentId    int
	sequenceType string
	elements     []SequenceElement
}

type SequencesController struct {
	log           util.Logger
	sequenceFlows []*SequenceFlow
}

func NewSequencesController(log util.Logger) *SequencesController {
	s := &SequencesController{
		log:           log,
		sequenceFlows: []*SequenceFlow{},
	}

	return s
}

func (s *SequencesController) Boot() {
	s.sequenceFlows = append(s.sequenceFlows, s.setupEVIdentificationSequences())
	s.sequenceFlows = append(s.sequenceFlows, s.setupEVConfigurationSequences())
	s.sequenceFlows = append(s.sequenceFlows, s.setupEVMeasurementSequences())
	s.sequenceFlows = append(s.sequenceFlows, s.setupEVLoadControlSequences())
	s.sequenceFlows = append(s.sequenceFlows, s.setupEVCoordinatedChargingSequences())
}

func (s *SequencesController) newElement(featureType model.FeatureTypeEnumType, functionType model.FunctionEnumType, cmdClassifier model.CmdClassifierType) SequenceElement {
	element := SequenceElement{
		featureType:   featureType,
		functionType:  functionType,
		cmdClassifier: cmdClassifier,
	}
	return element
}

// sequenceEVIdentification
func (s *SequencesController) setupEVIdentificationSequences() *SequenceFlow {
	newSequenceFlow := SequenceFlow{}
	newSequenceFlow.sequenceType = SequenceEnumTypeEVIdentification

	{
		element := s.newElement(model.FeatureTypeEnumTypeIdentification, model.FunctionEnumTypeIdentificationListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	return &newSequenceFlow
}

// sequenceEVCommissioningAndConfiguration
func (s *SequencesController) setupEVConfigurationSequences() *SequenceFlow {
	newSequenceFlow := SequenceFlow{}
	newSequenceFlow.sequenceType = SequenceEnumTypeEVDeviceConfiguration

	{
		element := s.newElement(model.FeatureTypeEnumTypeDeviceConfiguration, model.FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeDeviceConfiguration, model.FunctionEnumTypeDeviceConfigurationKeyValueListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	return &newSequenceFlow
}

// sequenceEVChargingElectricityMeasurement
func (s *SequencesController) setupEVMeasurementSequences() *SequenceFlow {
	newSequenceFlow := SequenceFlow{}
	newSequenceFlow.sequenceType = SequenceEnumTypeEVMeasurement

	{
		element := s.newElement(model.FeatureTypeEnumTypeMeasurement, model.FunctionEnumTypeMeasurementDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeMeasurement, model.FunctionEnumTypeMeasurementListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeElectricalConnection, model.FunctionEnumTypeElectricalConnectionParameterDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeElectricalConnection, model.FunctionEnumTypeElectricalConnectionDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeElectricalConnection, model.FunctionEnumTypeElectricalConnectionPermittedValueSetListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	return &newSequenceFlow
}

// SequenceEnumTypeEVLoadControl
func (s *SequencesController) setupEVLoadControlSequences() *SequenceFlow {
	newSequenceFlow := SequenceFlow{}
	newSequenceFlow.sequenceType = SequenceEnumTypeEVLoadControl

	{
		element := s.newElement(model.FeatureTypeEnumTypeLoadControl, model.FunctionEnumTypeLoadControlLimitDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeLoadControl, model.FunctionEnumTypeLoadControlLimitListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	return &newSequenceFlow
}

// SequenceEnumTypeEVCoordinatedCharging
func (s *SequencesController) setupEVCoordinatedChargingSequences() *SequenceFlow {
	newSequenceFlow := SequenceFlow{}
	newSequenceFlow.sequenceType = SequenceEnumTypeEVCoordinatedCharging

	// Scenario 1 + 4
	{
		element := s.newElement(model.FeatureTypeEnumTypeTimeSeries, model.FunctionEnumTypeTimeSeriesDescriptionListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeTimeSeries, model.FunctionEnumTypeTimeSeriesListData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	// Scenario 3
	{
		element := s.newElement(model.FeatureTypeEnumTypeIncentiveTable, model.FunctionEnumTypeIncentiveTableDescriptionData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeIncentiveTable, model.FunctionEnumTypeIncentiveTableConstraintsData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}
	{
		element := s.newElement(model.FeatureTypeEnumTypeIncentiveTable, model.FunctionEnumTypeIncentiveTableData, model.CmdClassifierTypeRead)
		newSequenceFlow.elements = append(newSequenceFlow.elements, element)
	}

	return &newSequenceFlow
}

func (s *SequencesController) processStepInSequenceFlow(ctx spine.Context, sequenceFlow *SequenceFlow) error {
	sequenceElement := sequenceFlow.elements[sequenceFlow.currentId]

	msgCounter, err := ctx.ProcessSequenceFlowRequest(sequenceElement.featureType, sequenceElement.functionType, sequenceElement.cmdClassifier)
	if err != nil {
		return err
	}
	if msgCounter == nil {
		return errors.New("No error returned with msgCounter as nil")
	}
	sequenceFlow.elements[sequenceFlow.currentId].msgCounter = *msgCounter

	return nil
}

// start a sequence of a specific type
func (s *SequencesController) StartSequenceFlow(ctx spine.Context, sequenceType string) error {
	for _, sequenceFlow := range s.sequenceFlows {
		if sequenceFlow.sequenceType == sequenceType {
			return s.processStepInSequenceFlow(ctx, sequenceFlow)
		}
	}

	return nil
}

// invoke the next step in a sequence
func (s *SequencesController) ProcessResponseInSequences(ctx spine.Context, msgCounter *model.MsgCounterType) error {
	if msgCounter == nil {
		return errors.New("invalid msgCounter")
	}

	var sequenceFlow *SequenceFlow

	for _, flow := range s.sequenceFlows {
		currentSequenceElement := flow.elements[flow.currentId]
		if currentSequenceElement.msgCounter != 0 && currentSequenceElement.msgCounter == *msgCounter {
			sequenceFlow = flow
		}
	}

	if sequenceFlow != nil {
		// if sequenceFlow.sequenceType == SequenceEnumTypeEVMeasurement {
		sequenceFlow.currentId += 1
		if sequenceFlow.currentId < len(sequenceFlow.elements) {
			return s.processStepInSequenceFlow(ctx, sequenceFlow)
		} else {
			sequenceFlow.currentId = 0
		}
		// }
	}

	return nil
}
