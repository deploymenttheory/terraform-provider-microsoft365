package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type PolicyValidationResult struct {
	PolicyNotFound       bool
	InvalidProvisionType bool
	ProvisioningType     string
}

func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyID string, reservePercentageSet bool) (*PolicyValidationResult, error) {
	result := &PolicyValidationResult{
		PolicyNotFound:       false,
		InvalidProvisionType: false,
	}

	tflog.Debug(ctx, "Validating provisioning policy", map[string]any{"policy_id": policyID})

	policy, err := client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(policyID).
		Get(ctx, nil)

	if err != nil {
		graphErr := errors.GraphError(ctx, err)
		if graphErr.StatusCode == 404 {
			result.PolicyNotFound = true
			tflog.Warn(ctx, "Provisioning policy not found", map[string]any{"policy_id": policyID})
			return result, nil
		}
		return nil, fmt.Errorf("failed to validate provisioning policy %s: %w", policyID, err)
	}

	// Only validate provisioning type if reserve_percentage is set
	if reservePercentageSet {
		provisioningType := policy.GetProvisioningType()
		if provisioningType != nil {
			provisioningTypeStr := provisioningType.String()
			result.ProvisioningType = provisioningTypeStr

			isFrontline := provisioningTypeStr == "shared" ||
				provisioningTypeStr == "sharedByUser" ||
				provisioningTypeStr == "sharedByEntraGroup"

			if !isFrontline {
				result.InvalidProvisionType = true
				tflog.Warn(ctx, "Policy is not Frontline type", map[string]any{
					"policy_id":         policyID,
					"provisioning_type": provisioningTypeStr,
				})
			}
		}
	}

	tflog.Debug(ctx, "Provisioning policy validated successfully", map[string]any{"policy_id": policyID})
	return result, nil
}
