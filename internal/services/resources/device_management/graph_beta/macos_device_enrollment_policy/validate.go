package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
