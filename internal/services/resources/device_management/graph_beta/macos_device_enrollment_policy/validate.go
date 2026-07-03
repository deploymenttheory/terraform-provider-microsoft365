package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// intuneProvisioningClientAppID is the App ID for the "Intune Provisioning Client" service
// principal (shown as "Intune Autopilot ConfidentialClient" in some tenants), which must own the
// device_security_group for enrollment time grouping to work.
const intuneProvisioningClientAppID = "f1346770-5b25-470b-88bd-d5744ab7952c"

// validateSecurityGroupOwnership validates that the specified security group has the Intune
// Provisioning Client as an owner, as required for enrollment time grouping.
func validateSecurityGroupOwnership(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, groupID string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Validating security group %s has Intune Provisioning Client as owner", groupID))

	owners, err := client.
		Groups().
		ByGroupId(groupID).
		Owners().
		Get(ctx, nil)
	if err != nil {
		diags.AddError(
			"Failed to validate security group ownership",
			fmt.Sprintf("Could not retrieve owners for security group %s: %s", groupID, err.Error()),
		)
		return diags
	}

	for _, owner := range owners.GetValue() {
		if servicePrincipal, ok := owner.(models.ServicePrincipalable); ok {
			if appID := servicePrincipal.GetAppId(); appID != nil && *appID == intuneProvisioningClientAppID {
				return diags
			}
		}
	}

	diags.AddError(
		"Invalid security group ownership",
		fmt.Sprintf(
			"Security group %s must have the Intune Provisioning Client (AppID: %s) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'.",
			groupID,
			intuneProvisioningClientAppID,
		),
	)
	return diags
}

// ConfigValidators implements resource.ResourceWithConfigValidators. It encodes the settings
// catalog business rule that the ade_accountsettings_createlocaladmin subtree is only ever present
// when ade_macos_awaitconfiguration is enabled, and that ade_accountsettings_prefillaccountinfo is
// only ever present when a primary account is being created.
func (r *MacOSDeviceEnrollmentPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		requireAdminAccountWhenAwaitConfigured{},
	}
}

// requireAdminAccountWhenAwaitConfigured enforces:
//   - admin_account must be set when await_device_configured is true, and must be omitted when it
//     is false.
//   - admin_account.primary_account must be omitted when admin_account.create_local_primary_account
//     is false.
type requireAdminAccountWhenAwaitConfigured struct{}

func (v requireAdminAccountWhenAwaitConfigured) Description(_ context.Context) string {
	return "admin_account must be set if and only if await_device_configured is true, and admin_account.primary_account must be omitted when admin_account.create_local_primary_account is false"
}

func (v requireAdminAccountWhenAwaitConfigured) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requireAdminAccountWhenAwaitConfigured) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if req.Config.Raw.IsNull() {
		return
	}

	var awaitConfigured types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("await_device_configured"), &awaitConfigured)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var adminAccount *AdminAccountModel
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("admin_account"), &adminAccount)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !awaitConfigured.IsUnknown() {
		switch {
		case !awaitConfigured.IsNull() && awaitConfigured.ValueBool() && adminAccount == nil:
			resp.Diagnostics.AddAttributeError(
				path.Root("admin_account"),
				"admin_account is required",
				"admin_account must be configured when await_device_configured is true.",
			)
		case (awaitConfigured.IsNull() || !awaitConfigured.ValueBool()) && adminAccount != nil:
			resp.Diagnostics.AddAttributeError(
				path.Root("admin_account"),
				"admin_account must not be set",
				"admin_account must be omitted when await_device_configured is false.",
			)
		}
	}

	if adminAccount == nil {
		return
	}

	createLocalPrimary := adminAccount.CreateLocalPrimaryAccount
	if !createLocalPrimary.IsUnknown() && (createLocalPrimary.IsNull() || !createLocalPrimary.ValueBool()) &&
		adminAccount.PrimaryAccount != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("admin_account").AtName("primary_account"),
			"primary_account must not be set",
			"admin_account.primary_account must be omitted when admin_account.create_local_primary_account is false.",
		)
	}
}
