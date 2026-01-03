package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceProactiveRemediation) *devicemanagement.ManagedDevicesItemInitiateOnDemandProactiveRemediationPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemInitiateOnDemandProactiveRemediationPostRequestBody()

	convert.FrameworkToGraphString(device.ScriptPolicyID, requestBody.SetScriptPolicyId)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device on-demand proactive remediation request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceProactiveRemediation) *devicemanagement.ComanagedDevicesItemInitiateOnDemandProactiveRemediationPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemInitiateOnDemandProactiveRemediationPostRequestBody()

	convert.FrameworkToGraphString(device.ScriptPolicyID, requestBody.SetScriptPolicyId)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device on-demand proactive remediation request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
