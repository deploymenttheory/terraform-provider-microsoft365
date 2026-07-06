package graphBetaIOSiPadOSDeviceEnrollmentPolicy

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
// catalog business rules around the ade_useraffinity subtree: the authentication method selectors
// are only meaningful with user affinity enabled, and await final configuration only exists under
// Setup Assistant with modern authentication.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		requireUserAuthenticationForAuthenticationOptions{},
	}
}

// requireUserAuthenticationForAuthenticationOptions enforces:
//   - enable_authentication_via_company_portal and require_setup_assistant_with_modern_authentication
//     may only be true when requires_user_authentication is true.
//   - await_final_configuration may only be true when
//     require_setup_assistant_with_modern_authentication is true.
type requireUserAuthenticationForAuthenticationOptions struct{}

func (v requireUserAuthenticationForAuthenticationOptions) Description(_ context.Context) string {
	return "authentication method options require requires_user_authentication, and await_final_configuration requires require_setup_assistant_with_modern_authentication"
}

func (v requireUserAuthenticationForAuthenticationOptions) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requireUserAuthenticationForAuthenticationOptions) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if req.Config.Raw.IsNull() {
		return
	}

	var requiresUserAuth, companyPortal, modernAuth, awaitFinalConfiguration types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("requires_user_authentication"), &requiresUserAuth)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("enable_authentication_via_company_portal"), &companyPortal)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("require_setup_assistant_with_modern_authentication"), &modernAuth)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("await_final_configuration"), &awaitFinalConfiguration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	companyPortalTrue := !companyPortal.IsUnknown() && !companyPortal.IsNull() && companyPortal.ValueBool()
	modernAuthTrue := !modernAuth.IsUnknown() && !modernAuth.IsNull() && modernAuth.ValueBool()
	awaitTrue := !awaitFinalConfiguration.IsUnknown() && !awaitFinalConfiguration.IsNull() && awaitFinalConfiguration.ValueBool()

	if !requiresUserAuth.IsUnknown() && (requiresUserAuth.IsNull() || !requiresUserAuth.ValueBool()) {
		if companyPortalTrue {
			resp.Diagnostics.AddAttributeError(
				path.Root("enable_authentication_via_company_portal"),
				"requires_user_authentication is required",
				"enable_authentication_via_company_portal may only be true when requires_user_authentication is true.",
			)
		}
		if modernAuthTrue {
			resp.Diagnostics.AddAttributeError(
				path.Root("require_setup_assistant_with_modern_authentication"),
				"requires_user_authentication is required",
				"require_setup_assistant_with_modern_authentication may only be true when requires_user_authentication is true.",
			)
		}
	}

	if awaitTrue && !modernAuth.IsUnknown() && (modernAuth.IsNull() || !modernAuth.ValueBool()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("await_final_configuration"),
			"require_setup_assistant_with_modern_authentication is required",
			"await_final_configuration may only be true when require_setup_assistant_with_modern_authentication is true - "+
				"the underlying ade_modernauth_awaitfinalconfiguration setting only exists under the Setup Assistant with "+
				"modern authentication choice.",
		)
	}
}
