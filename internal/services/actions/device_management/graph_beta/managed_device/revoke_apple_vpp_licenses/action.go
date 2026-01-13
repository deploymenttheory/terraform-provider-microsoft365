package graphBetaRevokeAppleVppLicenses

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &RevokeAppleVppLicensesAction{}
	_ action.ActionWithConfigure      = &RevokeAppleVppLicensesAction{}
	_ action.ActionWithValidateConfig = &RevokeAppleVppLicensesAction{}
)

func NewRevokeAppleVppLicensesAction() action.Action {
	return &RevokeAppleVppLicensesAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RevokeAppleVppLicensesAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RevokeAppleVppLicensesAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *RevokeAppleVppLicensesAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *RevokeAppleVppLicensesAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Revokes all Apple Volume Purchase Program (VPP) licenses from devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/revokeAppleVppLicenses` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/revokeAppleVppLicenses` endpoints. " +
			"This action reclaims all VPP-purchased app licenses assigned to iOS/iPadOS devices, making them available " +
			"for reassignment to other devices or users.\n\n" +
			"**What This Action Does:**\n" +
			"- Revokes all VPP app licenses from device\n" +
			"- Returns licenses to available pool\n" +
			"- Makes licenses available for reassignment\n" +
			"- Removes apps from device (if enforced)\n" +
			"- Updates license inventory\n" +
			"- Audits license revocation\n\n" +
			"**When to Use:**\n" +
			"- Device retirement or decommissioning\n" +
			"- Device lost or stolen\n" +
			"- User departure from organization\n" +
			"- License reallocation needed\n" +
			"- Device platform change\n" +
			"- License optimization\n" +
			"- Compliance requirements\n\n" +
			"**Platform Support:**\n" +
			"- **iOS**: Full support (VPP apps)\n" +
			"- **iPadOS**: Full support (VPP apps)\n" +
			"- **Other platforms**: Not applicable (no VPP)\n\n" +
			"**Important Considerations:**\n" +
			"- Only affects VPP-purchased apps\n" +
			"- Device may need to be online\n" +
			"- Apps may be removed from device\n" +
			"- User data typically preserved\n" +
			"- Licenses immediately available\n" +
			"- Cannot be undone easily\n\n" +
			"**Reference:** [Microsoft Graph API - Revoke Apple VPP Licenses](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices fully managed by Intune only. " +
					"Each ID must be a valid GUID format. All VPP licenses will be revoked from these devices. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"You can provide both to revoke licenses from different types of devices in one action.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"comanaged_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of co-managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices managed by both Intune and " +
					"Configuration Manager (SCCM). Each ID must be a valid GUID format. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\"]`\n\n" +
					"**Note:** Co-management is rare for iOS/iPadOS devices but supported by this action. " +
					"At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some operations fail. " +
					"Failed operations will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any operation fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist and are iOS/iPadOS devices before attempting to revoke licenses. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent or non-Apple devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
