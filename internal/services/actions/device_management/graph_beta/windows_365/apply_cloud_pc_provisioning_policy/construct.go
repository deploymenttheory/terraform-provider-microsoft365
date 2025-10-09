package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func constructRequest(ctx context.Context, data *ApplyCloudPcProvisioningPolicyActionModel) (*devicemanagement.VirtualEndpointProvisioningPoliciesItemApplyPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s request for policy ID: %s", ActionName, data.ProvisioningPolicyID.ValueString()))

	requestBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemApplyPostRequestBody()

	if !data.PolicySettings.IsNull() && !data.PolicySettings.IsUnknown() {
		policySettings := data.PolicySettings.ValueString()
		requestBody.GetAdditionalData()["policySettings"] = policySettings
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for action %s", ActionName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s request for %s", ActionName, data.ProvisioningPolicyID.ValueString()))
	return requestBody, nil
}
