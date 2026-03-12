package graphBetaCrossTenantAccessPartnerUserSyncSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a CrossTenantIdentitySyncPolicyPartner SDK request body.
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-put-identitysynchronization?view=graph-rest-beta
func constructResource(ctx context.Context, data *CrossTenantAccessPartnerUserSyncSettingsResourceModel) (graphmodels.CrossTenantIdentitySyncPolicyPartnerable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewCrossTenantIdentitySyncPolicyPartner()

	convert.FrameworkToGraphString(data.TenantID, requestBody.SetTenantId)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if data.UserSyncInbound != nil {
		userSyncInbound := graphmodels.NewCrossTenantUserSyncInbound()
		convert.FrameworkToGraphBool(data.UserSyncInbound.IsSyncAllowed, userSyncInbound.SetIsSyncAllowed)
		requestBody.SetUserSyncInbound(userSyncInbound)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
