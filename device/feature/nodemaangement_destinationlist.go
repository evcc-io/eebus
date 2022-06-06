package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

func (f *NodeManagement) readDestinationListData(ctrl spine.Context, data model.NodeManagementDestinationListDataType, isPartialForCmd bool) error {
	localDevice := f.Entity.GetDevice()

	deviceAddress := localDevice.GetAddress()
	deviceType := localDevice.GetType()
	featureSet := model.NetworkManagementFeatureSetTypeSmart

	res := model.CmdType{
		NodeManagementDestinationListData: &model.NodeManagementDestinationListDataType{
			NodeManagementDestinationData: []model.NodeManagementDestinationDataType{
				{
					DeviceDescription: &model.NetworkManagementDeviceDescriptionDataType{
						DeviceAddress: &model.DeviceAddressType{
							Device: &deviceAddress,
						},
						DeviceType:        &deviceType,
						NetworkFeatureSet: &featureSet,
					},
				},
			},
		},
	}

	return ctrl.Reply(model.CmdClassifierTypeReply, res)

}

func (f *NodeManagement) handleNodeManagementDestinationListData(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementDestinationListDataType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeRead:
		return f.readDestinationListData(ctrl, *data, isPartialForCmd)
	default:
		return fmt.Errorf("nodemanagement.handleNodeManagementDestinationListData: NodeManagementDestinationListDataType CmdClassifierType not implemented: %s", op)
	}
}
