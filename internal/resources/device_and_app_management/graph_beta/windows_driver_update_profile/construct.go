package graphBetaWindowsDriverUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows driver update profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsDriverUpdateProfileResourceModel, forUpdate bool) (graphmodels.WindowsDriverUpdateProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsDriverUpdateProfile()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	// Immutable field once created. Excluded from update req construction.
	if !forUpdate {
		if err := constructors.SetEnumProperty(data.ApprovalType, graphmodels.ParseDriverUpdateProfileApprovalType, requestBody.SetApprovalType); err != nil {
			return nil, fmt.Errorf("invalid approval type: %s", err)
		}
	}

	constructors.SetInt32Property(data.DeploymentDeferralInDays, requestBody.SetDeploymentDeferralInDays)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
