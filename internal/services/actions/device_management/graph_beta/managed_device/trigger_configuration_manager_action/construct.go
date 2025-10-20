package graphBetaTriggerConfigurationManagerActionManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceConfigManagerAction) *devicemanagement.ManagedDevicesItemTriggerConfigurationManagerActionPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemTriggerConfigurationManagerActionPostRequestBody()

	// Create the configurationManagerAction object
	configManagerAction := msgraphbetamodels.NewConfigurationManagerAction()

	// Convert string action to enum
	actionValue := device.Action.ValueString()
	var actionEnum msgraphbetamodels.ConfigurationManagerActionType

	switch actionValue {
	case "refreshMachinePolicy":
		actionEnum = msgraphbetamodels.REFRESHMACHINEPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	case "refreshUserPolicy":
		actionEnum = msgraphbetamodels.REFRESHUSERPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	case "wakeUpClient":
		actionEnum = msgraphbetamodels.WAKEUPCLIENT_CONFIGURATIONMANAGERACTIONTYPE
	case "appEvaluation":
		actionEnum = msgraphbetamodels.APPEVALUATION_CONFIGURATIONMANAGERACTIONTYPE
	case "quickScan":
		actionEnum = msgraphbetamodels.QUICKSCAN_CONFIGURATIONMANAGERACTIONTYPE
	case "fullScan":
		actionEnum = msgraphbetamodels.FULLSCAN_CONFIGURATIONMANAGERACTIONTYPE
	case "windowsDefenderUpdateSignatures":
		actionEnum = msgraphbetamodels.WINDOWSDEFENDERUPDATESIGNATURES_CONFIGURATIONMANAGERACTIONTYPE
	default:
		// This shouldn't happen due to schema validation, but default to refreshMachinePolicy
		actionEnum = msgraphbetamodels.REFRESHMACHINEPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	}

	configManagerAction.SetAction(&actionEnum)
	requestBody.SetConfigurationManagerAction(configManagerAction)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device Configuration Manager action request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceConfigManagerAction) *devicemanagement.ComanagedDevicesItemTriggerConfigurationManagerActionPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemTriggerConfigurationManagerActionPostRequestBody()

	// Create the configurationManagerAction object
	configManagerAction := msgraphbetamodels.NewConfigurationManagerAction()

	// Convert string action to enum
	actionValue := device.Action.ValueString()
	var actionEnum msgraphbetamodels.ConfigurationManagerActionType

	switch actionValue {
	case "refreshMachinePolicy":
		actionEnum = msgraphbetamodels.REFRESHMACHINEPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	case "refreshUserPolicy":
		actionEnum = msgraphbetamodels.REFRESHUSERPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	case "wakeUpClient":
		actionEnum = msgraphbetamodels.WAKEUPCLIENT_CONFIGURATIONMANAGERACTIONTYPE
	case "appEvaluation":
		actionEnum = msgraphbetamodels.APPEVALUATION_CONFIGURATIONMANAGERACTIONTYPE
	case "quickScan":
		actionEnum = msgraphbetamodels.QUICKSCAN_CONFIGURATIONMANAGERACTIONTYPE
	case "fullScan":
		actionEnum = msgraphbetamodels.FULLSCAN_CONFIGURATIONMANAGERACTIONTYPE
	case "windowsDefenderUpdateSignatures":
		actionEnum = msgraphbetamodels.WINDOWSDEFENDERUPDATESIGNATURES_CONFIGURATIONMANAGERACTIONTYPE
	default:
		// This shouldn't happen due to schema validation, but default to refreshMachinePolicy
		actionEnum = msgraphbetamodels.REFRESHMACHINEPOLICY_CONFIGURATIONMANAGERACTIONTYPE
	}

	configManagerAction.SetAction(&actionEnum)
	requestBody.SetConfigurationManagerAction(configManagerAction)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device Configuration Manager action request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
