package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUseCaseInformationReply(t *testing.T) {
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var actor UseCaseActorType = "CEM"
	var useCaseVersion SpecificationVersionType = "1.0.1"
	var useCaseName1 UseCaseNameType = "evseCommissioningAndConfiguration"
	var useCaseName2 UseCaseNameType = "evChargingSummary"
	var useCaseName3 UseCaseNameType = "measurementOfElectricityDuringEvCharging"
	var useCaseName4 UseCaseNameType = "optimizationOfSelfConsumptionDuringEvCharging"
	var useCaseName5 UseCaseNameType = "coordinatedEvCharging"
	var useCaseName6 UseCaseNameType = "overloadProtectionByEvChargingCurrentCurtailment"
	var useCaseName7 UseCaseNameType = "evCommissioningAndConfiguration"

	payload := PayloadType{
		Cmd: []CmdType{{
			NodeManagementUseCaseData: &NodeManagementUseCaseDataType{
				UseCaseInformation: []UseCaseInformationDataType{
					{
						Address: &FeatureAddressType{Device: &deviceHems},
						Actor:   &actor,
						UseCaseSupport: []UseCaseSupportType{
							{
								UseCaseName:     &useCaseName1,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2},
							},
							{
								UseCaseName:     &useCaseName2,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1},
							},
							{
								UseCaseName:     &useCaseName3,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2, 3},
							},
							{
								UseCaseName:     &useCaseName4,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2, 3},
							},
							{
								UseCaseName:     &useCaseName5,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
							},
							{
								UseCaseName:     &useCaseName6,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2, 3},
							},
							{
								UseCaseName:     &useCaseName7,
								UseCaseVersion:  &useCaseVersion,
								ScenarioSupport: []UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
							},
						},
					},
				},
			},
		}},
	}

	json, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("TestUseCaseInformationReply() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `[{"cmd":[[{"nodeManagementUseCaseData":[{"useCaseInformation":[[{"address":[{"device":"d:_i:3210_HEMS"}]},{"actor":"CEM"},{"useCaseSupport":[[{"useCaseName":"evseCommissioningAndConfiguration"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2]}],[{"useCaseName":"evChargingSummary"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1]}],[{"useCaseName":"measurementOfElectricityDuringEvCharging"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"optimizationOfSelfConsumptionDuringEvCharging"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"coordinatedEvCharging"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2,3,4,5,6,7,8]}],[{"useCaseName":"overloadProtectionByEvChargingCurrentCurtailment"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2,3]}],[{"useCaseName":"evCommissioningAndConfiguration"},{"useCaseVersion":"1.0.1"},{"scenarioSupport":[1,2,3,4,5,6,7,8]}]]}]]}]}]]}]`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestUseCaseInformationReply() actual json string doesn't match expected result")
	}
}
