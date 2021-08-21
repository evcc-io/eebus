package feature

import (
	"errors"
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

func (f *NodeManagement) readUseCaseData(ctrl spine.Context, data model.NodeManagementUseCaseDataType) error {
	// TODO: generate this!

	deviceAddress := f.GetEntity().GetDevice().GetAddress()
	actor := model.UseCaseActorType("CEM")
	var useCaseSupport []model.UseCaseSupportType
	useCaseVersion := model.SpecificationVersionType("1.0.1")

	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeEVSECommissioningAndConfiguration)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeEVCommissioningAndConfiguration)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	// {
	// 	useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeEVChargingSummary)
	// 	useCaseItem := model.UseCaseSupportType{
	// 		UseCaseVersion:  &useCaseVersion,
	// 		UseCaseName:     &useCaseName,
	// 		ScenarioSupport: []model.UseCaseScenarioSupportType{1},
	// 	}
	// 	useCaseSupport = append(useCaseSupport, useCaseItem)
	// }
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeMeasurementOfElectricityDuringEVCharging)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeCoordinatedEVCharging)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 3, 4, 5, 6, 7, 8},
			// ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeOptimizationOfSelfConsumptionDuringEVCharging)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeEVStateOfCharge)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}
	{
		useCaseName := model.UseCaseNameType(model.UseCaseNameEnumTypeOverloadProtectionByEVChargingCurrentCurtailment)
		useCaseItem := model.UseCaseSupportType{
			UseCaseVersion:  &useCaseVersion,
			UseCaseName:     &useCaseName,
			ScenarioSupport: []model.UseCaseScenarioSupportType{1, 2, 3},
		}
		useCaseSupport = append(useCaseSupport, useCaseItem)
	}

	res := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{
			UseCaseInformation: []model.UseCaseInformationDataType{
				{
					Address:        &model.FeatureAddressType{Device: &deviceAddress},
					Actor:          &actor,
					UseCaseSupport: useCaseSupport,
				},
			},
		},
	}

	return ctrl.Reply(model.CmdClassifierTypeReply, res)
}

func (f *NodeManagement) updateSupportedUseCases(ctrl spine.Context, remoteDevice spine.Device, data model.NodeManagementUseCaseDataType) error {
	useCaseInformation := data.UseCaseInformation

	for _, actorItem := range useCaseInformation {
		useCaseActor := actorItem.Actor
		useCaseSupport := actorItem.UseCaseSupport
		remoteDevice.SetUseCaseActor(string(*useCaseActor), useCaseSupport)

		if *actorItem.Actor == model.UseCaseActorType(model.UseCaseActorEnumTypeEV) {
			for _, item := range actorItem.UseCaseSupport {
				if f.Delegate != nil {
					f.Delegate.UpdateUseCaseSupportData(f, *item.UseCaseName, *item.UseCaseAvailable)
				}
			}
		}
	}

	return nil
}

func (f *NodeManagement) replyUseCaseData(ctrl spine.Context, data model.NodeManagementUseCaseDataType) error {
	remoteDevice := ctrl.GetDevice()

	// Exmaple EV: {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:EVSE"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":13484},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"nodeManagementUseCaseData":[{"useCaseInformation":[[{"actor":"EV"},{"useCaseSupport":[[{"useCaseName":"measurementOfElectricityDuringEvCharging"},{"useCaseAvailable":true},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"optimizationOfSelfConsumptionDuringEvCharging"},{"useCaseAvailable":true},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"overloadProtectionByEvChargingCurrentCurtailment"},{"useCaseAvailable":true},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"coordinatedEvCharging"},{"useCaseAvailable":true},{"scenarioSupport":[1,2,3,4,5,6,7,8]}],[{"useCaseName":"evCommissioningAndConfiguration"},{"useCaseAvailable":true},{"scenarioSupport":[1,2,3,4,5,6,7,8]}],[{"useCaseName":"evseCommissioningAndConfiguration"},{"useCaseAvailable":true},{"scenarioSupport":[1,2]}],[{"useCaseName":"evChargingSummary"},{"useCaseAvailable":true},{"scenarioSupport":[1]}],[{"useCaseName":"evStateOfCharge"},{"useCaseAvailable":false},{"scenarioSupport":[1]}]]}]]}]}]]}]}]}}]}

	useCaseInformation := data.UseCaseInformation
	if useCaseInformation == nil {
		return errors.New("nodemanagement.replyUseCaseData: invalid UseCaseInformation")
	}

	return f.updateSupportedUseCases(ctrl, remoteDevice, data)
}

func (f *NodeManagement) handleUseCaseData(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementUseCaseDataType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeRead:
		return f.readUseCaseData(ctrl, *data)

	case model.CmdClassifierTypeReply:
		return f.replyUseCaseData(ctrl, *data)

	case model.CmdClassifierTypeNotify:
		return f.replyUseCaseData(ctrl, *data)

	default:
		return fmt.Errorf("nodemanagement.handleUseCaseData: NodeManagementUseCaseData CmdClassifierType not implemented: %s", op)
	}
}
