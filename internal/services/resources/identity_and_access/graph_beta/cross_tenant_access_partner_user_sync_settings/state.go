package graphBetaCrossTenantAccessPartnerUserSyncSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote CrossTenantIdentitySyncPolicyPartner API response to Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *CrossTenantAccessPartnerUserSyncSettingsResourceModel, remoteResource graphmodels.CrossTenantIdentitySyncPolicyPartnerable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetTenantId()).ValueString(),
	})

	tenantID := convert.GraphToFrameworkString(remoteResource.GetTenantId())
	data.ID = tenantID
	data.TenantID = tenantID
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	if userSyncInbound := remoteResource.GetUserSyncInbound(); userSyncInbound != nil {
		data.UserSyncInbound = &CrossTenantUserSyncInbound{
			IsSyncAllowed: convert.GraphToFrameworkBool(userSyncInbound.GetIsSyncAllowed()),
		}
	} else {
		data.UserSyncInbound = nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
